package clicommands

import "fmt"

var (
	client APIClient
)

func Init(c APIClient) error {
	if c == nil {
		return fmt.Errorf("cannot init clicommands package: provided API client is null")
	}
	client = c
	return nil
}
