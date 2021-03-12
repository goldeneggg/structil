package decoder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

func TestTypeJSON(t *testing.T) {
	t.Parallel()

	if d := cmp.Diff(TypeJSON.String(), "json"); d != "" {
		t.Fatalf("mismatch TypeJSON.String(): (-got +want)\n%s", d)
	}

	var i interface{}
	err := TypeJSON.Unmarshal(singleJSON, &i)
	if err != nil {
		t.Fatalf("unexpected error is returned from TypeJSON.Unmarshal: %v", err)
	}
}

func TestTypeYAML(t *testing.T) {
	t.Parallel()

	if d := cmp.Diff(TypeYAML.String(), "yaml"); d != "" {
		t.Fatalf("mismatch TypeYAML.String(): (-got +want)\n%s", d)
	}

	var i interface{}
	err := TypeYAML.Unmarshal(singleYAML, &i)
	if err != nil {
		t.Fatalf("unexpected error is returned from TypeYAML.Unmarshal: %v", err)
	}
}

func TestTypeInvalid(t *testing.T) {
	t.Parallel()

	if d := cmp.Diff(TypeInvalid.String(), ""); d != "" {
		t.Fatalf("mismatch TypeInvalid.String(): (-got +want)\n%s", d)
	}

	var i interface{}
	err := TypeInvalid.Unmarshal(singleJSON, &i)
	if err == nil {
		t.Fatalf("expected error did not occur from TypeInvalid.Unmarshal: %v", err)
	}
}
