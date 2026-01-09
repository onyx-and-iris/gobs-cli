// Package util provides utility functions for the application.

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
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

func parseTimeStringToMilliseconds(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	var durationStr string

	switch len(parts) {
	case 1:
		// Format: SS -> "SSs"
		durationStr = parts[0] + "s"
	case 2:
		// Format: MM:SS -> "MMmSSs"
		durationStr = parts[0] + "m" + parts[1] + "s"
	case 3:
		// Format: HH:MM:SS -> "HHhMMmSSs"
		durationStr = parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
	default:
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration.Seconds() * 1000, nil
}

func formatMillisecondsToTimeString(ms float64) string {
	totalSeconds := int(ms / 1000)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
