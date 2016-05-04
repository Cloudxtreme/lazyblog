package model

import (
	"os"
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	t.Parallel()
	p := NewPost()

	if p.DateCreated.After(time.Now()) {
		t.Errorf("DateCreated must be before %s, , but is %s\n", time.Now().String(), p.DateCreated.String())
	}
}

func TestPost_Set(t *testing.T) {
	t.Parallel()
	s := NewBolt("test.db")
	p := NewPost()
	id, err := p.Set(s)

	if err != nil {
		t.Errorf("Error when setting new post: %s\n", err.Error())
	}
	if len(id) < 8 {
		t.Errorf("Saved Post ID doesn't meet length requirement of more than 8 characters: %s\n", id)
	}

	if err = os.Remove("test.db"); err != nil {
		t.Logf("Info: Error removing test database: %s\n", err.Error())
	}
}
