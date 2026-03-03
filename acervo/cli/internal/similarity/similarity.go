package similarity

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math/bits"
	"os"

	"golang.org/x/image/draw"
)

// DHash generates a 64-bit difference hash for an image.
// It resizes the image to 9x8 and compares adjacent pixels.
func DHash(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, fmt.Errorf("decode error: %w", err)
	}

	// Resize to 9x8
	// 9 columns so we can compare 8 pairs horizontally
	resized := image.NewRGBA(image.Rect(0, 0, 9, 8))
	draw.BiLinear.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Convert to grayscale and compute bits
	var hash uint64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			p1 := grayscale(resized.At(x, y))
			p2 := grayscale(resized.At(x+1, y))
			if p1 < p2 {
				hash |= 1 << (uint(y*8 + x))
			}
		}
	}

	return hash, nil
}

func grayscale(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	// Traditional grayscale conversion
	return (r*299 + g*587 + b*114) / 1000
}

// HammingDistance calculates the number of bits that differ between two hashes.
func HammingDistance(h1, h2 uint64) int {
	return bits.OnesCount64(h1 ^ h2)
}

// Probability returns a 0.0-1.0 similarity score given Hamming Distance.
func Probability(distance int) float64 {
	return 1.0 - (float64(distance) / 64.0)
}
