// Package model blah blah blah I'm wondering if this stuff should be private
// since it ends up being called by the model code.
package model

import "golang.org/x/crypto/bcrypt"

const cost = bcrypt.DefaultCost

// Store is the interface containing the methods needed for interaction
// between our models and any database.
type Store interface {
	SetPost(p *Post) ([]byte, error)
	GetPostHTML(id []byte) ([]byte, error)
	GetPostJSON(id []byte) ([]byte, error)
	GetPosts(num, offset int) ([]*Post, error)
	SetUser(username, password []byte) error
	GetUser(username, password []byte) error
}
