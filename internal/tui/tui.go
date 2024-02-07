package tui

import (
	"github.com/pterm/pterm"
)

func ShowMessage(message string) {
	pterm.Info.Printfln("Selected option: %s", pterm.Green(message))
}

func ShowOptions(options []string) string {
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		Show()
	return selectedOption
}
