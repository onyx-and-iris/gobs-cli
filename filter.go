package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/aquasecurity/table"
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
	SourceName string `arg:"" help:"Name of the source to list filters from."`
}

// Run executes the command to list all filters in a scene.
func (cmd *FilterListCmd) Run(ctx *context) error {
	filters, err := ctx.Client.Filters.GetSourceFilterList(
		filters.NewGetSourceFilterListParams().WithSourceName(cmd.SourceName),
	)
	if err != nil {
		return err
	}

	if len(filters.Filters) == 0 {
		fmt.Fprintf(ctx.Out, "No filters found for source %s.\n", cmd.SourceName)
		return nil
	}

	t := table.New(ctx.Out)
	t.SetPadding(3)
	t.SetAlignment(table.AlignLeft, table.AlignLeft, table.AlignCenter, table.AlignLeft)
	t.SetHeaders("Filter Name", "Kind", "Enabled", "Settings")

	for _, filter := range filters.Filters {
		var lines []string
		for k, v := range filter.FilterSettings {
			lines = append(lines, fmt.Sprintf("%s %v", k, v))
		}
		sort.Slice(lines, func(i, j int) bool {
			return strings.ToLower(lines[i]) < strings.ToLower(lines[j])
		})

		t.AddRow(
			filter.FilterName,
			snakeCaseToTitleCase(filter.FilterKind),
			getEnabledMark(filter.FilterEnabled),
			strings.Join(lines, "\n"),
		)
	}
	t.Render()
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
			cmd.FilterName, cmd.SourceName, err)
	}
	fmt.Fprintf(ctx.Out, "Filter %s enabled on source %s.\n",
		cmd.FilterName, cmd.SourceName)
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
			cmd.FilterName, cmd.SourceName, err)
	}
	fmt.Fprintf(ctx.Out, "Filter %s disabled on source %s.\n",
		cmd.FilterName, cmd.SourceName)
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
			cmd.FilterName, cmd.SourceName, err)
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
			cmd.FilterName, cmd.SourceName, err)
	}

	if newStatus {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is now enabled.\n",
			cmd.FilterName, cmd.SourceName)
	} else {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is now disabled.\n",
			cmd.FilterName, cmd.SourceName)
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
			cmd.FilterName, cmd.SourceName, err)
	}
	if filter.FilterEnabled {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is enabled.\n",
			cmd.FilterName, cmd.SourceName)
	} else {
		fmt.Fprintf(ctx.Out, "Filter %s on source %s is disabled.\n",
			cmd.FilterName, cmd.SourceName)
	}
	return nil
}
