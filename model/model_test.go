package model

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNewPost() {
	p := NewPost("Hello, world!", "Body")
	fmt.Println(p.Title)
	// Output: Hello, world!
}

func TestNewPost(t *testing.T) {
	t.Parallel()
	p := NewPost("Title", "Body")

	if p.DateCreated > time.Now().Unix() {
		t.Errorf("DateCreated must be before %d, but is %d\n", time.Now().Unix(), p.DateCreated)
	}
}

func TestPost_urlify(t *testing.T) {
	t.Parallel()
	px := &Post{
		ID:    []byte("572b7220"),
		Title: "mytest$1234+===((()",
	}
	px.urlify()
	py := &Post{
		ID:    []byte("572b7220"),
		Title: "my test with spaces",
	}
	py.urlify()
	pz := &Post{}
	err := pz.urlify()

	if px.Path != "mytest1234-572b7220" {
		t.Errorf("Didn't pass, got %s, expected mytest1234\n", px.Path)
	}
	if py.Path != "my-test-with-spaces-572b7220" {
		t.Errorf("Didn't pass, got %s, expected my-test-with-spaces\n", py.Path)
	}
	if err == nil {
		t.Error("Error must be returned if ID or Title are blank\n")
	}
}
