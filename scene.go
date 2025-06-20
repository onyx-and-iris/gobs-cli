package main

import (
	"fmt"
	"slices"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/aquasecurity/table"
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
func (cmd *SceneListCmd) Run(ctx *context) error {
	scenes, err := ctx.Client.Scenes.GetSceneList()
	if err != nil {
		return err
	}

	currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
	if err != nil {
		return err
	}

	t := table.New(ctx.Out)
	t.SetPadding(3)
	if cmd.UUID {
		t.SetAlignment(table.AlignLeft, table.AlignCenter, table.AlignLeft)
		t.SetHeaders("Scene Name", "Active", "UUID")
	} else {
		t.SetAlignment(table.AlignLeft, table.AlignCenter)
		t.SetHeaders("Scene Name", "Active")
	}

	slices.Reverse(scenes.Scenes)
	for _, scene := range scenes.Scenes {
		var activeMark string
		if scene.SceneName == currentScene.SceneName {
			activeMark = getEnabledMark(true)
		}
		if cmd.UUID {
			t.AddRow(scene.SceneName, activeMark, scene.SceneUuid)
		} else {
			t.AddRow(scene.SceneName, activeMark)
		}
	}
	t.Render()
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
		fmt.Fprintln(ctx.Out, scene.SceneName)
	} else {
		scene, err := ctx.Client.Scenes.GetCurrentProgramScene()
		if err != nil {
			return err
		}
		fmt.Fprintln(ctx.Out, scene.SceneName)
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

		fmt.Fprintln(ctx.Out, "Switched to preview scene:", cmd.NewScene)
	} else {
		_, err := ctx.Client.Scenes.SetCurrentProgramScene(scenes.NewSetCurrentProgramSceneParams().
			WithSceneName(cmd.NewScene))
		if err != nil {
			return err
		}

		fmt.Fprintln(ctx.Out, "Switched to program scene:", cmd.NewScene)
	}
	return nil
}
