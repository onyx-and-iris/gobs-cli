package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/ui"
	"github.com/aquasecurity/table"
)

// ProjectorCmd provides a command to manage projectors in OBS.
type ProjectorCmd struct {
	ListMonitors ProjectorListMonitorsCmd `cmd:"" help:"List available monitors."                                        aliases:"ls-m"`
	Open         ProjectorOpenCmd         `cmd:"" help:"Open a fullscreen projector for a source on a specific monitor." aliases:"o"`
}

// ProjectorListMonitorsCmd provides a command to list all monitors available for projectors.
type ProjectorListMonitorsCmd struct{} // size = 0x0

// Run executes the command to list all monitors available for projectors.
func (cmd *ProjectorListMonitorsCmd) Run(ctx *context) error {
	monitors, err := ctx.Client.Ui.GetMonitorList()
	if err != nil {
		return err
	}

	if len(monitors.Monitors) == 0 {
		ctx.Out.Write([]byte("No monitors found for projectors.\n"))
		return nil
	}

	t := table.New(ctx.Out)
	t.SetPadding(3)
	t.SetAlignment(table.AlignCenter, table.AlignLeft)
	t.SetHeaders("Monitor ID", "Monitor Name")

	for _, monitor := range monitors.Monitors {
		t.AddRow(fmt.Sprintf("%d", monitor.MonitorIndex), monitor.MonitorName)
	}

	t.Render()
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

	ctx.Client.Ui.OpenSourceProjector(ui.NewOpenSourceProjectorParams().
		WithSourceName(cmd.SourceName).
		WithMonitorIndex(cmd.MonitorIndex))

	fmt.Fprintf(ctx.Out, "Opened projector for source '%s' on monitor index %d.\n", cmd.SourceName, cmd.MonitorIndex)
	return nil
}
