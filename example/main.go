package main

import (
	"os"

	"github.com/mdwhatcott/tui"
)

func main() {
	ui := tui.New(os.Stdin, os.Stdout)
	color := ui.Prompt("What is your favorite color? ")
	ui.Println("You chose:", color)
}
