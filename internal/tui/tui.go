package tui

import (
	"github.com/pterm/pterm"
)

func ShowMessage(message string) {
	pterm.Printfln(message)
}

func ShowInfo(message string) {
	pterm.Info.Printfln(pterm.Green(message))
}

func ShowOptions(title string, options []string) string {
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText(title).
		WithOptions(options).
		Show()
	return selectedOption
}

func ShowSpinner(message string) func(string) {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	return func(result string) {
		if result == "success" {
			spinner.Success()
		} else if result == "fail" {
			spinner.Fail()
		} else {
			spinner.Info()
		}
	}
}
