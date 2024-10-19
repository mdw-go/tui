package main

import "github.com/mdw-go/tui/v2"

func main() {
	ui := tui.New()
	ui.Println("You like:", ui.Confirm("Favorite color", "blue"))
	ui.Println(ui.Prompt("If you would like to type something, please do so now: "))
	ui.Println(ui.MultilinePrompt("Why?"))
	ui.Println("You like blue?", ui.YesNo("You like blue?"))
	ui.Println("You like red?", ui.NoYes("You like red?"))
	ui.Println("You chose:", ui.Prompt("What is your favorite color? "))
	ui.Println("You chose:", ui.Select("your favorite color", "red", "blue", "yellow"))
	ui.Println("You chose:", ui.Suggest("your favorite color", "red", "blue", "yellow"))
	ui.Println("You chose:", ui.SelectMany("your favorite colors", "red", "blue", "yellow"))
}
