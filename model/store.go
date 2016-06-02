// Package model blah blah blah I'm wondering if this stuff should be private
// since it ends up being called by the model code.
package model

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

var (
	boltRaw    = []byte("boltRaw")
	boltCached = []byte("boltCached")
)

// Store is the interface containing the methods needed for interaction
// between our models and any database.
type Store interface {
	Set(p *Post) (string, error)
	Get(id string) (*Post, error)
	GetAll() ([]*Post, error)
}

// Bolt is a store satisfying the `Store` interface. It's used for communicating
// with BoltDB.
type Bolt struct {
	db *bolt.DB
}

// NewBolt returns a new instance of the Bolt struct.
func NewBolt(name string) *Bolt {
	db, err := bolt.Open(name, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	// we can't start without connecting to our db, so panic
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltRaw)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(boltCached)
		if err != nil {
			return err
		}

		return nil
	})
	// some errors can be ignored, so we break out if we see those
	if err != nil {
		switch err {
		case bolt.ErrBucketExists:
			break
		default:
			panic(err)
		}
	}

	return &Bolt{
		db: db,
	}
}

// Set sets a new post in a BoltDB instance. It returns the post ID and any
// errors that occurred.
func (b *Bolt) Set(p *Post) (string, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		rawBucket := tx.Bucket(boltRaw)
		// cachedBucket := tx.Bucket(boltCached)

		post, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return rawBucket.Put([]byte(p.ID), post)
	})
	return p.ID, err
}

// Get retrieves a post and marshals it into a struct.
func (b *Bolt) Get(id string) (*Post, error) {
	var p *Post
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltRaw)
		err := json.Unmarshal(bucket.Get([]byte(id)), &p)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetAll returns evert post in the database, and returns them as an array of
// `Post` structs.
func (b *Bolt) GetAll() ([]*Post, error) {
	var posts []*Post
	err := b.db.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(boltRaw)
		c := raw.Cursor()

		// Posts need to be ordered by the trailing base16 string, since it's
		// actually in descending byte-order
		for id, postJSON := c.Last(); id != nil; id, postJSON = c.Prev() {
			var post *Post
			json.Unmarshal(postJSON, &post)
			posts = append(posts, post)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return posts, nil
}
