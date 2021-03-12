package dynamicstruct_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct"
)

func BenchmarkAddString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddString("StringField")
	}
}

func BenchmarkAddStringWithTag(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStringWithTag("StringField", stringFieldTag)
	}
}

func BenchmarkAddInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInt("IntField")
	}
}

func BenchmarkAddFloa32(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFloat32("Float32Field")
	}
}

func BenchmarkAddFloa64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFloat64("Float64Field")
	}
}

func BenchmarkAddBool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddBool("BoolField")
	}

}

func BenchmarkAddMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddMap("MapField", SampleString, SampleFloat32)
	}
}

func BenchmarkAddFunc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddFunc("FuncField", []interface{}{SampleInt, SampleInt}, []interface{}{SampleBool, ErrSample})
	}
}

func BenchmarkAddChanBoth(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanBoth("ChanBothField", SampleInt)
	}
}

func BenchmarkAddChanRecv(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanRecv("ChanRecvField", SampleInt)
	}
}

func BenchmarkAddChanSend(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddChanSend("ChanSendField", SampleInt)
	}
}

func BenchmarkAddStruct(b *testing.B) {
	st := newDynamicTestStruct() // See: dynamicstruct_test.go

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStruct("StructField", st, false)
	}
}

func BenchmarkAddStructPtr(b *testing.B) {
	st := newDynamicTestStructPtr() // See: dynamicstruct_test.go

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddStructPtr("StructPtrField", st)
	}
}

func BenchmarkAddSlice(b *testing.B) {
	st := newDynamicTestStructPtr() // See: dynamicstruct_test.go

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddSlice("SliceField", st)
	}
}

func BenchmarkAddInterface(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInterface("InterfaceField", false)
	}
}

func BenchmarkAddInterfacePtr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder().AddInterface("InterfacePtrField", true)
	}
}

func BenchmarkBuild(b *testing.B) {
	builder := newDynamicTestBuilder() // See: dynamicstruct_test.go

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = builder.Build()
	}
}

func BenchmarkBuildNonPtr(b *testing.B) {
	builder := newDynamicTestBuilder() // See: dynamicstruct_test.go

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = builder.BuildNonPtr()
	}
}

func BenchmarkDefinition(b *testing.B) {
	builder := newDynamicTestBuilder() // See: dynamicstruct_test.go
	ds, _ := builder.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ds.Definition()
	}
}
