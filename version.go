package main

import (
	"fmt"
)

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
