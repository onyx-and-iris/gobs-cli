package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestStreamStart(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &StreamStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get stream status: %v", err)
	}
	var active bool
	if strings.Contains(out.String(), "Output active: true") {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStart := &StreamStartCmd{}
	err = cmdStart.Run(context)
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for the stream to start

	if active {
		if out.String() != "Stream is already active.\n" {
			t.Fatalf("Expected 'Stream is already active.', got: %s", out.String())
		}
	} else {
		if out.String() != "Streaming started successfully.\n" {
			t.Fatalf("Expected 'Streaming started successfully.', got: %s", out.String())
		}
	}
}

func TestStreamStop(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &StreamStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get stream status: %v", err)
	}
	var active bool
	if strings.Contains(out.String(), "Output active: true") {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStop := &StreamStopCmd{}
	err = cmdStop.Run(context)
	if err != nil {
		t.Fatalf("Failed to stop stream: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for the stream to stop

	if active {
		if out.String() != "Streaming stopped successfully.\n" {
			t.Fatalf("Expected 'Streaming stopped successfully.', got: %s", out.String())
		}
	} else {
		if out.String() != "Stream is already inactive.\n" {
			t.Fatalf("Expected 'Stream is already inactive.', got: %s", out.String())
		}
	}
}

func TestStreamToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &StreamStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get stream status: %v", err)
	}
	var active bool
	if strings.Contains(out.String(), "Output active: true") {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdToggle := &StreamToggleCmd{}
	err = cmdToggle.Run(context)
	if err != nil {
		t.Fatalf("Failed to toggle stream: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for the stream to toggle

	if active {
		if out.String() != "Streaming stopped successfully.\n" {
			t.Fatalf("Expected 'Streaming stopped successfully.', got: %s", out.String())
		}
	} else {
		if out.String() != "Streaming started successfully.\n" {
			t.Fatalf("Expected 'Streaming started successfully.', got: %s", out.String())
		}
	}
}
