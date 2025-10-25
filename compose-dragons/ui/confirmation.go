package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type model struct {
	textInput textinput.Model
	err       error
}

// initialModel creates a new text input model with the specified prompt
func initialModel(prompt string) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 255
	ti.Width = 80

	ti.Prompt = prompt

	return model{
		textInput: ti,
		err:       nil,
	}
}

// Init initializes the model and returns the initial command
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// Handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var promptStyle lipgloss.Style

// View renders the model as a string for display
func (m model) View() string {
	return promptStyle.Render(m.textInput.View() + "\n")
}

// Input creates a colored text input prompt and returns the user's input
func Input(color, prompt string) (string, error) {
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	p := tea.NewProgram(initialModel(prompt))
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	if m, ok := m.(model); ok {
		return strings.TrimSpace(m.textInput.Value()), nil
	}
	return "", fmt.Errorf("😡 unable to get input")
}

// GetConfirmation prompts the user for yes/no confirmation with a default value and specified color
func GetConfirmation(color string, message string, defaultYes bool) bool {
	defaultText := "y"
	if !defaultYes {
		defaultText = "n"
	}

	prompt := fmt.Sprintf("%s (y/n) [%s]: ", message, defaultText)

	input, err := Input(color, prompt)
	if err != nil {
		return defaultYes
	}

	input = strings.ToLower(strings.TrimSpace(input))

	// If empty input, return default
	if input == "" {
		return defaultYes
	}

	// Check for yes/no responses
	return input == "y" || input == "yes"
}

// GetChoice prompts the user to choose from multiple options.
// It displays the available choices in the format "message (choice1/choice2/choice3) [default]: "
// and validates the input against the provided choices. If an invalid choice is entered,
// it will prompt again until a valid choice is made. The comparison is case-insensitive.
// If an empty input is provided, it returns the defaultChoice.
//
// Parameters:
//   - color: the color to use for the prompt styling
//   - message: the message to display to the user
//   - choices: a slice of valid choices the user can select from
//   - defaultChoice: the default choice if user presses enter without input
//
// Returns:
//   - string: the selected choice from the choices slice
func GetChoice(color string, message string, choices []string, defaultChoice string) string {
	// Validate that defaultChoice is in choices
	defaultValid := false
	for _, choice := range choices {
		if choice == defaultChoice {
			defaultValid = true
			break
		}
	}
	if !defaultValid && len(choices) > 0 {
		defaultChoice = choices[0]
	}

	choicesStr := strings.Join(choices, "/")
	prompt := fmt.Sprintf("%s (%s) [%s]: ", message, choicesStr, defaultChoice)

	for {
		input, err := Input(color, prompt)
		if err != nil {
			return defaultChoice
		}

		input = strings.ToLower(strings.TrimSpace(input))

		// If empty input, return default
		if input == "" {
			return defaultChoice
		}

		// Check if input matches any choice
		for _, choice := range choices {
			if strings.ToLower(choice) == input {
				return choice
			}
		}

		// Invalid choice, prompt again
		fmt.Printf("Invalid choice. Please choose from: %s\n", choicesStr)
	}
}
