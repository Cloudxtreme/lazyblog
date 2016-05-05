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
	DateCreated int64
}

// NewPost returns a new post. It should be noted that `SavePost` must be
// called to save the post to the DB.
func NewPost() *Post {
	return &Post{
		ID:          util.NewID(),
		DateCreated: time.Now().Unix(),
	}
}

// Set persists a post to the chosen database.
func (p *Post) Set(s Store) (string, error) {
	return s.Set(p)
}

// Get retrieves a post from the chosen database, and returns the `Post` struct
// for it.
func Get(id string, s Store) (*Post, error) {
	p, err := s.Get(id)
	return p, err
}
