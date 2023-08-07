package tui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mdwhatcott/tui/internal/just"
)

type TUI struct {
	stdin  *bufio.Scanner
	stdout io.Writer
}

func New(stdin io.Reader, stdout io.Writer) *TUI {
	return &TUI{
		stdin:  bufio.NewScanner(stdin),
		stdout: stdout,
	}
}
func (this *TUI) Print(args ...any)                 { _, _ = fmt.Fprint(this.stdout, args...) }
func (this *TUI) Printf(format string, args ...any) { _, _ = fmt.Fprintf(this.stdout, format, args...) }
func (this *TUI) Println(args ...any)               { _, _ = fmt.Fprintln(this.stdout, args...) }

// Confirm displays a default selection but allows the user to overwrite it.
func (this *TUI) Confirm(label, value string) string {
	prompt := fmt.Sprintf("%s: [%s] <ENTER> to continue or type another value here: ", label, value)
	return just.Coalesce(this.Prompt(prompt), value)
}

// Prompt returns one line of text submitted via stdin.
func (this *TUI) Prompt(message string) string {
	this.Print(message)
	_ = this.stdin.Scan()
	return this.stdin.Text()
}

// MultilinePrompt returns newline-separated, trimmed text from stdin until 3+ blank lines are submitted.
func (this *TUI) MultilinePrompt(intro string) string {
	blanks := 0
	var result strings.Builder
	for text := this.Prompt(intro + multilinePrompt); blanks < 3; text = this.Prompt("") {
		if text == "" {
			blanks++
		} else {
			blanks = 0
		}
		result.WriteString(text + "\n")
	}
	return strings.TrimSpace(result.String())
}

// Yes returns whether the user indicates a 'yes' answer.
func (this *TUI) Yes(question string, defaultYes bool) bool {
	defaults := "[Y/n]"
	if !defaultYes {
		defaults = "[y/N]"
	}
	answer := strings.ToLower(this.Prompt(fmt.Sprintf("%s %s ", question, defaults)))
	return answer == "y" || answer == "" && defaultYes
}

// Select displays a set of numbered options for the user to choose from.
func (this *TUI) Select(label string, options ...string) string {
	if len(options) == 0 {
		return this.Confirm(label, "")
	}
	if len(options) == 1 {
		return this.Confirm(label, options[0])
	}
	for {
		this.Printf("Choose the %s:\n", label)
		for n, option := range options {
			this.Printf("%d. %s\n", n+1, option)
		}
		this.Printf("%d. None of the above\n", len(options)+1)

		choice := just.Value(strconv.Atoi(this.Prompt("Enter the number of your choice: "))) - 1
		if choice < 0 || choice >= len(options)+1 {
			this.Println("Invalid choice, try again.")
			continue
		}
		if choice == len(options) {
			return this.Confirm(label, "")
		}
		return options[choice]
	}
}

const multilinePrompt = " (type or paste multiline text below; 3+ consecutive blank lines signals EOF)\n"
