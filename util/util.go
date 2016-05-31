package util

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	defaultFlake = &flake{}

	// ErrTimeIsGoingBackwards occurs if the system clock appears to run backwards.
	ErrTimeIsGoingBackwards = errors.New("Time is running backwards on your machine.")
)

const nano = 1000 ^ 2

type flake struct {
	time uint64
	seq  uint32
	lock sync.Mutex
}

func (f *flake) next() (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	ts := uint64(time.Now().UnixNano() / nano)
	if ts == f.time {
		f.seq = f.seq + 1
	} else {
		f.seq = 0
	}

	if ts < f.time {
		return "", ErrTimeIsGoingBackwards
	}
	f.time = ts
	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, ts)
	return hex.EncodeToString(id), nil
}

// NewID base16 encodes the current time, to be used along with the Urlify'd
// Post title in order generate a usable ID for each post. BoltDB orders its
// data by byte order, so this ID is used to order posts from oldest to newest.
func NewID() string {
	id, err := defaultFlake.next()
	if err != nil {
		// Should we really panic if the clock goes backwards?
		panic(err)
	}
	return id
}

// RandStr returns a random string. Used for generating test database names.
func RandStr() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
