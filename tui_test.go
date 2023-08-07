package tui_test

import (
	"io"
	"strings"
	"testing"

	"github.com/mdwhatcott/tui"
	"github.com/mdwhatcott/tui/internal/should"
)

func TestPrompt(t *testing.T) {
	ui := tui.New(strings.NewReader("Hello, world!\n"), io.Discard)
	should.So(t, ui.Prompt(""), should.Equal, "Hello, world!")
}
func TestMultilinePrompt(t *testing.T) {
	ui := tui.New(strings.NewReader("1\n2\n3\n\n\n\n"), io.Discard)
	should.So(t, ui.MultilinePrompt(""), should.Equal, "1\n2\n3")
}
func TestYes_NoMeansNo(t *testing.T) {
	ui := tui.New(strings.NewReader("n"), io.Discard)
	should.So(t, ui.Yes("", false), should.BeFalse)
}
func TestYes_NoStillMeansNo(t *testing.T) {
	ui := tui.New(strings.NewReader("n"), io.Discard)
	should.So(t, ui.Yes("", true), should.BeFalse)
}
func TestYes_WhateverMeansNo(t *testing.T) {
	ui := tui.New(strings.NewReader("whatever"), io.Discard)
	should.So(t, ui.Yes("", true), should.BeFalse)
}
func TestYes_NothingMeansNo(t *testing.T) {
	ui := tui.New(strings.NewReader(""), io.Discard)
	should.So(t, ui.Yes("", false), should.BeFalse)
}
func TestYes_NothingMeansNoUnlessYesIsDefault(t *testing.T) {
	ui := tui.New(strings.NewReader(""), io.Discard)
	should.So(t, ui.Yes("", true), should.BeTrue)
}
func TestYes_YesMeansYes(t *testing.T) {
	ui := tui.New(strings.NewReader("y"), io.Discard)
	should.So(t, ui.Yes("", false), should.BeTrue)
}
func TestYes_YESMeansYes(t *testing.T) {
	ui := tui.New(strings.NewReader("Y"), io.Discard)
	should.So(t, ui.Yes("", false), should.BeTrue)
}
func TestConfirm_Default(t *testing.T) {
	ui := tui.New(strings.NewReader(""), io.Discard)
	should.So(t, ui.Confirm("label", "default"), should.Equal, "default")
}
func TestConfirm_Override(t *testing.T) {
	ui := tui.New(strings.NewReader("override"), io.Discard)
	should.So(t, ui.Confirm("label", "default"), should.Equal, "override")
}
func TestSelect_ZeroOptions(t *testing.T) {
	ui := tui.New(strings.NewReader(""), io.Discard)
	should.So(t, ui.Select("Nada"), should.Equal, "")
}
func TestSelect_OneOption(t *testing.T) {
	ui := tui.New(strings.NewReader(""), io.Discard)
	should.So(t, ui.Select("Nada", "only"), should.Equal, "only")
}
func TestSelect_MultipleOptions(t *testing.T) {
	ui := tui.New(strings.NewReader("2\n"), io.Discard)
	should.So(t, ui.Select("Nada", "first", "second"), should.Equal, "second")
}
func TestSelect_MultipleOptions_InvalidChoicePromptsRetry(t *testing.T) {
	ui := tui.New(strings.NewReader("invalid\n2\n"), io.Discard)
	should.So(t, ui.Select("Nada", "first", "second"), should.Equal, "second")
}
func TestSelect_CustomOption(t *testing.T) {
	ui := tui.New(strings.NewReader("3\ncustom-option\n"), io.Discard)
	should.So(t, ui.Select("Nada", "first", "second"), should.Equal, "custom-option")
}
