package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/mediainputs"
)

// MediaCmd represents a collection of commands to control media inputs.
type MediaCmd struct {
	Cursor  MediaCursorCmd  `cmd:"" help:"Get/set the cursor position of a media input." aliases:"c"`
	Play    MediaPlayCmd    `cmd:"" help:"Plays a media input."                          aliases:"p"`
	Pause   MediaPauseCmd   `cmd:"" help:"Pauses a media input."                         aliases:"pa"`
	Stop    MediaStopCmd    `cmd:"" help:"Stops a media input."                          aliases:"s"`
	Restart MediaRestartCmd `cmd:"" help:"Restarts a media input."                       aliases:"r"`
}

// MediaCursorCmd represents the command to get or set the cursor position of a media input.
type MediaCursorCmd struct {
	InputName  string `arg:"" help:"Name of the media input."`
	TimeString string `arg:"" help:"Time position to set the cursor to (e.g., '00:01:30' for 1 minute 30 seconds). If not provided, the current cursor position will be displayed." optional:""`
}

// Run executes the command to set the cursor position of the media input.
func (cmd *MediaCursorCmd) Run(ctx *context) error {
	if cmd.TimeString == "" {
		resp, err := ctx.Client.MediaInputs.GetMediaInputStatus(
			mediainputs.NewGetMediaInputStatusParams().
				WithInputName(cmd.InputName))
		if err != nil {
			return fmt.Errorf("failed to get media input cursor: %w", err)
		}

		fmt.Fprintf(
			ctx.Out,
			"%s cursor position: %s\n",
			ctx.Style.Highlight(cmd.InputName),
			formatMillisecondsToTimeString(resp.MediaCursor),
		)
		return nil
	}

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

// MediaPlayCmd represents the command to play a media input.
type MediaPlayCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to play the media input.
func (cmd *MediaPlayCmd) Run(ctx *context) error {
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

// MediaPauseCmd represents the command to pause a media input.
type MediaPauseCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to pause the media input.
func (cmd *MediaPauseCmd) Run(ctx *context) error {
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

// MediaStopCmd represents the command to stop a media input.
type MediaStopCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to stop the media input.
func (cmd *MediaStopCmd) Run(ctx *context) error {
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

// MediaRestartCmd represents the command to restart a media input.
type MediaRestartCmd struct {
	InputName string `arg:"" help:"Name of the media input."`
}

// Run executes the command to restart the media input.
func (cmd *MediaRestartCmd) Run(ctx *context) error {
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
