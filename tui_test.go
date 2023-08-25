package tui_test

import (
	"io"
	"strings"
	"testing"

	"github.com/mdwhatcott/tui/v2"
	"github.com/mdwhatcott/tui/v2/internal/should"
)

func ui(input string) *tui.TUI {
	ui := tui.New()
	ui.Reader = strings.NewReader(input)
	ui.Writer = io.Discard
	return ui
}
func TestPrompt(t *testing.T) {
	should.So(t, ui("Hello, world!\n").Prompt(""), should.Equal, "Hello, world!")
}
func TestMultilinePrompt(t *testing.T) {
	should.So(t, ui("1\n2\n3\n\n\n\n").MultilinePrompt(""), should.Equal, "1\n2\n3")
}
func TestNoYes(t *testing.T) {
	should.So(t, ui("n").NoYes(""), should.BeFalse)
	should.So(t, ui("n").NoYes(""), should.BeFalse)
	should.So(t, ui("whatever").NoYes(""), should.BeFalse)
	should.So(t, ui("").NoYes(""), should.BeFalse)
}
func TestYesNo(t *testing.T) {
	should.So(t, ui("").YesNo(""), should.BeTrue)
	should.So(t, ui("y").YesNo(""), should.BeTrue)
	should.So(t, ui("Y").YesNo(""), should.BeTrue)
}
func TestConfirm(t *testing.T) {
	should.So(t, ui("").Confirm("label", "default"), should.Equal, "default")
	should.So(t, ui("override").Confirm("label", "default"), should.Equal, "override")
}
func TestSelect(t *testing.T) {
	should.So(t, func() { ui("").Select("-") }, should.Panic)
	should.So(t, func() { ui("").Select("-", "only") }, should.Panic)
	should.So(t, func() { ui(strings.Repeat("invalid\n", 100)+"1 3\n").Select("-", "A", "B", "C") }, should.Panic)
	should.So(t, ui("2\n").Select("-", "A", "B"), should.Equal, "B")
}
func TestSuggest(t *testing.T) {
	should.So(t, func() { ui("").Suggest("-") }, should.Panic)
	should.So(t, func() { ui("").Suggest("-", "only") }, should.Panic)
	should.So(t, func() { ui(strings.Repeat("invalid\n", 100)+"1 3\n").Suggest("-", "A", "B", "C") }, should.Panic)
	should.So(t, ui("invalid\n2\n").Suggest("-", "A", "B"), should.Equal, "B")
	should.So(t, ui("3\nC\n").Suggest("-", "A", "B"), should.Equal, "C")
}
func TestSelectMany(t *testing.T) {
	should.So(t, func() { ui("").SelectMany("-") }, should.Panic)
	should.So(t, func() { ui("").SelectMany("-", "only") }, should.Panic)
	should.So(t, func() { ui(strings.Repeat("invalid\n", 100)+"1 3\n").SelectMany("-", "A", "B", "C") }, should.Panic)
	should.So(t, ui("invalid\n1 3\n").SelectMany("-", "A", "B", "C"), should.Equal, []string{"A", "C"})
	should.So(t, ui("1 3\n").SelectMany("-", "A", "B", "C"), should.Equal, []string{"A", "C"})
}
