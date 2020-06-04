package cli

import "fmt"

// Error wraps errors for nice formatting
type Error struct {
	message   string
	origError error
}

func (e *Error) formatError(template string) string {
	if e.message == "" {
		return e.origError.Error()
	}
	if e.origError == nil {
		return e.message
	}
	return fmt.Sprintf(template, e.message, e.origError)
}

func (e *Error) Error() string {
	return e.formatError("%s: %v")
}

// NiceError returns well-formatted error message formatted
func (e *Error) NiceError() string {
	return e.formatError("%s\nError message: %v")
}

func (e *Error) Unwrap() error {
	return e.origError
}

func newError(format string, args ...interface{}) *Error {
	var origError error
	if len(args) > 0 {
		lastArg := args[len(args)-1]
		err, ok := lastArg.(error)
		if ok {
			origError = err
			args = args[:len(args)-1]
		}
	}
	message := fmt.Sprintf(format, args...)
	return &Error{message, origError}
}
