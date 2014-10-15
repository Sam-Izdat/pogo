package gtspec

import (
    "io/ioutil"
    "os"
    "errors"
    "strings"
    "github.com/Sam-Izdat/pogo/deps/toml"
)
// Classifaction of a comment
type CommentSpec struct {
    Key, Prefix string
}

// CommentsPack groups comments keyed by classification (i.e. CommentSpec.key)
type CommentPack map[string][]string

type Msg struct {    
    Ctxt        string          // msgctxt: message context, if any
    Id          string          // msgid: untranslated singular string
    IdPlural    string          // msgid_plural: untranslated plural string
    Str         string          // msgstr: translated singular string
    StrPlural   []string        // msgstr[n]: translated plural strings
    Comments    CommentPack 
    Filename    string          // Name of file extracted from (to be shoved into comments)
    Line        int             // Line number within file (to be shoved into comment)
}

type Config struct {
    General confGeneral
    Parsing confParsing
    Po confPo
    Meta confMeta
}

type confGeneral struct {
    ProjectName string  `toml:"project_name"`
    ProjectFN string    `toml:"project_filename"`
    Targets []string
    DirProject string
    DirLocale string    `toml:"dir_locale"`
    DirMessages string  `toml:"dir_messages"`
}

type confParsing struct {
    TmplExts []string   `toml:"extensions_template"`
    FuncG string        `toml:"function_gettext"`
    FuncNG string       `toml:"function_ngettext"`
    FuncPG string       `toml:"function_pgettext"`
    FuncNPG string      `toml:"function_npgettext"`
    DelimL string       `toml:"delimiter_left"`
    DelimR string       `toml:"delimiter_right"`
}

type confPo struct {
    ReportBugs string   `toml:"report_bugs_to"`
    Comment string
}

type confMeta struct {
    Version string      `toml:"pogo_version"`
    fnTOML string
    fnTOMLTmpl string
}

var ps = string(os.PathSeparator)
var cfgFN = "POGO.toml"

func LoadOptions() (Config, error) {
    var options Config
    path, err := os.Getwd()
    if err != nil { 
        return Config{}, err
    }

    // Traverse up the directory tree until finding config file or hitting root
    var dir string
    c := strings.Split(path, ps[:1])
    for k := range c {
        dir = strings.Join(c[:len(c)-k], ps)
        if _, err := os.Stat(dir+ps+cfgFN); os.IsNotExist(err) { // file does not exist
            continue 
        } else if err == nil { // file exists
            data, err := ioutil.ReadFile(dir+ps+cfgFN)
            if err != nil {
                return Config{}, err
            }
            if _, err := toml.Decode(string(data), &options); err != nil {
                return Config{}, err
            }
            options.General.DirProject = dir
            ldir := strings.Replace(options.General.DirLocale, "/", ps, -1)
            options.General.DirLocale = strings.Replace(ldir, "%PROJECT%", dir, -1)
            return options, nil
        }
    }

    return Config{}, errors.New("config file not found")
}

func LoadOptionsGOPATH(path string) (Config, error) {
    var options Config
    gopath := os.Getenv("GOPATH")
    path = gopath+ps+"src"+ps+path
    data, err := ioutil.ReadFile(path+ps+cfgFN)
    if err != nil {
        return Config{}, err
    }
    if _, err := toml.Decode(string(data), &options); err != nil {
        return Config{}, err
    }
    options.General.DirProject = path
    ldir := strings.Replace(options.General.DirLocale, "/", ps, -1)
    options.General.DirLocale = strings.Replace(ldir, "%PROJECT%", path, -1)
    return options, nil
}

func init() {
}