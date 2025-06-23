package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/requests/record"
)

// RecordCmd handles the recording commands.
type RecordCmd struct {
	Start     RecordStartCmd     `cmd:"" help:"Start recording."                   aliases:"s"`
	Stop      RecordStopCmd      `cmd:"" help:"Stop recording."                    aliases:"st"`
	Toggle    RecordToggleCmd    `cmd:"" help:"Toggle recording."                  aliases:"tg"`
	Status    RecordStatusCmd    `cmd:"" help:"Show recording status."             aliases:"ss"`
	Pause     RecordPauseCmd     `cmd:"" help:"Pause recording."                   aliases:"p"`
	Resume    RecordResumeCmd    `cmd:"" help:"Resume recording."                  aliases:"r"`
	Directory RecordDirectoryCmd `cmd:"" help:"Get/Set recording directory."       aliases:"d"`
	Split     RecordSplitCmd     `cmd:"" help:"Split recording."                   aliases:"sp"`
	Chapter   RecordChapterCmd   `cmd:"" help:"Create a chapter in the recording." aliases:"c"`
}

// RecordStartCmd starts the recording.
type RecordStartCmd struct{} // size = 0x0

// Run executes the command to start recording.
func (cmd *RecordStartCmd) Run(ctx *context) error {
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if status.OutputActive {
		if status.OutputPaused {
			return fmt.Errorf("recording is already in progress and paused")
		}
		return fmt.Errorf("recording is already in progress")
	}

	_, err = ctx.Client.Record.StartRecord()
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
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if !status.OutputActive {
		return fmt.Errorf("recording is not in progress")
	}

	resp, err := ctx.Client.Record.StopRecord()
	if err != nil {
		return err
	}
	fmt.Fprintf(
		ctx.Out,
		"%s",
		fmt.Sprintf("Recording stopped successfully. Output file: %s\n", ctx.Style.Highlight(resp.OutputPath)),
	)
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

// RecordDirectoryCmd sets the recording directory.
type RecordDirectoryCmd struct {
	RecordDirectory string `arg:"" help:"Directory to save recordings." default:""`
}

// Run executes the command to set the recording directory.
func (cmd *RecordDirectoryCmd) Run(ctx *context) error {
	if cmd.RecordDirectory == "" {
		resp, err := ctx.Client.Config.GetRecordDirectory()
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Current recording directory: %s\n", ctx.Style.Highlight(resp.RecordDirectory))
		return nil
	}

	_, err := ctx.Client.Config.SetRecordDirectory(
		config.NewSetRecordDirectoryParams().WithRecordDirectory(cmd.RecordDirectory),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out, "Recording directory set to: %s\n", ctx.Style.Highlight(cmd.RecordDirectory))
	return nil
}

// RecordSplitCmd splits the current recording.
type RecordSplitCmd struct{} // size = 0x0

// Run executes the command to split the recording.
func (cmd *RecordSplitCmd) Run(ctx *context) error {
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if !status.OutputActive {
		return fmt.Errorf("recording is not in progress")
	}

	_, err = ctx.Client.Record.SplitRecordFile()
	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.Out, "Recording split successfully.")
	return nil
}

// RecordChapterCmd creates a chapter in the recording.
type RecordChapterCmd struct {
	ChapterName string `arg:"" help:"Name of the chapter to create." default:""`
}

// Run executes the command to create a chapter in the recording.
func (cmd *RecordChapterCmd) Run(ctx *context) error {
	status, err := ctx.Client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if !status.OutputActive {
		return fmt.Errorf("recording is not in progress")
	}

	var params *record.CreateRecordChapterParams
	if cmd.ChapterName == "" {
		params = record.NewCreateRecordChapterParams()
	} else {
		params = record.NewCreateRecordChapterParams().WithChapterName(cmd.ChapterName)
	}

	_, err = ctx.Client.Record.CreateRecordChapter(params)
	if err != nil {
		return err
	}

	if cmd.ChapterName == "" {
		cmd.ChapterName = "unnamed"
	}

	fmt.Fprintf(ctx.Out, "Chapter %s created successfully.\n", ctx.Style.Highlight(cmd.ChapterName))
	return nil
}
