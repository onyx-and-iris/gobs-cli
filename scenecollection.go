package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/config"
)

// SceneCollectionCmd provides commands to manage scene collections in OBS Studio.
type SceneCollectionCmd struct {
	List    ListSceneCollectionCmd    `help:"List scene collections."       cmd:"" aliases:"ls"`
	Current CurrentSceneCollectionCmd `help:"Get current scene collection." cmd:"" aliases:"c"`
	Switch  SwitchSceneCollectionCmd  `help:"Switch scene collection."      cmd:"" aliases:"sw"`
	Create  CreateSceneCollectionCmd  `help:"Create scene collection."      cmd:"" aliases:"cr"`
}

// ListSceneCollectionCmd provides a command to list all scene collections.
type ListSceneCollectionCmd struct{} // size = 0x0

// Run executes the command to list all scene collections.
func (cmd *ListSceneCollectionCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return fmt.Errorf("failed to get scene collection list: %w", err)
	}

	for _, collection := range collections.SceneCollections {
		fmt.Fprintln(ctx.Out, collection)
	}

	return nil
}

// CurrentSceneCollectionCmd provides a command to get the current scene collection.
type CurrentSceneCollectionCmd struct{} // size = 0x0

// Run executes the command to get the current scene collection.
func (cmd *CurrentSceneCollectionCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return fmt.Errorf("failed to get scene collection list: %w", err)
	}
	fmt.Fprintln(ctx.Out, collections.CurrentSceneCollectionName)

	return nil
}

// SwitchSceneCollectionCmd provides a command to switch to a different scene collection.
type SwitchSceneCollectionCmd struct {
	Name string `arg:"" help:"Name of the scene collection to switch to." required:""`
}

// Run executes the command to switch to a different scene collection.
func (cmd *SwitchSceneCollectionCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return err
	}
	current := collections.CurrentSceneCollectionName

	if current == cmd.Name {
		return fmt.Errorf("scene collection %s is already active", cmd.Name)
	}

	_, err = ctx.Client.Config.SetCurrentSceneCollection(
		config.NewSetCurrentSceneCollectionParams().WithSceneCollectionName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to switch scene collection: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Switched to scene collection: %s\n", cmd.Name)

	return nil
}

// CreateSceneCollectionCmd provides a command to create a new scene collection.
type CreateSceneCollectionCmd struct {
	Name string `arg:"" help:"Name of the scene collection to create." required:""`
}

// Run executes the command to create a new scene collection.
func (cmd *CreateSceneCollectionCmd) Run(ctx *context) error {
	_, err := ctx.Client.Config.CreateSceneCollection(
		config.NewCreateSceneCollectionParams().WithSceneCollectionName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to create scene collection: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Created scene collection: %s\n", cmd.Name)
	return nil
}
