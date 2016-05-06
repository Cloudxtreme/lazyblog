package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bentranter/lazyblog/util"
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

func TestPost_Set(t *testing.T) {
	t.Parallel()
	dbStr := util.RandStr() + ".db"
	s := NewBolt(dbStr)
	p := NewPost("Title", "Body")
	id, err := p.Set(s)
	if err != nil {
		t.Errorf("Error when setting new post: %s\n", err.Error())
	}
	if len(id) < 8 {
		t.Errorf("Saved Post ID doesn't meet length requirement of more than 8 characters: %s\n", id)
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error removing test database: %s\n", err.Error())
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	dbStr := util.RandStr() + ".db"
	s := NewBolt(dbStr)
	p := NewPost("Title", "Body")
	p.Set(s)
	px, err := Get(p.ID, s)

	if err != nil {
		t.Errorf("Error while getting post: %s\n", err.Error())
	}
	if !reflect.DeepEqual(p, px) {
		t.Errorf("Posts do not match: %s %s\n", p.DateCreated, px.DateCreated)
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error removing test database: %s\n", err.Error())
	}
}

func TestPost_urlify(t *testing.T) {
	t.Parallel()
	px := &Post{
		ID:    "572b7220",
		Title: "mytest$1234+===((()",
	}
	px.urlify()
	py := &Post{
		ID:    "572b7220",
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
