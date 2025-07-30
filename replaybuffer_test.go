package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func skipIfSkipReplayBufferTests(t *testing.T) {
	if os.Getenv("GOBS_TEST_SKIP_REPLAYBUFFER_TESTS") != "" {
		t.Skip("Skipping replay buffer tests due to GOBS_TEST_SKIP_REPLAYBUFFER_TESTS environment variable")
	}
}

func TestReplayBufferStart(t *testing.T) {
	skipIfSkipReplayBufferTests(t)

	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &ReplayBufferStartCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to start replay buffer: %v", err)
	}
	if out.String() != "Replay buffer started.\n" {
		t.Fatalf("Expected output to be 'Replay buffer started', got '%s'", out.String())
	}
	time.Sleep(500 * time.Millisecond) // Wait for the replay buffer to start
}

func TestReplayBufferStop(t *testing.T) {
	skipIfSkipReplayBufferTests(t)

	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &ReplayBufferStopCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to stop replay buffer: %v", err)
	}
	if out.String() != "Replay buffer stopped.\n" {
		t.Fatalf("Expected output to be 'Replay buffer stopped.', got '%s'", out.String())
	}
	time.Sleep(500 * time.Millisecond) // Wait for the replay buffer to stop
}

func TestReplayBufferToggle(t *testing.T) {
	skipIfSkipReplayBufferTests(t)

	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

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
	time.Sleep(500 * time.Millisecond) // Wait for the toggle to take effect
}
