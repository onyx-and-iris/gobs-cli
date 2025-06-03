package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var version = "unknown"

// VersionFlag is a custom flag type for displaying version information.
type VersionFlag string

// Decode implements the kong.Flag interface for VersionFlag.
func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }

// IsBool implements the kong.Flag interface for VersionFlag.
func (v VersionFlag) IsBool() bool { return true }

// BeforeApply implements the kong.Flag interface for VersionFlag.
func (v VersionFlag) BeforeApply(app *kong.Kong, _ kong.Vars) error { // nolint: unparam
	fmt.Printf("gobs-cli version: %s\n", version)
	app.Exit(0) // Exit the application after printing the version
	return nil
}

// ObsVersionCmd handles the version command.
type ObsVersionCmd struct{} // size = 0x0

// Run executes the command to get the OBS client version.
func (cmd *ObsVersionCmd) Run(ctx *context) error {
	version, err := ctx.Client.General.GetVersion()
	if err != nil {
		return err
	}
	fmt.Fprintf(
		ctx.Out,
		"OBS Client Version: %s with Websocket Version: %s\n",
		version.ObsVersion,
		version.ObsWebSocketVersion,
	)

	return nil
}
