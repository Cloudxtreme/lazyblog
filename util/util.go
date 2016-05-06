package util

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

// NewID base16 encodes the current time, to be used along with the Urlify'd
// Post title in order generate a usable ID for each post. BoltDB orders its
// data by byte order, so this ID is used to order posts from oldest to newest.
func NewID() string {
	now := big.NewInt(time.Now().Unix()).Bytes()
	return hex.EncodeToString(now)
}

// RandStr returns a random string. Used for generating test database names.
func RandStr() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
