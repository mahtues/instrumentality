package zerrors

import (
	"errors"
	"fmt"
)

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Unwrap returns nil if the Unwrap method returns []error.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, Is examines err followed by a
// depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, As examines err followed by a
// depth-first traversal of its children.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return &errWrap{
		wrapped: err,
		message: fmt.Sprintf(format, args...),
	}
}

func Newf(format string, args ...interface{}) error {
	return &errString{
		message: fmt.Sprintf(format, args...),
	}
}

type errString struct {
	message string
}

func (err *errString) Error() string {
	return err.message
}

type errWrap struct {
	wrapped error
	message string
}

func (err *errWrap) Error() string {
	if err.wrapped == nil {
		return e.message
	}
	return err.message + ": " + err.wrapped.Error()
}

func (err *errWrap) Unwrap() error {
	return err.wrapped
}
