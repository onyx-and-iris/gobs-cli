package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, "")

	cmd := &ObsVersionCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if !strings.Contains(out.String(), "OBS Client Version:") {
		t.Fatalf("Expected output to contain 'OBS Client Version:', got '%s'", out.String())
	}
	if !strings.Contains(out.String(), "with Websocket Version:") {
		t.Fatalf("Expected output to contain 'with Websocket Version:', got '%s'", out.String())
	}
}
