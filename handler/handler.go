package handler

import (
	"encoding/json"

	"github.com/bentranter/lazyblog/model"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var s = model.NewBolt("dev_bolt.db")

// Info displays info about the available API routes.
func Info(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	routes := `{
  "posts_url": "/api/posts{/id}"
}`
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.WriteString(routes)
}

// SetPost is the API method for creating a new post.
func SetPost(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	v := &struct {
		Title string
		Body  string
	}{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	p := model.NewPost(v.Title, v.Body)
	id, err := p.Set(s)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	resp, err := json.MarshalIndent(map[string]string{
		"id": string(id),
	}, "", "  ")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	ctx.Write(resp)
}

// GetPostJSON is the API method for getting a post's JSON.
func GetPostJSON(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	id := ps.ByName("id")
	postJSON, err := model.GetJSON([]byte(id), s)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(postJSON)
}

// GetAllPosts is a method for getting every post
func GetAllPosts(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	posts, err := model.GetAll(s)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}

	resp, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(resp)
}
