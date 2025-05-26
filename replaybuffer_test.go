package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplayBufferStart(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmd := &ReplayBufferStartCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to start replay buffer: %v", err)
	}
	if out.String() != "Replay buffer started.\n" {
		t.Fatalf("Expected output to be 'Replay buffer started', got '%s'", out.String())
	}
}

func TestReplayBufferStop(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmd := &ReplayBufferStopCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to stop replay buffer: %v", err)
	}
	if out.String() != "Replay buffer stopped.\n" {
		t.Fatalf("Expected output to be 'Replay buffer stopped.', got '%s'", out.String())
	}
}

func TestReplayBufferToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &ReplayBufferStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get replay buffer status: %v", err)
	}
	var active bool
	if strings.Contains(out.String(), "Replay buffer is active") {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdToggle := &ReplayBufferToggleCmd{}
	err = cmdToggle.Run(context)
	if err != nil {
		t.Fatalf("Failed to toggle replay buffer: %v", err)
	}
	if active {
		if out.String() != "Replay buffer stopped.\n" {
			t.Fatalf("Expected output to be 'Replay buffer stopped.', got '%s'", out.String())
		}
	} else {
		if out.String() != "Replay buffer started.\n" {
			t.Fatalf("Expected output to be 'Replay buffer started.', got '%s'", out.String())
		}
	}
}
