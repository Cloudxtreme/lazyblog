package main

import (
	"github.com/bentranter/lazyblog/handler"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	r := fasthttprouter.New()
	r.GET("/api", handler.Info)
	r.GET("/api/post/:id", handler.GetPost)
	r.POST("/api/post", handler.SetPost)

	fasthttp.ListenAndServe(":3000", r.Handler)
}
