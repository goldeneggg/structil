package structil

import (
	"fmt"
	"runtime"
)

func recoverToError(r interface{}) (err error) {
	if r != nil {
		msg := fmt.Sprintf("\n%v\n", r) + stackTrace()
		err = fmt.Errorf("unexpected panic occured: %s", msg)
	}

	return
}

func stackTrace() string {
	msg := ""

	for d := 0; ; d++ {
		pc, file, line, ok := runtime.Caller(d)
		if !ok {
			break
		}
		msg = msg + fmt.Sprintf(" -> %d: %s: %s:%d\n", d, runtime.FuncForPC(pc).Name(), file, line)
	}

	return msg
}
