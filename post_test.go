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
	str := Urlify("mytest$1234+===((()")
	str2 := Urlify("my test with spaces")

	if !strings.Contains(string(str), "mytest1234") {
		t.Errorf("Didn't pass, got %s, expected mytest1234\n", str)
	}
	if !strings.Contains(string(str2), "my-test-with-spaces") {
		t.Errorf("Didn't pass, got %s, expected my-test-with-spaces\n", str2)
	}
}

func BenchmarkUrlfify(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Urlify("mytestif$56334+===(((")
	}
}
