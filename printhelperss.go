package main

import "github.com/fatih/color"

func PrintErrorf(s string) {
	printer := color.New(color.FgRed)
	printer.Printf(s)
}

func PrintError(s string) {
	printer := color.New(color.FgRed)
	printer.Printf(s + "\n")
}

func PrintWarningf(s string) {
	printer := color.New(color.FgYellow)
	printer.Printf(s)
}

func PrintWarning(s string) {
	printer := color.New(color.FgYellow)
	printer.Printf(s + "\n")
}

func PrintSuccessf(s string) {
	printer := color.New(color.FgGreen)
	printer.Printf(s)
}

func PrintSuccess(s string) {
	printer := color.New(color.FgGreen)
	printer.Printf(s + "\n")
}
