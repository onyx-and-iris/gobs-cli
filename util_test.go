package main

import "testing"

func TestSnakeCaseToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "Hello World"},
		{"snake_case_to_title_case", "Snake Case To Title Case"},
	}

	for _, test := range tests {
		result := snakeCaseToTitleCase(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s' but got '%s'", test.expected, result)
		}
	}
}
