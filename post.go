package lazyblog

import (
	"encoding/base64"
	"math/rand"
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

// NewID is the math one
func NewID() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	return base64.URLEncoding.EncodeToString(buf)
}
