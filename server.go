package lazyblog

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidSigningMethod means that method used to sign the token doesn't
	// match the method stated in the token header.
	ErrInvalidSigningMethod = errors.New("Invalid signing method")

	// ErrInvalidToken means the token isn't valid.
	ErrInvalidToken = errors.New("Invalid token")

	// ErrExpiredToken means that token has expired.
	ErrExpiredToken = errors.New("Expired token")

	// Router is the router for our application.
	Router = NewDefaultMux()

	templatePath = os.Getenv("GOPATH") + "/src/github.com/bentranter/lazyblog/cmd/layout/*"
	assetPath    = os.Getenv("GOPATH") + "/src/github.com/bentranter/lazyblog/cmd/static/*"
	t            = template.Must(template.ParseGlob(templatePath))
	signingKey   = genRandBytes()
	cookieName   = "_lazyblog_token"
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

// APIGetPostHandler returns the post as JSON.
func APIGetPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("id")
	w.Write(GetPostForAPI(id))
}

// AdminGetPostHandler returns the post as JSON to the admin panel so that post
// may be edited.
func AdminGetPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var post PostJSON
	postJSON := GetPostForAPI(ps.ByName("id"))
	err := json.Unmarshal(postJSON, &post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t.ExecuteTemplate(w, "edit", post)
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
	id := NewID()
	post := &PostJSON{
		ID:          id,
		Path:        Urlify(title) + id,
		Title:       title,
		Body:        r.FormValue("body"),
		DateCreated: time.Now(),
	}
	SetPost(w, post)
}

// EditPostSubmitHandler handles editing existing posts.
func EditPostSubmitHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	post := &PostJSON{
		ID:          r.FormValue("id"),
		Path:        r.FormValue("path"),
		Title:       r.FormValue("title"),
		Body:        r.FormValue("body"),
		DateCreated: time.Now(), // should be switched to lastUpdated
	}
	SetPost(w, post)
}

// DeletePostSubmitHandler deletes an existing post.
func DeletePostSubmitHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	err := DeletePost(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

type httprouterHandler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

// AuthenticatedRoute protects the route.
func AuthenticatedRoute(next httprouterHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			// If there isn't a cookie, send the user to the sign in page. We
			// should probably tell the user that they must be signed in to
			// view anything under /admin.
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		err = verifyToken(cookie.Value)
		if err != nil {
			// If the token can't be verified, then redirect them to the sign
			// in page. In the future, we should show a message that informs
			// the user that tokens are invalidated between server restarts,
			// since that's what typically causes the error.
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next(w, r, ps)
	})
}

// AdminHandler serves the admin page.
func AdminHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := t.ExecuteTemplate(w, "admin", GetAll())
	if err != nil {
		log.Printf("Error rendering login template: ", err.Error())
	}
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
		log.Printf("Error logging in: %s\n", err.Error())
		return
	}

	tok, err := genToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error generating token: %s\n", err.Error())
		return
	}

	isSecure := false
	if r.URL.Scheme == "https" {
		isSecure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
	})
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// NewDefaultMux returns the router with its routes already initialized.
func NewDefaultMux() *httprouter.Router {
	// Create a new serve mux
	r := httprouter.New()

	// Routes
	r.GET("/", HomeHandler)
	r.GET("/posts/:id", GetPostHandler)
	r.GET("/admin/login", LoginHandler)
	r.GET("/admin/posts/:id", AdminGetPostHandler)
	r.POST("/admin/login", LoginPostHandler)

	// API
	r.GET("/api/posts/:id", APIGetPostHandler)

	// Authenticated routes
	r.GET("/admin", AuthenticatedRoute(AdminHandler))
	r.GET("/admin/new", AuthenticatedRoute(NewPostHandler))
	r.POST("/admin/new", AuthenticatedRoute(NewPostSubmitHandler))
	r.POST("/admin/edit", AuthenticatedRoute(EditPostSubmitHandler))
	r.POST("/admin/delete", AuthenticatedRoute(DeletePostSubmitHandler))

	// Serve static assets
	r.ServeFiles("/static/*filepath", http.Dir(assetPath))

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

	if int64(tok.Claims["exp"].(float64)) < time.Now().Unix() {
		return ErrExpiredToken
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
