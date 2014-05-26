/*
Palettize creates a composite image using the brightness of one image and the
color palette of another. Only works with PNGs.

Example syntax:
    ./palettize original.png palette.png result.png

The alogorithm used gets a list of the colors from each input file and sorts
them by brightness. The color of each pixel in the first image is mapped onto
the color at the corresponding point in the second image's color list in order
to produce the result image.
*/
package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
)

// Prints an error message to stderr and exits with a non-zero status.
func die(err error) {
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 4 {
		die(errors.New(fmt.Sprintf("Usage: %s original palette result",
			os.Args[0])))
	}

	valueImg := readImage(os.Args[1])

	oldPalette := getPalette(valueImg)
	newPalette := getPalette(readImage(os.Args[2]))

	ratio := float64(len(newPalette)) / float64(len(oldPalette))

	b := valueImg.Bounds()
	imgOut := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	width := b.Max.X - b.Min.X
	for x := b.Min.X; x < b.Max.X; x++ {

		// Progress display
		colNumber := x - b.Min.X + 1
		fmt.Printf("\rConverting column %d of %d (%d%%)", colNumber, width,
			100*colNumber/width)
		os.Stdout.Sync()

		for y := b.Min.Y; y < b.Max.Y; y++ {
			index := indexOf(valueImg.At(x, y), oldPalette)
			if index != -1 {
				imgOut.Set(x, y, newPalette[int(float64(index)*ratio)])
			} else {
				imgOut.Set(x, y, valueImg.At(x, y))
			}
		}
	}

	// Erase progress display
	print("\r")
	os.Stdout.Sync()

	writeImage(imgOut, os.Args[3])
}

// Gets a slice of colors from an image, sorted from least to most brightness.
func getPalette(img image.Image) []color.Color {

	// Get colors from image
	allColors := make([]color.Color, 0)
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if !transparent(img.At(x, y)) {
				allColors = append(allColors, img.At(x, y))
			}
		}
	}

	// Convert slice of colors into sorted set of (unique) colors
	sort.Sort(ByBrightness(allColors))
	palette := make([]color.Color, 0)
	for _, c := range allColors {
		if len(palette) == 0 || palette[len(palette)-1] != c {
			palette = append(palette, c)
		}
	}

	return palette
}

// Gets the index of a color in a slice of colors.
func indexOf(c color.Color, colors []color.Color) int {
	for i := 0; i < len(colors); i++ {
		if colors[i] == c {
			return i
		}
	}

	return -1
}

// ByBrightness implements sort.Interface for []color.Color based on value
// (brightness).
type ByBrightness []color.Color

func (a ByBrightness) Len() int      { return len(a) }
func (a ByBrightness) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByBrightness) Less(i, j int) bool {
	ri, gi, bi, _ := a[i].RGBA()
	rj, gj, bj, _ := a[j].RGBA()
	return (ri + gi + bi) < (rj + gj + bj)
}

// Returns true if the color is transparent, false if it is opaque.
func transparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

// Gets an image from a PNG file.
func readImage(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		die(err)
	}
	image, err := png.Decode(file)
	if err != nil {
		die(err)
	}
	return image
}

// Writes an image to a PNG file.
func writeImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		die(err)
	}
	png.Encode(file, img)
}

// vim: ts=4 sw=0
