package lazyblog

import (
	"strings"
	"testing"
)

func BenchmarkNewID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewID()
	}
}

func TestUrlify(t *testing.T) {
	strA := Urlify("mytest$1234+===((()")
	strB := Urlify("my test with spaces")

	if !strings.Contains(string(strA), "mytest1234") {
		t.Errorf("Didn't pass, got %s, expected mytest1234\n", strA)
	}
	if !strings.Contains(string(strB), "my-test-with-spaces") {
		t.Errorf("Didn't pass, got %s, expected my-test-with-spaces\n", strB)
	}
}

func BenchmarkUrlfify(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Urlify("mytestif$56334+===(((")
	}
}
