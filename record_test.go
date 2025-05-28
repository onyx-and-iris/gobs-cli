package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRecordStart(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &RecordStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get recording status: %v", err)
	}
	var active bool
	if out.String() == "Recording is in progress.\n" {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStart := &RecordStartCmd{}
	err = cmdStart.Run(context)
	if err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}
	time.Sleep(1 * time.Second) // Wait for a second to ensure recording has started
	if active {
		if out.String() != "Recording is already in progress.\n" {
			t.Fatalf("Expected output to be 'Recording is already in progress.', got '%s'", out.String())
		}
	} else {
		if !strings.Contains(out.String(), "Recording started successfully.") {
			t.Fatalf("Expected output to contain 'Recording started successfully.', got '%s'", out.String())
		}
	}
}

func TestRecordStop(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &RecordStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get recording status: %v", err)
	}
	var active bool
	if out.String() == "Recording is in progress.\n" {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStop := &RecordStopCmd{}
	err = cmdStop.Run(context)
	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}
	time.Sleep(1 * time.Second) // Wait for a second to ensure recording has stopped
	if !active {
		if out.String() != "No recording is currently in progress.\n" {
			t.Fatalf("Expected output to be 'No recording is currently in progress.', got '%s'", out.String())
		}
	} else {
		if !strings.Contains(out.String(), "Recording stopped successfully. Output file:") {
			t.Fatalf("Expected output to contain 'Recording stopped successfully. Output file:', got '%s'", out.String())
		}
	}
}

func TestRecordToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmdStatus := &RecordStatusCmd{}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get recording status: %v", err)
	}
	var active bool
	if out.String() == "Recording is in progress.\n" {
		active = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdToggle := &RecordToggleCmd{}
	err = cmdToggle.Run(context)
	if err != nil {
		t.Fatalf("Failed to toggle recording: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for a second to ensure toggle has taken effect

	if active {
		if out.String() != "Recording stopped successfully.\n" {
			t.Fatalf("Expected output to be 'Recording stopped successfully.', got '%s'", out.String())
		}
	} else {
		if out.String() != "Recording started successfully.\n" {
			t.Fatalf("Expected output to be 'Recording started successfully.', got '%s'", out.String())
		}
	}
}
