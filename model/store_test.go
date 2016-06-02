package model

import (
	"os"
	"testing"

	"github.com/bentranter/lazyblog/util"
)

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

func TestBolt_Set(t *testing.T) {
	t.Parallel()
	dbStr := util.RandStr() + ".db"
	b := NewBolt(dbStr)
	p := NewPost("Title", "Body")
	_, err := b.Set(p)

	if err != nil {
		t.Errorf("Error while creating new post: %s\n", err.Error())
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error while removing test database: %s\n", err.Error())
	}

}

func TestBolt_Get(t *testing.T) {
	t.Parallel()
	dbStr := util.RandStr() + ".db"
	b := NewBolt(dbStr)
	p := NewPost("Title", "Body")
	id, _ := b.Set(p)
	px, err := b.Get(id)

	if err != nil {
		t.Errorf("Error while retrieving saved post: %s\n", err.Error())
	}
	if len(px.Bytes()) == 0 {
		t.Errorf("Post doesn't exist")
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error while removing test database: %s\n", err.Error())
	}
}

func TestBolt_GetAll(t *testing.T) {
	t.Skip()
	t.Parallel()
	dbStr := util.RandStr() + ".db"
	b := NewBolt(dbStr)
	px := NewPost("TitleX", "BodyX")
	b.Set(px)
	py := NewPost("TitleY", "BodyY")
	b.Set(py)
	posts, err := b.GetAll()

	t.Logf("%s\n", posts)

	if err != nil {
		t.Errorf("Error while retrieving all posts: %s\n", err.Error())
	}
	if len(posts) != 2 {
		t.Errorf("Didn't get all posts, expected %d, got %d\n", 2, len(posts))
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error while removing test database: %s\n", err.Error())
	}
}
