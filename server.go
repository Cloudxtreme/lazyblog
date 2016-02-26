package lazyblog

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

const (
	templatePath = "/Users/bentranter/Go/src/github.com/bentranter/lazyblog/cmd/templates/*"
	assetPath    = "/Users/bentranter/Go/src/github.com/bentranter/lazyblog/cmd/assets"
)

var (
	// ErrInvalidSigningMethod means that method used to sign the token doesn't
	// match the method stated in the token header.
	ErrInvalidSigningMethod = errors.New("Invalid signing method")

	// ErrInvalidToken means the token isn't valid.
	ErrInvalidToken = errors.New("Invalid token")

	// Router is the router for our application.
	Router = NewDefaultMux()

	t          = template.Must(template.ParseGlob(templatePath))
	signingKey = genRandBytes()
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
	// @TODO
}

// LoginHandler serves the admin login page.
func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := t.ExecuteTemplate(w, "login", nil)
	if err != nil {
		log.Printf("Error rendering login template: ", err.Error())
	}
}

// LoginPostHandler handles the user's login request. If their password is
// incorrect, they're redirected to the login page with a flash message
// informiang them what went wrong. If their password is correct, they're given
// a JSON web token and redirected to the admin page.
func LoginPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword := GetUser(username)

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		// @TODO: should redirect with flash message for v0.1.0
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	// Issue the authenticated user a token
	w.Write([]byte("Good credentials"))
}

// NewDefaultMux returns the router with its routes already initialized.
func NewDefaultMux() *httprouter.Router {
	// Create a new serve mux
	r := httprouter.New()

	// Routes
	r.GET("/", HomeHandler)
	r.GET("/new", NewPostHandler)
	r.GET("/posts/:id", GetPostHandler)
	r.GET("/admin/login", LoginHandler)

	r.POST("/admin/login", LoginPostHandler)

	// Authenticated routes
	r.POST("/new", NewPostSubmitHandler)

	// Server static files
	r.ServeFiles("/assets/*filepath", http.Dir(assetPath))

	return r
}

func genToken() (string, error) {
	tok := jwt.New(jwt.SigningMethodHS256)

	tok.Claims["sub"] = "admin"
	tok.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tok.Claims["iat"] = time.Now().Unix()

	return tok.SignedString(signingKey)
}

func verifyToken(tokStr string) error {
	tok, err := jwt.Parse(tokStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return signingKey, nil
	})
	if err != nil {
		return err
	}

	if !tok.Valid {
		return ErrInvalidToken
	}

	return nil
}

func genRandBytes() []byte {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return []byte(base64.URLEncoding.EncodeToString(b))
}
