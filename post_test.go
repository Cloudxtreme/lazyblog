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
	if !strings.Contains(string(str), "mytest1234") {
		t.Errorf("Didn't pass, got %s\n", str)
	}

}

func BenchmarkUrlfify(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Urlify("mytestif$56334+===(((")
	}
}
