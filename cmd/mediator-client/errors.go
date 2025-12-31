package main

import "errors"

var (
	ErrSeveralInteractiveFlags error = errors.New("only one of the --scripted-condition --pre-assignment --scripted-task and --risk-analysis flags can be used at a time")
	ErrNoInteractiveFlags      error = errors.New("one of the --scripted-condition --pre-assignment --scripted-task or --risk-analysis flags must be selected")
)
