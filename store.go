package lazyblog

import (
	"time"

	"github.com/boltdb/bolt"
)

// DefaultStore @Todo
var DefaultStore = NewDefaultStore()

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
		_, err := tx.CreateBucketIfNotExists([]byte("raw"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("rendered"))
		if err != nil {
			return err
		}

		return nil
	})

	return db
}
