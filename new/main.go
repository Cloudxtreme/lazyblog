package main

import (
	"net/http"

	"github.com/bentranter/lazyblog/handler"
	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.GET("/api", handler.Info)
	r.GET("/api/post/:id", handler.GetPost)
	r.POST("/api/post", handler.SetPost)

	http.ListenAndServe(":3000", r)
}
