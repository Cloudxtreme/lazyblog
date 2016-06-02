package main

import (
	"github.com/bentranter/lazyblog/handler"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	r := fasthttprouter.New()
	r.GET("/api", handler.Info)
	r.GET("/api/posts", handler.GetAllPosts)
	r.GET("/api/posts/:id", handler.GetPost)

	r.POST("/api/posts", handler.SetPost)

	fasthttp.ListenAndServe(":3000", r.Handler)
}
