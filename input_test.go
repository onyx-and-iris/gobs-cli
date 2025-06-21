package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInputList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &InputListCmd{}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list inputs: %v", err)
	}

	expectedInputs := []string{
		"Desktop Audio",
		"Mic/Aux",
		"Colour Source",
		"Colour Source 2",
		"Colour Source 3",
	}
	output := out.String()
	for _, input := range expectedInputs {
		if !strings.Contains(output, input) {
			t.Fatalf("Expected output to contain '%s', got '%s'", input, output)
		}
	}
}

func TestInputListFilterInput(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &InputListCmd{Input: true}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list inputs with filter: %v", err)
	}

	expectedInputs := []string{
		"Mic/Aux",
	}
	expectedFilteredOut := []string{
		"Desktop Audio",
		"Colour Source",
		"Colour Source 2",
		"Colour Source 3",
	}
	for _, input := range expectedInputs {
		if !strings.Contains(out.String(), input) {
			t.Fatalf("Expected output to contain '%s', got '%s'", input, out.String())
		}
	}
	for _, filteredOut := range expectedFilteredOut {
		if strings.Contains(out.String(), filteredOut) {
			t.Fatalf("Expected output to NOT contain '%s', got '%s'", filteredOut, out.String())
		}
	}
}

func TestInputListFilterOutput(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &InputListCmd{Output: true}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list outputs with filter: %v", err)
	}

	expectedInputs := []string{
		"Desktop Audio",
	}
	expectedFilteredOut := []string{
		"Mic/Aux",
		"Colour Source",
		"Colour Source 2",
		"Colour Source 3",
	}
	for _, input := range expectedInputs {
		if !strings.Contains(out.String(), input) {
			t.Fatalf("Expected output to contain '%s', got '%s'", input, out.String())
		}
	}
	for _, filteredOut := range expectedFilteredOut {
		if strings.Contains(out.String(), filteredOut) {
			t.Fatalf("Expected output to NOT contain '%s', got '%s'", filteredOut, out.String())
		}
	}
}

func TestInputListFilterColour(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := newContext(client, &out, StyleConfig{})

	cmd := &InputListCmd{Colour: true}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list colour inputs with filter: %v", err)
	}

	expectedInputs := []string{
		"Colour Source",
		"Colour Source 2",
		"Colour Source 3",
	}
	for _, input := range expectedInputs {
		if !strings.Contains(out.String(), input) {
			t.Fatalf("Expected output to contain '%s', got '%s'", input, out.String())
		}
	}
}
