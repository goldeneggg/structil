package reflectil_test

import (
	"testing"

	. "github.com/goldeneggg/structil/reflectil"
)

func TestRecoverToError(t *testing.T) {
	t.Parallel()

	t.Run("recover to error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err := RecoverToError(r)
				if err == nil {
					t.Errorf("RecoverToError() did not return error. r: %v", r)
					return
				}
				t.Logf("%v", err)
			}
		}()
		panic("panic for test")
	})
}
