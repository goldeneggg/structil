package reflectil

import (
	"fmt"
	"runtime"
	"strings"
)

// RecoverToError returns an error converted from recoverd panic information.
func RecoverToError(r interface{}) (err error) {
	if r != nil {
		msg := fmt.Sprintf("\n%v\n", r) + stackTrace()
		err = fmt.Errorf("unexpected panic occured: %s", msg)
	}

	return
}

func stackTrace() string {
	msgs := make([]string, 0, 10)

	for d := 0; ; d++ {
		pc, file, line, ok := runtime.Caller(d)
		if !ok {
			break
		}
		msgs = append(msgs, fmt.Sprintf(" -> %d: %s: %s:%d", d, runtime.FuncForPC(pc).Name(), file, line))
	}

	return strings.Join(msgs, "\n")
}
