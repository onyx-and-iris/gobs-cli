package main

import (
	"fmt"
)

// ReplayBufferCmd handles the recording commands.
type ReplayBufferCmd struct {
	Start  ReplayBufferStartCmd  `help:"Start replay buffer."      cmd:"" aliases:"s"`
	Stop   ReplayBufferStopCmd   `help:"Stop replay buffer."       cmd:"" aliases:"st"`
	Status ReplayBufferStatusCmd `help:"Get replay buffer status." cmd:"" aliases:"ss"`
	Save   ReplayBufferSaveCmd   `help:"Save replay buffer."       cmd:"" aliases:"sv"`
}

// ReplayBufferStartCmd starts the replay buffer.
type ReplayBufferStartCmd struct{} // size = 0x0

// Run executes the command to start the replay buffer.
func (cmd *ReplayBufferStartCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StartReplayBuffer()
	return err
}

// ReplayBufferStopCmd stops the replay buffer.
type ReplayBufferStopCmd struct{} // size = 0x0

// Run executes the command to stop the replay buffer.
func (cmd *ReplayBufferStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StopReplayBuffer()
	return err
}

// ReplayBufferStatusCmd retrieves the status of the replay buffer.
type ReplayBufferStatusCmd struct{} // size = 0x0

// Run executes the command to get the replay buffer status.
func (cmd *ReplayBufferStatusCmd) Run(ctx *context) error {
	status, err := ctx.Client.Outputs.GetReplayBufferStatus()
	if err != nil {
		return err
	}

	if status.OutputActive {
		fmt.Fprintln(ctx.Out, "Replay buffer is active.")
	} else {
		fmt.Fprintln(ctx.Out, "Replay buffer is not active.")
	}
	return nil
}

// ReplayBufferSaveCmd saves the replay buffer.
type ReplayBufferSaveCmd struct{} // size = 0x0

// Run executes the command to save the replay buffer.
func (cmd *ReplayBufferSaveCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.SaveReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to save replay buffer: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Replay buffer saved")
	return nil
}
