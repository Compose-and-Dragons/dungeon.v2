package ui

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer holds the glamour renderer instance
var markdownRenderer *glamour.TermRenderer

// InitMarkdownRenderer initializes the global markdown renderer with terminal-optimized settings
func InitMarkdownRenderer() error {
	var err error
	markdownRenderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(200), // Increased width to prevent aggressive line breaking
	)
	return err
}

// RenderMarkdown converts markdown content to styled terminal output and prints it
func RenderMarkdown(content string) error {
	if markdownRenderer == nil {
		if err := InitMarkdownRenderer(); err != nil {
			return err
		}
	}

	rendered, err := markdownRenderer.Render(content)
	if err != nil {
		return err
	}

	fmt.Print(rendered)
	return nil
}

// PrintMarkdown safely renders and prints markdown content with automatic fallback to plain text
func PrintMarkdown(content string) {
	if err := RenderMarkdown(content); err != nil {
		// Fallback to plain text if markdown rendering fails
		fmt.Print(content)
	}
}
