package model

import (
	"bytes"
	"testing"
)

func TestNewMemoryStore(t *testing.T) {
	_ = NewMemoryStore()
}

func TestMemoryStore_SetPost(t *testing.T) {
	s := NewMemoryStore()
	p := NewPost("Test", "Test")
	_, err := s.SetPost(p)

	if err != nil {
		t.Errorf("Failed to set post: %s\n", err.Error())
	}
}

func TestMemoryStore_GetPostHTML(t *testing.T) {
	s := NewMemoryStore()
	p := NewPost("Test", "Test")
	id, err := s.SetPost(p)

	if err != nil {
		t.Errorf("Failed to set post: %s\n", err.Error())
	}

	html, err := s.GetPostHTML(id)
	if err != nil {
		t.Errorf("Failed to get post: %s\n", err.Error())
	}

	if html == nil {
		t.Errorf("Failed to retrieve HTML from set post\n")
	}

	if bytes.Equal(html, []byte("Test\n")) {
		t.Errorf("Expected html to be %s but got %s\n", html, []byte("Test\n"))
	}
}

func TestMemoryStore_GetPostJSON(t *testing.T) {
	s := NewMemoryStore()
	p := NewPost("Test", "Test")
	id, err := s.SetPost(p)

	if err != nil {
		t.Errorf("Failed to set post: %s\n", err.Error())
	}

	j, err := s.GetPostJSON(id)
	if err != nil {
		t.Errorf("Failed to get post: %s\n", err.Error())
	}

	if j == nil {
		t.Errorf("Failed to retrieve JSON from set post\n")
	}

	if bytes.Equal(j, []byte("Test\n")) {
		t.Errorf("Expected json to be %s but got %s\n", j, []byte("Test\n"))
	}
}

func TestMemoryStore_GetPosts(t *testing.T) {
	s := NewMemoryStore()
	px := NewPost("TestX", "TestX")
	py := NewPost("TestY", "TestY")

	_, err := s.SetPost(px)
	if err != nil {
		t.Errorf("Failed to set post: %s\n", err.Error())
	}

	_, err = s.SetPost(py)
	if err != nil {
		t.Errorf("Failed to set post: %s\n", err.Error())
	}

	posts, err := s.GetPosts(0, 0)
	if err != nil {
		t.Errorf("Failed to get posts: %s\n", err.Error())
	}

	f := func(a []*Post, b *Post) bool {
		for _, post := range a {
			if post.Title == b.Title {
				return true
			}
		}
		return false
	}

	if !f(posts, px) {
		t.Errorf("Expected %s to exist in %s\n", px.Title, posts)
	}

	if !f(posts, py) {
		t.Errorf("Expected %s to exist in %s\n", py.Title, posts)
	}
}

func TestMemoryStore_SetUser(t *testing.T) {
	t.Skip("SetUser not yet implemented")
}

func TestMemoryStore_GetUser(t *testing.T) {
	t.Skip("GetUser not yet implemented")
}
