package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// SettingsCmd handles settings management.
type SettingsCmd struct {
	Show          SettingsShowCmd          `help:"Show settings."                     cmd:"" aliases:"s"`
	Profile       SettingsProfileCmd       `help:"Get/Set profile parameter setting." cmd:"" aliases:"p"`
	StreamService SettingsStreamServiceCmd `help:"Get/Set stream service setting."    cmd:"" aliases:"ss"`
	Video         SettingsVideoCmd         `help:"Get/Set video setting."             cmd:"" aliases:"v"`
}

// SettingsShowCmd shows the video settings.
type SettingsShowCmd struct {
	Video   bool `flag:"" help:"Show video settings."`
	Record  bool `flag:"" help:"Show record directory."`
	Profile bool `flag:"" help:"Show profile parameters."`
}

// Run executes the show command.
// nolint: misspell
func (cmd *SettingsShowCmd) Run(ctx *context) error {
	if !cmd.Video && !cmd.Record && !cmd.Profile {
		cmd.Video = true
		cmd.Record = true
		cmd.Profile = true
	}

	// Get video settings
	videoResp, err := ctx.Client.Config.GetVideoSettings()
	if err != nil {
		return fmt.Errorf("failed to get video settings: %w", err)
	}

	vt := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Video Setting", "Value").
		StyleFunc(func(row, _ int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
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

	vt.Row("Base Width", fmt.Sprintf("%.0f", videoResp.BaseWidth))
	vt.Row("Base Height", fmt.Sprintf("%.0f", videoResp.BaseHeight))
	vt.Row("Output Width", fmt.Sprintf("%.0f", videoResp.OutputWidth))
	vt.Row("Output Height", fmt.Sprintf("%.0f", videoResp.OutputHeight))
	vt.Row("FPS Numerator", fmt.Sprintf("%.0f", videoResp.FpsNumerator))
	vt.Row("FPS Denominator", fmt.Sprintf("%.0f", videoResp.FpsDenominator))

	// Get record directory
	dirResp, err := ctx.Client.Config.GetRecordDirectory()
	if err != nil {
		return fmt.Errorf("failed to get record directory: %w", err)
	}

	rt := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Record Setting", "Value").
		StyleFunc(func(row, _ int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
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

	rt.Row("Directory", dirResp.RecordDirectory)

	// Get profile prameters
	pt := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Profile Parameter", "Value").
		StyleFunc(func(row, _ int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
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

	// Common profile parameters to display
	params := []struct {
		category string
		name     string
		label    string
	}{
		{"Output", "Mode", "Output Mode"},

		{"SimpleOutput", "StreamEncoder", "Simple Streaming Encoder"},
		{"SimpleOutput", "RecEncoder", "Simple Recording Encoder"},
		{"SimpleOutput", "RecFormat2", "Simple Recording Video Format"},
		{"SimpleOutput", "RecAudioEncoder", "Simple Recording Audio Format"},
		{"SimpleOutput", "RecQuality", "Simple Recording Quality"},

		{"AdvOut", "Encoder", "Advanced Streaming Encoder"},
		{"AdvOut", "RecEncoder", "Advanced Recording Encoder"},
		{"AdvOut", "RecType", "Advanced Recording Type"},
		{"AdvOut", "RecFormat2", "Advanced Recording Video Format"},
		{"AdvOut", "RecAudioEncoder", "Advanced Recording Audio Format"},
	}

	for _, param := range params {
		resp, err := ctx.Client.Config.GetProfileParameter(
			config.NewGetProfileParameterParams().
				WithParameterCategory(param.category).
				WithParameterName(param.name),
		)
		if err == nil && resp.ParameterValue != "" {
			pt.Row(param.label, resp.ParameterValue)
		}
	}

	if cmd.Video {
		fmt.Fprintln(ctx.Out, vt.Render())
	}

	if cmd.Record {
		fmt.Fprintln(ctx.Out, rt.Render())
	}

	if cmd.Profile {
		fmt.Fprintln(ctx.Out, pt.Render())
	}

	return nil
}

// SettingsProfileCmd gets/ sets a profile parameter.
type SettingsProfileCmd struct {
	Category string `arg:"" help:"Parameter category (e.g., AdvOut, SimpleOutput, Output)." required:""`
	Name     string `arg:"" help:"Parameter name (e.g., RecFormat2, RecEncoder)."           required:""`
	Value    string `arg:"" help:"Parameter value to set."                                              optional:""`
}

// Run executes the set command.
func (cmd *SettingsProfileCmd) Run(ctx *context) error {
	if cmd.Value == "" {
		resp, err := ctx.Client.Config.GetProfileParameter(
			config.NewGetProfileParameterParams().
				WithParameterCategory(cmd.Category).
				WithParameterName(cmd.Name),
		)
		if err != nil {
			return fmt.Errorf("failed to get parameter %s.%s: %w", cmd.Category, cmd.Name, err)
		}

		fmt.Fprintf(ctx.Out, "%s.%s = %s\n", cmd.Category, cmd.Name, resp.ParameterValue)
		return nil
	}

	_, err := ctx.Client.Config.SetProfileParameter(
		config.NewSetProfileParameterParams().
			WithParameterCategory(cmd.Category).
			WithParameterName(cmd.Name).
			WithParameterValue(cmd.Value),
	)
	if err != nil {
		return fmt.Errorf("failed to set parameter %s.%s: %w", cmd.Category, cmd.Name, err)
	}

	fmt.Fprintf(ctx.Out, "Set %s.%s = %s\n", cmd.Category, cmd.Name, cmd.Value)
	return nil
}

// SettingsStreamServiceCmd gets/ sets stream service settings.
type SettingsStreamServiceCmd struct {
	Type   string `arg:"" help:"Stream type (e.g., rtmp_common, rtmp_custom)." required:""`
	Key    string `       help:"Stream key."                                               flag:""`
	Server string `       help:"Stream server URL."                                        flag:""`
}

// Run executes the set stream service command.
// nolint: misspell
func (cmd *SettingsStreamServiceCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Config.GetStreamServiceSettings()
	if err != nil {
		return fmt.Errorf("failed to get stream service settings: %w", err)
	}

	if cmd.Key == "" && cmd.Server == "" {
		t := table.New().Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
			Headers("Stream Service Setting", "Value").
			StyleFunc(func(row, _ int) lipgloss.Style {
				style := lipgloss.NewStyle().Padding(0, 3)
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

		t.Row("Type", cmd.Type)
		t.Row("Key", resp.StreamServiceSettings.Key)
		t.Row("Server", resp.StreamServiceSettings.Server)

		fmt.Fprintln(ctx.Out, t.Render())
		return nil
	}

	if cmd.Key == "" {
		cmd.Key = resp.StreamServiceSettings.Key
	}
	if cmd.Server == "" {
		cmd.Server = resp.StreamServiceSettings.Server
	}

	_, err = ctx.Client.Config.SetStreamServiceSettings(
		config.NewSetStreamServiceSettingsParams().
			WithStreamServiceSettings(&typedefs.StreamServiceSettings{
				Key:    cmd.Key,
				Server: cmd.Server,
			}).
			WithStreamServiceType(cmd.Type),
	)
	if err != nil {
		return fmt.Errorf("failed to set stream service settings: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Stream service settings updated successfully.")
	return nil
}

// SettingsVideoCmd gets/ sets video settings.
type SettingsVideoCmd struct {
	Show         bool `flag:"" help:"Show video settings."`
	BaseWidth    int  `flag:"" help:"Base (canvas) width."           min:"8"`
	BaseHeight   int  `flag:"" help:"Base (canvas) height."          min:"8"`
	OutputWidth  int  `flag:"" help:"Output (scaled) width."         min:"8"`
	OutputHeight int  `flag:"" help:"Output (scaled) height."        min:"8"`
	FPSNum       int  `flag:"" help:"Frames per second numerator."   min:"1"`
	FPSDen       int  `flag:"" help:"Frames per second denominator." min:"1"`
}

// Run executes the gets/ set video command.
// nolint: misspell
func (cmd *SettingsVideoCmd) Run(ctx *context) error {
	resp, err := ctx.Client.Config.GetVideoSettings()
	if err != nil {
		return fmt.Errorf("failed to get video settings: %w", err)
	}

	if cmd.Show {
		t := table.New().Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
			Headers("Video Setting", "Value").
			StyleFunc(func(row, _ int) lipgloss.Style {
				style := lipgloss.NewStyle().Padding(0, 3)
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

		t.Row("Base Width", fmt.Sprintf("%.0f", resp.BaseWidth))
		t.Row("Base Height", fmt.Sprintf("%.0f", resp.BaseHeight))
		t.Row("Output Width", fmt.Sprintf("%.0f", resp.OutputWidth))
		t.Row("Output Height", fmt.Sprintf("%.0f", resp.OutputHeight))
		t.Row("FPS Numerator", fmt.Sprintf("%.0f", resp.FpsNumerator))
		t.Row("FPS Denominator", fmt.Sprintf("%.0f", resp.FpsDenominator))

		fmt.Fprintln(ctx.Out, t.Render())
		return nil
	}

	if cmd.BaseWidth == 0 {
		cmd.BaseWidth = int(resp.BaseWidth)
	}
	if cmd.BaseHeight == 0 {
		cmd.BaseHeight = int(resp.BaseHeight)
	}
	if cmd.OutputWidth == 0 {
		cmd.OutputWidth = int(resp.OutputWidth)
	}
	if cmd.OutputHeight == 0 {
		cmd.OutputHeight = int(resp.OutputHeight)
	}
	if cmd.FPSNum == 0 {
		cmd.FPSNum = int(resp.FpsNumerator)
	}
	if cmd.FPSDen == 0 {
		cmd.FPSDen = int(resp.FpsDenominator)
	}

	_, err = ctx.Client.Config.SetVideoSettings(
		config.NewSetVideoSettingsParams().
			WithBaseWidth(float64(cmd.BaseWidth)).
			WithBaseHeight(float64(cmd.BaseHeight)).
			WithOutputWidth(float64(cmd.OutputWidth)).
			WithOutputHeight(float64(cmd.OutputHeight)).
			WithFpsNumerator(float64(cmd.FPSNum)).
			WithFpsDenominator(float64(cmd.FPSDen)),
	)
	if err != nil {
		return fmt.Errorf("failed to set video settings: %w", err)
	}

	fmt.Fprintln(ctx.Out, "Video settings updated successfully.")
	return nil
}
