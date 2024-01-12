package console

import (
	"fmt"
	"strconv"
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
