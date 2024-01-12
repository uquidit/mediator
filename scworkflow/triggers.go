package scworkflow

import (
	"bytes"
	"encoding/json"
	"strings"
)

type SecurechangeTrigger int

// WARNING!!!
// Order is VERY IMPORTANT here!
// triggers below ADVANCE have only one script
// Other triggers need a step to get the appropriate script
const (
	NO_TRIGGER SecurechangeTrigger = iota
	CREATE
	CLOSE
	CANCEL
	REJECT
	RESUBMIT
	RESOLVE
	ADVANCE
	REDO
	REOPEN
	AUTOMATION_FAILED
	LAST_TRIGGER
)

const (
	NEED_STEP_THRESHOLD = ADVANCE
)

var scTriggersToString = map[SecurechangeTrigger]string{
	NO_TRIGGER:        "",
	CREATE:            "Create",
	CLOSE:             "Close",
	CANCEL:            "Cancel",
	REJECT:            "Reject",
	ADVANCE:           "Advance",
	REDO:              "Redo",
	RESUBMIT:          "Resubmit",
	REOPEN:            "Reopen",
	RESOLVE:           "Resolve",
	AUTOMATION_FAILED: "Automatic step failed",
}

var scTriggersToID = map[string]SecurechangeTrigger{
	"Create":                CREATE,
	"Close":                 CLOSE,
	"Cancel":                CANCEL,
	"Reject":                REJECT,
	"Advance":               ADVANCE,
	"Redo":                  REDO,
	"Resubmit":              RESUBMIT,
	"Reopen":                REOPEN,
	"Resolve":               RESOLVE,
	"Automatic step failed": AUTOMATION_FAILED,
	"":                      NO_TRIGGER,
}

func (t SecurechangeTrigger) String() string {
	if s, ok := scTriggersToString[t]; ok {
		return s
	} else {
		return "unknown"
	}
}

func (t SecurechangeTrigger) Slug() string {
	if t == AUTOMATION_FAILED {
		return "AUTOMATION_FAILED"
	}
	if s, ok := scTriggersToString[t]; ok {
		return strings.ToUpper(s)
	} else {
		return "UNKNOWN"
	}
}

func (t SecurechangeTrigger) NeedStepToGetScript() bool {
	return t >= NEED_STEP_THRESHOLD
}

func (t SecurechangeTrigger) UseNextStep() bool {
	return t == ADVANCE /*|| t == CREATE*/
}

func (t SecurechangeTrigger) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(scTriggersToString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *SecurechangeTrigger) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*t = scTriggersToID[j]
	return nil
}

func GetTriggerFromString(t string) SecurechangeTrigger {
	for trigger_name, trigger_id := range scTriggersToID {
		if strings.EqualFold(t, trigger_name) {
			return trigger_id
		}
	}
	return NO_TRIGGER
}
