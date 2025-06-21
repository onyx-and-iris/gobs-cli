package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/general"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Hotkey Name").
		StyleFunc(func(row, _ int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch {
			case row == table.HeaderRow:
				style = style.Bold(true).Align(lipgloss.Center) // nolint: misspell
			case row%2 == 0:
				style = style.Foreground(ctx.Style.evenRows)
			default:
				style = style.Foreground(ctx.Style.oddRows)
			}
			return style
		})

	for _, hotkey := range resp.Hotkeys {
		t.Row(hotkey)
	}
	fmt.Fprintln(ctx.Out, t.Render())
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
