package main

import (
	"net/http"
)

func main() {
	// Create a new serve mux
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", IndexHandler)

	// Server static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":3000", mux)
}
