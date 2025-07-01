package main

import "github.com/fatih/color"

func PrintError(s string) {
	printer := color.New(color.FgRed)
	printer.Printf(s + "\n")
}

func PrintWarning(s string) {
	printer := color.New(color.FgYellow)
	printer.Printf(s + "\n")
}

func PrintSuccess(s string) {
	printer := color.New(color.FgGreen)
	printer.Printf(s + "\n")
}
