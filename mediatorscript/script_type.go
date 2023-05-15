package mediatorscript

import (
	"bytes"
	"encoding/json"
)

type ScriptType int

const (
	ScriptTrigger ScriptType = iota //this will be default value as first call to iota returns 0
	ScriptCondition
	ScriptTask
	ScriptAssignment
	ScriptAll
)

var toString = map[ScriptType]string{
	ScriptTrigger:    "Trigger script",
	ScriptCondition:  "Scripted Condition script",
	ScriptTask:       "Scripted Task script",
	ScriptAssignment: "Pre-Assignment script",
}

var toSlug = map[ScriptType]string{
	ScriptTrigger:    "trigger",
	ScriptCondition:  "scripted-condition",
	ScriptTask:       "scripted-task",
	ScriptAssignment: "pre-assignment",
}

var toID = map[string]ScriptType{
	"Trigger script":            ScriptTrigger,
	"Scripted Condition script": ScriptCondition,
	"Scripted Task script":      ScriptTask,
	"Pre-Assignment script":     ScriptAssignment,
}

var fromSlugToID = map[string]ScriptType{
	"trigger":            ScriptTrigger,
	"scripted-condition": ScriptCondition,
	"scripted-task":      ScriptTask,
	"pre-assignment":     ScriptAssignment,
}

func (ft ScriptType) String() string {
	if s, ok := toString[ft]; ok {
		return s
	} else {
		return "unknown"
	}
}

func (ft ScriptType) Slug() string {
	if s, ok := toSlug[ft]; ok {
		return s
	} else {
		return "unknown"
	}
}

func GetTypeFromSlug(s string) (ScriptType, error) {
	if t, ok := fromSlugToID[s]; ok {
		return t, nil
	}
	return ScriptAll, ErrUnknownScriptType
}

func IsScriptTypeSlug(s string) bool {
	for _, slug := range toSlug {
		if slug == s {
			return true
		}
	}
	return false
}

func (ct ScriptType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[ct])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (ct *ScriptType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*ct = toID[j]
	return nil
}
