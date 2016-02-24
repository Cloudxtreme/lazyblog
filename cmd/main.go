package main

import (
	"net/http"

	"github.com/bentranter/lazyblog"
)

func main() {
	// Create a new serve mux
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", lazyblog.IndexHandler)
	mux.HandleFunc("/new", lazyblog.NewPostHandler)
	mux.HandleFunc("/posts", lazyblog.GetPostHandler)

	// Server static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":3000", mux)
}
