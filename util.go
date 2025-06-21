// Package util provides utility functions for the application.

package main

import (
	"os"
	"strings"
)

func snakeCaseToTitleCase(snake string) string {
	words := strings.Split(snake, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func getEnabledMark(enabled bool) string {
	if enabled {
		if os.Getenv("NO_COLOR") != "" { // nolint: misspell
			return "✓"
		}
		return "✅"
	}
	if os.Getenv("NO_COLOR") != "" { // nolint: misspell
		return "✗"
	}
	return "❌"
}

func trimPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}
