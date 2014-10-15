package main
    
import (
    "fmt"
    "net/http"
    "html/template"
    "path/filepath"
    pogo "github.com/Sam-Izdat/pogo/translate"
)

func handler(w http.ResponseWriter, r *http.Request) {
    var T = pogo.New("ru")

    var bottles []int
    for i := 99; i >= 0; i-- { bottles = append(bottles, i) }

    data := struct {
        T pogo.Translator
        Title string
        Bottles []int
    } {
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
    pogo.LoadCfg("github.com/Sam-Izdat/pogo/example")
    port := ":8383"
    fmt.Println("Serving on port", port)
    http.HandleFunc("/", handler)
    http.ListenAndServe(port, nil)
}