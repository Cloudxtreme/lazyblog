package main

import (
	"github.com/bentranter/lazyblog/handler"
	"github.com/dinever/golf"
)

func main() {
	app := golf.New()

	app.Get("/api", handler.Info)
	app.Get("/api/posts", handler.GetAllPosts)
	app.Get("/api/posts/:id", handler.GetPostJSON)

	app.Post("/api/posts", handler.SetPost)

	app.Run(":3000")
}
