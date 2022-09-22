package main

import "github.com/fatih/color"

func PrintErrorf(s string) {
	errorPrinter := color.New(color.FgRed)
	errorPrinter.Printf(s)
}

func PrintError(s string) {
	errorPrinter := color.New(color.FgRed)
	errorPrinter.Printf(s + "\n")
}

func PrintWarningf(s string) {
	errorPrinter := color.New(color.FgYellow)
	errorPrinter.Printf(s)
}

func PrintWarning(s string) {
	errorPrinter := color.New(color.FgYellow)
	errorPrinter.Printf(s + "\n")
}

func PrintSuccessf(s string) {
	errorPrinter := color.New(color.FgGreen)
	errorPrinter.Printf(s)
}

func PrintSuccess(s string) {
	errorPrinter := color.New(color.FgGreen)
	errorPrinter.Printf(s + "\n")
}
