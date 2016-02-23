package lazyblog

import (
	"testing"
)

func BenchmarkNewID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewID()
	}
}
