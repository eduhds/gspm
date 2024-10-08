package tui

import (
	"github.com/pterm/pterm"
)

func ShowMessage(message string) {
	pterm.Printfln(message)
}

func ShowInfo(message string) {
	pterm.Info.Printfln(pterm.LightBlue(message))
}

func ShowSuccess(message string) {
	pterm.Success.Printfln(pterm.Green(message))
}

func ShowError(message string) {
	pterm.Error.Printfln(pterm.Red(message))
}

func ShowWarning(message string) {
	pterm.Warning.Printfln(pterm.Yellow(message))
}

func ShowOptions(title string, options []string) string {
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText(title).
		WithOptions(options).
		WithMaxHeight(10).
		Show()
	pterm.Println()
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

func ShowTextInput(message string, multiline bool, defaultValue string) string {
	if multiline {
		input, _ := pterm.DefaultInteractiveTextInput.
			WithDefaultText(message).
			WithDefaultValue(defaultValue).
			WithMultiLine().
			Show()
		return input
	} else {
		input, _ := pterm.DefaultInteractiveTextInput.
			WithDefaultText(message).
			WithDefaultValue(defaultValue).
			Show()
		return input
	}
}

func ShowConfirm(message string) bool {
	result, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultText(message).
		Show()
	pterm.Println()
	return result
}

func ShowBox(text string) {
	pterm.DefaultBox.Println(text)
}

func ShowLine() {
	pterm.Println()
}

func TextSuccess(message string) {
	pterm.Printfln(pterm.Green(message))
}

func TextError(message string) {
	pterm.Printfln(pterm.Red(message))
}

func TextInfo(message string) {
	pterm.Printfln(pterm.LightBlue(message))
}

func TextWarning(message string) {
	pterm.Printfln(pterm.Yellow(message))
}
