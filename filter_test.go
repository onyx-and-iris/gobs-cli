package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFilterList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &FilterListCmd{
		SourceName: "Mic/Aux",
	}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list filters: %v", err)
	}
	if !strings.Contains(out.String(), "test_filter") {
		t.Fatalf("Expected output to contain 'test_filter', got '%s'", out.String())
	}
}

func TestFilterListScene(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &FilterListCmd{
		SourceName: "gobs-test-scene",
	}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list filters in scene: %v", err)
	}
	if !strings.Contains(out.String(), "test_filter") {
		t.Fatalf("Expected output to contain 'test_filter', got '%s'", out.String())
	}
}

func TestFilterListEmpty(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &FilterListCmd{
		SourceName: "NonExistentSource",
	}
	err := cmd.Run(context)
	if err == nil {
		t.Fatal("Expected error for non-existent source, but got none")
	}
	if !strings.Contains(err.Error(), "No source was found by the name of `NonExistentSource`.") {
		t.Fatalf(
			"Expected error to contain 'No source was found by the name of `NonExistentSource`.', got '%s'",
			err.Error(),
		)
	}
}
