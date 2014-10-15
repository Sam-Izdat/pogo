
![pogo](doc/pogo-logo.png)

...provides gettext-like internationalization for golang 

##Why?

...because most of the major, recurrent challenges in internationalizing software have been addressed by GNU almost twenty years ago and golang is badly lacking in this functionality. Gettext enjoys extremely wide adoption, established toolchains and great popularity among both translators and programmers. It's simple. It provides for better separation of localization from the rest of your architecture, allows non-programmers to translate your software, eliminates ugly, unmanageable JSON files and, when used properly, effectively solves the problems of context and grammatical differences in pluralization. Until something better comes along, it just doesn't make much sense to reinvent the square wheel. This project is meant to be a loose interpretation replicating the most useful functionality and making use of common po/mo standards.

##Features 

- i18n in your templates, not crammed into them as an afterthought
- Simple CLI package to extract strings in lieu of xgettext
- Built-in support for pluralization rules, from Alcholi to Yoruba
- Built-in support for context fields

##Install

    $ go get github.com/Sam-Izdat/pogo
    $ go install github.com/Sam-Izdat/pogo

##Status

It's still very raw. A pre-release is available. Contributions and bug reports are most welcome. If you find this helpful and feel generous you can throw me some loose change.

[![Build Status](http://drone.io/github.com/Sam-Izdat/pogo/status.png)](https://drone.io/github.com/Sam-Izdat/pogo/latest) 
[![License MIT](http://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](http://opensource.org/licenses/MIT)
[![GoDoc](http://img.shields.io/badge/doc-REFERENCE-blue.svg?style=flat-square)](https://godoc.org/github.com/Sam-Izdat/pogo/translate)
[![Donate](http://img.shields.io/badge/needfood-GIVEMEMONEYS-yellow.svg?style=flat-square)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=CY6VMGLZ7XA64)

##Getting started
If your `GOPATH` and `GOBIN` environment variables were set correctly, you should now be able to run `pogo` from any directory. This little CLI package will scan your project for calls pogo's gettext-ish functions and compile string literals into po files. The first thing is to initialize it in the main project directory and set up a configuration file.

    $ cd github.com/MyStuff/MyProject
    $ pogo init

Now, edit the top few settings in the generated POGO.toml file. It should be enough to name your project and give it a list of target locales. That's it. There are no domains, as such. When pogo walks the directory it'll search until it bumps into another POGO.toml file somewhere; if it does, that subdirectory will be ignored and left to another configuration and collection of catalogs. If needed, an application's translation files can be compartmentalized by packages, or just by separate directories of views.

###Putting it to use
Here's a basic webserver that can be found in the example folder.

```pogo
package main
    
import (
    "fmt"
    "net/http"
    "html/template"
    "path/filepath"
    pogo "github.com/Sam-Izdat/pogo/translate"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Grab a translator in your request flow
    var T = pogo.New("ru") // "ru" is Russian

    // Set up some data for the template
    var bottles []int
    for i := 99; i >= 0; i-- { bottles = append(bottles, i) }

    data := struct {
        T pogo.Translator       // Throw a translator at the template
        Title string
        Bottles []int
    } {
        T,                      // The Russian translator above
        T.G("Internationalization Example"), // Translate outside template
        bottles,
    }

    lp := filepath.Join("views", "layout.html")
    fp := filepath.Join("views", "index.html")
    templates := template.Must(template.ParseFiles(lp, fp))
    templates.ExecuteTemplate(w, "layout", data)
}

func main() {
    // Load your POGO.toml configuration file *before* processing the request
    pogo.LoadCfg("github.com/MyStuff/MyProject")

    // Point to a handler and serve
    port := ":8383"
    fmt.Println("Serving on port", port)
    http.HandleFunc("/", handler)
    http.ListenAndServe(port, nil)
}
```
###Translating strings
By default, a pogo "translator" exports four methods called:
- `G()` - for basic translation (roughly equivalent to `gettext()` or `_()`)
- `PG()` - for translation with context (`pgettext()`)
- `NG()`- for translation with quantities (`ngettext()`)
- `NPG()` - for translation with both quantities and context (`npgettext()`)

It will be expected by the CLI program that the above are reserved for pogo in .go files (as function and method names) and in templates (as function, method *and* variable names). If you would like something more verbose with less chance conflict/collision, just alias these methods and edit the POGO.toml file for the scanner, to reflect the new names. 

String literals can be queued up for translation directly in your go files but, chances are, most of the content to be translated will reside in templates. Passing the translator to a template as above now lets you do this:

####G - just translate
```
{{.T.G "A robot?!"}}
```
All of these functions are variadic and will take an arbitrary number of arguments. `G()` will pass along all subsequent arguments to `fmt.Sprintf()`, often eliminating the need to nest it in another function:

```
{{.T.G "I seem to have lost %s" "my shoes"}}
```
The words "my shoes" above will be translated if a translation is available. 

####PG - context
`PG()` allows you to provide a context for the translator, which is passed as the first argument:

```
{{.T.PG "childrens' game" "tag"}}
{{.T.PG "category or identifier" "tag"}}
{{.T.PG "the price on an item" "tag"}}
```

####NG - plurals
`NG()` is for anything that may vary in quantity:
```
{{.T.NG "You have %d new message" "You have %d new messages" .ctMsg}}
```
The first argument must be the singular form of the string to be translated; the second must be the plural. The very *last* argument must be the quantity of the thing, which will be used to pull up a locale-specific pluralization rule. Anything in between will be interpolated.

```
{{.T.NG "I ate %[2]d muffin %[1]s" "I ate %[2]d muffins %[1]s" "yesterday" .ctMuffins}}
```

####NPG - plurals & context
`NPG()` takes at least four arguments, combining the purposes of `PG()` and `NG()`:
```
{{.T.NPG "sternly!" "Delete %d file?" "Delete %d files?" .ctFiles}}
```
The first argument is context, then singular, then plural, then zero or more other arguments, then the quantity.

Use a sigil (`$`) to access `T` while ranging over something or within conditional statements:

```html
{{range .Bottles}}
<p>{{$.T.NG "%[2]d bottle of beer %[1]s!" "%[2]d bottles of beer %[1]s!" "on the wall" (.)}}</p>
{{end}}
```
There's really not much more to it.

###Mos, pos and pots and other things

Now that there's some stuff to be translated you can compile the po files. One file with the extension "pot" will serve as the original template and every locale will have its own "po" file (catalog) containing the actual translations. Editors like Poedit can merge these catalogs with any new messages added to the template. 

    $ cd github.com/MyStuff/MyProject
    $ mkdir locale
    $ pogo build pot

...will compile the template...
        
    $ pogo build po

...will produce individual po files for all your targets with some meta-data already in place. Whenever you have new strings to translate, just run `pogo build -o pot` again. It does roughly what xgettext does. Currently pogo does not compile mo files and leaves that up the fancy editors. 

#On the to-do list

- [ ] Unit tests
- [ ] Better documentation
- [ ] General code cleanup
- [ ] "pogo build mo" CLI command
- [ ] "pogo status" CLI command
- [ ] HTML UI stats, previews

#Attributions
built upon:
- [gorilla web toolkit](https://github.com/gorilla)'s mo reader & writer, template parser
- [odin](https://github.com/jwaldrip/odin) CLI library
- [TOML parser](TOML parser)

#License

MIT for pogo; see dependencies' docs for their respective licenses
