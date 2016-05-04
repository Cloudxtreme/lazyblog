package model

import (
	"os"
	"testing"
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
