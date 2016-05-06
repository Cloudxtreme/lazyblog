package model

import (
	"errors"
	"strings"
	"time"

	"github.com/bentranter/lazyblog/util"
)

var (
	ErrMissingPostID    = errors.New("Post ID is required to save a post")
	ErrMissingPostTitle = errors.New("Post Title is required to save a post")
)

// Post is the struct that represents each post.
type Post struct {
	ID          string
	Path        string // @TODO: Decide if path should be private
	Title       string
	Body        string
	DateCreated int64
}

// NewPost returns a new post. It should be noted that `SavePost` must be
// called to save the post to the DB.
func NewPost(title string, body string) *Post {
	return &Post{
		ID:          util.NewID(),
		Title:       title,
		Body:        body,
		DateCreated: time.Now().Unix(),
	}
}

// Set persists a post to the chosen database.
func (p *Post) Set(s Store) (string, error) {
	err := p.urlify()
	if err != nil {
		return "", err
	}
	return s.Set(p)
}

// Get retrieves a post from the chosen database, and returns the `Post` struct
// for it.
func Get(id string, s Store) (*Post, error) {
	p, err := s.Get(id)
	return p, err
}

// urlify is a utility for making strings URL safe. It removes anything that
// isn't a number or letter, and replaces each with a `-`. It then appends a
// single "-" to the end, followed by the post ID.
func (p *Post) urlify() error {
	if p.ID == "" {
		return ErrMissingPostID
	}
	if p.Title == "" {
		return ErrMissingPostTitle
	}
	if p.Path != "" {
		return nil
	}

	title := strings.ToLower(p.Title)
	buf := []byte(title)

	var bytebuf []byte
	for _, b := range buf {
		if (b >= 48 && b <= 57) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122) {
			bytebuf = append(bytebuf, b)
		} else if b == 32 {
			bytebuf = append(bytebuf, '-')
		}
	}

	bytebuf = append(bytebuf, '-')
	p.Path = string(bytebuf) + p.ID
	return nil
}
