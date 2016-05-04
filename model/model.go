package model

import (
	"time"

	"github.com/bentranter/lazyblog/util"
	// "github.com/boltdb/bolt"
)

// Post is the struct that represents each post.
type Post struct {
	ID          string
	Path        string
	Title       string
	Body        string
	DateCreated time.Time
}

// PostCached is the struct that holds a rendered version of a post.
type PostCached struct {
	ID   string
	Body []byte
}

// NewPost returns a new post. It should be noted that `SavePost` must be
// called to save the post to the DB.
func NewPost() *Post {
	return &Post{
		ID:          util.NewID(),
		DateCreated: time.Now(),
	}
}
