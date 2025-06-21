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
	context := newContext(client, &out, "")

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
	if active {
		if err == nil {
			t.Fatalf("Expected error when starting recording while active, got nil")
		}
		if !strings.Contains(err.Error(), "recording is already in progress") {
			t.Fatalf("Expected error message to contain 'recording is already in progress', got '%s'", err.Error())
		}
		return
	}

	if err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}
	if out.String() != "Recording started successfully.\n" {
		t.Fatalf("Expected output to contain 'Recording started successfully.', got '%s'", out.String())
	}
	time.Sleep(1 * time.Second) // Wait for the recording to start
}

func TestRecordStop(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

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
	if !active {
		if err == nil {
			t.Fatalf("Expected error when stopping recording while inactive, got nil")
		}
		if !strings.Contains(err.Error(), "recording is not in progress") {
			t.Fatalf("Expected error message to contain 'recording is not in progress', got '%s'", err.Error())
		}
		return
	}

	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}
	if !strings.Contains(out.String(), "Recording stopped successfully. Output file: ") {
		t.Fatalf("Expected output to contain 'Recording stopped successfully. Output file: ', got '%s'", out.String())
	}
	time.Sleep(1 * time.Second) // Wait for the recording to stop
}

func TestRecordToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

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
