package model

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

var (
	boltJSON  = []byte("boltJSON")
	boltHTML  = []byte("boltHTML")
	boltUsers = []byte("boltUsers")
)

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
		_, err := tx.CreateBucketIfNotExists(boltJSON)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(boltHTML)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(boltUsers)
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

// SetPost sets a new post in a BoltDB instance. It returns the post ID and any
// errors that occurred.
func (b *Bolt) SetPost(p *Post) ([]byte, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		JSONBucket := tx.Bucket(boltJSON)
		HTMLBucket := tx.Bucket(boltHTML)

		post, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return err
		}

		err = JSONBucket.Put(p.ID, post)
		if err != nil {
			return err
		}
		return HTMLBucket.Put(p.ID, []byte(p.Body))
	})
	return p.ID, err
}

// GetPostHTML returns the rendered HTML for the associated post.
func (b *Bolt) GetPostHTML(id []byte) ([]byte, error) {
	var p []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltHTML)
		p = bucket.Get(id)
		return nil
	})
	return p, err
}

// GetPostJSON returns the JSON for the associated post.
func (b *Bolt) GetPostJSON(id []byte) ([]byte, error) {
	var p []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltJSON)
		p = bucket.Get(id)
		return nil
	})
	return p, err
}

// GetPosts returns every post in the database, and returns them as an array of
// `Post` structs.
func (b *Bolt) GetPosts(num, offset int) ([]*Post, error) {
	var posts []*Post
	err := b.db.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(boltJSON)
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
	return posts, err
}

// SetUser sets the username and password for the given user in BoltDB. It
// hashed and salts the password before saving it.
func (b *Bolt) SetUser(username, password []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		rawBucket := tx.Bucket(boltUsers)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, cost)
		if err != nil {
			return err
		}
		return rawBucket.Put(username, hashedPassword)
	})
}

// GetUser compares the password of the given user with the hashed password for
// that user, returning an error if the user doesn't exist, or has provided an
// incorrect password.
func (b *Bolt) GetUser(username, password []byte) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltUsers)
		hashedPassword := bucket.Get(username)
		return bcrypt.CompareHashAndPassword(hashedPassword, password)
	})
}
