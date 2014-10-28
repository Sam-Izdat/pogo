package main

import (
	"fmt"
	"github.com/Sam-Izdat/pogo/translate"
	"html/template"
	"net/http"
	"path/filepath"
)

var POGO = translate.LoadCfg("github.com/Sam-Izdat/pogo/example")

func handler(w http.ResponseWriter, r *http.Request) {
	var T = POGO.New("ru")

	var bottles []int
	for i := 99; i >= 0; i-- {
		bottles = append(bottles, i)
	}

	data := struct {
		T       translate.Translator
		Title   string
		Bottles []int
	}{
		T,
		T.G("Internationalization Example"),
		bottles,
	}

	lp := filepath.Join("views", "layout.html")
	fp := filepath.Join("views", "index.html")
	templates := template.Must(template.ParseFiles(lp, fp))
	templates.ExecuteTemplate(w, "layout", data)
}

func main() {
	port := ":8383"
	fmt.Println("Serving on port", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
