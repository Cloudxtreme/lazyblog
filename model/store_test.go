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
	p := NewPost()
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
	p := NewPost()
	id, _ := b.Set(p)
	px, err := b.Get(id)

	if err != nil {
		t.Errorf("Error while retrieving saved post: %s\n", err.Error())
	}
	if id != px.ID {
		t.Errorf("Retrieved post contains incorrect id: %d vs %d\n", id, px.ID)
	}

	if err = os.Remove(dbStr); err != nil {
		t.Logf("Info: Error while removing test database: %s\n", err.Error())
	}
}
