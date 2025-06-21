// nolint: misspell
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Style defines colours for the table styles.
type Style struct {
	name      string
	border    lipgloss.Color
	oddRows   lipgloss.Color
	evenRows  lipgloss.Color
	highlight lipgloss.Color
}

// Highlight applies the highlight style to the given text.
func (s *Style) Highlight(text string) string {
	return lipgloss.NewStyle().Foreground(s.highlight).Render(text)
}

func (s *Style) Error(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(text) // Red for errors
}

func newRedStyle() *Style {
	return &Style{
		name:      "red",
		border:    lipgloss.Color("#D32F2F"), // Strong red for border
		oddRows:   lipgloss.Color("#FFCDD2"), // Very light red for odd rows
		evenRows:  lipgloss.Color("#EF9A9A"), // Light red for even rows
		highlight: lipgloss.Color("#EF9A9A"),
	}
}

func newMagentaStyle() *Style {
	return &Style{
		name:      "magenta",
		border:    lipgloss.Color("#C2185B"), // Strong magenta for border
		oddRows:   lipgloss.Color("#F8BBD0"), // Very light magenta/pink for odd rows
		evenRows:  lipgloss.Color("#F48FB1"), // Light magenta/pink for even rows
		highlight: lipgloss.Color("#F48FB1"),
	}
}

func newPurpleStyle() *Style {
	return &Style{
		name:      "purple",
		border:    lipgloss.Color("#7B1FA2"), // Strong purple for border
		oddRows:   lipgloss.Color("#E1BEE7"), // Very light purple for odd rows
		evenRows:  lipgloss.Color("#CE93D8"), // Light purple for even rows
		highlight: lipgloss.Color("#CE93D8"),
	}
}

func newBlueStyle() *Style {
	return &Style{
		name:      "blue",
		border:    lipgloss.Color("#1976D2"), // Medium blue for border
		oddRows:   lipgloss.Color("#E3F2FD"), // Very light blue for odd rows
		evenRows:  lipgloss.Color("#BBDEFB"), // Light blue for even rows
		highlight: lipgloss.Color("#1976D2"),
	}
}

func newCyanStyle() *Style {
	return &Style{
		name:      "cyan",
		border:    lipgloss.Color("#00BFCF"), // A strong cyan for border
		oddRows:   lipgloss.Color("#E0F7FA"), // Very light cyan for odd rows
		evenRows:  lipgloss.Color("#B2EBF2"), // Slightly darker light cyan for even rows
		highlight: lipgloss.Color("#00BFCF"),
	}
}

func newGreenStyle() *Style {
	return &Style{
		name:      "green",
		border:    lipgloss.Color("#43A047"), // Medium green for border
		oddRows:   lipgloss.Color("#E8F5E9"), // Very light green for odd rows
		evenRows:  lipgloss.Color("#C8E6C9"), // Light green for even rows
		highlight: lipgloss.Color("#43A047"),
	}
}

func newYellowStyle() *Style {
	return &Style{
		name:      "yellow",
		border:    lipgloss.Color("#FBC02D"), // Strong yellow for border
		oddRows:   lipgloss.Color("#FFF9C4"), // Very light yellow for odd rows
		evenRows:  lipgloss.Color("#FFF59D"), // Light yellow for even rows
		highlight: lipgloss.Color("#FBC02D"),
	}
}

func newOrangeStyle() *Style {
	return &Style{
		name:      "orange",
		border:    lipgloss.Color("#F57C00"), // Strong orange for border
		oddRows:   lipgloss.Color("#FFF3E0"), // Very light orange for odd rows
		evenRows:  lipgloss.Color("#FFE0B2"), // Light orange for even rows
		highlight: lipgloss.Color("#F57C00"),
	}
}

func newWhiteStyle() *Style {
	return &Style{
		name:      "white",
		border:    lipgloss.Color("#FFFFFF"), // White for border
		oddRows:   lipgloss.Color("#F0F0F0"), // Very light grey for odd rows
		evenRows:  lipgloss.Color("#E0E0E0"), // Light grey for even rows
		highlight: lipgloss.Color("#FFFFFF"), // Highlight in white
	}
}

func newGreyStyle() *Style {
	return &Style{
		name:      "grey",
		border:    lipgloss.Color("#9E9E9E"), // Medium grey for border
		oddRows:   lipgloss.Color("#F5F5F5"), // Very light grey for odd rows
		evenRows:  lipgloss.Color("#EEEEEE"), // Light grey for even rows
		highlight: lipgloss.Color("#9E9E9E"), // Highlight in medium grey
	}
}

func newNavyBlueStyle() *Style {
	return &Style{
		name:      "navy",
		border:    lipgloss.Color("#001F3F"), // Navy blue for border
		oddRows:   lipgloss.Color("#CFE2F3"), // Very light blue for odd rows
		evenRows:  lipgloss.Color("#A9CCE3"), // Light blue for even rows
		highlight: lipgloss.Color("#001F3F"), // Highlight in navy blue
	}
}

func newBlackStyle() *Style {
	return &Style{
		name:      "black",
		border:    lipgloss.Color("#000000"), // Black for border
		oddRows:   lipgloss.Color("#333333"), // Dark grey for odd rows
		evenRows:  lipgloss.Color("#444444"), // Slightly lighter dark grey for even rows
		highlight: lipgloss.Color("#000000"), // Highlight in black
	}
}

func styleFromFlag(colour string) *Style {
	switch colour {
	case "red":
		return newRedStyle()
	case "magenta":
		return newMagentaStyle()
	case "purple":
		return newPurpleStyle()
	case "blue":
		return newBlueStyle()
	case "cyan":
		return newCyanStyle()
	case "green":
		return newGreenStyle()
	case "yellow":
		return newYellowStyle()
	case "orange":
		return newOrangeStyle()
	case "white":
		return newWhiteStyle()
	case "grey":
		return newGreyStyle()
	case "navy":
		return newNavyBlueStyle()
	case "black":
		return newBlackStyle()
	default:
		err := os.Setenv("NO_COLOR", "1") // nolint: misspell
		if err != nil {
			// If we can't set NO_COLOR, we just log the error and continue
			// This is a fallback to ensure that the application can still run
			fmt.Fprintf(os.Stderr, "Error setting NO_COLOR: %v\n", err)
		}
		return &Style{}
	}
}
