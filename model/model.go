package model

import (
	"time"

	"github.com/bentranter/lazyblog/util"
)

// Post is the struct that represents each post.
type Post struct {
	ID          string
	Path        string
	Title       string
	Body        string
	DateCreated time.Time
}

// NewPost returns a new post. It should be noted that `SavePost` must be
// called to save the post to the DB.
func NewPost() *Post {
	return &Post{
		ID:          util.NewID(),
		DateCreated: time.Now(),
	}
}

// Set persists a post to the chosen database.
func (p *Post) Set(s Store) (string, error) {
	return s.Set(p)
}
