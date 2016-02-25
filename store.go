package lazyblog

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

var (
	_raw      = []byte("raw")
	_rendered = []byte("rendered")
)

// DefaultStore @Todo
var DefaultStore = NewDefaultStore()

// SetPost creates a new post. It will save two versions of the post: a "raw"
// version, containg the id and the body, and a "rendered" version, containing
// the id and the rendered and compressed HTML for that post.
func SetPost(w http.ResponseWriter, post *PostJSON) {
	// immediately write the response so the client feels snappy
	// this should be a buffer pool, and response should be gzipped first!
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, "post", post)
	if err != nil {
		log.Fatalln("Error executing: ", err)
	}

	DefaultStore.Update(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		rendered := tx.Bucket(_rendered)

		postJSON, err := json.Marshal(post)
		if err != nil {
			return err
		}

		err = raw.Put([]byte(post.ID), postJSON)
		if err != nil {
			return err
		}

		return rendered.Put([]byte(post.ID), buf.Bytes())
	})

	// Ideally this would happen before we touch the buckets, but:
	//     1. It's good to check for errors
	//     2. Calling `writeTo` drains the buffer.
	// So idk what to yet
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error: ", err)
	}
}

// GetPost retrieves the post body given its id, and returns it as a byte
// slice.
func GetPost(id string) []byte {
	var buf []byte
	DefaultStore.View(func(tx *bolt.Tx) error {
		rendered := tx.Bucket(_rendered)
		buf = rendered.Get([]byte(id))
		return nil
	})

	return buf
}

// GetAll retrieves every post from the database. You probably don't want to
// use this -- use GetPosts instead.
func GetAll() []*PostJSON {
	var posts []*PostJSON
	DefaultStore.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		c := raw.Cursor()

		for id, postJSON := c.First(); id != nil; id, postJSON = c.Next() {
			var post *PostJSON
			json.Unmarshal(postJSON, &post)
			posts = append(posts, post)
		}
		return nil
	})
	return posts
}

// NewDefaultStore creates our store if it doesn't already exist. We also
// create two buckets along with our db: a "raw" bucket, for storing the
// raw post data, and a "rendered" bucket for storing the compiled and
// compressed HTML data.
func NewDefaultStore() *bolt.DB {
	db, err := bolt.Open("store.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	// we can't start without connecting to our db, so panic
	if err != nil {
		panic(err)
	}

	// make two buckets: one for raw post data, another for compiled templates
	db.Update(func(tx *bolt.Tx) error {
		// we don't care about the returned bucket, so ignore it
		_, err := tx.CreateBucketIfNotExists(_raw)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(_rendered)
		if err != nil {
			return err
		}

		return nil
	})

	return db
}
