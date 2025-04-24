package main

import (
	"fmt"
)

// VirtualCamCmd handles the virtual camera commands.
type VirtualCamCmd struct {
	Start  StartVirtualCamCmd  `help:"Start virtual camera."      cmd:"" aliases:"s"`
	Stop   StopVirtualCamCmd   `help:"Stop virtual camera."       cmd:"" aliases:"st"`
	Toggle ToggleVirtualCamCmd `help:"Toggle virtual camera."     cmd:"" aliases:"tg"`
	Status StatusVirtualCamCmd `help:"Get virtual camera status." cmd:"" aliases:"ss"`
}

// StartVirtualCamCmd starts the virtual camera.
type StartVirtualCamCmd struct{} // size = 0x0

// Run executes the command to start the virtual camera.
func (c *StartVirtualCamCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StartVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to start virtual camera: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Virtual camera started.")
	return nil
}

// StopVirtualCamCmd stops the virtual camera.
type StopVirtualCamCmd struct{} // size = 0x0

// Run executes the command to stop the virtual camera.
func (c *StopVirtualCamCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StopVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to stop virtual camera: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Virtual camera stopped.")
	return nil
}

// ToggleVirtualCamCmd toggles the virtual camera.
type ToggleVirtualCamCmd struct{} // size = 0x0

// Run executes the command to toggle the virtual camera.
func (c *ToggleVirtualCamCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.ToggleVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to toggle virtual camera: %w", err)
	}
	return nil
}

// StatusVirtualCamCmd retrieves the status of the virtual camera.
type StatusVirtualCamCmd struct{} // size = 0x0

// Run executes the command to get the status of the virtual camera.
func (c *StatusVirtualCamCmd) Run(ctx *context) error {
	status, err := ctx.Client.Outputs.GetVirtualCamStatus()
	if err != nil {
		return fmt.Errorf("failed to get virtual camera status: %w", err)
	}

	if status.OutputActive {
		fmt.Fprintln(ctx.Out, "Virtual camera is active.")
	} else {
		fmt.Fprintln(ctx.Out, "Virtual camera is inactive.")
	}
	return nil
}
