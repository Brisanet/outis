package outis

import (
	"fmt"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type wrappedError interface {
	Unwrap() error
}

// reconstructStackTrace reconstructs the stack trace
func reconstructStackTrace(err error) (output []string, traced bool) {
	var (
		wrapped wrappedError
		tracer  stackTracer
	)
	if errors.As(err, &wrapped) {
		if !traced && errors.As(err, &tracer) {
			stack := tracer.StackTrace()
			for _, frame := range stack {
				output = append(output, fmt.Sprintf("%+v", frame))
			}
			traced = true
		}
	}
	return
}
