package lazyblog

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

// counter is incremented every time an id generated
var (
	counter = 0
	prng    = seed()
)

// Post is the struct that represents our post data. It works for the two types
// of data we store: the raw data that is passed to an HTML template, and the
// rendered and compressed data that we show to the user.
//
// In the future, it's likely we'll add a comments field to this struct.
type Post struct {
	ID   []byte
	Body []byte
}

func seed() []byte {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		// couldn't read from PRNG, so panic
		panic(err)
	}
	return buf
}

// NewID generates a new id
func NewID() []byte {
	// some weird suff in here, may refactor
	id := append(prng, []byte(strconv.Itoa(counter))...)
	counter++
	return []byte(base64.URLEncoding.EncodeToString(id))
}
