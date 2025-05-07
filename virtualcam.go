package main

import (
	"fmt"
)

// VirtualCamCmd handles the virtual camera commands.
type VirtualCamCmd struct {
	Start  VirtualCamStartCmd  `help:"Start virtual camera."      cmd:"" aliases:"s"`
	Stop   VirtualCamStopCmd   `help:"Stop virtual camera."       cmd:"" aliases:"st"`
	Toggle VirtualCamToggleCmd `help:"Toggle virtual camera."     cmd:"" aliases:"tg"`
	Status VirtualCamStatusCmd `help:"Get virtual camera status." cmd:"" aliases:"ss"`
}

// VirtualCamStartCmd starts the virtual camera.
type VirtualCamStartCmd struct{} // size = 0x0

// Run executes the command to start the virtual camera.
func (c *VirtualCamStartCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StartVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to start virtual camera: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Virtual camera started.")
	return nil
}

// VirtualCamStopCmd stops the virtual camera.
type VirtualCamStopCmd struct{} // size = 0x0

// Run executes the command to stop the virtual camera.
func (c *VirtualCamStopCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.StopVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to stop virtual camera: %w", err)
	}
	fmt.Fprintln(ctx.Out, "Virtual camera stopped.")
	return nil
}

// VirtualCamToggleCmd toggles the virtual camera.
type VirtualCamToggleCmd struct{} // size = 0x0

// Run executes the command to toggle the virtual camera.
func (c *VirtualCamToggleCmd) Run(ctx *context) error {
	_, err := ctx.Client.Outputs.ToggleVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to toggle virtual camera: %w", err)
	}
	return nil
}

// VirtualCamStatusCmd retrieves the status of the virtual camera.
type VirtualCamStatusCmd struct{} // size = 0x0

// Run executes the command to get the status of the virtual camera.
func (c *VirtualCamStatusCmd) Run(ctx *context) error {
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
