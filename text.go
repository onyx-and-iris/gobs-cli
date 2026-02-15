package main

import (
	"fmt"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/inputs"
)

// TextCmd provides commands for managing text inputs in OBS.
type TextCmd struct {
	Current TextCurrentCmd `cmd:"current" help:"Display current text for a text input." aliases:"c"`
	Update  TextUpdateCmd  `cmd:"update"  help:"Update the text of a text input."       aliases:"u"`
}

// TextCurrentCmd provides a command to display the current text of a text input.
type TextCurrentCmd struct {
	InputName string `arg:"" help:"Name of the text source."`
}

// Run executes the command to display the current text of a text input.
func (cmd *TextCurrentCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Inputs.GetInputSettings(
		inputs.NewGetInputSettingsParams().WithInputName(cmd.InputName),
	)
	if err != nil {
		return fmt.Errorf("failed to get input settings: %w", err)
	}

	// Check if the input is a text input
	kind := resp.InputKind
	if !strings.HasPrefix(kind, "text_") {
		return fmt.Errorf("input %s is of %s", cmd.InputName, kind)
	}

	currentText, ok := resp.InputSettings["text"]
	if !ok {
		return fmt.Errorf("input %s does not have a 'text' setting", cmd.InputName)
	}
	if currentText == "" {
		currentText = "(empty)"
	}
	fmt.Fprintf(
		ctx.Out,
		"Current text for source %s: %s\n",
		ctx.Style.Highlight(cmd.InputName),
		currentText,
	)
	return nil
}

// TextUpdateCmd provides a command to update the text of a text input.
type TextUpdateCmd struct {
	InputName string `arg:"" help:"Name of the text source."`
	NewText   string `arg:"" help:"New text to set for the source." default:""`
}

// Run executes the command to update the text of a text input.
func (cmd *TextUpdateCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Inputs.GetInputSettings(
		inputs.NewGetInputSettingsParams().WithInputName(cmd.InputName),
	)
	if err != nil {
		return fmt.Errorf("failed to get input settings: %w", err)
	}

	// Check if the input is a text input
	kind := resp.InputKind
	if !strings.HasPrefix(kind, "text_") {
		return fmt.Errorf("input %s is of %s", cmd.InputName, kind)
	}

	if _, err := ctx.Client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName:     &cmd.InputName,
		InputSettings: map[string]any{"text": &cmd.NewText},
	}); err != nil {
		return fmt.Errorf("failed to update text for source %s: %w", cmd.InputName, err)
	}

	if cmd.NewText == "" {
		cmd.NewText = "(empty)"
	}
	fmt.Fprintf(
		ctx.Out,
		"Updated text for source %s to: %s\n",
		ctx.Style.Highlight(cmd.InputName),
		cmd.NewText,
	)
	return nil
}
