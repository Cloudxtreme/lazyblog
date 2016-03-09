package lazyblog

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

// Setup runs our setup process, creating a new DB,
func Setup() {
	var username []byte
	var password []byte

	for username == nil {
		fmt.Printf("Enter a username:\n")
		fmt.Scanf("%s", &username)
	}
	for password == nil {
		fmt.Printf("Enter a password:\n")
		fmt.Scanf("%s", &password)
	}

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

// NumUsers returns how many users are registered.
func NumUsers() (int, error) {
	var numUsers int
	err := DefaultStore.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(_users)
		numUsers = users.Stats().KeyN
		return nil
	})
	if err != nil {
		return 0, err
	}
	return numUsers, nil
}
