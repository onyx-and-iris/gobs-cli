package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// GroupCmd provides commands to manage groups in OBS Studio.
type GroupCmd struct {
	List   GroupListCmd   `cmd:"" help:"List all groups."    aliases:"ls"`
	Show   GroupShowCmd   `cmd:"" help:"Show group details." aliases:"sh"`
	Hide   GroupHideCmd   `cmd:"" help:"Hide group."         aliases:"h"`
	Toggle GroupToggleCmd `cmd:"" help:"Toggle group."       aliases:"tg"`
	Status GroupStatusCmd `cmd:"" help:"Get group status."   aliases:"ss"`
}

// GroupListCmd provides a command to list all groups in a scene.
type GroupListCmd struct {
	SceneName string `arg:"" help:"Name of the scene to list groups from." default:""`
}

// Run executes the command to list all groups in a scene.
// nolint: misspell
func (cmd *GroupListCmd) Run(ctx *context) error {
	if cmd.SceneName == "" {
		currentScene, err := ctx.Client.Scenes.GetCurrentProgramScene()
		if err != nil {
			return fmt.Errorf("failed to get current program scene: %w", err)
		}
		cmd.SceneName = currentScene.SceneName
	}

	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("ID", "Group Name", "Enabled").
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch col {
			case 0:
				style = style.Align(lipgloss.Center)
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

	var found bool
	for _, item := range resp.SceneItems {
		if item.IsGroup {
			t.Row(
				fmt.Sprintf("%d", item.SceneItemID),
				item.SourceName,
				getEnabledMark(item.SceneItemEnabled),
			)
			found = true
		}
	}

	if !found {
		fmt.Fprintf(ctx.Out, "No groups found in scene %s.\n", ctx.Style.Highlight(cmd.SceneName))
		return nil
	}

	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// GroupShowCmd provides a command to show a group in a scene.
type GroupShowCmd struct {
	SceneName string `arg:"" help:"Name of the scene to show group from."`
	GroupName string `arg:"" help:"Name of the group to show."`
}

// Run executes the command to show a group in a scene.
func (cmd *GroupShowCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}

	var found bool
	for _, item := range resp.SceneItems {
		if item.IsGroup && item.SourceName == cmd.GroupName {
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(
				sceneitems.NewSetSceneItemEnabledParams().
					WithSceneName(cmd.SceneName).
					WithSceneItemId(item.SceneItemID).
					WithSceneItemEnabled(true),
			)
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			fmt.Fprintf(ctx.Out, "Group %s is now shown.\n", ctx.Style.Highlight(cmd.GroupName))
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf(
			"group %s not found in scene %s",
			ctx.Style.Error(cmd.GroupName),
			ctx.Style.Error(cmd.SceneName),
		)
	}
	return nil
}

// GroupHideCmd provides a command to hide a group in a scene.
type GroupHideCmd struct {
	SceneName string `arg:"" help:"Name of the scene to hide group from."`
	GroupName string `arg:"" help:"Name of the group to hide."`
}

// Run executes the command to hide a group in a scene.
func (cmd *GroupHideCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}

	var found bool
	for _, item := range resp.SceneItems {
		if item.IsGroup && item.SourceName == cmd.GroupName {
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(
				sceneitems.NewSetSceneItemEnabledParams().
					WithSceneName(cmd.SceneName).
					WithSceneItemId(item.SceneItemID).
					WithSceneItemEnabled(false),
			)
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			fmt.Fprintf(ctx.Out, "Group %s is now hidden.\n", ctx.Style.Highlight(cmd.GroupName))
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf(
			"group %s not found in scene %s",
			ctx.Style.Error(cmd.GroupName),
			ctx.Style.Error(cmd.SceneName),
		)
	}
	return nil
}

// GroupToggleCmd provides a command to toggle a group in a scene.
type GroupToggleCmd struct {
	SceneName string `arg:"" help:"Name of the scene to toggle group from."`
	GroupName string `arg:"" help:"Name of the group to toggle."`
}

// Run executes the command to toggle a group in a scene.
func (cmd *GroupToggleCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}

	var found bool
	for _, item := range resp.SceneItems {
		if item.IsGroup && item.SourceName == cmd.GroupName {
			newState := !item.SceneItemEnabled
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(
				sceneitems.NewSetSceneItemEnabledParams().
					WithSceneName(cmd.SceneName).
					WithSceneItemId(item.SceneItemID).
					WithSceneItemEnabled(newState),
			)
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			if newState {
				fmt.Fprintf(ctx.Out, "Group %s is now shown.\n", ctx.Style.Highlight(cmd.GroupName))
			} else {
				fmt.Fprintf(
					ctx.Out,
					"Group %s is now hidden.\n",
					ctx.Style.Highlight(cmd.GroupName),
				)
			}
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf(
			"group %s not found in scene %s",
			ctx.Style.Error(cmd.GroupName),
			ctx.Style.Error(cmd.SceneName),
		)
	}

	return nil
}

// GroupStatusCmd provides a command to get the status of a group in a scene.
type GroupStatusCmd struct {
	SceneName string `arg:"" help:"Name of the scene to get group status from."`
	GroupName string `arg:"" help:"Name of the group to get status."`
}

// Run executes the command to get the status of a group in a scene.
func (cmd *GroupStatusCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}
	for _, item := range resp.SceneItems {
		if item.IsGroup && item.SourceName == cmd.GroupName {
			if item.SceneItemEnabled {
				fmt.Fprintf(ctx.Out, "Group %s is shown.\n", ctx.Style.Highlight(cmd.GroupName))
			} else {
				fmt.Fprintf(ctx.Out, "Group %s is hidden.\n", ctx.Style.Highlight(cmd.GroupName))
			}
			return nil
		}
	}
	return fmt.Errorf(
		"group %s not found in scene %s",
		ctx.Style.Error(cmd.GroupName),
		ctx.Style.Error(cmd.SceneName),
	)
}
