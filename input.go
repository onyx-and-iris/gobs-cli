package main

import (
	"fmt"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/aquasecurity/table"
)

// InputCmd provides commands to manage inputs in OBS Studio.
type InputCmd struct {
	List   InputListCmd   `cmd:"" help:"List all inputs." aliases:"ls"`
	Mute   InputMuteCmd   `cmd:"" help:"Mute input."      aliases:"m"`
	Unmute InputUnmuteCmd `cmd:"" help:"Unmute input."    aliases:"um"`
	Toggle InputToggleCmd `cmd:"" help:"Toggle input."    aliases:"tg"`
}

// InputListCmd provides a command to list all inputs.
type InputListCmd struct {
	Input  bool `flag:"" help:"List all inputs."         aliases:"i"`
	Output bool `flag:"" help:"List all outputs."        aliases:"o"`
	Colour bool `flag:"" help:"List all colour sources." aliases:"c"`
	Ffmpeg bool `flag:"" help:"List all ffmpeg sources." aliases:"f"`
	Vlc    bool `flag:"" help:"List all VLC sources."    aliases:"v"`
}

// Run executes the command to list all inputs.
func (cmd *InputListCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Inputs.GetInputList(inputs.NewGetInputListParams())
	if err != nil {
		return err
	}

	t := table.New(ctx.Out)
	t.SetPadding(3)
	t.SetAlignment(table.AlignLeft, table.AlignLeft, table.AlignCenter)
	t.SetHeaders("Input Name", "Kind", "Muted")

	for _, input := range resp.Inputs {
		var muteMark string
		for _, kind := range []string{"input", "output", "ffmpeg", "vlc"} {
			if strings.Contains(input.InputKind, kind) {
				resp, err := ctx.Client.Inputs.GetInputMute(
					inputs.NewGetInputMuteParams().WithInputName(input.InputName),
				)
				if err != nil {
					return fmt.Errorf("failed to get input mute state: %w", err)
				}
				muteMark = getEnabledMark(resp.InputMuted)
				break
			}
		}

		type filter struct {
			enabled bool
			keyword string
		}
		filters := []filter{
			{cmd.Input, "input"},
			{cmd.Output, "output"},
			{cmd.Colour, "color"}, // nolint: misspell
			{cmd.Ffmpeg, "ffmpeg"},
			{cmd.Vlc, "vlc"},
		}

		var added bool
		for _, f := range filters {
			if f.enabled && strings.Contains(input.InputKind, f.keyword) {
				t.AddRow(input.InputName, input.InputKind, muteMark)
				added = true
				break
			}
		}

		if !added && (!cmd.Input && !cmd.Output && !cmd.Colour && !cmd.Ffmpeg && !cmd.Vlc) {
			t.AddRow(input.InputName, input.InputKind, muteMark)
		}
	}
	t.Render()
	return nil
}

// InputMuteCmd provides a command to mute an input.
type InputMuteCmd struct {
	InputName string `arg:"" help:"Name of the input to mute."`
}

// Run executes the command to mute an input.
func (cmd *InputMuteCmd) Run(ctx *context) error {
	_, err := ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(true),
	)
	if err != nil {
		return fmt.Errorf("failed to mute input: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Muted input: %s\n", cmd.InputName)
	return nil
}

// InputUnmuteCmd provides a command to unmute an input.
type InputUnmuteCmd struct {
	InputName string `arg:"" help:"Name of the input to unmute."`
}

// Run executes the command to unmute an input.
func (cmd *InputUnmuteCmd) Run(ctx *context) error {
	_, err := ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(false),
	)
	if err != nil {
		return fmt.Errorf("failed to unmute input: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Unmuted input: %s\n", cmd.InputName)
	return nil
}

// InputToggleCmd provides a command to toggle the mute state of an input.
type InputToggleCmd struct {
	InputName string `arg:"" help:"Name of the input to toggle."`
}

// Run executes the command to toggle the mute state of an input.
func (cmd *InputToggleCmd) Run(ctx *context) error {
	// Get the current mute state of the input
	resp, err := ctx.Client.Inputs.GetInputMute(
		inputs.NewGetInputMuteParams().WithInputName(cmd.InputName),
	)
	if err != nil {
		return fmt.Errorf("failed to get input mute state: %w", err)
	}
	// Toggle the mute state
	newMuteState := !resp.InputMuted
	_, err = ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(newMuteState),
	)
	if err != nil {
		return fmt.Errorf("failed to toggle input mute state: %w", err)
	}

	if newMuteState {
		fmt.Fprintf(ctx.Out, "Muted input: %s\n", cmd.InputName)
	} else {
		fmt.Fprintf(ctx.Out, "Unmuted input: %s\n", cmd.InputName)
	}
	return nil
}
