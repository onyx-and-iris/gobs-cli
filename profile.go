package main

import (
	"fmt"
	"slices"

	"github.com/andreykaipov/goobs/api/requests/config"
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
func (cmd *ProfileListCmd) Run(ctx *context) error {
	profiles, err := ctx.Client.Config.GetProfileList()
	if err != nil {
		return err
	}

	for _, profile := range profiles.Profiles {
		fmt.Fprintln(ctx.Out, profile)
	}

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
	fmt.Fprintf(ctx.Out, "Current profile: %s\n", profiles.CurrentProfileName)

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
		return nil
	}

	_, err = ctx.Client.Config.SetCurrentProfile(config.NewSetCurrentProfileParams().WithProfileName(cmd.Name))
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out, "Switched from profile %s to %s\n", current, cmd.Name)

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
		return fmt.Errorf("profile %s already exists", cmd.Name)
	}

	_, err = ctx.Client.Config.CreateProfile(config.NewCreateProfileParams().WithProfileName(cmd.Name))
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out, "Created profile: %s\n", cmd.Name)

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
		return fmt.Errorf("profile %s does not exist", cmd.Name)
	}

	if profiles.CurrentProfileName == cmd.Name {
		return fmt.Errorf("cannot delete current profile %s", cmd.Name)
	}

	_, err = ctx.Client.Config.RemoveProfile(config.NewRemoveProfileParams().WithProfileName(cmd.Name))
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out, "Deleted profile: %s\n", cmd.Name)

	return nil
}
