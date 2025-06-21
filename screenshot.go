package main

import (
	"fmt"
	"path/filepath"

	"github.com/andreykaipov/goobs/api/requests/sources"
)

// ScreenshotCmd provides commands to manage screenshots in OBS Studio.
type ScreenshotCmd struct {
	Save ScreenshotSaveCmd `cmd:"" help:"Take a screenshot and save it to a file." aliases:"sv"`
}

// ScreenshotSaveCmd represents the command to save a screenshot of a source in OBS.
type ScreenshotSaveCmd struct {
	SourceName string  `arg:"" help:"Name of the source to take a screenshot of."`
	FilePath   string  `arg:"" help:"Path to the file where the screenshot will be saved."`
	Width      float64 `       help:"Width of the screenshot in pixels."                   flag:"" default:"1920"`
	Height     float64 `       help:"Height of the screenshot in pixels."                  flag:"" default:"1080"`
	Quality    float64 `       help:"Quality of the screenshot (1-100)."                   flag:"" default:"-1"`
}

// Run executes the command to take a screenshot and save it to a file.
func (cmd *ScreenshotSaveCmd) Run(ctx *context) error {
	_, err := ctx.Client.Sources.SaveSourceScreenshot(
		sources.NewSaveSourceScreenshotParams().
			WithSourceName(cmd.SourceName).
			WithImageFormat(trimPrefix(filepath.Ext(cmd.FilePath), ".")).
			WithImageFilePath(cmd.FilePath).
			WithImageWidth(cmd.Width).
			WithImageHeight(cmd.Height).
			WithImageCompressionQuality(cmd.Quality),
	)
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %w", err)
	}

	fmt.Fprintf(ctx.Out, "Screenshot saved to %s.\n", ctx.Style.Highlight(cmd.FilePath))
	return nil
}
