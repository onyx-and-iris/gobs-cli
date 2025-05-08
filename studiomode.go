package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/ui"
)

// StudioModeCmd provides commands to manage studio mode in OBS Studio.
type StudioModeCmd struct {
	Enable  StudioModeEnableCmd  `cmd:"enable"  help:"Enable studio mode."     aliases:"on"`
	Disable StudioModeDisableCmd `cmd:"disable" help:"Disable studio mode."    aliases:"off"`
	Toggle  StudioModeToggleCmd  `cmd:"toggle"  help:"Toggle studio mode."     aliases:"tg"`
	Status  StudioModeStatusCmd  `cmd:"status"  help:"Get studio mode status." aliases:"ss"`
}

// StudioModeEnableCmd provides a command to enable studio mode.
type StudioModeEnableCmd struct{} // size = 0x0

// Run executes the command to enable studio mode.
func (cmd *StudioModeEnableCmd) Run(ctx *context) error {
	_, err := ctx.Client.Ui.SetStudioModeEnabled(ui.NewSetStudioModeEnabledParams().WithStudioModeEnabled(true))
	if err != nil {
		return fmt.Errorf("failed to enable studio mode: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Studio mode is now enabled")
	return nil
}

// StudioModeDisableCmd provides a command to disable studio mode.
type StudioModeDisableCmd struct{} // size = 0x0

// Run executes the command to disable studio mode.
func (cmd *StudioModeDisableCmd) Run(ctx *context) error {
	_, err := ctx.Client.Ui.SetStudioModeEnabled(ui.NewSetStudioModeEnabledParams().WithStudioModeEnabled(false))
	if err != nil {
		return fmt.Errorf("failed to disable studio mode: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Studio mode is now disabled")
	return nil
}

// StudioModeToggleCmd provides a command to toggle studio mode.
type StudioModeToggleCmd struct{} // size = 0x0

// Run executes the command to toggle studio mode.
func (cmd *StudioModeToggleCmd) Run(ctx *context) error {
	status, err := ctx.Client.Ui.GetStudioModeEnabled(&ui.GetStudioModeEnabledParams{})
	if err != nil {
		return fmt.Errorf("failed to get studio mode status: %w", err)
	}

	newStatus := !status.StudioModeEnabled
	_, err = ctx.Client.Ui.SetStudioModeEnabled(ui.NewSetStudioModeEnabledParams().WithStudioModeEnabled(newStatus))
	if err != nil {
		return fmt.Errorf("failed to toggle studio mode: %w", err)
	}

	if newStatus {
		fmt.Fprintln(ctx.Out, "Studio mode is now enabled")
	} else {
		fmt.Fprintln(ctx.Out, "Studio mode is now disabled")
	}

	return nil
}

// StudioModeStatusCmd provides a command to get the status of studio mode.
type StudioModeStatusCmd struct{} // size = 0x0

// Run executes the command to get the status of studio mode.
func (cmd *StudioModeStatusCmd) Run(ctx *context) error {
	status, err := ctx.Client.Ui.GetStudioModeEnabled(&ui.GetStudioModeEnabledParams{})
	if err != nil {
		return fmt.Errorf("failed to get studio mode status: %w", err)
	}
	if status.StudioModeEnabled {
		fmt.Fprintln(ctx.Out, "Studio mode is enabled")
	} else {
		fmt.Fprintln(ctx.Out, "Studio mode is disabled")
	}
	return nil
}
