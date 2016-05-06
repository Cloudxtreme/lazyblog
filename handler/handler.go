package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bentranter/lazyblog/model"
	"github.com/julienschmidt/httprouter"
)

var s = model.NewBolt("prod.db")

// SetPost is the API method for creating a new post.
func SetPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := struct {
		Title string
		Body  string
	}{}
	d := json.NewDecoder(r.Body)
	d.Decode(&v)

	p := model.NewPost(v.Title, v.Body)
	id, err := p.Set(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(id))
}

// GetPost is the API method for getting a post.
func GetPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	post, err := model.Get(id, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
