package lazyblog

import (
	"encoding/base64"
	"math/rand"
	"strings"
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

// Urlify is a utility for making strings URL safe. It removes anything that
// isn't a number of letter, and replaces each with a `-`. It then generates
// an id and appends it to the end of the string.
func Urlify(id string) []byte {
	id = strings.ToLower(id)
	buf := []byte(id)
	var bytebuf []byte
	for _, b := range buf {
		if (b >= 48 && b <= 57) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122) {
			bytebuf = append(bytebuf, b)
		} else if b == 32 {
			bytebuf = append(bytebuf, '-')
		}
	}
	bytebuf = append(bytebuf, []byte(NewID())...) // worst syntax ever lol
	return bytebuf
}
