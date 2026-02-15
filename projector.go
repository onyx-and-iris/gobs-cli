package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// ProjectorCmd provides a command to manage projectors in OBS.
type ProjectorCmd struct {
	ListMonitors ProjectorListMonitorsCmd `cmd:"" help:"List available monitors."                                        aliases:"ls-m"`
	Open         ProjectorOpenCmd         `cmd:"" help:"Open a fullscreen projector for a source on a specific monitor." aliases:"o"`
}

// ProjectorListMonitorsCmd provides a command to list all monitors available for projectors.
type ProjectorListMonitorsCmd struct{} // size = 0x0

// Run executes the command to list all monitors available for projectors.
// nolint: misspell
func (cmd *ProjectorListMonitorsCmd) Run(ctx *context) error {
	monitors, err := ctx.Client.Ui.GetMonitorList()
	if err != nil {
		return err
	}

	if len(monitors.Monitors) == 0 {
		fmt.Fprintf(ctx.Out, "No monitors found.\n")
		return nil
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Monitor ID", "Monitor Name").
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch col {
			case 0:
				style = style.Align(lipgloss.Center)
			case 1:
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

	for _, monitor := range monitors.Monitors {
		t.Row(fmt.Sprintf("%d", monitor.MonitorIndex), monitor.MonitorName)
	}

	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// ProjectorOpenCmd provides a command to open a fullscreen projector for a specific source.
type ProjectorOpenCmd struct {
	MonitorIndex int    `flag:"" help:"Index of the monitor to open the projector on." default:"0"`
	SourceName   string `        help:"Name of the source to project."                 default:""  arg:""`
}

// Run executes the command to show details of a specific projector.
func (cmd *ProjectorOpenCmd) Run(ctx *context) error {
	if cmd.SourceName == "" {
		currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
		if err != nil {
			return fmt.Errorf("failed to get current program scene: %w", err)
		}
		cmd.SourceName = currentScene.SceneName
	}

	monitors, err := ctx.Client.Ui.GetMonitorList()
	if err != nil {
		return err
	}

	var monitorName string
	for _, monitor := range monitors.Monitors {
		if monitor.MonitorIndex == cmd.MonitorIndex {
			monitorName = monitor.MonitorName
			break
		}
	}

	if monitorName == "" {
		return fmt.Errorf(
			"monitor with index %s not found. use %s to list available monitors",
			ctx.Style.Error(fmt.Sprintf("%d", cmd.MonitorIndex)),
			ctx.Style.Error("gobs-cli prj ls-m"),
		)
	}

	_, err = ctx.Client.Ui.OpenSourceProjector(ui.NewOpenSourceProjectorParams().
		WithSourceName(cmd.SourceName).
		WithMonitorIndex(cmd.MonitorIndex))
	if err != nil {
		return fmt.Errorf("failed to open projector: %w", err)
	}

	fmt.Fprintf(
		ctx.Out,
		"Opened projector for source %s on monitor %s.\n",
		ctx.Style.Highlight(cmd.SourceName),
		ctx.Style.Highlight(monitorName),
	)
	return nil
}
