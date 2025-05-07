package main

import (
	"fmt"
)

// StreamCmd handles the streaming commands.
type StreamCmd struct {
	Start  StreamStartCmd  `cmd:"" help:"Start streaming."      aliases:"s"`
	Stop   StreamStopCmd   `cmd:"" help:"Stop streaming."       aliases:"st"`
	Toggle StreamToggleCmd `cmd:"" help:"Toggle streaming."     aliases:"tg"`
	Status StreamStatusCmd `cmd:"" help:"Get streaming status." aliases:"ss"`
}

// StreamStartCmd starts the stream.
type StreamStartCmd struct{} // size = 0x0

// Run executes the command to start streaming.
func (cmd *StreamStartCmd) Run(ctx *context) error {
	_, err := ctx.Client.Stream.StartStream()
	if err != nil {
		return err
	}
	return nil
}

// StreamStopCmd stops the stream.
type StreamStopCmd struct{} // size = 0x0

// Run executes the command to stop streaming.
func (cmd *StreamStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.Stream.StopStream()
	if err != nil {
		return err
	}
	return nil
}

// StreamToggleCmd toggles the stream status.
type StreamToggleCmd struct{} // size = 0x0

// Run executes the command to toggle streaming.
func (cmd *StreamToggleCmd) Run(ctx *context) error {
	status, err := ctx.Client.Stream.ToggleStream()
	if err != nil {
		return err
	}

	if status.OutputActive {
		fmt.Fprintln(ctx.Out, "Streaming started successfully.")
	} else {
		fmt.Fprintln(ctx.Out, "Streaming stopped successfully.")
	}
	return nil
}

// StreamStatusCmd retrieves the status of the stream.
type StreamStatusCmd struct{} // size = 0x0

// Run executes the command to get the stream status.
func (cmd *StreamStatusCmd) Run(ctx *context) error {
	status, err := ctx.Client.Stream.GetStreamStatus()
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Output active: %v\n", status.OutputActive)
	if status.OutputActive {
		seconds := status.OutputDuration / 1000
		minutes := int(seconds / 60)
		secondsInt := int(seconds) % 60
		if minutes > 0 {
			fmt.Fprintf(ctx.Out, "Output duration: %d minutes and %d seconds\n", minutes, secondsInt)
		} else {
			fmt.Fprintf(ctx.Out, "Output duration: %d seconds\n", secondsInt)
		}
	}
	return nil
}
