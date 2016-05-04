package model

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	t.Parallel()
	p := NewPost()

	if p.DateCreated.After(time.Now()) {
		t.Errorf("DateCreated must be less than or equal to now")
	}
}
