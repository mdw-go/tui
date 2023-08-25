package tui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mdwhatcott/tui/v2/internal/just"
)

type TUI struct {
	io.Reader
	io.Writer
}

func New() *TUI {
	return &TUI{
		Reader: os.Stdin,
		Writer: os.Stderr,
	}
}

func (this *TUI) Print(args ...any)                 { _, _ = fmt.Fprint(this.Writer, args...) }
func (this *TUI) Printf(format string, args ...any) { _, _ = fmt.Fprintf(this.Writer, format, args...) }
func (this *TUI) Println(args ...any)               { _, _ = fmt.Fprintln(this.Writer, args...) }

// Confirm displays a default selection but allows the user to overwrite it.
func (this *TUI) Confirm(label, value string) string {
	prompt := fmt.Sprintf("%s: [%s] Press <ENTER> or type another %s: ", label, value, label)
	return just.Coalesce(this.Prompt(prompt), value)
}

// Prompt returns one line of text submitted via stdin.
func (this *TUI) Prompt(message string) string {
	this.Print(message)
	return line(bufio.NewScanner(this.Reader))
}

// MultilinePrompt returns newline-separated, trimmed text from stdin until 3+ blank lines are submitted.
func (this *TUI) MultilinePrompt(intro string) string {
	this.Println(intro + multilinePrompt)
	var result strings.Builder
	scanner := bufio.NewScanner(this.Reader)
	for blanks := 0; blanks < 3; {
		text := line(scanner)
		if text == "" {
			blanks++
		} else {
			blanks = 0
		}
		result.WriteString(text + "\n")
	}
	return strings.TrimSpace(result.String())
}

// YesNo returns whether the user indicates a 'yes' answer (and 'yes' is default).
func (this *TUI) YesNo(question string) bool {
	return strings.ToLower(this.Prompt(fmt.Sprintf("%s [Y/n] ", question))) != "n"
}

// NoYes returns whether the user indicates a 'yes' answer (but 'no' is default).
func (this *TUI) NoYes(question string) bool {
	return strings.ToLower(this.Prompt(fmt.Sprintf("%s [y/N] ", question))) == "y"
}

// Select displays a set of numbered options and allows the user to choose one of them.
func (this *TUI) Select(label string, options ...string) string {
	if len(options) <= 1 {
		panic("not enough options provided")
	}
	scanner := bufio.NewScanner(this.Reader)
	for x := 0; x < 100; x++ {
		this.Printf("Choose %s:\n", label)
		for n, option := range options {
			this.Printf("%d. %s\n", n+1, option)
		}
		this.Print("Enter the number of your choice: ")
		raw := line(scanner)
		choice := just.Value(strconv.Atoi(raw)) - 1
		if choice < 0 || choice >= len(options) {
			this.Println("Invalid choice, try again.")
			continue
		}
		return options[choice]
	}
	panic("user failed to enter a valid choice")
}

// Suggest is like Select but allows the user to input something other than what is presented.
func (this *TUI) Suggest(label string, options ...string) string {
	if len(options) <= 1 {
		panic("not enough options provided")
	}
	scanner := bufio.NewScanner(this.Reader)
	for x := 0; x < 100; x++ {
		this.Printf("Choose %s:\n", label)
		for n, option := range options {
			this.Printf("%d. %s\n", n+1, option)
		}
		this.Printf("%d. None of the above\n", len(options)+1)

		this.Print("Enter the number of your choice: ")
		raw := line(scanner)
		choice := just.Value(strconv.Atoi(raw)) - 1
		if choice < 0 || choice >= len(options)+1 {
			this.Println("Invalid choice, try again.")
			continue
		}
		if choice == len(options) {
			this.Printf("Enter %s: ", label)
			return line(scanner)
		}
		return options[choice]
	}
	panic("user failed to enter a valid choice")
}

// SelectMany displays a set of numbered options and allows the user to choose zero or more of them.
func (this *TUI) SelectMany(label string, options ...string) (results []string) {
	if len(options) <= 1 {
		panic("not enough options provided")
	}
	scanner := bufio.NewScanner(this.Reader)
	for x := 0; x < 100; x++ {
		this.Printf("Choose one or more of %s:\n", label)
		for n, option := range options {
			this.Printf("%d. %s\n", n+1, option)
		}
		this.Print("Enter the numbers of your choice (separated by ' '): ")
		raw := line(scanner)
		choices := strings.Fields(raw)
		for _, raw := range choices {
			choice := just.Value(strconv.Atoi(raw)) - 1
			if 0 <= choice && choice < len(options) {
				results = append(results, options[choice])
			}
		}
		if len(results) < len(choices) {
			this.Println("Invalid choices, try again.")
			continue
		}
		return results
	}
	panic("user failed to enter a valid choice")
}

func line(scanner *bufio.Scanner) string {
	_ = scanner.Scan()
	return scanner.Text()
}

const multilinePrompt = " (type or paste multiline text below; several consecutive blank lines will signal EOF)\n"
