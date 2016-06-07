package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bentranter/lazyblog/model"
	"github.com/julienschmidt/httprouter"
)

var s = model.NewBolt("dev_bolt.db")

// Info displays info about the available API routes.
//
// @TODO: Don't marshal JSON for every request, just hard code this as a byte
//        slice.
func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	routes := map[string]string{
		"posts_url": "/api/posts{/id}",
	}
	resp, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// SetPost is the API method for creating a new post.
func SetPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := &struct {
		Title string
		Body  string
	}{}
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := model.NewPost(v.Title, v.Body)
	id, err := p.Set(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.MarshalIndent(map[string]string{
		"id": string(id),
	}, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(resp)
}

// GetPostJSON is the API method for getting a post's JSON.
func GetPostJSON(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	postJSON, err := model.GetJSON([]byte(id), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postJSON)
}

// GetAllPosts is a method for getting every post
func GetAllPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := model.GetAll(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
