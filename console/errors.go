package console

import "errors"

var (
	ErrUnknownChoice error = errors.New("unkown choice")
	ErrBadAnwser     error = errors.New("bad answer")
	ErrEmptyLabel    error = errors.New("empty label")
)
