package dumper

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"text/tabwriter"
)

type DumpWriter interface {
	io.Writer
	Flush() error
	Dump(rvs ...reflect.Value) error
}

type dwImpl struct {
	tw *tabwriter.Writer
}

type DumpWriterParam struct {
	MinWidth int
	TabWidth int
	Padding  int
	PadChar  byte
	Flags    uint
}

func New() DumpWriter {
	dwp := &DumpWriterParam{
		MinWidth: 0,
		TabWidth: 4,
		Padding:  4,
		PadChar:  ' ',
		Flags:    0,
	}

	return NewWithSetupInfo(dwp, os.Stdout)
}

func NewWithSetupInfo(dwp *DumpWriterParam, wrap io.Writer) DumpWriter {
	dw := &dwImpl{}

	dw.tw = tabwriter.NewWriter(
		wrap,
		dwp.MinWidth,
		dwp.TabWidth,
		dwp.Padding,
		dwp.PadChar,
		dwp.Flags,
	)

	return dw
}

func (dw *dwImpl) Write(b []byte) (int, error) {
	return dw.tw.Write(b)
}

func (dw *dwImpl) Flush() error {
	return dw.tw.Flush()
}

func (dw *dwImpl) Dump(rvs ...reflect.Value) error {
	var t interface{}

	ds := make([][]interface{}, len(rvs))

	for i, rv := range rvs {
		if rv.IsValid() {
			t = rv.Type()
		} else {
			t = rv.Kind()
		}
		ds[i] = []interface{}{
			t,  // Type
			rv, // Value
		}
	}

	dw.tw.Write([]byte(fmt.Sprintf("%s\t%s\n", "Type", "Value")))
	dw.tw.Write([]byte(fmt.Sprintf("%s\t%s\n", "-----", "-----")))

	for _, d := range ds {
		dw.tw.Write([]byte(fmt.Sprintf(
			"%v\t%+v\n", d[0], d[1],
		)))
	}
	err := dw.tw.Flush()

	return err
}
