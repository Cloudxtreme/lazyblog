package model

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/bentranter/lazyblog/util"
)

var sBolt *Bolt

func TestMain(m *testing.M) {
	dbStr := util.RandStr() + ".db"
	sBolt = NewBolt(dbStr)

	exitCode := m.Run()

	os.Remove(dbStr)
	os.Exit(exitCode)
}

func TestNewBolt(t *testing.T) {
	t.Parallel()
	_ = NewBolt("test_new_bolt.db")

	if _, err := os.Stat("test_new_bolt.db"); os.IsNotExist(err) {
		t.Errorf("No DB was created: %s\n", err.Error())
	}

	if err := os.Remove("test_new_bolt.db"); err != nil {
		t.Logf("Info: Error removing test database: %s\n", err.Error())
	}
}

func TestBolt_SetPost(t *testing.T) {
	t.Parallel()

	p := NewPost("Title", "Body")
	_, err := sBolt.SetPost(p)

	if err != nil {
		t.Errorf("Error while setting post: %s\n", err.Error())
	}
}

func TestBolt_GetPostHTML(t *testing.T) {
	t.Parallel()

	p := NewPost("Title", "Body")
	id, _ := sBolt.SetPost(p)
	postHTML, err := sBolt.GetPostHTML(id)

	if err != nil {
		t.Errorf("Error while retrieving saved post: %s\n", err.Error())
	}

	if len(postHTML) == 0 {
		t.Errorf("Post doesn't exist\n")
	}

	if bytes.Equal(postHTML, []byte("Body\n")) {
		t.Errorf(`Expected "%s" to match "Body"`, postHTML)
	}
}

func TestBolt_GetPostJSON(t *testing.T) {
	t.Parallel()

	p := NewPost("Title", "Body")
	id, _ := sBolt.SetPost(p)
	postJSON, err := sBolt.GetPostJSON(id)

	if err != nil {
		t.Errorf("Error while retrieving saved post: %s\n", err.Error())
	}

	if len(postJSON) == 0 {
		t.Errorf("Post doesn't exist\n")
	}

	var unmarshaledPostJSON *Post
	err = json.Unmarshal(postJSON, &unmarshaledPostJSON)
	if err != nil {
		t.Errorf("Error marshalling post JSON: %s\n", err.Error())
	}

	if unmarshaledPostJSON.Title != p.Title {
		t.Errorf("Expected post titles to match, but got %s and %s\n", unmarshaledPostJSON.Title, p.Title)
	}
}

// func TestBolt_GetAll(t *testing.T) {
// 	t.Parallel()

// 	px := NewPost("TitleX", "BodyX")
// 	sBolt.Set(px)
// 	py := NewPost("TitleY", "BodyY")
// 	sBolt.Set(py)
// 	posts, err := sBolt.GetPosts(num, offset)

// 	t.Logf("%s\n", posts)

// 	if err != nil {
// 		t.Errorf("Error while retrieving all posts: %s\n", err.Error())
// 	}
// 	if len(posts) != 2 {
// 		t.Errorf("Didn't get all posts, expected %d, got %d\n", 2, len(posts))
// 	}
// }

func TestBolt_SetUser(t *testing.T) {
	t.Parallel()

	err := sBolt.SetUser([]byte("Test"), []byte("Test"))
	if err != nil {
		t.Errorf("Error creating new user: %s\n", err.Error())
	}
}

func TestBolt_GetUser(t *testing.T) {
	t.Parallel()

	err := sBolt.SetUser([]byte("Test"), []byte("Test"))
	if err != nil {
		t.Errorf("Error creating new user: %s\n", err.Error())
	}

	err = sBolt.GetUser([]byte("Test"), []byte("Test"))
	if err != nil {
		t.Errorf("User should exist and passwords should match, but got: %s\n", err.Error())
	}

	err = sBolt.GetUser([]byte("Test"), []byte("TestWrongPassword"))
	if err == nil {
		t.Errorf("Password should not match, but appears to match\n")
	}

	err = sBolt.GetUser([]byte("TestWrongUser"), []byte("Test"))
	if err == nil {
		t.Errorf("Tried to get non-existent user, but user appears to exist\n")
	}

	err = sBolt.GetUser(nil, nil)
	if err == nil {
		t.Errorf("Tried to get nil user and password, but didn't encounter error\n")
	}

	// @TODO: Add cases where user changes password
}
