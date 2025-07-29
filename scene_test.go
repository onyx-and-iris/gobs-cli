package main

import (
	"bytes"
	"testing"
)

func TestSceneList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &SceneListCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list scenes: %v", err)
	}
	if out.String() == "Current program scene: gobs-test-scene\n" {
		t.Fatalf("Expected output to be 'Current program scene: gobs-test-scene', got '%s'", out.String())
	}
}

func TestSceneCurrent(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	// Set the current scene to "gobs-test-scene"
	cmdSwitch := &SceneSwitchCmd{
		NewScene: "gobs-test-scene",
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
	if out.String() != "Current program scene: gobs-test-scene\n" {
		t.Fatalf("Expected output to be 'Current program scene: gobs-test-scene', got '%s'", out.String())
	}
}
