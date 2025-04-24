package main

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
)

// SceneItemCmd provides commands to manage scene items in OBS Studio.
type SceneItemCmd struct {
	List    SceneItemListCmd    `cmd:"" help:"List all scene items."      aliases:"ls"`
	Show    SceneItemShowCmd    `cmd:"" help:"Show scene item."           aliases:"sh"`
	Hide    SceneItemHideCmd    `cmd:"" help:"Hide scene item."           aliases:"h"`
	Toggle  SceneItemToggleCmd  `cmd:"" help:"Toggle scene item."         aliases:"tg"`
	Visible SceneItemVisibleCmd `cmd:"" help:"Get scene item visibility." aliases:"v"`
}

// SceneItemListCmd provides a command to list all scene items in a scene.
type SceneItemListCmd struct {
	SceneName string `arg:"" help:"Scene name."`
}

// Run executes the command to list all scene items in a scene.
func (cmd *SceneItemListCmd) Run(ctx *context) error {
	resp, err := ctx.Client.SceneItems.GetSceneItemList(sceneitems.NewGetSceneItemListParams().
		WithSceneName(cmd.SceneName))
	if err != nil {
		return fmt.Errorf("failed to get scene item list: %w", err)
	}
	for _, item := range resp.SceneItems {
		fmt.Fprintf(ctx.Out, "Item ID: %d, Source Name: %s\n", item.SceneItemID, item.SourceName)
	}
	return nil
}

func getSceneNameAndItemID(
	client *goobs.Client,
	sceneName string,
	itemName string,
	parent string,
) (string, int, error) {
	if parent != "" {
		resp, err := client.SceneItems.GetGroupSceneItemList(sceneitems.NewGetGroupSceneItemListParams().
			WithSceneName(parent))
		if err != nil {
			return "", 0, err
		}
		for _, item := range resp.SceneItems {
			if item.SourceName == itemName {
				return parent, int(item.SceneItemID), nil
			}
		}
		return "", 0, fmt.Errorf("item '%s' not found in scene '%s'", itemName, sceneName)
	}

	itemID, err := client.SceneItems.GetSceneItemId(sceneitems.NewGetSceneItemIdParams().
		WithSceneName(sceneName).
		WithSourceName(itemName))
	if err != nil {
		return "", 0, err
	}
	return sceneName, int(itemID.SceneItemId), nil
}

// SceneItemShowCmd provides a command to show a scene item.
type SceneItemShowCmd struct {
	Parent string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to show a scene item.
func (cmd *SceneItemShowCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Parent)
	if err != nil {
		return err
	}

	_, err = ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
		WithSceneName(sceneName).
		WithSceneItemId(sceneItemID).
		WithSceneItemEnabled(true))
	if err != nil {
		return err
	}
	return nil
}

// SceneItemHideCmd provides a command to hide a scene item.
type SceneItemHideCmd struct {
	Parent string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to hide a scene item.
func (cmd *SceneItemHideCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Parent)
	if err != nil {
		return err
	}

	_, err = ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
		WithSceneName(sceneName).
		WithSceneItemId(sceneItemID).
		WithSceneItemEnabled(false))
	if err != nil {
		return err
	}
	return nil
}

// getItemEnabled retrieves the enabled status of a scene item.
func getItemEnabled(client *goobs.Client, sceneName string, itemID int) (bool, error) {
	item, err := client.SceneItems.GetSceneItemEnabled(sceneitems.NewGetSceneItemEnabledParams().
		WithSceneName(sceneName).
		WithSceneItemId(itemID))
	if err != nil {
		return false, err
	}
	return item.SceneItemEnabled, nil
}

// SceneItemToggleCmd provides a command to toggle the visibility of a scene item.
type SceneItemToggleCmd struct {
	Parent string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to toggle the visibility of a scene item.
func (cmd *SceneItemToggleCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Parent)
	if err != nil {
		return err
	}

	itemEnabled, err := getItemEnabled(ctx.Client, sceneName, sceneItemID)
	if err != nil {
		return err
	}

	_, err = ctx.Client.SceneItems.SetSceneItemEnabled(sceneitems.NewSetSceneItemEnabledParams().
		WithSceneName(sceneName).
		WithSceneItemId(sceneItemID).
		WithSceneItemEnabled(!itemEnabled))
	if err != nil {
		return err
	}
	return nil
}

// SceneItemVisibleCmd provides a command to check the visibility of a scene item.
type SceneItemVisibleCmd struct {
	Parent string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to check the visibility of a scene item.
func (cmd *SceneItemVisibleCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Parent)
	if err != nil {
		return err
	}

	itemEnabled, err := getItemEnabled(ctx.Client, sceneName, sceneItemID)
	if err != nil {
		return err
	}

	if itemEnabled {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is visible.\n", cmd.ItemName, cmd.SceneName)
	} else {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is hidden.\n", cmd.ItemName, cmd.SceneName)
	}
	return nil
}
