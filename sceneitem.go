package main

import (
	"fmt"
	"sort"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/aquasecurity/table"
)

// SceneItemCmd provides commands to manage scene items in OBS Studio.
type SceneItemCmd struct {
	List      SceneItemListCmd      `cmd:"" help:"List all scene items."      aliases:"ls"`
	Show      SceneItemShowCmd      `cmd:"" help:"Show scene item."           aliases:"sh"`
	Hide      SceneItemHideCmd      `cmd:"" help:"Hide scene item."           aliases:"h"`
	Toggle    SceneItemToggleCmd    `cmd:"" help:"Toggle scene item."         aliases:"tg"`
	Visible   SceneItemVisibleCmd   `cmd:"" help:"Get scene item visibility." aliases:"v"`
	Transform SceneItemTransformCmd `cmd:"" help:"Transform scene item."      aliases:"t"`
}

// SceneItemListCmd provides a command to list all scene items in a scene.
type SceneItemListCmd struct {
	UUID      bool   `flag:"" help:"Display UUIDs of scene items."`
	SceneName string `        help:"Name of the scene to list items from." arg:"" default:""`
}

// Run executes the command to list all scene items in a scene.
func (cmd *SceneItemListCmd) Run(ctx *context) error {
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

	if len(resp.SceneItems) == 0 {
		fmt.Fprintf(ctx.Out, "No scene items found in scene '%s'.\n", cmd.SceneName)
		return nil
	}

	t := table.New(ctx.Out)
	t.SetPadding(3)
	if cmd.UUID {
		t.SetAlignment(table.AlignCenter, table.AlignLeft, table.AlignCenter, table.AlignCenter, table.AlignCenter)
		t.SetHeaders("Item ID", "Item Name", "In Group", "Enabled", "UUID")
	} else {
		t.SetAlignment(table.AlignCenter, table.AlignLeft, table.AlignCenter, table.AlignCenter)
		t.SetHeaders("Item ID", "Item Name", "In Group", "Enabled")
	}

	sort.Slice(resp.SceneItems, func(i, j int) bool {
		return resp.SceneItems[i].SceneItemID < resp.SceneItems[j].SceneItemID
	})

	for _, item := range resp.SceneItems {
		if item.IsGroup {
			resp, err := ctx.Client.SceneItems.GetGroupSceneItemList(sceneitems.NewGetGroupSceneItemListParams().
				WithSceneName(item.SourceName))
			if err != nil {
				return fmt.Errorf("failed to get group scene item list for '%s': %w", item.SourceName, err)
			}

			sort.Slice(resp.SceneItems, func(i, j int) bool {
				return resp.SceneItems[i].SceneItemID < resp.SceneItems[j].SceneItemID
			})

			for _, groupItem := range resp.SceneItems {
				if cmd.UUID {
					t.AddRow(
						fmt.Sprintf("%d", groupItem.SceneItemID),
						groupItem.SourceName,
						item.SourceName,
						getEnabledMark(item.SceneItemEnabled && groupItem.SceneItemEnabled),
						groupItem.SourceUuid,
					)
				} else {
					t.AddRow(
						fmt.Sprintf("%d", groupItem.SceneItemID),
						groupItem.SourceName,
						item.SourceName,
						getEnabledMark(item.SceneItemEnabled && groupItem.SceneItemEnabled),
					)
				}
			}
		} else {
			if cmd.UUID {
				t.AddRow(fmt.Sprintf("%d", item.SceneItemID), item.SourceName, "",
					getEnabledMark(item.SceneItemEnabled), item.SourceUuid)
			} else {
				t.AddRow(fmt.Sprintf("%d", item.SceneItemID), item.SourceName, "", getEnabledMark(item.SceneItemEnabled))
			}
		}
	}
	t.Render()
	return nil
}

// getSceneNameAndItemID retrieves the scene name and item ID for a given item in a scene or group.
func getSceneNameAndItemID(
	client *goobs.Client,
	sceneName string,
	itemName string,
	group string,
) (string, int, error) {
	if group != "" {
		resp, err := client.SceneItems.GetGroupSceneItemList(sceneitems.NewGetGroupSceneItemListParams().
			WithSceneName(group))
		if err != nil {
			return "", 0, err
		}
		for _, item := range resp.SceneItems {
			if item.SourceName == itemName {
				return group, int(item.SceneItemID), nil
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
	Group string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to show a scene item.
func (cmd *SceneItemShowCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Group)
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

	if cmd.Group != "" {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in group '%s' is now visible.\n", cmd.ItemName, cmd.Group)
	} else {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is now visible.\n", cmd.ItemName, cmd.SceneName)
	}

	return nil
}

// SceneItemHideCmd provides a command to hide a scene item.
type SceneItemHideCmd struct {
	Group string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to hide a scene item.
func (cmd *SceneItemHideCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Group)
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

	if cmd.Group != "" {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in group '%s' is now hidden.\n", cmd.ItemName, cmd.Group)
	} else {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is now hidden.\n", cmd.ItemName, cmd.SceneName)
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
	Group string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to toggle the visibility of a scene item.
func (cmd *SceneItemToggleCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Group)
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

	if itemEnabled {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is now hidden.\n", cmd.ItemName, cmd.SceneName)
	} else {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' is now visible.\n", cmd.ItemName, cmd.SceneName)
	}

	return nil
}

// SceneItemVisibleCmd provides a command to check the visibility of a scene item.
type SceneItemVisibleCmd struct {
	Group string `flag:"" help:"Parent group name."`

	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`
}

// Run executes the command to check the visibility of a scene item.
func (cmd *SceneItemVisibleCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Group)
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

// SceneItemTransformCmd provides a command to transform a scene item.
type SceneItemTransformCmd struct {
	SceneName string `arg:"" help:"Scene name."`
	ItemName  string `arg:"" help:"Item name."`

	Group string `flag:"" help:"Parent group name."`

	Alignment       float64 `flag:"" help:"Alignment of the scene item."`
	BoundsAlignment float64 `flag:"" help:"Bounds alignment of the scene item."`
	BoundsHeight    float64 `flag:"" help:"Bounds height of the scene item."          default:"1.0"`
	BoundsType      string  `flag:"" help:"Bounds type of the scene item."            default:"OBS_BOUNDS_NONE"`
	BoundsWidth     float64 `flag:"" help:"Bounds width of the scene item."           default:"1.0"`
	CropToBounds    bool    `flag:"" help:"Whether to crop the scene item to bounds."`
	CropBottom      float64 `flag:"" help:"Crop bottom value of the scene item."`
	CropLeft        float64 `flag:"" help:"Crop left value of the scene item."`
	CropRight       float64 `flag:"" help:"Crop right value of the scene item."`
	CropTop         float64 `flag:"" help:"Crop top value of the scene item."`
	PositionX       float64 `flag:"" help:"X position of the scene item."`
	PositionY       float64 `flag:"" help:"Y position of the scene item."`
	Rotation        float64 `flag:"" help:"Rotation of the scene item."`
	ScaleX          float64 `flag:"" help:"X scale of the scene item."`
	ScaleY          float64 `flag:"" help:"Y scale of the scene item."`
}

// Run executes the command to transform a scene item.
func (cmd *SceneItemTransformCmd) Run(ctx *context) error {
	sceneName, sceneItemID, err := getSceneNameAndItemID(ctx.Client, cmd.SceneName, cmd.ItemName, cmd.Group)
	if err != nil {
		return err
	}

	// Get the current transform of the scene item
	resp, err := ctx.Client.SceneItems.GetSceneItemTransform(sceneitems.NewGetSceneItemTransformParams().
		WithSceneName(sceneName).
		WithSceneItemId(sceneItemID))
	if err != nil {
		return err
	}

	// Update the transform with the provided values
	transform := resp.SceneItemTransform

	if cmd.Alignment != 0 {
		transform.Alignment = cmd.Alignment
	}
	if cmd.BoundsAlignment != 0 {
		transform.BoundsAlignment = cmd.BoundsAlignment
	}

	if cmd.BoundsHeight != 0 {
		transform.BoundsHeight = cmd.BoundsHeight
	}
	if cmd.BoundsType != "" {
		transform.BoundsType = cmd.BoundsType
	}
	if cmd.BoundsWidth != 0 {
		transform.BoundsWidth = cmd.BoundsWidth
	}

	if cmd.CropToBounds {
		transform.CropToBounds = cmd.CropToBounds
	}
	if cmd.CropBottom != 0 {
		transform.CropBottom = cmd.CropBottom
	}
	if cmd.CropLeft != 0 {
		transform.CropLeft = cmd.CropLeft
	}
	if cmd.CropRight != 0 {
		transform.CropRight = cmd.CropRight
	}
	if cmd.CropTop != 0 {
		transform.CropTop = cmd.CropTop
	}
	if cmd.PositionX != 0 {
		transform.PositionX = cmd.PositionX
	}
	if cmd.PositionY != 0 {
		transform.PositionY = cmd.PositionY
	}
	if cmd.Rotation != 0 {
		transform.Rotation = cmd.Rotation
	}
	if cmd.ScaleX != 0 {
		transform.ScaleX = cmd.ScaleX
	}
	if cmd.ScaleY != 0 {
		transform.ScaleY = cmd.ScaleY
	}

	_, err = ctx.Client.SceneItems.SetSceneItemTransform(sceneitems.NewSetSceneItemTransformParams().
		WithSceneName(sceneName).
		WithSceneItemId(sceneItemID).
		WithSceneItemTransform(transform))
	if err != nil {
		return err
	}

	if cmd.Group != "" {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in group '%s' transformed.\n", cmd.ItemName, cmd.Group)
	} else {
		fmt.Fprintf(ctx.Out, "Scene item '%s' in scene '%s' transformed.\n", cmd.ItemName, cmd.SceneName)
	}

	return nil
}
