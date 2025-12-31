package console

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Item interface {
	GetLabel() string
	GetValue() int
}

type ListItem struct {
	Label string
	ID    int
}

func (li ListItem) GetLabel() string { return li.Label }
func (li ListItem) GetValue() int    { return li.ID }

func SelectFromItemList(label string, options []Item, defaultvalue *int) (int, error) {
	label_list := make([]string, len(options))
	value_list := make(map[string]int)
	var defaultlabel string
	var defaultlabelptr *string

	for i, o := range options {
		if o == nil {
			continue //WTF?
		}
		label_list[i] = o.GetLabel()
		value_list[label_list[i]] = o.GetValue()

		if defaultvalue != nil && *defaultvalue == o.GetValue() {
			defaultlabel = o.GetLabel()
			defaultlabelptr = &defaultlabel
		}
	}

	if choice, err := SelectFromList(label, label_list, defaultlabelptr); err != nil {
		return 0, err
	} else if val, ok := value_list[choice]; !ok {
		return 0, fmt.Errorf("%w: %s", ErrUnknownChoice, choice)
	} else {
		return val, nil
	}
}

func SelectFromList(label string, options []string, defaultvalue *string) (string, error) {
	var text string
	text = fmt.Sprintf("%s\n", label)

	switch len(options) {
	case 0:
		return "", ErrEmptyOptionList
	case 1:
		// return if there is only one option in list
		fmt.Printf("%s :\n!! Unique option has been selected: %s\n\n", label, options[0])
		return options[0], nil
	}

	for i, opt := range options {
		if defaultvalue != nil {
			if *defaultvalue == opt {
				text = fmt.Sprintf("%s - %d*: %s\n", text, i+1, opt)
			} else {
				text = fmt.Sprintf("%s - %d : %s\n", text, i+1, opt)
			}
		} else {
			text = fmt.Sprintf("%s - %d: %s\n", text, i+1, opt)
		}
	}

	// loop until valid result
	for {
		if txt, err := GetText(text); err != nil {
			// return if getText() fails
			return "", err

		} else if txt == "" {
			// empty string
			// if use default, we return default value (used for updates)
			// otherwise we return empty string (used for creation)
			if defaultvalue != nil {
				return *defaultvalue, nil
			} else {
				return "", nil
			}

		} else if value, err := strconv.Atoi(txt); err != nil {
			// do not return if conv fails. stay in loop
			fmt.Printf("'%s' is not a valid choice.\n", txt)
		} else {
			if 0 < value && value <= len(options) {
				return options[value-1], nil
			} else {
				// do not return if out of bounds. stay in loop
				fmt.Printf("'%s' is not a valid choice.\n", txt)
			}
		}
	}
}

const (
	cmdALL    = "all"
	cmdNONE   = "none"
	cmdCANCEL = "cancel"
	cmdOK     = "ok"
)

func createPrintedList(options []string, selection []string) string {
	box := ""
	list := []string{}

	for i, opt := range options {
		if slices.Contains(selection, opt) {
			box = "[X]"
		} else {
			box = "[ ]"
		}
		list = append(list, fmt.Sprintf(" %s %d: %s", box, i+1, opt))
	}
	// add commands
	list = append(list, "------")

	if len(options) != len(selection) {
		list = append(list, fmt.Sprintf(" - %s    : Select all", cmdALL))
	}
	if len(selection) != 0 {
		list = append(list, fmt.Sprintf(" - %s   : Select none", cmdNONE))
	}
	list = append(list, fmt.Sprintf(" - %s     : Save selection", cmdOK))
	list = append(list, fmt.Sprintf(" - %s : Cancel\n", cmdCANCEL))

	return strings.Join(list, "\n")
}

func SelectManyFromList(label string, options []string, selection []string) ([]string, error) {
	var (
		text string
	)
	if selection == nil {
		selection = make([]string, 0)
	}
	for {
		text = label + "\n" + createPrintedList(options, selection)
		if txt, err := GetText(text); err != nil {
			// return if getText() fails
			return nil, err

		} else {
			switch txt {
			case cmdALL:
				selection = make([]string, len(options))
				copy(selection, options)
			case cmdNONE:
				selection = []string{}
			case cmdCANCEL:
				return nil, nil
			case cmdOK:
				return selection, nil
			default:
				if value, err := strconv.Atoi(txt); err != nil {
					// do not return if conv fails. stay in loop
					fmt.Printf("'%s' is not a valid choice.\n", txt)
				} else {
					if 0 < value && value <= len(options) {
						chosen_item := options[value-1]
						if index := slices.Index(selection, chosen_item); index == -1 {
							selection = append(selection, chosen_item)
						} else {
							selection = slices.Delete(selection, index, index+1)
						}
					} else {
						// do not return if out of bounds. stay in loop
						fmt.Printf("'%s' is not a valid choice.\n", txt)
					}
				}
			}
		}
		fmt.Println()
	}
}
