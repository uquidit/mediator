package clicommands

import (
	"fmt"
	"mediator/console"
	"mediator/mediatorscript"
	"mediator/mediatorsettings"
	"mediator/scworkflow"
)

func GetTriggerScripts() ([]string, error) {
	fmt.Print("Get list of registered trigger scripts from backend...")
	if script_list, err := getAllScriptNamesByType(); err != nil {
		fmt.Println("")
		return nil, err
	} else {
		list_with_noscript := []string{}
		for _, s := range script_list[mediatorscript.ScriptTrigger] {
			list_with_noscript = append(list_with_noscript, s.Name)
		}
		fmt.Println("    OK !")
		return list_with_noscript, nil
	}
}

func selectScript(default_script *string) string {
	var (
		selected_script string
		err             error
	)
	label := "Choose a script from list of registered trigger scripts"

	for {
		selected_script, err = console.SelectFromList(label, trigger_script_list, default_script)
		if err == nil {
			return selected_script
		}
	}
}

func selectRule(rules mediatorsettings.RulesSlice) (index int, new bool, exit bool, edit_desc bool) {
	options := []console.Item{}
	for i, s := range rules {
		item := console.ListItem{
			Label: s.String(),
			ID:    i,
		}
		options = append(options, item)
	}
	options = append(options, console.ListItem{Label: "Edit description", ID: len(options)})
	index_desc := len(options) - 1
	options = append(options, console.ListItem{Label: "New rule", ID: len(options)})
	index_new := len(options) - 1
	options = append(options, console.ListItem{Label: "Exit", ID: len(options)})
	index_exit := len(options) - 1

	var err error
	for {
		index, err = console.SelectFromItemList("\nDo you want to edit a rule, add a new one or edit description", options, &index_exit)
		if err == nil {
			switch index {
			case index_desc:
				edit_desc = true
				new = false
				exit = false
				index = -1
			case index_new:
				edit_desc = false
				new = true
				exit = false
				index = -1
			case index_exit:
				edit_desc = false
				new = false
				exit = true
				index = -1
			default:
				edit_desc = false
				new = false
				exit = false
			}
			return
		} else {
			fmt.Printf("%v\n", err)
		}
	}
}

func selectTrigger(default_trigger *scworkflow.SecurechangeTrigger) scworkflow.SecurechangeTrigger {
	options := []console.Item{}
	for i := 1; i < int(scworkflow.LAST_TRIGGER); i++ {
		item := console.ListItem{
			Label: scworkflow.SecurechangeTrigger(i).String(),
			ID:    i,
		}
		options = append(options, item)
	}

	var default_trigger_value *int = nil
	if default_trigger != nil {
		i := int(*default_trigger)
		default_trigger_value = &i
	}
	for {
		index, err := console.SelectFromItemList("Choose a trigger", options, default_trigger_value)
		if err == nil {
			return scworkflow.SecurechangeTrigger(index)
		}
	}
}

func selectStep(steps []string, default_step *string) string {
	var (
		step_list []string
		d         *string
	)
	if len(steps) < 1 { //sanity check
		return ""
	}
	step_list = steps
	d = default_step

	for {
		s, err := console.SelectFromList("Choose a step:", step_list, d)
		if err == nil {
			return s
		}
	}
}

func getNewRule(steps []string) *mediatorsettings.Rule {
	setting := mediatorsettings.Rule{}

	// select trigger
	trigger := selectTrigger(nil)
	setting.Trigger = trigger.String()

	// select a step if needed
	if trigger.NeedStepToGetScript() {
		var step string
		if trigger.UseNextStep() {
			step = selectStep(steps[1:], nil)
		} else {
			step = selectStep(steps, nil)
		}

		if step == "" {
			setting.Step = nil
		} else {
			setting.Step = &step
		}
	} else {
		setting.Step = nil
	}

	// select a script
	setting.Script = selectScript(nil)

	// get comment
	setting.Comment, _ = console.GetText("Rule comment")

	return &setting
}

func editRule(rule *mediatorsettings.Rule, steps []string) bool {
	if rule == nil {
		return false
	}

	// modify or delete?
	label := fmt.Sprintf("Do you want to modify or delete this rule: '%v'", rule)
	modify_label := "Modify"
	delete_label := "Delete"
	cancel_label := "Cancel"
	for {
		action, err := console.SelectFromList(label, []string{modify_label, delete_label, cancel_label}, &modify_label)
		if err == nil {
			if action == cancel_label {
				return false
			}
			if action == delete_label {
				for {
					if confirmation, err := console.GetBoolean(fmt.Sprintf("Are you sure you want to permanently delete rule: '%s'", rule), console.GetBooleanDefault_No); err == nil {
						return confirmation
					}
				}
			}
			break
		}
	}

	// select trigger
	default_trigger := scworkflow.GetTriggerFromString(rule.Trigger)
	trigger := selectTrigger(&default_trigger)
	rule.Trigger = trigger.String()

	// select a step if needed
	if trigger.NeedStepToGetScript() {
		var step string
		if trigger.UseNextStep() {
			step = selectStep(steps[1:], rule.Step)
		} else {
			step = selectStep(steps, rule.Step)
		}
		if step == "" {
			rule.Step = nil
		} else {
			rule.Step = &step
		}
	} else {
		rule.Step = nil
	}

	// select a script
	rule.Script = selectScript(&rule.Script)
	// get comment
	rule.Comment, _ = console.GetTextWithDefault("Rule comment", rule.Comment)
	return false
}

func selectWorkflow(SCworkflows *scworkflow.Workflows) (*scworkflow.WorkflowXML, bool) {
	options := []console.Item{}
	for _, scwf := range SCworkflows.Workflows {
		options = append(options, scwf)
	}
	exit := console.ListItem{
		Label: "Exit",
		ID:    -1,
	}
	options = append(options, exit)
	save_and_exit := console.ListItem{
		Label: "Save and exit",
		ID:    0,
	}
	options = append(options, save_and_exit)

	for {
		id, err := console.SelectFromItemList("Choose a workflow (! indicates workflows without settings):", options, nil)
		if err == nil {
			switch id {
			case -1: // exit, no save
				for {
					if confirmation, err := console.GetBoolean("Are you sure you want to exit without saving ?", console.GetBooleanDefault_No); err == nil {
						if confirmation {
							return nil, false
						} else {
							break
						}
					}
				}

			case 0: // exit and save
				return nil, true
			default:
				return SCworkflows.GetWorkflowByID(id), true
			}
		}
	}

}
