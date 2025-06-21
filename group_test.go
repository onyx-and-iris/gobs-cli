package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGroupList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmd := &GroupListCmd{
		SceneName: "Scene",
	}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list groups: %v", err)
	}
	if !strings.Contains(out.String(), "test_group") {
		t.Fatalf("Expected output to contain 'test_group', got '%s'", out.String())
	}
}

func TestGroupShow(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmd := &GroupShowCmd{
		SceneName: "Scene",
		GroupName: "test_group",
	}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to show group: %v", err)
	}
	if out.String() != "Group test_group is now shown.\n" {
		t.Fatalf("Expected output to be 'Group test_group is now shown.', got '%s'", out.String())
	}
}

func TestGroupToggle(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmdStatus := &GroupStatusCmd{
		SceneName: "Scene",
		GroupName: "test_group",
	}
	err := cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get group status: %v", err)
	}
	var enabled bool
	if strings.Contains(out.String(), "Group test_group is shown.") {
		enabled = true
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdToggle := &GroupToggleCmd{
		SceneName: "Scene",
		GroupName: "test_group",
	}
	err = cmdToggle.Run(context)
	if err != nil {
		t.Fatalf("Failed to toggle group: %v", err)
	}
	if enabled {
		if out.String() != "Group test_group is now hidden.\n" {
			t.Fatalf("Expected output to be 'Group test_group is now hidden.', got '%s'", out.String())
		}
	} else {
		if out.String() != "Group test_group is now shown.\n" {
			t.Fatalf("Expected output to be 'Group test_group is now shown.', got '%s'", out.String())
		}
	}
}

func TestGroupStatus(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmdShow := &GroupShowCmd{
		SceneName: "Scene",
		GroupName: "test_group",
	}
	err := cmdShow.Run(context)
	if err != nil {
		t.Fatalf("Failed to show group: %v", err)
	}
	// Reset output buffer for the next command
	out.Reset()

	cmdStatus := &GroupStatusCmd{
		SceneName: "Scene",
		GroupName: "test_group",
	}
	err = cmdStatus.Run(context)
	if err != nil {
		t.Fatalf("Failed to get group status: %v", err)
	}
	if out.String() != "Group test_group is shown.\n" {
		t.Fatalf("Expected output to be 'Group test_group is shown.', got '%s'", out.String())
	}
}
