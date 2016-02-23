package lazyblog

import (
	"html/template"
	"log"
	"net/http"
)

var t = template.Must(template.ParseGlob("templates/*"))

// IndexHandler serves the home page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "index", map[string]string{"Title": "Lazy Blog"})
	if err != nil {
		log.Fatalln("Couldn't render template for home page!")
	}
}
