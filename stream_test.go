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
	context := newContext(client, &out, StyleConfig{})

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
	if active {
		if err == nil {
			t.Fatalf("Expected error when starting stream while active, got nil")
		}
		if !strings.Contains(err.Error(), "stream is already in progress") {
			t.Fatalf("Expected error message to contain 'stream is already in progress', got '%s'", err.Error())
		}
		return
	}
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}
	if out.String() != "Stream started successfully.\n" {
		t.Fatalf("Expected output to contain 'Stream started successfully.', got '%s'", out.String())
	}
	time.Sleep(2 * time.Second) // Wait for the stream to start
}

func TestStreamStop(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

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
	if !active {
		if err == nil {
			t.Fatalf("Expected error when stopping stream while inactive, got nil")
		}
		if !strings.Contains(err.Error(), "stream is not in progress") {
			t.Fatalf("Expected error message to contain 'stream is not in progress', got '%s'", err.Error())
		}
		return
	}
	if err != nil {
		t.Fatalf("Failed to stop stream: %v", err)
	}
	if out.String() != "Stream stopped successfully.\n" {
		t.Fatalf("Expected output to contain 'Stream stopped successfully.', got '%s'", out.String())
	}
	time.Sleep(2 * time.Second) // Wait for the stream to stop
}

func TestStreamToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

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

	if active {
		if out.String() != "Stream stopped successfully.\n" {
			t.Fatalf("Expected 'Stream stopped successfully.', got: %s", out.String())
		}
	} else {
		if out.String() != "Stream started successfully.\n" {
			t.Fatalf("Expected 'Stream started successfully.', got: %s", out.String())
		}
	}
	time.Sleep(2 * time.Second) // Wait for the stream to toggle
}
