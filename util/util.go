package util

import (
	"encoding/hex"
	"math/big"
	"strings"
	"time"
)

// NewID base16 encodes the current time, to be used along with the Urlify'd
// Post title in order generate a usable ID for each post. BoltDB orders its
// data by byte order, so this ID is used to order posts from oldest to newest.
func NewID() string {
	now := big.NewInt(time.Now().Unix()).Bytes()
	return hex.EncodeToString(now)
}

// Urlify is a utility for making strings URL safe. It removes anything that
// isn't a number or letter, and replaces each with a `-`. It then appends a
// single "-" to the end.
func Urlify(id string) string {
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
	bytebuf = append(bytebuf, '-')
	return string(bytebuf)
}
