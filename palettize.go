/*
Palettize creates a composite image using the brightness of one image and the
color palette of another. Only works with PNGs.

Example syntax:
    ./palettizer original.png palette.png result.png

The alogorithm used gets a list of the colors from each input file and sorts
them by brightness. The color of each pixel in the first image is mapped onto
the color at the corresponding point in the second image's color list in order
to produce the result image.
*/
package main

import (
    "fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s original palette result\n",
                    os.Args[0])
        return
	}
	valueImg := readImage(os.Args[1])

	oldPalette := getPalette(valueImg)
	newPalette := getPalette(readImage(os.Args[2]))

	ratio := float64(len(newPalette)) / float64(len(oldPalette))

	b := valueImg.Bounds()
	imgOut := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for x := b.Min.X; x < b.Max.X; x++ {
        println(x)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			index := indexOf(valueImg.At(x, y), oldPalette)
			if index == -1 {
				imgOut.Set(x, y, valueImg.At(x, y))
			} else {
				imgOut.Set(x, y, newPalette[int(float64(index)*ratio)])
			}
		}
	}
    println("done")

	writeImage(imgOut, os.Args[3])
}

// Gets a slice of colors from an image, sorted from least to most brightness.
func getPalette(img image.Image) []color.Color {
	palette := make([]color.Color, 0)
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if !transparent(img.At(x, y)) {
				palette = append(palette, img.At(x, y))
			}
		}
	}
	sort.Sort(ByBrightness(palette))

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
	return brightness(a[i]) < brightness(a[j])
}

// Returns true if the color is transparent, false if it is opaque.
func transparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

// Gets the brightness of a color (the sum of its red, green, and blue values).
func brightness(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return r + g + b
}

// Gets an image from a PNG file.
func readImage(filename string) image.Image {
	file, _ := os.Open(filename)
	image, _ := png.Decode(file)
	return image
}

// Writes an image to a PNG file.
func writeImage(img image.Image, filename string) {
	file, _ := os.Create(filename)
	png.Encode(file, img)
}

// vim: ts=4 sw=0
