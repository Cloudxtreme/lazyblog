package lazyblog

import (
	"bytes"
	"encoding/binary"
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
func SetPost(w http.ResponseWriter, post *Post) {
	// immediately write the response so the client feels snappy
	// this should be a buffer pool, and response should be gzipped first!
	var buf bytes.Buffer
	t.ExecuteTemplate(&buf, "new", post)
	_, err := buf.WriteTo(w)
	if err != nil {
		// not sure what to do with this yet... I guess don't write it to the
		// store at least?
		log.Println("Error: ", err)
	}

	DefaultStore.Update(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		rendered := tx.Bucket(_rendered)

		// this will never error under an "update" transaction, so it's safe to
		// ignore this error
		id, _ := raw.NextSequence()
		post.ID = append(post.ID, itob(id)...)

		err = raw.Put(post.ID, post.Body)
		if err != nil {
			return err
		}

		return rendered.Put(post.ID, buf.Bytes())
	})

}

// GetPost retrieves the post body given its id, and returns it as a byte
// slice.
func GetPost(id []byte) []byte {
	var buf []byte
	DefaultStore.View(func(tx *bolt.Tx) error {
		rendered := tx.Bucket(_rendered)
		buf = rendered.Get(id)
		return nil
	})

	return buf
}

// GetAll retrieves every post from the database. You probably don't want to
// use this -- use GetPosts instead.
func GetAll() []*Post {
	var posts []*Post
	DefaultStore.View(func(tx *bolt.Tx) error {
		rendered := tx.Bucket(_rendered)
		c := rendered.Cursor()

		for id, body := c.First(); id != nil; id, body = c.Next() {
			posts = append(posts, &Post{
				ID:   id,
				Body: body,
			})
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

// itob is a utility for converting uint64s to []bytes. I *think* your
// system's endian-ness doesn't matter since Bolt will know what to do either
// way.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
