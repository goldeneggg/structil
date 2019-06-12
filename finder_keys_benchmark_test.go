package structil_test

import (
	"testing"

	. "github.com/goldeneggg/structil"
)

func BenchmarkNewFinderKeysFromConf_yml(b *testing.B) {
	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_test1_yml")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}

func BenchmarkNewFinderKeysFromConf_json(b *testing.B) {
	f, err := NewFinder(newTestStructPtr())
	if err != nil {
		b.Fatalf("NewFinder() occurs unexpected error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_test1_json")
		if err == nil {
			_ = f.FromKeys(fks)
			f.Reset()
		} else {
			b.Fatalf("abort benchmark because error %v occurd.", err)
		}
	}
}
