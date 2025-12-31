package main

type arguments struct {
	positional        []string
	data_filename     string
	scriptedCondition bool
	preAssignment     bool
	scriptedTask      bool
	riskAnalysis      bool
	trigger           string
	settings_filename string
}

func (args arguments) NPositional() int {
	return len(args.positional)
}
func (args arguments) isInteractiveScript() bool {
	return args.scriptedCondition || args.preAssignment || args.scriptedTask || args.riskAnalysis
}

func (args arguments) isUniqueInteractiveScriptFlag() error {
	var i uint8 = 0
	if args.scriptedCondition {
		i += 1
	}
	if args.preAssignment {
		i += 1
	}
	if args.scriptedTask {
		i += 1
	}
	if args.riskAnalysis {
		i += 1
	}
	if i > 1 {
		return ErrSeveralInteractiveFlags
	}
	if i == 0 {
		return ErrNoInteractiveFlags
	}
	return nil
}
