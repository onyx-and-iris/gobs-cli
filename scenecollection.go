package main

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// SceneCollectionCmd provides commands to manage scene collections in OBS Studio.
type SceneCollectionCmd struct {
	List    SceneCollectionListCmd    `help:"List scene collections."       cmd:"" aliases:"ls"`
	Current SceneCollectionCurrentCmd `help:"Get current scene collection." cmd:"" aliases:"c"`
	Switch  SceneCollectionSwitchCmd  `help:"Switch scene collection."      cmd:"" aliases:"sw"`
	Create  SceneCollectionCreateCmd  `help:"Create scene collection."      cmd:"" aliases:"new"`
}

// SceneCollectionListCmd provides a command to list all scene collections.
type SceneCollectionListCmd struct{} // size = 0x0

// Run executes the command to list all scene collections.
// nolint: misspell
func (cmd *SceneCollectionListCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return fmt.Errorf("failed to get scene collection list: %w", err)
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Scene Collection Name").
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch col {
			case 0:
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

	for _, collection := range collections.SceneCollections {
		t.Row(collection)
	}
	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// SceneCollectionCurrentCmd provides a command to get the current scene collection.
type SceneCollectionCurrentCmd struct{} // size = 0x0

// Run executes the command to get the current scene collection.
func (cmd *SceneCollectionCurrentCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return fmt.Errorf("failed to get scene collection list: %w", err)
	}
	fmt.Fprintln(ctx.Out, collections.CurrentSceneCollectionName)

	return nil
}

// SceneCollectionSwitchCmd provides a command to switch to a different scene collection.
type SceneCollectionSwitchCmd struct {
	Name string `arg:"" help:"Name of the scene collection to switch to." required:""`
}

// Run executes the command to switch to a different scene collection.
func (cmd *SceneCollectionSwitchCmd) Run(ctx *context) error {
	collections, err := ctx.Client.Config.GetSceneCollectionList()
	if err != nil {
		return err
	}
	current := collections.CurrentSceneCollectionName

	if current == cmd.Name {
		return fmt.Errorf("scene collection %s is already active", ctx.Style.Error(cmd.Name))
	}

	_, err = ctx.Client.Config.SetCurrentSceneCollection(
		config.NewSetCurrentSceneCollectionParams().WithSceneCollectionName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to switch scene collection %s: %w", ctx.Style.Error(cmd.Name), err)
	}

	fmt.Fprintf(ctx.Out, "Switched to scene collection: %s\n", ctx.Style.Highlight(cmd.Name))

	return nil
}

// SceneCollectionCreateCmd provides a command to create a new scene collection.
type SceneCollectionCreateCmd struct {
	Name string `arg:"" help:"Name of the scene collection to create." required:""`
}

// Run executes the command to create a new scene collection.
func (cmd *SceneCollectionCreateCmd) Run(ctx *context) error {
	_, err := ctx.Client.Config.CreateSceneCollection(
		config.NewCreateSceneCollectionParams().WithSceneCollectionName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to create scene collection %s: %w", ctx.Style.Error(cmd.Name), err)
	}

	fmt.Fprintf(ctx.Out, "Created scene collection: %s\n", ctx.Style.Highlight(cmd.Name))
	return nil
}
