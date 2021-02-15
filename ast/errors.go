package ast

import (
	"errors"
	"fmt"
)

// NotFounder interface is implemented by errors, that indicate that a record
// is not found.
type NotFounder interface {
	Error() string
	NotFound()
}

type notFoundError struct {
	err error
}

// NewNotFoundError wraps an error as a not found error.
func NewNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	return notFoundError{
		err: err,
	}
}

// NotFoundErrorf formats according to a format specifier and returns the string
// as a value that satisfies NotFounder.
func NotFoundErrorf(format string, a ...interface{}) error {
	return NewNotFoundError(fmt.Errorf(format, a...))
}

// NotFound implements the NotFounder interface to indicate, that the error is
// of type not found.
func (e notFoundError) NotFound() {}

func (e notFoundError) Error() string {
	return "not found: " + e.err.Error()
}

// Format implements the fmt.Formatter interface to print the error differently
// depending on the format verb.
func (e notFoundError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "not found: %+v", e.err)
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// IsNotFoundError returns true, if the provided error implements the NotFound
// interface.
func IsNotFoundError(e error) bool {
	var notFoundErr NotFounder
	return errors.As(e, &notFoundErr)
}
