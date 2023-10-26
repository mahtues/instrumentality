package zmisc

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type AppError struct {
	cause   error
	status  int
	message string
}

func (err *AppError) Error() string {
	if err.cause == nil {
		return err.message
	}
	return err.message + ": " + err.cause.Error()
}

var New = errors.New

var Wrapf = errors.Wrapf

func WithStatus(status int, message string) error {
	return errors.WithStack(&AppError{status: status, message: message})
}

// Although the stackTracer interface is not exported by the errors package, it
// is considered a part of its stable public interface.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Although the causer interface is not exported by the errors package, it is
// considered a part of its stable public interface.
type causer interface {
	Cause() error
}

func PrintStackTrace(err error) string {
	if err == nil {
		return ""
	}

	err = errors.Cause(err)

	lines := []string{}

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			lines = append(lines, fmt.Sprintf("%+s:%d", f, f))
		}
	}

	return strings.Join(lines, "\n")
}
