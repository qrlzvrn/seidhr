package errorz

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

type ErrStack struct {
	StackTrace []byte
	Err        error
}

func (e *ErrStack) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Error:\n %s\n", e.Err)
	fmt.Fprintf(&buf, "Trace:\n %s\n", e.StackTrace)
	return buf.String()
}

func NewErrStack(msg string) *ErrStack {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return &ErrStack{StackTrace: buf, Err: errors.New(msg)}
}
