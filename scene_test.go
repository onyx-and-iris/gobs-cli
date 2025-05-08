package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSceneList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmd := &SceneListCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list scenes: %v", err)
	}
	if !strings.Contains(out.String(), "gobs-test") {
		t.Fatalf("Expected output to contain 'gobs-test', got '%s'", out.String())
	}
}

func TestSceneCurrent(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	// Set the current scene to "gobs-test"
	cmdSwitch := &SceneSwitchCmd{
		NewScene: "gobs-test",
	}
	err := cmdSwitch.Run(context)
	if err != nil {
		t.Fatalf("Failed to switch to scene: %v", err)
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdCurrent := &SceneCurrentCmd{}
	err = cmdCurrent.Run(context)
	if err != nil {
		t.Fatalf("Failed to get current scene: %v", err)
	}
	if out.String() != "gobs-test\n" {
		t.Fatalf("Expected output to contain 'gobs-test', got '%s'", out.String())
	}
}
