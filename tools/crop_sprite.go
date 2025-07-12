package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run tools/crop_sprite.go <input_sprite_sheet> <output_sprite_sheet>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Open the sprite sheet
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal("Error opening sprite sheet:", err)
	}
	defer file.Close()

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal("Error decoding sprite sheet:", err)
	}

	bounds := img.Bounds()
	fmt.Printf("Original sprite sheet dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Find the bounds of actual content (same logic as analyzer)
	foundPixels := false
	contentBounds := image.Rectangle{}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 { // Non-transparent pixel
				if !foundPixels {
					contentBounds.Min.X = x
					contentBounds.Min.Y = y
					contentBounds.Max.X = x
					contentBounds.Max.Y = y
					foundPixels = true
				} else {
					if x < contentBounds.Min.X {
						contentBounds.Min.X = x
					}
					if y < contentBounds.Min.Y {
						contentBounds.Min.Y = y
					}
					if x > contentBounds.Max.X {
						contentBounds.Max.X = x
					}
					if y > contentBounds.Max.Y {
						contentBounds.Max.Y = y
					}
				}
			}
		}
	}

	if !foundPixels {
		fmt.Println("❌ No non-transparent pixels found!")
		return
	}

	fmt.Printf("Content bounds: %s\n", contentBounds)
	fmt.Printf("Content dimensions: %dx%d\n", contentBounds.Dx(), contentBounds.Dy())

	// Crop the image to just the content area
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(contentBounds)

	// Save the cropped image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, croppedImg)
	if err != nil {
		log.Fatal("Error encoding cropped image:", err)
	}

	fmt.Printf("✅ Cropped sprite sheet saved to: %s\n", outputPath)
	fmt.Printf("New dimensions: %dx%d\n", contentBounds.Dx(), contentBounds.Dy())

	// Suggest grid configurations for the cropped image
	fmt.Println("\nSuggested grid configurations for cropped image:")
	fmt.Println("================================================")

	width := contentBounds.Dx()
	height := contentBounds.Dy()

	configs := []struct {
		frameWidth, frameHeight int
		description             string
	}{
		{width / 2, height / 2, "2x2 grid"},
		{width / 4, height / 2, "4x2 grid"},
		{width / 2, height / 4, "2x4 grid"},
		{width / 4, height / 4, "4x4 grid"},
		{width, height / 2, "1x2 grid"},
		{width / 2, height, "2x1 grid"},
	}

	for _, config := range configs {
		if config.frameWidth > 0 && config.frameHeight > 0 {
			cols := width / config.frameWidth
			rows := height / config.frameHeight
			if cols > 0 && rows > 0 {
				fmt.Printf("Frame size %dx%d → %dx%d grid (%s)\n",
					config.frameWidth, config.frameHeight, cols, rows, config.description)
			}
		}
	}
}
