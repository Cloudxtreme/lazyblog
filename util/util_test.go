package util

import (
	"testing"
)

func TestNewID(t *testing.T) {
	t.Parallel()
	id := NewID()

	if len(id) < 8 {
		t.Errorf("Failed to generate usable id, length is too short")
	}
}
