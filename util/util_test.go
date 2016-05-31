package util

import (
	"testing"
)

func TestNewID(t *testing.T) {
	t.Parallel()
	idA := NewID()
	idB := NewID()

	if len(idA) < 8 {
		t.Errorf("Failed to generate usable id, length is too short: %s\n", idA)
	}

	if idA == idB {
		t.Errorf("IDs must be unique, but %s and %s are equal.\n", idA, idB)
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

func BenchmarkNewID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewID()
	}
}
