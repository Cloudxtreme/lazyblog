package lazyblog

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var t = template.Must(template.ParseGlob("templates/*"))

// IndexHandler serves the home page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	posts := GetAll()
	err := t.ExecuteTemplate(w, "index", posts)
	if err != nil {
		log.Println("Couldn't render template for home page!", err)
	}
}

// GetPostHandler Test a rino
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimLeft(r.URL.Path, "/posts")
	w.Write([]byte(id))
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := t.ExecuteTemplate(w, "new", nil)
		if err != nil {
			log.Println("Couldn't render template for home page!", err)
		}
	case "POST":
		r.ParseForm()
		post := &Post{
			ID:   []byte(r.FormValue("id")),
			Body: []byte(r.FormValue("body")),
		}
		SetPost(w, post)
	}
}
