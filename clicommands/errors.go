package clicommands

import "errors"

var (
	ErrDoNothing = errors.New("this command does nothing: use a sub-command")
)
