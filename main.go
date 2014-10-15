package main

import (
    "os"
    "io/ioutil"
    "fmt"
    spec "github.com/Sam-Izdat/pogo/gtspec"
    "github.com/Sam-Izdat/pogo/po"
    "github.com/Sam-Izdat/pogo/deps/odin/cli"
    "time"
    "strings"
    "runtime"
    "errors"
)

var (
    o spec.Config
    CLI = cli.New("0.0.3", "pogo command line utility", exec)
    ps = string(os.PathSeparator)
    cmdInit, cmdBuild *cli.SubCommand
)

func init() {
    cmdInit  = CLI.DefineSubCommand("init", "initialize pogo in this directory", pinit)
    cmdBuild = CLI.DefineSubCommand("build", "scan source and compile .pot or .po files", build, "filetype")
    cmdBuild.DefineBoolFlag("overwrite", false, "overwrite existing files")
    cmdBuild.AliasFlag('o', "overwrite")
}

func main() {
    CLI.Start()
}

func exec(c cli.Command) {
    fmt.Println(pWarn, "missing command \n", `See "pogo --help" for usage information`)
    os.Exit(1)
}


func pinit(c cli.Command) {
    wd, err := os.Getwd()
    clid := getCallerDir()
    fnFrom, fnTo := "tmpl.toml", "POGO.toml"

    if _, err := os.Stat(wd+ps+fnTo); err == nil {
        fmt.Println(pNotice, 
            fStr("aborting").s("bold"), "- configuration file already exists \n", wd+ps+fnTo)
        os.Exit(1)
    }

    tmpl, err := ioutil.ReadFile(clid+ps+fnFrom)
    if err != nil {
        fmt.Println(pWarn, "could not read default configuration file:\n", err)
        os.Exit(1)
    }
    err = ioutil.WriteFile(wd+ps+fnTo, tmpl, 0755)
    if err != nil {
        fmt.Println(pWarn, "could not write configuration file:\n", err)
        os.Exit(1)
    }        
    fmt.Println(pSuccess, "default configuration file created - edit POGO.toml to configure")
    os.Exit(1)
}

func build(c cli.Command) {
    loadOptions()
    if c.Params()["filetype"] == nil {
        fmt.Println("\n", pWarn, "missing filetype parameter \n", 
            `Specify what to build ("pot"/"po") - e.g. "pogo build pot"`)
        os.Exit(1)
    }

    switch fmt.Sprintf("%s", c.Param("filetype")) {
    case "pot":    
        verifyLocaleDir()

        // Verify no file or active overwrite flag
        fn := o.General.DirLocale+ps+o.General.ProjectFN+".pot"
        if _, err := os.Stat(fn); err == nil && c.Flag("overwrite").Get() == false {
            fmt.Println(pNotice, fStr("skipping").s("bold"),
                `template - file already exists; use "-o" flag to overwrite`)
            os.Exit(1)
        }

        // Scan
        defer un(trace("scan/build"))
        fmt.Println("Parsing...")
        pdir := o.General.DirProject
        msgs := []spec.Msg{}
        msgs = append(po.ScanGo(pdir), po.ScanTmpl(pdir)...)
        po.RemoveDuplicates(&msgs)
        fmt.Println(len(msgs), "unique message(s) extracted")
        fmt.Println("Compiling", o.General.ProjectFN+".pot...")

        // Write
        err := WritePOT(msgs)
        if err != nil {
            fmt.Println(pWarn, err)
        }
    case "po":
        verifyLocaleDir()

        defer un(trace("scan/build"))
        fmt.Println("Parsing...")
        pdir := o.General.DirProject
        msgs := []spec.Msg{}
        msgs = append(po.ScanGo(pdir), po.ScanTmpl(pdir)...)
        po.RemoveDuplicates(&msgs)
        fmt.Println(len(msgs), "unique message(s) extracted")

        if c.Flag("overwrite").Get() == true {
            fmt.Println(pNotice, fStr("WARNING!").s("bold"))
            fmt.Println(fStr("OVERWRITING EXISTING CATALOGS").s("red"))
            fmt.Println("You will be asked to confirm overwrites individually.")
            fmt.Println("This is a destructive and not an additive process.")
            fmt.Println("Any completed translations in these files will be", fStr("LOST!").s("bold"))
        }
        
        for _, target := range o.General.Targets {
            path := o.General.DirLocale+ps+target+ps+o.General.DirMessages
            dir, err := os.Stat(path)
            if err != nil || !dir.IsDir(){
                err = os.MkdirAll(path, os.FileMode(0755))
                if err != nil {
                    fmt.Println(pWarn, fStr("skipping").s("bold"), target,
                        `could not create locale directory`)
                }
            }
            file := o.General.ProjectFN+"."+target+".po"
            fn := path+ps+file
            if _, err := os.Stat(fn); err == nil {
                if c.Flag("overwrite").Get() == true {
                    var confirm string
                    fmt.Print("Are you sure you want to overwrite ", file, "? (y/N) : ")
                    fmt.Scanln(&confirm)
                    if confirm != "y" && confirm != "Y" {
                        continue
                    }
                } else {
                    fmt.Println(pNotice, fStr("skipping").s("bold"), file,
                        `- file already exists; use "-o" flag to overwrite`)
                    continue                    
                }
            }
            fmt.Println("Compiling", file+"...")

            // Write
            err = WritePO(msgs, target, path)
            if err != nil {
                fmt.Println(pWarn, "ERROR compiling", file, "-", err)
            }
        }

    default: 
        fmt.Println(pWarn, "invalid filetype parameter \n", 
            `Expecting "pot" or "po" - e.g. "pogo build pot"`)
        os.Exit(1)
    }
}

func loadOptions() {
    var err error
    o, err = spec.LoadOptions()
    if err != nil {
        fmt.Println(pWarn, "could not locate or parse POGO.toml configuration file \n" + 
            `Run "pogo init" to intialize i18n in this directory`)
        os.Exit(1)
    }
}

func verifyLocaleDir() {
    dir, err := os.Stat(o.General.DirLocale)
    if err != nil || !dir.IsDir(){
        fmt.Println(pWarn, 
            "locale directory does not exist")
        os.Exit(1)
    }
}

func WritePOT(msgs []spec.Msg) error {
    pofile := po.Compile(msgs, "", "", "")    

    // open output file
    fn := o.General.DirLocale+ps+o.General.ProjectFN+".pot"
    fo, err := os.Create(fn)
    if err != nil {
        return err
    }

    // close fo on exit and check for its returned error
    defer func() error {
        if err := fo.Close(); err != nil {
            return err
        }
        return nil
    }()

    // write a chunk
    if _, err := fo.Write([]byte(pofile)); err != nil {
        return err
    }
    return nil
}

func WritePO(msgs []spec.Msg, target string, path string) error {
    prule := spec.Plurals[target]
    if prule == nil {
        prule = spec.Plurals[strings.Split(target, "_")[0]] 
        if prule == nil {
            return errors.New("unknown locale")
        }
    }

    name := prule.Name()
    pf := prule.Header()
    pofile := po.Compile(msgs, target, name, pf)
   
    // open output file
    fn := path+ps+o.General.ProjectFN+"."+target+".po"
    fo, err := os.Create(fn)
    if err != nil {
        return err
    }

    // close fo on exit and check for its returned error
    defer func() error {
        if err := fo.Close(); err != nil {
            return err
        }
        return nil
    }()

    // write a chunk
    if _, err := fo.Write([]byte(pofile)); err != nil {
        return err
    }
    return nil
}

func getCallerDir() string {
    _, fn, _, _ := runtime.Caller(1)
    tmp := strings.Split(fn, "/") // always "/" -- do NOT use OS-specific path separator
    dir := strings.Join(tmp[:len(tmp)-1], ps)
    return dir
}

func trace(s string) (string, time.Time) {
    fmt.Println(fStr("Starting").s("green"), fStr(s).s("green"))
    return s, time.Now()
}

func un(s string, startTime time.Time) {
    endTime := time.Now()
    fmt.Println(
        fStr("Ending").s("green"), 
        fStr(s).s("green"), 
        fStr("- took").s("green"), 
        fStr(endTime.Sub(startTime).String()).s("green"),
    )
}

func isWindows() bool { // <!--[if IE 6]> lalalalalala
    return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}