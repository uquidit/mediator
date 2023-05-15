package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type GetBooleanDefault int

const (
	GetBooleanDefault_No GetBooleanDefault = iota //this will be default value as first call to iota returns 0
	GetBooleanDefault_Yes
	GetBooleanDefault_None
)

func isYes(v string) bool {
	return v == "y" || v == "yes"
}
func isNo(v string) bool {
	return v == "n" || v == "no"
}

func checkBoolean(answer string, default_value GetBooleanDefault) (bool, error) {
	switch default_value {
	case GetBooleanDefault_No:
		return isYes(answer), nil
	case GetBooleanDefault_Yes:
		return !isNo(answer), nil
	default:
		if isNo(answer) {
			return false, nil
		}
		if isYes(answer) {
			return true, nil
		}
		// any other value? return an error
		return false, fmt.Errorf("bad answer")
	}
}

// Prompt user for a yes/no question.
// If default_value is
// * GetBooleanDefault_No, returns true if answer is 'y' or 'yes' (case insensitive), false otherwise
// * GetBooleanDefault_Yes, returns false if answer is 'n' or 'no' (case insensitive), true otherwise
// * GetBooleanDefault_None, expects 'y', 'yes', 'n', 'no' case insensitive and return accordingly. Loops otherwise.
// Function will add "[y/n]" suffix to provided label. 'y' or 'n' will be capitalized depending on default_value
func GetBoolean(label string, default_value GetBooleanDefault) (bool, error) {
	if label == "" {
		return false, fmt.Errorf("empty label in GetBoolean()")
	}

	// add [y/n] suffix to label. Capitalize according to default value
	switch default_value {
	case GetBooleanDefault_No:
		label = fmt.Sprintf("%s [y/N]: ", label)
	case GetBooleanDefault_Yes:
		label = fmt.Sprintf("%s [Y/n]: ", label)
	case GetBooleanDefault_None:
		label = fmt.Sprintf("%s [y/n]: ", label)
	}

	for {
		fmt.Print(label)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return false, err
		}
		answer := strings.ToLower(scanner.Text())

		if a, err := checkBoolean(answer, default_value); err == nil {
			// if no error, exit
			return a, nil
		}
		// else loop: ask question one more time
	}
}

func GetText(label string) (string, error) {
	if label != "" {
		fmt.Printf("%s : ", label)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return scanner.Text(), nil

}

func GetPassword(label string) (string, error) {
	if label != "" {
		fmt.Printf("%s : ", label)
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	// defer commands order is important
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	defer fmt.Printf("\n")

	if b, err := term.ReadPassword(int(os.Stdin.Fd())); err != nil {
		return "", err
	} else {
		return string(b), nil
	}

}
