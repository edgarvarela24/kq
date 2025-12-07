// Package ui provides interactive terminal prompts using the huh library.
package ui

import "github.com/charmbracelet/huh"

// SelectOne presents a filterable list of options and returns the user's selection.
func SelectOne(label string, options []string) (string, error) {
	var selection string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(label).
				Options(huh.NewOptions(options...)...).
				Value(&selection).
				Filtering(true).
				Height(15),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return selection, nil
}

// SelectLogOptions prompts the user for log streaming options.
func SelectLogOptions() (follow, timestamps, previous bool, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Follow logs? (-f)").
				Value(&follow),
			huh.NewConfirm().
				Title("Show timestamps?").
				Value(&timestamps),
			huh.NewConfirm().
				Title("Previous container logs? (-p)").
				Value(&previous),
		),
	)
	if err = form.Run(); err != nil {
		return false, false, false, err
	}
	return follow, timestamps, previous, nil
}
