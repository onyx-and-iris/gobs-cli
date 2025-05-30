// Package util provides utility functions for the application.

package main

import "strings"

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
		return "\u2713" // green check mark
	}
	return "\u274c" // red cross mark
}
