package lazyblog

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var (
	_raw      = []byte("raw")
	_rendered = []byte("rendered")
	_users    = []byte("users")
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
	//   1. It's good to check for errors
	//   2. Calling `writeTo` drains the buffer.
	// So idk what to do yet
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
		buf = rendered.Get(GetID([]byte(id)))
		return nil
	})

	return buf
}

// DeletePost deletes a post given its id. It deletes a post from both buckets.
func DeletePost(id string) error {
	return DefaultStore.Update(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		rendered := tx.Bucket(_rendered)

		err := raw.Delete([]byte(id))
		if err != nil {
			return err
		}
		return rendered.Delete([]byte("id"))
	})
}

// GetPostForAPI gets the JSON for the post the given id. For use with the
// API endpoints.
func GetPostForAPI(id string) []byte {
	var buf []byte
	DefaultStore.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		buf = raw.Get(GetID([]byte(id)))
		return nil
	})
	return buf
}

// GetAll retrieves every post from the database, in byte order.
func GetAll() []*PostJSON {
	var posts []*PostJSON
	DefaultStore.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(_raw)
		c := raw.Cursor()

		// Posts need to be ordered by the trailing base16 string, since it's
		// actually in descending byte-order
		for id, postJSON := c.Last(); id != nil; id, postJSON = c.Prev() {
			var post *PostJSON
			json.Unmarshal(postJSON, &post)
			posts = append(posts, post)
		}
		return nil
	})
	return posts
}

// GetUser returns the hashed password for the given username.
func GetUser(username string) []byte {
	var buf []byte
	DefaultStore.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(_users)
		buf = users.Get([]byte(username))
		return nil
	})

	return buf
}

// GetID gets trailing 8 bytes of each post title that represent the ID of each
// post. Since these bytes are in "order", they're used for ordering the
// results by date created.
func GetID(id []byte) []byte {
	// This is my fav line of code in this repo. It looks pretty confusing, but
	// it's pretty simple once you figure out Go's weird slice operations.
	// Anyway, this is what's happening:
	//   1. We know that each ID is the last 8 chars in the ID string, so
	//   2. We find the index of the first char by subtracting 8 from the
	//      length of our id.
	//   3. We return all the bytes from that index to the end of the string.
	// This generates no garbage and runs at 13 ns/op!
	return id[(len(id) - 8):]
}
