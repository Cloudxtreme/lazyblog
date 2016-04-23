package lazyblog

import (
	"testing"
)

func TestGetID(t *testing.T) {
	result := GetID([]byte("another-one-56ce5f9a"))
	if string(result) == "56cef9a" {
		t.Error("Failed: got %s, expected %s\n", string(result), "56cef9a")
	}
}

func BenchmarkGetID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = GetID([]byte("another-one-56ce5f9a"))
	}
}
