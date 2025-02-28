package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

const (
	verticalStep   = 4
	horizontalStep = 2
	blockWidth     = verticalStep
	blockHeight    = horizontalStep
	asciiRamp      = " .:-=+*#%@"
)

type ImageProcessor struct {
	img image.Image
}

func NewImageProcessor(image image.Image) *ImageProcessor {
	return &ImageProcessor{img: image}
}

func (ip *ImageProcessor) CalculateBlockAverage(startX, startY int) int {
	var sum, count int
	bounds := ip.img.Bounds()
	maxX, maxY := bounds.Max.X, bounds.Max.Y

	for x := startX; x < startX+blockWidth && x < maxX; x++ {
		for y := startY; y < startY+blockHeight && y < maxY; y++ {
			sum += grayscale(ip.img.At(x, y))
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return sum / count
}

func (ip *ImageProcessor) Generate() string {
	var output string
	bounds := ip.img.Bounds()
	max := bounds.Max

	for y := 0; y < max.Y; y += verticalStep {
		for x := 0; x < max.X; x += horizontalStep {
			avg := ip.CalculateBlockAverage(x, y)
			output += mapBrightnessToChar(avg)
		}
		output += "\n"
	}

	return output
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func grayscale(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}

func mapBrightnessToChar(avgBrightness int) string {
	index := len(asciiRamp) * avgBrightness / 65536
	return string(asciiRamp[index])
}

func main() {
	img, err := loadImage("images.png")
	if err != nil {
		panic(err)
	}

	processor := NewImageProcessor(img)
	fmt.Print(processor.Generate())
}
