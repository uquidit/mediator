package securechangeapi

import "errors"

var (
	ErrNoSecurechangeWorkflows error = errors.New("no activated Securechange workflows")
)
