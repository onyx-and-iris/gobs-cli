package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

// Mediainput represents a collection of commands to control media inputs.
type Mediainput struct {
	SetCursor MediainputSetCursorCmd `cmd:"" help:"Sets the cursor position of a media input."`
	Play      MediainputPlayCmd      `cmd:"" help:"Plays a media input."`
	Pause     MediainputPauseCmd     `cmd:"" help:"Pauses a media input."`
	Stop      MediainputStopCmd      `cmd:"" help:"Stops a media input."`
	Restart   MediainputRestartCmd   `cmd:"" help:"Restarts a media input."`
}

// MediainputSetCursorCmd represents the command to set the cursor position of a media input.
type MediainputSetCursorCmd struct {
	InputName  string `arg:"" help:"Name of the media input."`
	TimeString string `arg:"" help:"Time position to set the cursor to (e.g., '00:01:30' for 1 minute 30 seconds)."`
}

// Run executes the command to set the cursor position of the media input.
func (cmd *MediainputSetCursorCmd) Run(ctx *context) error {
	position, err := parseTimeStringToMilliseconds(cmd.TimeString)
	if err != nil {
		return fmt.Errorf("failed to parse time string: %w", err)
	}

	_, err = ctx.Client.MediaInputs.SetMediaInputCursor(
		mediainputs.NewSetMediaInputCursorParams().
			WithInputName(cmd.InputName).
			WithMediaCursor(position))
	if err != nil {
		return fmt.Errorf("failed to set media input cursor: %w", err)
	}

	fmt.Fprintf(
		ctx.Out,
		"Set %s cursor to %s (%.0f ms)\n",
		ctx.Style.Highlight(cmd.InputName),
		ctx.Style.Highlight(cmd.TimeString),
		position,
	)
	return nil
}

// MediainputPlayCmd represents the command to play a media input.
type MediainputPlayCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to play the media input.
func (cmd *MediainputPlayCmd) Run(ctx *context) error {
	_, err := ctx.Client.MediaInputs.TriggerMediaInputAction(
		mediainputs.NewTriggerMediaInputActionParams().
			WithInputName(cmd.InputName).
			WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PLAY"))
	if err != nil {
		return fmt.Errorf("failed to play media input: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Playing media input:", cmd.InputName)
	return nil
}

// MediainputPauseCmd represents the command to pause a media input.
type MediainputPauseCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to pause the media input.
func (cmd *MediainputPauseCmd) Run(ctx *context) error {
	_, err := ctx.Client.MediaInputs.TriggerMediaInputAction(
		mediainputs.NewTriggerMediaInputActionParams().
			WithInputName(cmd.InputName).
			WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PAUSE"))
	if err != nil {
		return fmt.Errorf("failed to pause media input: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Pausing media input:", cmd.InputName)
	return nil
}

// MediainputStopCmd represents the command to stop a media input.
type MediainputStopCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to stop the media input.
func (cmd *MediainputStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.MediaInputs.TriggerMediaInputAction(
		mediainputs.NewTriggerMediaInputActionParams().
			WithInputName(cmd.InputName).
			WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_STOP"))
	if err != nil {
		return fmt.Errorf("failed to stop media input: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Stopping media input:", cmd.InputName)
	return nil
}

// MediainputRestartCmd represents the command to restart a media input.
type MediainputRestartCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to restart the media input.
func (cmd *MediainputRestartCmd) Run(ctx *context) error {
	_, err := ctx.Client.MediaInputs.TriggerMediaInputAction(
		mediainputs.NewTriggerMediaInputActionParams().
			WithInputName(cmd.InputName).
			WithMediaAction("OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART"))
	if err != nil {
		return fmt.Errorf("failed to restart media input: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Restarting media input:", cmd.InputName)
	return nil
}
