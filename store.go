package lazyblog

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
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

// Setup runs our setup process. If the user hasn't run the code before, it
// will create a new user account with the given username and password. If the
// user has, it will requre that username and password before starting. Idk if
// this is a good idea or not, so I'm open to feedback.
//
// @TODO before v0.1.0:
//
// On second thought, this sucks for users who plan to use systemd or launchd
// or upstart or whatever windows servers use since you can't use the default
// script to start your app... so I need to change this.
func Setup(username string, password string) {
	if username == "" || password == "" {
		log.Fatalln("You must provide a username and password") // should use os.Exit
		return
	}

	var hashedPassword []byte
	DefaultStore.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(_users)
		hashedPassword = users.Get([]byte(username))
		return nil
	})

	if hashedPassword == nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalln("There was an error hashing and salting your password: ", err.Error())
			return
		}

		DefaultStore.Update(func(tx *bolt.Tx) error {
			users := tx.Bucket(_users)
			return users.Put([]byte(username), hashed)
		})
		log.Println("Congrats on setting up your Lazyblog! You're all set, try visiting /admin to begin.")
		return
	}

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Fatalln("Password incorrect, please try again: ", err.Error())
		return
	}

	log.Println("Started Lazyblog successfully. Welcome back,", username)
}

// NewDefaultStore creates our store if it doesn't already exist. We also
// create three buckets along with our db: a "raw" bucket, for storing the
// raw post data, and a "rendered" bucket for storing the compiled and
// compressed HTML data, and a "users" bucket, which is used for authenticating
// the user.
func NewDefaultStore() *bolt.DB {
	db, err := bolt.Open("store.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	// we can't start without connecting to our db, so panic
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(_raw)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(_rendered)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(_users)
		if err != nil {
			return err
		}

		return nil
	})

	return db
}
