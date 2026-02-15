package main

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// FilterCmd provides commands to manage filters in OBS Studio.
type FilterCmd struct {
	List    FilterListCmd    `cmd:"" help:"List all filters."  aliases:"ls"`
	Enable  FilterEnableCmd  `cmd:"" help:"Enable filter."     aliases:"on"`
	Disable FilterDisableCmd `cmd:"" help:"Disable filter."    aliases:"off"`
	Toggle  FilterToggleCmd  `cmd:"" help:"Toggle filter."     aliases:"tg"`
	Status  FilterStatusCmd  `cmd:"" help:"Get filter status." aliases:"ss"`
}

// FilterListCmd provides a command to list all filters in a scene.
type FilterListCmd struct {
	SourceName string `arg:"" help:"Name of the source to list filters from." default:""`
}

// Run executes the command to list all filters in a scene.
// nolint: misspell
func (cmd *FilterListCmd) Run(ctx *context) error {
	if cmd.SourceName == "" {
		currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
		if err != nil {
			return fmt.Errorf("failed to get current program scene: %w", err)
		}
		cmd.SourceName = currentScene.SceneName
	}

	sourceFilters, err := ctx.Client.Filters.GetSourceFilterList(
		filters.NewGetSourceFilterListParams().WithSourceName(cmd.SourceName),
	)
	if err != nil {
		return err
	}

	if len(sourceFilters.Filters) == 0 {
		fmt.Fprintf(
			ctx.Out,
			"No filters found for source %s.\n",
			ctx.Style.Highlight(cmd.SourceName),
		)
		return nil
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Filter Name", "Kind", "Enabled", "Settings").
		StyleFunc(func(row, col int) lipgloss.Style {
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

	for _, filter := range sourceFilters.Filters {
		defaultSettings, err := ctx.Client.Filters.GetSourceFilterDefaultSettings(
			filters.NewGetSourceFilterDefaultSettingsParams().
				WithFilterKind(filter.FilterKind),
		)
		if err != nil {
			return fmt.Errorf("failed to get default settings for filter %s: %w",
				ctx.Style.Error(filter.FilterName), err)
		}
		maps.Insert(defaultSettings.DefaultFilterSettings, maps.All(filter.FilterSettings))

		var lines []string
		for k, v := range defaultSettings.DefaultFilterSettings {
			lines = append(lines, fmt.Sprintf("%s: %v", snakeCaseToTitleCase(k), v))
		}
		sort.Slice(lines, func(i, j int) bool {
			return strings.ToLower(lines[i]) < strings.ToLower(lines[j])
		})

		t.Row(
			filter.FilterName,
			snakeCaseToTitleCase(filter.FilterKind),
			getEnabledMark(filter.FilterEnabled),
			strings.Join(lines, "\n"),
		)
	}
	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// FilterEnableCmd provides a command to enable a filter in a scene.
type FilterEnableCmd struct {
	SourceName string `arg:"" help:"Name of the source to enable filter from."`
	FilterName string `arg:"" help:"Name of the filter to enable."`
}

// Run executes the command to enable a filter in a scene.
func (cmd *FilterEnableCmd) Run(ctx *context) error {
	_, err := ctx.Client.Filters.SetSourceFilterEnabled(
		filters.NewSetSourceFilterEnabledParams().
			WithSourceName(cmd.SourceName).
			WithFilterName(cmd.FilterName).
			WithFilterEnabled(true),
	)
	if err != nil {
		return fmt.Errorf("failed to enable filter %s on source %s: %w",
			ctx.Style.Error(cmd.FilterName), ctx.Style.Error(cmd.SourceName), err)
	}
	fmt.Fprintf(ctx.Out, "Filter %s enabled on source %s.\n",
		ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	return nil
}

// FilterDisableCmd provides a command to disable a filter in a scene.
type FilterDisableCmd struct {
	SourceName string `arg:"" help:"Name of the source to disable filter from."`
	FilterName string `arg:"" help:"Name of the filter to disable."`
}

// Run executes the command to disable a filter in a scene.
func (cmd *FilterDisableCmd) Run(ctx *context) error {
	_, err := ctx.Client.Filters.SetSourceFilterEnabled(
		filters.NewSetSourceFilterEnabledParams().
			WithSourceName(cmd.SourceName).
			WithFilterName(cmd.FilterName).
			WithFilterEnabled(false),
	)
	if err != nil {
		return fmt.Errorf("failed to disable filter %s on source %s: %w",
			ctx.Style.Error(cmd.FilterName), ctx.Style.Error(cmd.SourceName), err)
	}
	fmt.Fprintf(ctx.Out, "Filter %s disabled on source %s.\n",
		ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	return nil
}

// FilterToggleCmd provides a command to toggle a filter in a scene.
type FilterToggleCmd struct {
	SourceName string `arg:"" help:"Name of the source to toggle filter from."`
	FilterName string `arg:"" help:"Name of the filter to toggle."`
}

// Run executes the command to toggle a filter in a scene.
func (cmd *FilterToggleCmd) Run(ctx *context) error {
	filter, err := ctx.Client.Filters.GetSourceFilter(
		filters.NewGetSourceFilterParams().
			WithSourceName(cmd.SourceName).
			WithFilterName(cmd.FilterName),
	)
	if err != nil {
		return fmt.Errorf("failed to get filter %s on source %s: %w",
			ctx.Style.Error(cmd.FilterName), ctx.Style.Error(cmd.SourceName), err)
	}

	newStatus := !filter.FilterEnabled
	_, err = ctx.Client.Filters.SetSourceFilterEnabled(
		filters.NewSetSourceFilterEnabledParams().
			WithSourceName(cmd.SourceName).
			WithFilterName(cmd.FilterName).
			WithFilterEnabled(newStatus),
	)
	if err != nil {
		return fmt.Errorf("failed to toggle filter %s on source %s: %w",
			ctx.Style.Error(cmd.FilterName), ctx.Style.Error(cmd.SourceName), err)
	}

	if newStatus {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is now enabled.\n",
			ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	} else {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is now disabled.\n",
			ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	}
	return nil
}

// FilterStatusCmd provides a command to get the status of a filter in a scene.
type FilterStatusCmd struct {
	SourceName string `arg:"" help:"Name of the source to get filter status from."`
	FilterName string `arg:"" help:"Name of the filter to get status."`
}

// Run executes the command to get the status of a filter in a scene.
func (cmd *FilterStatusCmd) Run(ctx *context) error {
	filter, err := ctx.Client.Filters.GetSourceFilter(
		filters.NewGetSourceFilterParams().
			WithSourceName(cmd.SourceName).
			WithFilterName(cmd.FilterName),
	)
	if err != nil {
		return fmt.Errorf("failed to get status of filter %s on source %s: %w",
			ctx.Style.Error(cmd.FilterName), ctx.Style.Error(cmd.SourceName), err)
	}
	if filter.FilterEnabled {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is enabled.\n",
			ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	} else {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is disabled.\n",
			ctx.Style.Highlight(cmd.FilterName), ctx.Style.Highlight(cmd.SourceName))
	}
	return nil
}
