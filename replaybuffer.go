package main

import (
	"fmt"
)

// ReplayBufferCmd handles the recording commands.
type ReplayBufferCmd struct {
	Start  ReplayBufferStartCmd  `cmd:"" help:"Start replay buffer."      aliases:"s"  completion-command-alias-enabled:"false"`
	Stop   ReplayBufferStopCmd   `cmd:"" help:"Stop replay buffer."       aliases:"st" completion-command-alias-enabled:"false"`
	Toggle ReplayBufferToggleCmd `cmd:"" help:"Toggle replay buffer."     aliases:"tg" completion-command-alias-enabled:"false"`
	Status ReplayBufferStatusCmd `cmd:"" help:"Get replay buffer status." aliases:"ss" completion-command-alias-enabled:"false"`
	Save   ReplayBufferSaveCmd   `cmd:"" help:"Save replay buffer."       aliases:"sv" completion-command-alias-enabled:"false"`
}

// ReplayBufferStartCmd starts the replay buffer.
type ReplayBufferStartCmd struct{} // size = 0x0

// Run executes the command to start the replay buffer.
func (cmd *ReplayBufferStartCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StartReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to start replay buffer: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Replay buffer started.")
	return nil
}

// ReplayBufferStopCmd stops the replay buffer.
type ReplayBufferStopCmd struct{} // size = 0x0

// Run executes the command to stop the replay buffer.
func (cmd *ReplayBufferStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StopReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to stop replay buffer: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Replay buffer stopped.")
	return nil
}

// ReplayBufferToggleCmd toggles the replay buffer state.
type ReplayBufferToggleCmd struct{} // size = 0x0

// Run executes the command to toggle the replay buffer.
func (cmd *ReplayBufferToggleCmd) Run(ctx *context) error {
	status, err := ctx.Client.Outputs.ToggleReplayBuffer()
	if err != nil {
		return err
	}

	if status.OutputActive {
		fmt.Fprintln(ctx.Out, "Replay buffer started.")
	} else {
		fmt.Fprintln(ctx.Out, "Replay buffer stopped.")
	}
	return nil
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
