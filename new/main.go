package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/bentranter/lazyblog/handler"
	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()

	r.GET("/api", handler.Info)
	r.GET("/api/posts", handler.GetAllPosts)
	r.GET("/api/posts/:id", handler.GetPostJSON)

	r.POST("/api/posts", handler.SetPost)

	r.ServeFiles("/static/*filepath", http.Dir("./static"))

	// pprof server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("go to /pprof/debug"))
	})
	go func() {
		log.Fatalln(http.ListenAndServe(":8080", nil))
	}()
	log.Fatalln(http.ListenAndServe(":3000", r))
}
