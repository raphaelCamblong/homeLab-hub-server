package cache

import "fmt"

type Error struct {
	err string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.err)
}

func NewError(err string) *Error {
	return &Error{err: err}
}

var ErrNoConnection Error = Error{"no connection"}
