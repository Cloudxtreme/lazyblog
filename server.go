package lazyblog

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	templatePath = "/Users/bentranter/Go/src/github.com/bentranter/lazyblog/cmd/templates/*"
	assetPath    = "/Users/bentranter/Go/src/github.com/bentranter/lazyblog/cmd/assets"
)

var (
	// Router is the router for our application.
	Router = NewDefaultMux()

	t = template.Must(template.ParseGlob(templatePath))
)

// HomeHandler serves the home page.
func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts := GetAll()
	err := t.ExecuteTemplate(w, "index", posts)
	if err != nil {
		log.Println("Couldn't render template for home page!", err)
	}
}

// GetPostHandler returns a post with the given id.
func GetPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	buf := GetPost(id)
	_, err := w.Write(buf)
	if err != nil {
		log.Printf("Error rendering for id %s: %s\n", id, err.Error())
	}
}

// NewPostHandler shows the page that allows you to create a new post.
func NewPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := t.ExecuteTemplate(w, "new", nil)
	if err != nil {
		log.Println("Couldn't render template for home page!", err)
	}
}

// NewPostSubmitHandler handles the post submission.
func NewPostSubmitHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	title := r.FormValue("title")
	post := &PostJSON{
		ID:          Urlify(title) + NewID(),
		Title:       title,
		Body:        r.FormValue("body"),
		DateCreated: time.Now(),
	}
	SetPost(w, post)
}

type httprouterHandler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

// AuthenticatedRoute protects the route
func AuthenticatedRoute(next httprouterHandler) httprouter.Handle {
	// check if user is authenticated
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		next(w, r, ps)
	})
}

// AdminHandler serves the admin page.
func AdminHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// NewDefaultMux returns the router with its routes already initialized.
func NewDefaultMux() *httprouter.Router {
	// Create a new serve mux
	r := httprouter.New()

	// Routes
	r.GET("/", HomeHandler)
	r.GET("/new", NewPostHandler)
	r.GET("/posts/:id", GetPostHandler)

	r.POST("/new", NewPostSubmitHandler)
	r.GET("/what", AuthenticatedRoute(GetPostHandler)) // it works y'all

	// Server static files
	r.ServeFiles("/assets/*filepath", http.Dir(assetPath))

	return r
}
