package po

import (
	prsTmpl "github.com/Sam-Izdat/pogo/deps/template/parse"
	spec "github.com/Sam-Izdat/pogo/gtspec"
	"go/ast"
	prsGo "go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var gf, ngf, pgf, npgf, lDelim, rDelim string
var cfgFN = "POGO.toml"

func init() {
	gf = o.Parsing.FuncG     // "gettext" function
	ngf = o.Parsing.FuncNG   // "ngettext" function
	pgf = o.Parsing.FuncPG   // "pgettext" function
	npgf = o.Parsing.FuncNPG // "npgettext" function
	lDelim = o.Parsing.DelimL
	rDelim = o.Parsing.DelimR
}

func ScanGo(path string) (res []spec.Msg) {
	filepath.Walk(path, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !!fi.IsDir() {
			switch fi.Name() {
			case ".", "..":
				return nil
			default:
				tmp := scanGoDir(fp)
				if tmp != nil {
					res = append(res, tmp...)
				}
			}
		}
		return nil
	})
	prepMsg(&res)
	return
}

func scanGoDir(path string) (res []spec.Msg) {
	// skip any subdirectory with its own POGO.toml config
	conf := filepath.Join(path, cfgFN)
	if _, err := os.Stat(conf); err == nil && path != o.General.DirProject {
		return nil
	}

	fset := token.NewFileSet()
	pkgs, err := prsGo.ParseDir(fset, path, nil, 0)
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		for fn, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				var funcName string
				switch x := n.(type) {
				case *ast.CallExpr:
					switch y := x.Fun.(type) {
					case *ast.Ident: // function call
						funcName = y.Name
					case *ast.SelectorExpr: // method call
						funcName = y.Sel.Name
					}
					switch funcName {
					case gf, ngf, pgf, npgf: // do nothing
					default:
						return true
					}
					msgs := []spec.Msg{}
					for k, arg := range x.Args {
						switch y := arg.(type) {
						case *ast.BasicLit:
							linePos := fset.Position(y.ValuePos).Line
							switch funcName {
							case gf: // just singular
								msgs = append(msgs, spec.Msg{Filename: fn, Line: linePos, Id: y.Value})
							case ngf: // singular and plural
								switch k { // 0 - singular, 1 - plural, subsequent ignored
								case 0:
									msgs = append(msgs, spec.Msg{Filename: fn, Line: linePos, Id: y.Value})
								case 1:
									msgs[0].IdPlural = y.Value
								}
							case pgf: // context and text
								switch k { // 0 - context, 1 - text, subsequent ignored
								case 0:
									msgs = append(msgs, spec.Msg{Filename: fn, Line: linePos, Ctxt: y.Value})
								case 1:
									msgs[0].Id = y.Value
								}
							case npgf:
								switch k { // 0 - context, 1 - singular, 2 - plural, subsequent ignored
								case 0:
									msgs = append(msgs, spec.Msg{Filename: fn, Line: linePos, Ctxt: y.Value})
								case 1:
									msgs[0].Id = y.Value
								case 2:
									msgs[0].IdPlural = y.Value
								}
							}
						}
					}
					if len(msgs) > 0 {
						res = append(res, msgs...)
					}
				}
				return true
			})
		}
	}
	return
}

func ScanTmpl(path string) (res []spec.Msg) {
	filepath.Walk(path, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// skip any subdirectory with its own POGO.toml config
		dir := filepath.Base(fp)
		conf := filepath.Join(dir, cfgFN)
		if _, err := os.Stat(conf); err == nil && dir != o.General.DirProject {
			return nil
		}

		if !!fi.IsDir() {
			return nil
		}

		matched := false
		for _, ext := range o.Parsing.TmplExts {
			var err error
			matched, err = filepath.Match("*."+ext, fi.Name())
			if err != nil {
				panic(err)
			}
			if matched {
				break
			}
		}

		if matched {
			buf, err := ioutil.ReadFile(fp)
			if err != nil {
				panic(err)
			}
			tmpl := string(buf)
			scn := scanTmplString(tmpl)
			for k := range scn {
				scn[k].Filename = fp
			}
			res = append(res, scn...)
		}
		return nil
	})
	prepMsg(&res)
	return
}

func scanTmplString(s string) (res []spec.Msg) {
	t, err := prsTmpl.Parse("p", s, lDelim, rDelim, map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	for _, va := range t {
		res = append(res, scanNodes(va.List.Nodes)...)
	}
	return
}

func scanNodes(nodes []prsTmpl.Node) (res []spec.Msg) {
	for _, node := range nodes {
		nt := node.Type().Type()
		switch nt {
		case prsTmpl.NodeAction:
			an := node.(*prsTmpl.ActionNode)
			res = append(res, scanActionNode(an)...)
		case prsTmpl.NodePipe:
			pn := node.(*prsTmpl.PipeNode)
			res = append(res, scanPipeNode(pn)...)
		case prsTmpl.NodeRange:
			rn := node.(*prsTmpl.RangeNode)
			res = append(res, scanNodes(rn.List.Nodes)...)
		case prsTmpl.NodeIf:
			in := node.(*prsTmpl.IfNode)
			res = append(res, scanNodes(in.List.Nodes)...)
			if in.ElseList != nil {
				res = append(res, scanNodes(in.ElseList.Nodes)...)
			}
		default:
			continue
		}
	}
	return
}

func scanActionNode(an *prsTmpl.ActionNode) []spec.Msg {
	pn := an.Pipe
	return scanPipeNode(pn)
}

func scanPipeNode(pn *prsTmpl.PipeNode) (res []spec.Msg) {
	cmds := pn.Cmds    // PipeNode.[]CommandNode
	linePos := pn.Line // line position
	for _, cmd := range cmds {
		ok, fun, msgs := false, "", []spec.Msg{}
		for k, arg := range cmd.Args {
			switch arg.Type() {
			case prsTmpl.NodeField, prsTmpl.NodeIdentifier, prsTmpl.NodeVariable: // func, method, var
				ok = true // allow for loop to roll through
				call := strings.Split(arg.String(), ".")
				switch call[len(call)-1] { // lock in for subsequent passes
				case gf:
					fun = "G" // "gettext"
				case ngf:
					fun = "NG" // "ngettext"
				case pgf:
					fun = "PG" // "ngettext"
				case npgf:
					fun = "NPG" // "npgettext"
				default:
					fun = ""
				}
			case prsTmpl.NodeString:
				ok = true
			}
			if !ok {
				continue
			}
			switch arg.Type() {
			case prsTmpl.NodePipe:
				msgs = append(msgs, scanPipeNode(arg.(*prsTmpl.PipeNode))...)
			case prsTmpl.NodeString:
				switch fun {
				case "G":
					msgs = append(msgs, spec.Msg{Line: linePos, Id: arg.String()})
				case "NG":
					switch k { // 1 - singular, 2 - plural, subsequent ignored
					case 1:
						msgs = append(msgs, spec.Msg{Line: linePos, Id: arg.String()})
					case 2:
						msgs[0].IdPlural = arg.String()
					default:
						msgs = append(msgs, spec.Msg{Line: linePos, Id: arg.String()})
					}
				case "PG":
					switch k { // 1 - context, 2 - text, subsequent ignored
					case 1:
						msgs = append(msgs, spec.Msg{Line: linePos, Ctxt: arg.String()})
					case 2:
						msgs[0].Id = arg.String()
					default:
						msgs = append(msgs, spec.Msg{Line: linePos, Id: arg.String()})
					}
				case "NPG":
					switch k { // 1 - context, 2 - singular, 3 - plural, subsequent ignored
					case 1:
						msgs = append(msgs, spec.Msg{Line: linePos, Ctxt: arg.String()})
					case 2:
						msgs[0].Id = arg.String()
					case 3:
						msgs[0].IdPlural = arg.String()
					default:
						msgs = append(msgs, spec.Msg{Line: linePos, Id: arg.String()})
					}
				}
			}
		}
		res = append(res, msgs...)
	}
	return
}

func prepMsg(msgs *[]spec.Msg) {
	ps := string(os.PathSeparator)
	for k, v := range *msgs {
		// prep meta
		(*msgs)[k].Comments = make(spec.CommentPack)
		ref := strings.Join(strings.Split(v.Filename, o.General.DirProject+ps), "")
		ref += ":" + strconv.Itoa(v.Line)
		(*msgs)[k].Comments["reference"] = append(v.Comments["reference"], ref)

		// escape tilde-`quoted` strings (note that v is assigned to)
		switch {
		case len(v.Ctxt) > 0 && v.Ctxt[:1] == "`":
			v.Ctxt = escapeString(v.Ctxt)
		case len(v.Id) > 0 && v.Id[:1] == "`":
			v.Id = escapeString(v.Id)
		case len(v.IdPlural) > 0 && v.IdPlural[:1] == "`":
			v.IdPlural = escapeString(v.IdPlural)
		}

		// remove quotes
		if len(v.Ctxt) > 0 {
			(*msgs)[k].Ctxt = v.Ctxt[1 : len(v.Ctxt)-1]
		}
		if len(v.Id) > 0 {
			(*msgs)[k].Id = v.Id[1 : len(v.Id)-1]
		}
		if len(v.IdPlural) > 0 {
			(*msgs)[k].IdPlural = v.IdPlural[1 : len(v.IdPlural)-1]
		}
	}
}

func escapeString(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\t", `\t`, "\r", `\r`, "\n", `\n`)
	return r.Replace(s)
}
