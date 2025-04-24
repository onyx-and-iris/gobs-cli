package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/sceneitems"
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
	SceneName string `arg:"" help:"Name of the scene to list groups from."`
}

// Run executes the command to list all groups in a scene.
func (cmd *GroupListCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}
	for _, item := range resp.SceneItems {
		if item.IsGroup {
			fmt.Fprintf(ctx.Out, "Group ID: %d, Source Name: %s\n", item.SceneItemID, item.SourceName)
		}
	}
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
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
				WithSceneName(cmd.SceneName).
				WithSceneItemId(item.SceneItemID).
				WithSceneItemEnabled(true))
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			fmt.Fprintf(ctx.Out, "Group %s is now shown.\n", cmd.GroupName)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("group '%s' not found", cmd.GroupName)
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
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
				WithSceneName(cmd.SceneName).
				WithSceneItemId(item.SceneItemID).
				WithSceneItemEnabled(false))
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			fmt.Fprintf(ctx.Out, "Group %s is now hidden.\n", cmd.GroupName)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("group '%s' not found", cmd.GroupName)
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
			_, err := ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
				WithSceneName(cmd.SceneName).
				WithSceneItemId(item.SceneItemID).
				WithSceneItemEnabled(newState))
			if err != nil {
				return fmt.Errorf("failed to set scene item enabled: %w", err)
			}
			if newState {
				fmt.Fprintf(ctx.Out, "Group %s is now shown.\n", cmd.GroupName)
			} else {
				fmt.Fprintf(ctx.Out, "Group %s is now hidden.\n", cmd.GroupName)
			}
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("group '%s' not found", cmd.GroupName)
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
				fmt.Fprintf(ctx.Out, "Group %s is shown.\n", cmd.GroupName)
			} else {
				fmt.Fprintf(ctx.Out, "Group %s is hidden.\n", cmd.GroupName)
			}
			return nil
		}
	}
	return fmt.Errorf("group '%s' not found", cmd.GroupName)
}
