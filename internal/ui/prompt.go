// internal/ui/prompt.go - Interactive prompts using huh
//
// This package wraps the huh library to provide simple selection prompts.
//
// huh is from Charm (charmbracelet) — same folks behind bubbletea and lipgloss.
// It provides a nice declarative API for building interactive prompts.
//
// Docs: https://github.com/charmbracelet/huh
package ui

import "github.com/charmbracelet/huh"

// SelectOne presents a list of options and returns the user's selection.
//
// TODO: Implement this function
//
// This should:
// 1. Create a huh.Select with the given options
// 2. Set a title/label for the prompt
// 3. Run the form and return the selected value
//
// Docs to explore:
// - huh.NewSelect[string]() — creates a select prompt
// - .Title() — sets the prompt title
// - .Options() — sets the options (use huh.NewOptions())
// - huh.NewForm() — wraps fields into a runnable form
//
// Example structure:
//
//	var result string
//	form := huh.NewForm(
//	    huh.NewGroup(
//	        huh.NewSelect[string]().
//	            Title("...").
//	            Options(...).
//	            Value(&result),
//	    ),
//	)
//	err := form.Run()
//	return result, err
func SelectOne(label string, options []string) (string, error) {
	var selection string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(label).
				Options(huh.NewOptions(options...)...).
				Value(&selection),
		),
	)
	err := form.Run()
	if err != nil {
		return "", err
	}

	return selection, nil
}
