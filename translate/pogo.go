package translate

import (
	"bytes"
	"fmt"
	gt "github.com/Sam-Izdat/pogo/deps/gettext"
	spec "github.com/Sam-Izdat/pogo/gtspec"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Translator delivers translation methods for a particular locale.
// It is exported for reference, but the New constructor should be
// used to initialize every translator.
type Translator struct {
	Locale string
	Ctrl   POGOCtrl
}

type collection map[string]gt.Catalog

// POGOCtrl is a configured handler for constructing translators
type POGOCtrl struct {
	o        spec.Config
	Catalogs collection
}

var LangDefault string
var LangsSupported = map[string]bool{}

// LoadCfg takes the path of the project directory
// (relative to $GOPATH/src/) containing the POGO.toml
// configuration file and loads the configuration variables.
// Normally, this will be the main directory of your package.
func LoadCfg(path string) POGOCtrl {
	var err error
	o, err := spec.LoadOptionsGOPATH(path)
	if err != nil {
		panic(err)
	}
	LangDefault = "UNSUPPORTED"
	for _, v := range o.General.Targets {
		LangsSupported[v] = true
	}
	return POGOCtrl{o, make(collection)}
}

// New takes a locale string and creates a new translator
func (p POGOCtrl) New(locale string) Translator {
	if p.o.General.ProjectFN == "" {
		panic("no pogo configuration loaded")
	}
	if supported, ok := LangsSupported[locale]; !ok || !supported {
		locale = LangDefault
	} else {
		p.readMo(locale)
	}
	return Translator{locale, p}
}

// NewQV takes a slice of locale strings, sorted by quality value and creates
// the best available translator, falling back on default (first) language if
// no match is found.
func (p POGOCtrl) NewQV(locales []string) Translator {
	if p.o.General.ProjectFN == "" {
		panic("no pogo configuration loaded")
	}
	var locale string
	for _, v := range locales {
		if supported, ok := LangsSupported[v]; ok && supported {
			locale = v
			break
		}
	}
	if locale == "" {
		locale = LangDefault
	}
	p.readMo(locale)
	return Translator{locale, p}
}

func (p *POGOCtrl) readMo(locale string) {
	if _, ok := p.Catalogs[locale]; ok {
		return
	}
	fn := strings.Join([]string{p.o.General.ProjectFN, ".", locale, ".mo"}, "")
	path := filepath.Join(p.o.General.DirLocale, locale, p.o.General.DirMessages, fn)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	c := gt.NewCatalog()
	if err := c.ReadMo(bytes.NewReader(data)); err != nil {
		panic(err)
	}
	p.Catalogs[locale] = *c
}

// G translates a string. The first argument must be
// the string to be translated. Any subsequent arguments
// will be translated, if possible, and considered arguments
// for sprintf.
func (t Translator) G(input ...interface{}) string {
	if len(input) < 1 {
		return ""
	}
	id := input[0].(string)

	// translate chain of remaining string inputs
	for k, v := range input[1:] {
		if s, ok := v.(string); ok {
			input[k+1] = t.G(s)
		}
	}

	c := t.Ctrl.Catalogs[t.Locale]
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
func (t Translator) NG(input ...interface{}) string {
	if len(input) < 3 {
		return ""
	}
	ct := input[len(input)-1].(int)

	// translate chain of remaining string inputs
	for k, v := range input[2:] {
		if s, ok := v.(string); ok {
			input[k+2] = t.G(s)
		}
	}

	idx, err := spec.GetPluralIdx(t.Locale, ct)
	if err == nil {
		c := t.Ctrl.Catalogs[t.Locale]
		if msg, ok := c.Msgs[input[0].(string)]; ok {
			if text := msg.StrPlural[idx]; text != nil {
				if len(input) < 3 {
					return fmt.Sprintf(string(text), ct)
				}
				return fmt.Sprintf(string(text), input[2:]...)
			}
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
func (t Translator) PG(input ...interface{}) string {
	if len(input) < 2 {
		return ""
	}

	// translate chain of remaining string inputs
	for k, v := range input[2:] {
		if s, ok := v.(string); ok {
			input[k+2] = t.G(s)
		}
	}

	c := t.Ctrl.Catalogs[t.Locale]
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
func (t Translator) NPG(input ...interface{}) string {
	if len(input) < 4 {
		return ""
	}

	// translate chain of remaining string inputs
	for k, v := range input[3:] {
		if s, ok := v.(string); ok {
			input[k+3] = t.G(s)
		}
	}

	ct := input[len(input)-1].(int)
	idx, err := spec.GetPluralIdx(t.Locale, ct)
	if err == nil {
		c := t.Ctrl.Catalogs[t.Locale]
		key := strings.Join([]string{input[0].(string), "\x04", input[1].(string)}, "")
		if msg, ok := c.Msgs[key]; ok {
			if text := msg.StrPlural[idx]; text != nil {
				if len(input) < 5 {
					return fmt.Sprintf(string(text), ct)
				}
				return fmt.Sprintf(string(text), input[3:]...)
			}
		}
	}

	if ct == 1 {
		return fmt.Sprintf(input[1].(string), input[3:]...)
	} else {
		return fmt.Sprintf(input[2].(string), input[3:]...)
	}
}
