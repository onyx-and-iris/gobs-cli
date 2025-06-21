package main

import (
	"fmt"
	"slices"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// SceneCmd provides commands to manage scenes in OBS Studio.
type SceneCmd struct {
	List    SceneListCmd    `cmd:"" help:"List all scenes."       aliases:"ls"`
	Current SceneCurrentCmd `cmd:"" help:"Get the current scene." aliases:"c"`
	Switch  SceneSwitchCmd  `cmd:"" help:"Switch to a scene."     aliases:"sw"`
}

// SceneListCmd provides a command to list all scenes.
type SceneListCmd struct {
	UUID bool `flag:"" help:"Display UUIDs of scenes."`
}

// Run executes the command to list all scenes.
// nolint: misspell
func (cmd *SceneListCmd) Run(ctx *context) error {
	scenes, err := ctx.Client.Scenes.GetSceneList()
	if err != nil {
		return err
	}

	currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
	if err != nil {
		return err
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border))
	if cmd.UUID {
		t.Headers("Scene Name", "Active", "UUID")
	} else {
		t.Headers("Scene Name", "Active")
	}
	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 3)
		switch col {
		case 0:
			style = style.Align(lipgloss.Left)
		case 1:
			style = style.Align(lipgloss.Center)
		case 2:
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

	slices.Reverse(scenes.Scenes)
	for _, scene := range scenes.Scenes {
		var activeMark string
		if scene.SceneName == currentScene.SceneName {
			activeMark = getEnabledMark(true)
		}
		if cmd.UUID {
			t.Row(scene.SceneName, activeMark, scene.SceneUuid)
		} else {
			t.Row(scene.SceneName, activeMark)
		}
	}
	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// SceneCurrentCmd provides a command to get the current scene.
type SceneCurrentCmd struct {
	Preview bool `flag:"" help:"Preview scene."`
}

// Run executes the command to get the current scene.
func (cmd *SceneCurrentCmd) Run(ctx *context) error {
	if cmd.Preview {
		scene, err := ctx.Client.Scenes.GetCurrentPreviewScene()
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Current preview scene: %s\n", ctx.Style.Highlight(scene.SceneName))
	} else {
		scene, err := ctx.Client.Scenes.GetCurrentProgramScene()
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Current program scene: %s\n", ctx.Style.Highlight(scene.SceneName))
	}
	return nil
}

// SceneSwitchCmd provides a command to switch to a different scene.
type SceneSwitchCmd struct {
	Preview  bool   `flag:"" help:"Preview scene."`
	NewScene string `        help:"Scene name to switch to." arg:""`
}

// Run executes the command to switch to a different scene.
func (cmd *SceneSwitchCmd) Run(ctx *context) error {
	if cmd.Preview {
		_, err := ctx.Client.Scenes.SetCurrentPreviewScene(scenes.NewSetCurrentPreviewSceneParams().
			WithSceneName(cmd.NewScene))
		if err != nil {
			return err
		}

		fmt.Fprintf(ctx.Out, "Switched to preview scene: %s\n", ctx.Style.Highlight(cmd.NewScene))
	} else {
		_, err := ctx.Client.Scenes.SetCurrentProgramScene(scenes.NewSetCurrentProgramSceneParams().
			WithSceneName(cmd.NewScene))
		if err != nil {
			return err
		}

		fmt.Fprintf(ctx.Out, "Switched to program scene: %s\n", ctx.Style.Highlight(cmd.NewScene))
	}
	return nil
}
