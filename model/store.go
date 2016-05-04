package model

import (
	"github.com/boltdb/bolt"
)

// Store is the interface containing the methods needed for interaction
// between our models and any database.
type Store interface {
	Create(p *Post) (string, error)
}

// Bolt is a store saisfying the `Store` interface. It's used for communicating
// with BoltDB.
type Bolt struct {
	db *bolt.DB
}

// Create creates a new post in a BoltDB instance.
func (b *Bolt) Create(p *Post) (string, error) {
	return "", nil
}
