package decoder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func TestTypeJSON(t *testing.T) {
	t.Parallel()

	t.Run("TestTypeJSON", func(t *testing.T) {
		if d := cmp.Diff(TypeJSON.String(), "json"); d != "" {
			t.Errorf("mismatch TypeJSON.String(): (-got +want)\n%s", d)
			return
		}

		var i interface{}
		err := TypeJSON.Unmarshal(singleJSON, &i)
		if err != nil {
			t.Errorf("unexpected error is returned from TypeJSON.Unmarshal: %v", err)
			return
		}
	})
}

func TestTypeYAML(t *testing.T) {
	t.Parallel()

	t.Run("TestTypeYAML", func(t *testing.T) {
		if d := cmp.Diff(TypeYAML.String(), "yaml"); d != "" {
			t.Errorf("mismatch TypeYAML.String(): (-got +want)\n%s", d)
			return
		}

		var i interface{}
		err := TypeYAML.Unmarshal(singleYAML, &i)
		if err != nil {
			t.Errorf("unexpected error is returned from TypeYAML.Unmarshal: %v", err)
			return
		}
	})
}
