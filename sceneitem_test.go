package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSceneItemList(t *testing.T) {
	client, disconnect := getClient(t)
	defer disconnect()

	var out bytes.Buffer
	context := &context{
		Client: client,
		Out:    &out,
	}

	cmd := &SceneItemListCmd{
		SceneName: "gobs-test",
	}
	err := cmd.Run(context)
	if err != nil {
		t.Fatalf("Failed to list scene items: %v", err)
	}
	if !strings.Contains(out.String(), "gobs-test-input") {
		t.Fatalf("Expected output to contain 'gobs-test-input', got '%s'", out.String())
	}
	if !strings.Contains(out.String(), "gobs-test-input-2") {
		t.Fatalf("Expected output to contain 'gobs-test-input-2', got '%s'", out.String())
	}
}
