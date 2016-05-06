package util

import (
	"testing"
)

func TestNewID(t *testing.T) {
	t.Parallel()
	id := NewID()

	if len(id) < 8 {
		t.Errorf("Failed to generate usable id, length is too short: %s\n", id)
	}
}

func TestRandStr(t *testing.T) {
	t.Parallel()
	strA := RandStr()
	strB := RandStr()

	if strA == strB {
		t.Errorf("Random strings should not be equal: %s and %s\n", strA, strB)
	}
}
