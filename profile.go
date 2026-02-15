package main

import (
	"fmt"
	"slices"

	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// ProfileCmd provides commands to manage profiles in OBS Studio.
type ProfileCmd struct {
	List    ProfileListCmd    `help:"List profiles."       cmd:"" aliases:"ls"`
	Current ProfileCurrentCmd `help:"Get current profile." cmd:"" aliases:"c"`
	Switch  ProfileSwitchCmd  `help:"Switch profile."      cmd:"" aliases:"sw"`
	Create  ProfileCreateCmd  `help:"Create profile."      cmd:"" aliases:"new"`
	Remove  ProfileRemoveCmd  `help:"Remove profile."      cmd:"" aliases:"rm"`
}

// ProfileListCmd provides a command to list all profiles.
type ProfileListCmd struct{} // size = 0x0

// Run executes the command to list all profiles.
// nolint: misspell
func (cmd *ProfileListCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}

	t := table.New().Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ctx.Style.border)).
		Headers("Profile Name", "Current").
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 3)
			switch col {
			case 0:
				style = style.Align(lipgloss.Left)
			case 1:
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

	for _, profile := range profiles.Profiles {
		var enabledMark string
		if profile == profiles.CurrentProfileName {
			enabledMark = getEnabledMark(true)
		}

		t.Row(profile, enabledMark)
	}
	fmt.Fprintln(ctx.Out, t.Render())
	return nil
}

// ProfileCurrentCmd provides a command to get the current profile.
type ProfileCurrentCmd struct{} // size = 0x0

// Run executes the command to get the current profile.
func (cmd *ProfileCurrentCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Current profile: %s\n", ctx.Style.Highlight(profiles.CurrentProfileName))

	return nil
}

// ProfileSwitchCmd provides a command to switch to a different profile.
type ProfileSwitchCmd struct {
	Name string `arg:"" help:"Name of the profile to switch to." required:""`
}

// Run executes the command to switch to a different profile.
func (cmd *ProfileSwitchCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}
	current := profiles.CurrentProfileName

	if current == cmd.Name {
		return fmt.Errorf("already using profile %s", ctx.Style.Error(cmd.Name))
	}

	_, err = ctx.Client.Config.SetCurrentProfile(
		config.NewSetCurrentProfileParams().WithProfileName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to switch to profile %s: %w", ctx.Style.Error(cmd.Name), err)
	}

	fmt.Fprintf(
		ctx.Out,
		"Switched from profile %s to %s\n",
		ctx.Style.Highlight(current),
		ctx.Style.Highlight(cmd.Name),
	)

	return nil
}

// ProfileCreateCmd provides a command to create a new profile.
type ProfileCreateCmd struct {
	Name string `arg:"" help:"Name of the profile to create." required:""`
}

// Run executes the command to create a new profile.
func (cmd *ProfileCreateCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}

	if slices.Contains(profiles.Profiles, cmd.Name) {
		return fmt.Errorf("profile %s already exists", ctx.Style.Error(cmd.Name))
	}

	_, err = ctx.Client.Config.CreateProfile(
		config.NewCreateProfileParams().WithProfileName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to create profile %s: %w", ctx.Style.Error(cmd.Name), err)
	}

	fmt.Fprintf(ctx.Out, "Created profile: %s\n", ctx.Style.Highlight(cmd.Name))

	return nil
}

// ProfileRemoveCmd provides a command to remove an existing profile.
type ProfileRemoveCmd struct {
	Name string `arg:"" help:"Name of the profile to delete." required:""`
}

// Run executes the command to remove an existing profile.
func (cmd *ProfileRemoveCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}

	if !slices.Contains(profiles.Profiles, cmd.Name) {
		return fmt.Errorf("profile %s does not exist", ctx.Style.Error(cmd.Name))
	}

	// Prevent deletion of the current profile
	// This is allowed in OBS Studio (with a confirmation prompt), but we want to prevent it here
	if profiles.CurrentProfileName == cmd.Name {
		return fmt.Errorf("cannot delete current profile %s", ctx.Style.Error(cmd.Name))
	}

	_, err = ctx.Client.Config.RemoveProfile(
		config.NewRemoveProfileParams().WithProfileName(cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to delete profile %s: %w", ctx.Style.Error(cmd.Name), err)
	}

	fmt.Fprintf(ctx.Out, "Deleted profile: %s\n", ctx.Style.Highlight(cmd.Name))

	return nil
}
