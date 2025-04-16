package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// GenerateThumbnail creates a thumbnail for the given image file.
func GenerateThumbnail(inputPath, outputPath string, width, height int) error {
	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open the source image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Resize the image to the specified dimensions
	thumbnail := imaging.Resize(src, width, height, imaging.Lanczos)

	// Save the thumbnail to the output path
	err = imaging.Save(thumbnail, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return nil
}
