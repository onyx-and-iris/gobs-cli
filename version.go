package main

import (
	"fmt"
)

// VersionCmd handles the version command.
type VersionCmd struct{} // size = 0x0

// Run executes the command to get the OBS client version.
func (cmd *VersionCmd) Run(ctx *context) error {
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
