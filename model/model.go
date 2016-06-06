package model

import (
	// "encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bentranter/lazyblog/util"
	// "github.com/russross/blackfriday"
)

var (
	ErrMissingPostID    = errors.New("Post ID is required to save a post")
	ErrMissingPostTitle = errors.New("Post Title is required to save a post")
)

// Post is the struct that represents each post.
type Post struct {
	ID          []byte `json: "id"`
	Path        string `json: "path"`
	Title       string `json: "title"`
	Body        string `json: "body"`
	DateCreated int64  `json: "dateCreated"`
}

// A User represents a user.
type User struct {
	Username string
	Password string
}

// NewPost returns a new post. It should be noted that `SavePost` must be
// called to save the post to the DB.
func NewPost(title string, body string) *Post {
	return &Post{
		ID:          []byte(util.NewID()),
		Title:       title,
		Body:        body,
		DateCreated: time.Now().Unix(),
	}
}

// Set persists a post to the chosen database.
func (p *Post) Set(s Store) ([]byte, error) {
	err := p.urlify()
	if err != nil {
		return nil, err
	}
	return s.SetPost(p)
}

func Get(id []byte, s Store) (*Post, error) {
	return nil, nil
}

// GetHTML returns the post HTML.
func GetHTML(id []byte, s Store) ([]byte, error) {
	return s.GetPostHTML(id)
}

// GetJSON retrieves a post from the chosen database, and returns the `Post` struct
// for it.
func GetJSON(id []byte, s Store) ([]byte, error) {
	return s.GetPostJSON(id)
}

// GetAll retrieves every post from the chosen database, and returns every
// `Post` struct in there.
func GetAll(s Store) ([]*Post, error) {
	return s.GetPosts(0, 0)
}

// urlify is a utility for making strings URL safe. It removes anything that
// isn't a number or letter, and replaces each with a `-`. It then appends a
// single "-" to the end, followed by the post ID.
func (p *Post) urlify() error {
	if p.ID == nil {
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
	p.Path = string(bytebuf) + string(p.ID)
	return nil
}
