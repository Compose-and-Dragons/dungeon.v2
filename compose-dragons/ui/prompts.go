package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
)

// UserCommand represents a parsed user command
type UserCommand struct {
	Input      string
	SkipTools  bool
	ShouldExit bool
}

// SimplePrompt creates an interactive text input prompt and returns a parsed UserCommand
func SimplePrompt(promptTitle, placeHolder string) (*UserCommand, error) {

	parseCommand := func(input string) *UserCommand {
		// Trim whitespace
		input = strings.TrimSpace(input)
		// Check for empty input
		if input == "" {
			return &UserCommand{
				Input:      "",
				SkipTools:  false,
				ShouldExit: false,
			}
		}
		// Check for /bye command
		if input == "/bye" {
			return &UserCommand{
				Input:      input,
				ShouldExit: true,
			}
		}

		return &UserCommand{
			Input:      input,
			ShouldExit: false,
		}
	}
	// Create a new form with a single text input field
	var userInput string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(promptTitle).
				Placeholder(placeHolder).
				Value(&userInput).
				ExternalEditor(false),
		),
	)

	// Run the form
	if err := form.Run(); err != nil {
		return nil, err
	}

	// Parse the command
	return parseCommand(userInput), nil
}
