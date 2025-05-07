package main

import (
	"fmt"
)

// RecordCmd handles the recording commands.
type RecordCmd struct {
	Start  RecordStartCmd  `cmd:"" help:"Start recording."       aliases:"s"`
	Stop   RecordStopCmd   `cmd:"" help:"Stop recording."        aliases:"st"`
	Toggle RecordToggleCmd `cmd:"" help:"Toggle recording."      aliases:"tg"`
	Status RecordStatusCmd `cmd:"" help:"Show recording status." aliases:"ss"`
	Pause  RecordPauseCmd  `cmd:"" help:"Pause recording."       aliases:"p"`
	Resume RecordResumeCmd `cmd:"" help:"Resume recording."      aliases:"r"`
}

// RecordStartCmd starts the recording.
type RecordStartCmd struct{} // size = 0x0

// Run executes the command to start recording.
func (cmd *RecordStartCmd) Run(ctx *context) error {
	_, err := ctx.Client.Record.StartRecord()
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Out, "Recording started successfully.")
	return nil
}

// RecordStopCmd stops the recording.
type RecordStopCmd struct{} // size = 0x0

// Run executes the command to stop recording.
func (cmd *RecordStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.Record.StopRecord()
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Out, "Recording stopped successfully.")
	return nil
}

// RecordToggleCmd toggles the recording state.
type RecordToggleCmd struct{} // size = 0x0

// Run executes the command to toggle recording.
func (cmd *RecordToggleCmd) Run(ctx *context) error {
	status, err := ctx.Client.Record.ToggleRecord()
	if err != nil {
		return err
	}

	if status.OutputActive {
		fmt.Fprintln(ctx.Out, "Recording started successfully.")
	} else {
		fmt.Fprintln(ctx.Out, "Recording stopped successfully.")
	}
	return nil
}

// RecordStatusCmd shows the recording status.
type RecordStatusCmd struct{} // size = 0x0

// Run executes the command to show recording status.
func (cmd *RecordStatusCmd) Run(ctx *context) error {
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if status.OutputActive {
		if status.OutputPaused {
			fmt.Fprintln(ctx.Out, "Recording is paused.")
		} else {
			fmt.Fprintln(ctx.Out, "Recording is in progress.")
		}
	} else {
		fmt.Fprintln(ctx.Out, "Recording is not in progress.")
	}

	return nil
}

// RecordPauseCmd pauses the recording.
type RecordPauseCmd struct{} // size = 0x0

// Run executes the command to pause recording.
func (cmd *RecordPauseCmd) Run(ctx *context) error {
	// Check if recording in progress and not already paused
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}
	if !status.OutputActive {
		return fmt.Errorf("recording is not in progress")
	}
	if status.OutputPaused {
		return fmt.Errorf("recording is already paused")
	}

	_, err = ctx.Client.Record.PauseRecord()
	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.Out, "Recording paused successfully.")
	return nil
}

// RecordResumeCmd resumes the recording.
type RecordResumeCmd struct{} // size = 0x0

// Run executes the command to resume recording.
func (cmd *RecordResumeCmd) Run(ctx *context) error {
	// Check if recording in progress and not already resumed
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}
	if !status.OutputActive {
		return fmt.Errorf("recording is not in progress")
	}
	if !status.OutputPaused {
		return fmt.Errorf("recording is not paused")
	}

	_, err = ctx.Client.Record.ResumeRecord()
	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.Out, "Recording resumed successfully.")
	return nil
}
