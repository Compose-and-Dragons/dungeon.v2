package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	Red     string = "#FF0000"
	Green   string = "#00FF00"
	Blue    string = "#0000FF"
	Yellow  string = "#FFFF00"
	Orange  string = "#FFA500"
	Purple  string = "#800080"
	Pink    string = "#FFC0CB"
	Brown   string = "#A52A2A"
	Black   string = "#000000"
	White   string = "#FFFFFF"
	Gray    string = "#808080"
	Cyan    string = "#00FFFF"
	Magenta string = "#FF00FF"
)

// Println prints the provided arguments with specified color styling followed by a newline
func Println(color string, strs ...any) {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	// Convert all arguments to strings
	strSlice := make([]string, len(strs))
	for i, v := range strs {
		strSlice[i] = fmt.Sprint(v)
	}

	// Join all strings and render with the style
	renderedString := textStyle.Render(strings.Join(strSlice, " "))

	// Print the rendered string with a newline
	fmt.Println(renderedString)
}

// Print prints the provided arguments with specified color styling without a newline
func Print(color string, strs ...any) {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	// Convert all arguments to strings
	strSlice := make([]string, len(strs))
	for i, v := range strs {
		strSlice[i] = fmt.Sprint(v)
	}

	// Join all strings and render with the style
	renderedString := textStyle.Render(strings.Join(strSlice, " "))

	// Print the rendered string with a newline
	fmt.Print(renderedString)
}

// Printf formats and prints text with specified color styling using printf-style formatting
func Printf(color string, format string, a ...interface{}) {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	
	// Format the string using fmt.Sprintf
	formattedString := fmt.Sprintf(format, a...)
	
	// Handle newlines properly by splitting and rendering each line
	if strings.Contains(formattedString, "\n") {
		lines := strings.Split(formattedString, "\n")
		for i, line := range lines {
			if line != "" {
				renderedLine := textStyle.Render(line)
				fmt.Print(renderedLine)
			}
			if i < len(lines)-1 {
				fmt.Print("\n")
			}
		}
	} else {
		renderedString := textStyle.Render(formattedString)
		fmt.Print(renderedString)
	}
}
