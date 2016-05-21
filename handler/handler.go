package handler

import (
	"encoding/json"

	"github.com/bentranter/lazyblog/model"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var s = model.NewBolt("prod.db")

// Info displays info about the available API routes.
//
// @TODO: Don't marshal JSON for every request, just hard code this as a byte
//        slice.
func Info(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	routes := map[string]string{
		"posts_url": "/api/post/:id",
	}
	resp, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// SetPost is the API method for creating a new post.
func SetPost(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	v := struct {
		Title string
		Body  string
	}{}
	err := json.Unmarshal(ctx.PostBody(), v)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer ctx.Request.ConnectionClose()

	p := model.NewPost(v.Title, v.Body)
	id, err := p.Set(s)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	resp, err := json.MarshalIndent(map[string]string{
		"id": id,
	}, "", "  ")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	ctx.Write(resp)
}

// GetPost is the API method for getting a post.
func GetPost(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	id := ps.ByName("id")
	post, err := model.Get(id, s)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}

	resp, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Write(resp)
}
