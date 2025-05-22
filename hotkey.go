package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/general"
	"github.com/andreykaipov/goobs/api/typedefs"
)

// HotkeyCmd provides commands to manage hotkeys in OBS Studio.
type HotkeyCmd struct {
	List            HotkeyListCmd            `cmd:"" help:"List all hotkeys."             aliases:"ls"`
	Trigger         HotkeyTriggerCmd         `cmd:"" help:"Trigger a hotkey by name."     aliases:"tr"`
	TriggerSequence HotkeyTriggerSequenceCmd `cmd:"" help:"Trigger a hotkey by sequence." aliases:"trs"`
}

// HotkeyListCmd provides a command to list all hotkeys.
type HotkeyListCmd struct{} // size = 0x0

// Run executes the command to list all hotkeys.
func (cmd *HotkeyListCmd) Run(ctx *context) error {
	resp, err := ctx.Client.General.GetHotkeyList()
	if err != nil {
		return err
	}

	for _, hotkey := range resp.Hotkeys {
		fmt.Fprintln(ctx.Out, hotkey)
	}
	return nil
}

// HotkeyTriggerCmd provides a command to trigger a hotkey.
type HotkeyTriggerCmd struct {
	Hotkey string `help:"Hotkey name to trigger." arg:""`
}

// Run executes the command to trigger a hotkey.
func (cmd *HotkeyTriggerCmd) Run(ctx *context) error {
	_, err := ctx.Client.General.TriggerHotkeyByName(
		general.NewTriggerHotkeyByNameParams().WithHotkeyName(cmd.Hotkey),
	)
	if err != nil {
		return err
	}
	return nil
}

// HotkeyTriggerSequenceCmd provides a command to trigger a hotkey sequence.
type HotkeyTriggerSequenceCmd struct {
	Shift bool   `flag:"" help:"Shift modifier."`
	Ctrl  bool   `flag:"" help:"Control modifier."`
	Alt   bool   `flag:"" help:"Alt modifier."`
	Cmd   bool   `flag:"" help:"Command modifier."`
	KeyID string `        help:"Key ID to trigger." arg:""`
}

// Run executes the command to trigger a hotkey sequence.
func (cmd *HotkeyTriggerSequenceCmd) Run(ctx *context) error {
	_, err := ctx.Client.General.TriggerHotkeyByKeySequence(
		general.NewTriggerHotkeyByKeySequenceParams().
			WithKeyId(cmd.KeyID).
			WithKeyModifiers(&typedefs.KeyModifiers{
				Shift:   cmd.Shift,
				Control: cmd.Ctrl,
				Alt:     cmd.Alt,
				Command: cmd.Cmd,
			}),
	)
	if err != nil {
		return err
	}
	return nil
}
