/*
	Creates a composite image using the brightness of one image and the color
	palette of another. Only works with PNGs.

	Example syntax:
		./palettizer brightness.png palette.png result.png
	
	The alogorithm used gets a list of the colors from each input file and
	sorts them by brightness. The color of each pixel in the first image is
	mapped onto the color at the corresponding point in the second image's
	color list in order to produce the result image.
*/
package main

import (
	"os"
	"image"
	"image/png"
	"image/color"
)

func main() {
	if len(os.Args) != 4 {
		panic("Command format: palettize <image file> <pallete file> <output file>")
	}
	imgOriginal := readImage(os.Args[1])
	imgPalette := readImage(os.Args[2])

	oldPalette := getPalette(imgOriginal)
	newPalette := getPalette(imgPalette)

	ratio := float64(len(newPalette)) / float64(len(oldPalette))

	b := imgOriginal.Bounds()
	imgOut := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			index := indexOf(imgOriginal.At(x, y), oldPalette)
			if index == -1 {
				imgOut.Set(x, y, imgOriginal.At(x, y))
			} else {
				imgOut.Set(x, y, newPalette[int(float64(index) * ratio)])
			}
		}
	}

	writeImage(imgOut, os.Args[3])
}

// Get a slice of colors from an image, sorted from least to most brightness
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
	palette = sort(palette)

	return palette
}

// Get the index of a color in a slice of colors
func indexOf(c color.Color, colors []color.Color) int {
	for i := 0; i < len(colors); i++ {
		if colors[i] == c {
			return i
		}
	}

	return -1
}

// Sorts a slice of colors from least to most brightness
// There must be a better way to do this :)
func sort(unsorted []color.Color) []color.Color {
	sorted := make([]color.Color, 0)
	for _, col := range unsorted {
		for i, color2 := range sorted {
			if col == color2 {
				break
			} else if brightness(col) > brightness(color2) {
				oldSorted := sorted
				sorted = make([]color.Color, len(oldSorted) + 1)
				for j := 0; j < len(sorted); j++ {
					if j < i {
						sorted[j] = oldSorted[j]
					} else if j == i {
						sorted[j] = col
					} else {
						sorted[j] = oldSorted[j - 1]
					}
				}
				break
			}
		}

		if len(sorted) == 0 {
			sorted = append(sorted, col)
		}
	}

	return sorted
}

// Returns true if the color is transparent, false if it is opaque
func transparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

// Gets the brightness of a color (the sum of its red, green, and blue values)
func brightness(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return r + g + b
}

// Get an image from a PNG file
func readImage(filename string) image.Image {
	file, _ := os.Open(filename)
	image, _ := png.Decode(file)
	return image
}

// Write an image to a PNG file
func writeImage(img image.Image, filename string) {
	file, _ := os.Create(filename)
	png.Encode(file, img)
}
