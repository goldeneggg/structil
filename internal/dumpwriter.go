package internal

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"text/tabwriter"
)

// DumpWriter is the interface that wraps the basic Write, Flush and Dump method.
type DumpWriter interface {
	io.Writer
	Flush() error
	Dump(rvs ...reflect.Value) error
}

type dwImpl struct {
	tw *tabwriter.Writer
}

// Param has configurations for DumpWriter.
type Param struct {
	MinWidth int
	TabWidth int
	Padding  int
	PadChar  byte
	Flags    uint
}

// NewDumpWriter returns a new default DumpWriter that wraps tabwriter.
func NewDumpWriter() DumpWriter {
	dwp := &Param{
		MinWidth: 0,
		TabWidth: 4,
		Padding:  4,
		PadChar:  ' ',
		Flags:    0,
	}

	return NewDumpWriterWithSetupInfo(dwp, os.Stdout)
}

// NewDumpWriterWithSetupInfo returns a new default DumpWriter that wraps the Writer assigned by "wrap" arg.
func NewDumpWriterWithSetupInfo(dwp *Param, wrap io.Writer) DumpWriter {
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

// Write writes data to a dump target.
func (dw *dwImpl) Write(b []byte) (int, error) {
	return dw.tw.Write(b)
}

// Flush should be called after the last call to Write to ensure
// that any data buffered in the Writer is written to output.
func (dw *dwImpl) Flush() error {
	return dw.tw.Flush()
}

// Dump writes reflection values to a dump target with automation Flush.
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

	if _, err := dw.tw.Write([]byte(fmt.Sprintf("%s\t%s\n", "Type", "Value"))); err != nil {
		return err
	}
	if _, err := dw.tw.Write([]byte(fmt.Sprintf("%s\t%s\n", "-----", "-----"))); err != nil {
		return err
	}

	for _, d := range ds {
		if _, err := dw.tw.Write([]byte(fmt.Sprintf("%v\t%+v\n", d[0], d[1]))); err != nil {
			return err
		}
	}
	err := dw.tw.Flush()

	return err
}
