package util

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	t.Parallel()
	id := NewID()

	if len(id) < 8 {
		t.Errorf("Failed to generate usable id, length is too short: %s\n", id)
	}
}

func TestUrlify(t *testing.T) {
	t.Parallel()
	strA := Urlify("mytest$1234+===((()")
	strB := Urlify("my test with spaces")

	if !strings.Contains(string(strA), "mytest1234") {
		t.Errorf("Didn't pass, got %s, expected mytest1234\n", strA)
	}
	if !strings.Contains(string(strB), "my-test-with-spaces") {
		t.Errorf("Didn't pass, got %s, expected my-test-with-spaces\n", strB)
	}
}

func TestRandStr(t *testing.T) {
	strA := RandStr()
	strB := RandStr()

	if strA == strB {
		t.Errorf("Random strings should not be equal: %s and %s\n", strA, strB)
	}
}
