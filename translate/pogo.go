package translate

import (
    "fmt"
    "bytes"
    "strings"
    "io/ioutil"
    "path/filepath"
    spec "github.com/Sam-Izdat/pogo/gtspec"
    gt "github.com/Sam-Izdat/pogo/deps/gettext"
)

// Translator creates locale-bound instances used for translation methods.
// It is exported for reference, but the translate.New constructor should be 
// used to initialize every translator.
type Translator struct {
    Locale string
}

type collection map[string]gt.Catalog
var Catalogs collection

var o spec.Config

// New takes a locale and creates a new translator
func New(locale string) Translator {
    if o.General.ProjectFN == "" {
        panic("no pogo configuration loaded")
    }
    readMo(locale)
    return Translator{locale}
}

// LoadCfg takes the path of the project directory 
// (relative to $GOPATH/src/) containing the POGO.toml 
// configuration file and loads the configuration variables. 
// Normally, this will be the main directory of your package.
func LoadCfg(path string) {
    var err error
    o, err = spec.LoadOptionsGOPATH(path)
    if err != nil {
        panic(err)
    }
}

func readMo(locale string) {
    if _, ok := Catalogs[locale]; ok { return }
    fn := strings.Join([]string{o.General.ProjectFN, ".", locale, ".mo"}, "")
    path := filepath.Join(o.General.DirLocale, locale, o.General.DirMessages, fn)
    data, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }
    c := gt.NewCatalog()
    if err := c.ReadMo(bytes.NewReader(data)); err != nil {
        panic(err)
    }
    Catalogs["ru"] = *c
}

func init() {
    Catalogs = make(collection)
}

// G translates a string. The first argument must be 
// the string to be translated. Any subsequent arguments 
// will be translated, if possible, and considered arguments
// for sprintf.
func (p Translator) G(input ...interface{}) string {
    if len(input) < 1 { return "" }
    id := input[0].(string)

    // translate chain of remaining string inputs
    for k, v := range input[1:] {
        if s, ok := v.(string); ok {
            input[k+1] = p.G(s)
        }
    }

    c := Catalogs[p.Locale]
    if msg, ok := c.Msgs[id]; ok {
        if text := msg.Str; text != nil {
            if len(input) < 2 {
                return string(text)
            }
            return fmt.Sprintf(string(text), input[1:]...)
        }
    }

    if len(input) == 1 {
        return id
    } else {
        return fmt.Sprintf(input[0].(string), input[1:]...)
    }
}

// NG translates and pluralizes a string, according to 
// a language's pluralization rules, if the quantity 
// calls for a plural form. The first argument must 
// be the singular form of the string to be translated; 
// the second must be the plural form; the *last* argument
// must be the quantity; any other arguments preceding it 
// will be translated, if possible, and considered arguments 
// for sprintf.
func (p Translator) NG(input ...interface{}) string {
    if len(input) < 3 { return "" }
    ct := input[len(input)-1].(int)

    // translate chain of remaining string inputs
    for k, v := range input[2:] {
        if s, ok := v.(string); ok {
            input[k+2] = p.G(s)
        }
    }

    c := Catalogs[p.Locale]
    idx := spec.Plurals[p.Locale].Idx(ct)
    if msg, ok := c.Msgs[input[0].(string)]; ok {
        if text := msg.StrPlural[idx]; text != nil {
            if len(input) < 3 {
                return fmt.Sprintf(string(text), ct)
            }
            return fmt.Sprintf(string(text), input[2:]...)
        }
    }

    if ct == 1 { 
        return fmt.Sprintf(input[0].(string), input[2:]...) 
    } else { 
        return fmt.Sprintf(input[1].(string), input[2:]...) 
    }
}

// PG translates a string with context/disambiguation. The first argument 
// must be the context; the second must be the string to be translated. 
// Any subsequent arguments will be translated, if possible, 
// and considered arguments for sprintf.
func (p Translator) PG(input ...interface{}) string {
    if len(input) < 2 { return "" }

    // translate chain of remaining string inputs
    for k, v := range input[2:] {
        if s, ok := v.(string); ok {
            input[k+2] = p.G(s)
        }
    }

    c := Catalogs[p.Locale]
    key := strings.Join([]string{input[0].(string), "\x04", input[1].(string)}, "")
    if msg, ok := c.Msgs[key]; ok {
        if text := msg.Str; text != nil {
            if len(input) < 3 {
                return string(text)
            }
            return fmt.Sprintf(string(text), input[2:]...)
        }
    }
    return fmt.Sprintf(input[1].(string), input[2:]...)
}

// NPG translates a string with context, and pluralizes it,
// according to a language's pluralization rules, if the 
// quantity calls for a plural form. The first argument must 
// be the context, the second must be the singular form 
// of the string to be translated; the third must be the plural form;
// the *last* argument must be the quantity; any other arguments 
// preceding it will be translated, if possible, and considered 
// arguments for sprintf.
func (p Translator) NPG(input ...interface{}) string {
    if len(input) < 4 { return "" }

    // translate chain of remaining string inputs
    for k, v := range input[3:] {
        if s, ok := v.(string); ok {
            input[k+3] = p.G(s)
        }
    }

    ct := input[len(input)-1].(int)

    c := Catalogs[p.Locale]
    idx := spec.Plurals[p.Locale].Idx(ct)
    key := strings.Join([]string{input[0].(string), "\x04", input[1].(string)}, "")
    if msg, ok := c.Msgs[key]; ok {
        if text := msg.StrPlural[idx]; text != nil {
            if len(input) < 5 {
                return fmt.Sprintf(string(text), ct)
            }
            return fmt.Sprintf(string(text), input[3:]...)
        }
    }
    if ct == 1 { 
        return fmt.Sprintf(input[1].(string), input[3:]...) 
    } else { 
        return fmt.Sprintf(input[2].(string), input[3:]...) 
    }
}
