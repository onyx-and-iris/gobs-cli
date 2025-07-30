// nolint: misspell
package main

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// InputCmd provides commands to manage inputs in OBS Studio.
type InputCmd struct {
	Create InputCreateCmd `cmd:"" help:"Create input." aliases:"c"`
	Delete InputDeleteCmd `cmd:"" help:"Delete input." aliases:"d"`
	Kinds  InputKindsCmd  `cmd:"" help:"List input kinds." aliases:"k"`
	List   InputListCmd   `cmd:"" help:"List all inputs." aliases:"ls"`
	Mute   InputMuteCmd   `cmd:"" help:"Mute input."      aliases:"m"`
	Unmute InputUnmuteCmd `cmd:"" help:"Unmute input."    aliases:"um"`
	Show   InputShowCmd   `cmd:"" help:"Show input details." aliases:"s"`
	Toggle InputToggleCmd `cmd:"" help:"Toggle input."    aliases:"tg"`
	Update InputUpdateCmd `cmd:"" help:"Update input settings." aliases:"up"`
}

// InputCreateCmd provides a command to create an input.
type InputCreateCmd struct {
	Name string `arg:"" help:"Name for the input." required:""`
	Kind string `arg:"" help:"Input kind (e.g., coreaudio_input_capture, macos-avcapture)." required:""`
}

// InputDeleteCmd provides a command to delete an input.
type InputDeleteCmd struct {
	Name string `arg:"" help:"Name of the input to delete." required:""`
}

// InputListCmd provides a command to list all inputs.
type InputListCmd struct {
	Input  bool `flag:"" help:"List all inputs."         aliases:"i"`
	Output bool `flag:"" help:"List all outputs."        aliases:"o"`
	Colour bool `flag:"" help:"List all colour sources." aliases:"c"`
	Ffmpeg bool `flag:"" help:"List all ffmpeg sources." aliases:"f"`
	Vlc    bool `flag:"" help:"List all VLC sources."    aliases:"v"`
	UUID   bool `flag:"" help:"Display UUIDs of inputs." aliases:"u"`
}

// Run executes the command to create an input.
func (cmd *InputCreateCmd) Run(ctx *context) error {
	currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
	if err != nil {
		return err
	}

	_, err = ctx.Client.Inputs.CreateInput(
		inputs.NewCreateInputParams().
			WithInputKind(cmd.Kind).
			WithInputName(cmd.Name).
			WithSceneName(currentScene.CurrentProgramSceneName),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out, "Created input: %s (%s) in scene %s\n",
		ctx.Style.Highlight(cmd.Name), cmd.Kind, ctx.Style.Highlight(currentScene.CurrentProgramSceneName))
	return nil
}

// Run executes the command to delete an input.
func (cmd *InputDeleteCmd) Run(ctx *context) error {
	_, err := ctx.Client.Inputs.RemoveInput(
		inputs.NewRemoveInputParams().WithInputName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to delete input: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Deleted %s\n", ctx.Style.Highlight(cmd.Name))
	return nil
}

// InputKindsCmd provides a command to list all input kinds.
type InputKindsCmd struct{}

// Run executes the command to list all input kinds.
func (cmd *InputKindsCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Inputs.GetInputKindList(
		inputs.NewGetInputKindListParams().WithUnversioned(false),
	)
	if err != nil {
		return fmt.Errorf("failed to get input kinds: %w", err)
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border))
	t.Headers("Kind")
	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 3)
		switch col {
		case 0:
			style = style.Align(lipgloss.Left)
		}
		switch {
		case row == table.HeaderRow:
			style = style.Bold(true).Align(lipgloss.Center)
		case row%2 == 0:
			style = style.Foreground(ctx.Style.evenRows)
		default:
			style = style.Foreground(ctx.Style.oddRows)
		}
		return style
	})

	for _, kind := range resp.InputKinds {
		t.Row(kind)
	}

	fmt.Fprintln(ctx.Out, t.Render())

	return nil
}

// Run executes the command to list all inputs.
func (cmd *InputListCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Inputs.GetInputList(inputs.NewGetInputListParams())
	if err != nil {
		return err
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border))
	if cmd.UUID {
		t.Headers("Input Name", "Kind", "Muted", "UUID")
	} else {
		t.Headers("Input Name", "Kind", "Muted")
	}
	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 3)
		switch col {
		case 0:
			style = style.Align(lipgloss.Left)
		case 1:
			style = style.Align(lipgloss.Left)
		case 2:
			style = style.Align(lipgloss.Center)
		case 3:
			style = style.Align(lipgloss.Left)
		}
		switch {
		case row == table.HeaderRow:
			style = style.Bold(true).Align(lipgloss.Center)
		case row%2 == 0:
			style = style.Foreground(ctx.Style.evenRows)
		default:
			style = style.Foreground(ctx.Style.oddRows)
		}
		return style
	})

	sort.Slice(resp.Inputs, func(i, j int) bool {
		return resp.Inputs[i].InputName < resp.Inputs[j].InputName
	})

	for _, input := range resp.Inputs {
		var muteMark string
		resp, err := ctx.Client.Inputs.GetInputMute(
			inputs.NewGetInputMuteParams().WithInputName(input.InputName),
		)
		if err != nil {
			if err.Error() == "request GetInputMute: InvalidResourceState (604): The specified input does not support audio." {
				muteMark = "N/A"
			} else {
				return fmt.Errorf("failed to get input mute state: %w", err)
			}
		} else {
			muteMark = getEnabledMark(resp.InputMuted)
		}

		type filter struct {
			enabled bool
			keyword string
		}
		filters := []filter{
			{cmd.Input, "input"},
			{cmd.Output, "output"},
			{cmd.Colour, "color"}, // nolint: misspell
			{cmd.Ffmpeg, "ffmpeg"},
			{cmd.Vlc, "vlc"},
		}

		var added bool
		for _, f := range filters {
			if f.enabled && strings.Contains(input.InputKind, f.keyword) {
				if cmd.UUID {
					t.Row(input.InputName, input.InputKind, muteMark, input.InputUuid)
				} else {
					t.Row(input.InputName, input.InputKind, muteMark)
				}
				added = true
				break
			}
		}

		if !added && (!cmd.Input && !cmd.Output && !cmd.Colour && !cmd.Ffmpeg && !cmd.Vlc) {
			if cmd.UUID {
				t.Row(input.InputName, snakeCaseToTitleCase(input.InputKind), muteMark, input.InputUuid)
			} else {
				t.Row(input.InputName, snakeCaseToTitleCase(input.InputKind), muteMark)
			}
		}
	}
	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// InputMuteCmd provides a command to mute an input.
type InputMuteCmd struct {
	InputName string `arg:"" help:"Name of the input to mute."`
}

// Run executes the command to mute an input.
func (cmd *InputMuteCmd) Run(ctx *context) error {
	_, err := ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(true),
	)
	if err != nil {
		return fmt.Errorf("failed to mute input: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Muted input: %s\n", ctx.Style.Highlight(cmd.InputName))
	return nil
}

// InputUnmuteCmd provides a command to unmute an input.
type InputUnmuteCmd struct {
	InputName string `arg:"" help:"Name of the input to unmute."`
}

// Run executes the command to unmute an input.
func (cmd *InputUnmuteCmd) Run(ctx *context) error {
	_, err := ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(false),
	)
	if err != nil {
		return fmt.Errorf("failed to unmute input: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Unmuted input: %s\n", ctx.Style.Highlight(cmd.InputName))
	return nil
}

// InputShowCmd provides a command to show input details.
type InputShowCmd struct {
	Name    string `arg:"" help:"Name of the input to show." required:""`
	Verbose bool   `flag:"" help:"Show all properties and their keys and values."`
}

// Run executes the command to show input details.
func (cmd *InputShowCmd) Run(ctx *context) error {
	lresp, err := ctx.Client.Inputs.GetInputList(inputs.NewGetInputListParams())
	if err != nil {
		return fmt.Errorf("failed to get input list: %w", err)
	}

	var inputKind string
	var found bool
	for _, input := range lresp.Inputs {
		if input.InputName == cmd.Name {
			inputKind = input.InputKind
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("input '%s' not found", cmd.Name)
	}

	prop, name, _, err := device(ctx.Client, cmd.Name)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border))
	t.Headers("Input Name", "Kind", "Device")
	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 3)
		switch col {
		case 0:
			style = style.Align(lipgloss.Left)
		case 1:
			style = style.Align(lipgloss.Left)
		case 2:
			style = style.Align(lipgloss.Center)
		}
		switch {
		case row == table.HeaderRow:
			style = style.Bold(true).Align(lipgloss.Center)
		case row%2 == 0:
			style = style.Foreground(ctx.Style.evenRows)
		default:
			style = style.Foreground(ctx.Style.oddRows)
		}
		return style
	})
	t.Row(cmd.Name, snakeCaseToTitleCase(inputKind), name)

	fmt.Fprintln(ctx.Out, t.Render())

	if cmd.Verbose {
		resp, err := ctx.Client.Inputs.GetInputPropertiesListPropertyItems(
			inputs.NewGetInputPropertiesListPropertyItemsParams().
				WithInputName(cmd.Name).
				WithPropertyName(prop),
		)
		if err != nil {
			return err
		}

		t := table.New().Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border))
		t.StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch col {
			case 0:
				style = style.Align(lipgloss.Left)
			}
			switch {
			case row == table.HeaderRow:
				style = style.Bold(true).Align(lipgloss.Center)
			case row%2 == 0:
				style = style.Foreground(ctx.Style.evenRows)
			default:
				style = style.Foreground(ctx.Style.oddRows)
			}
			return style
		})

		t.Headers("Devices")

		for _, item := range resp.PropertyItems {
			if item.ItemName != "" {
				t.Row(item.ItemName)
			}
		}

		fmt.Fprintln(ctx.Out, t.Render())
	}

	return nil
}

func device(c *goobs.Client, inputName string) (string, string, string, error) {
	propNames := []string{
		"device", "device_id",
	}

	for _, propName := range propNames {
		deviceListResp, err := c.Inputs.GetInputPropertiesListPropertyItems(
			inputs.NewGetInputPropertiesListPropertyItemsParams().
				WithInputName(inputName).
				WithPropertyName(propName),
		)
		if err == nil && len(deviceListResp.PropertyItems) > 0 {
			for _, item := range deviceListResp.PropertyItems {
				if item.ItemName != "" {
					return propName, item.ItemName, fmt.Sprint(item.ItemValue), nil
				}
			}
		}
	}

	return "", "", "", nil
}

// InputToggleCmd provides a command to toggle the mute state of an input.
type InputToggleCmd struct {
	InputName string `arg:"" help:"Name of the input to toggle."`
}

// Run executes the command to toggle the mute state of an input.
func (cmd *InputToggleCmd) Run(ctx *context) error {
	// Get the current mute state of the input
	resp, err := ctx.Client.Inputs.GetInputMute(
		inputs.NewGetInputMuteParams().WithInputName(cmd.InputName),
	)
	if err != nil {
		return fmt.Errorf("failed to get input mute state: %w", err)
	}
	// Toggle the mute state
	newMuteState := !resp.InputMuted
	_, err = ctx.Client.Inputs.SetInputMute(
		inputs.NewSetInputMuteParams().WithInputName(cmd.InputName).WithInputMuted(newMuteState),
	)
	if err != nil {
		return fmt.Errorf("failed to toggle input mute state: %w", err)
	}

	if newMuteState {
		fmt.Fprintf(ctx.Out, "Muted input: %s\n", ctx.Style.Highlight(cmd.InputName))
	} else {
		fmt.Fprintf(ctx.Out, "Unmuted input: %s\n", ctx.Style.Highlight(cmd.InputName))
	}
	return nil
}

// InputUpdateCmd provides a command to update input settings.
type InputUpdateCmd struct {
	InputName  string `arg:"" help:"Name of the input to update." required:""`
	DeviceName string `arg:"" help:"Name of the device to set." required:""`
}

// Run executes the command to update input settings.
func (cmd *InputUpdateCmd) Run(ctx *context) error {
	// Use the device helper to find the correct device property name
	prop, _, _, err := device(ctx.Client, cmd.InputName)
	if err != nil {
		return fmt.Errorf("failed to get device property: %w", err)
	}
	if prop == "" {
		return fmt.Errorf("no device property found for input '%s'", cmd.InputName)
	}

	resp, err := ctx.Client.Inputs.GetInputPropertiesListPropertyItems(
		inputs.NewGetInputPropertiesListPropertyItemsParams().
			WithInputName(cmd.InputName).
			WithPropertyName(prop),
	)
	if err != nil {
		return err
	}

	var deviceValue any
	var found bool
	for _, item := range resp.PropertyItems {
		if item.ItemName == cmd.DeviceName {
			deviceValue = item.ItemValue
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("device '%s' not found for input '%s'", cmd.DeviceName, cmd.InputName)
	}

	sresp, err := ctx.Client.Inputs.GetInputSettings(
		inputs.NewGetInputSettingsParams().WithInputName(cmd.InputName),
	)
	if err != nil {
		return err
	}

	settings := make(map[string]any)
	maps.Copy(settings, sresp.InputSettings)
	settings[prop] = deviceValue

	_, err = ctx.Client.Inputs.SetInputSettings(
		inputs.NewSetInputSettingsParams().
			WithInputName(cmd.InputName).
			WithInputSettings(settings),
	)
	if err != nil {
		return fmt.Errorf("failed to update input settings: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Input %s %s set to %s\n",
		ctx.Style.Highlight(cmd.InputName), prop, ctx.Style.Highlight(cmd.DeviceName))

	return nil
}
