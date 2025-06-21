package main

import (
	"bytes"
	"testing"
)

func TestStudioModeEnable(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmdEnable := &StudioModeEnableCmd{}
	err := cmdEnable.Run(context)
	if err != nil {
		t.Fatalf("failed to enable studio mode: %v", err)
	}
	if out.String() != "Studio mode is now enabled\n" {
		t.Fatalf("expected 'Studio mode is now enabled', got: %s", out.String())
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStatus := &StudioModeStatusCmd{}
	err = cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("failed to get studio mode status: %v", err)
	}
	if out.String() != "Studio mode is enabled\n" {
		t.Fatalf("expected 'Studio mode is enabled', got: %s", out.String())
	}
}

func TestStudioModeDisable(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmdDisable := &StudioModeDisableCmd{}
	err := cmdDisable.Run(context)
	if err != nil {
		t.Fatalf("failed to disable studio mode: %v", err)
	}
	if out.String() != "Studio mode is now disabled\n" {
		t.Fatalf("expected 'Studio mode is now disabled', got: %s", out.String())
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStatus := &StudioModeStatusCmd{}
	err = cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("failed to get studio mode status: %v", err)
	}
	if out.String() != "Studio mode is disabled\n" {
		t.Fatalf("expected 'Studio mode is disabled', got: %s", out.String())
	}
}
