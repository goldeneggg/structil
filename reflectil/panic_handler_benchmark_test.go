package reflectil_test

import (
	"testing"

	. "github.com/goldeneggg/structil/reflectil"
)

func BenchmarkRecoverToError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recoverToErrorDemo()
	}
}

func recoverToErrorDemo() {
	defer func() {
		if r := recover(); r != nil {
			_ = RecoverToError(r)
		}
	}()
	panic("panic for benchmark")
}
