package render

import (
	"fmt"
	"image"
	"os"

	"github.com/dubyte/crosspoint-deck/pkg/bmp"
	"github.com/dubyte/crosspoint-deck/pkg/card"
)

// Content dimensions.
const (
	LandscapeW = 800
	LandscapeH = 480
	PortraitW  = 480
	PortraitH  = 800
)

// Display orientation.
const (
	DisplayPortrait  = "portrait"
	DisplayLandscape = "landscape"
)

// ToFile renders a card and writes an uncompressed 24-bit BMP to path.
// If display orientation differs from content orientation, the image is
// rotated 90° clockwise so it fills the display correctly.
func ToFile(c card.Card, display, path string) error {
	img := c.Render()

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if !isValidDimension(w, h) {
		return fmt.Errorf("render: expected %dx%d or %dx%d, got %dx%d",
			LandscapeW, LandscapeH, PortraitW, PortraitH, w, h)
	}

	contentIsPortrait := w == PortraitW && h == PortraitH
	displayIsPortrait := display == DisplayPortrait

	if contentIsPortrait != displayIsPortrait {
		img = rotate90(img)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("render: create file: %w", err)
	}
	defer f.Close()

	if err := bmp.Encode(f, img); err != nil {
		return fmt.Errorf("render: encode bmp: %w", err)
	}
	return f.Close()
}

func isValidDimension(w, h int) bool {
	return (w == LandscapeW && h == LandscapeH) || (w == PortraitW && h == PortraitH)
}

// rotate90 rotates an image 90° clockwise.
func rotate90(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, srcH, srcW))

	for y := 0; y < srcH; y++ {
		for x := 0; x < srcW; x++ {
			dst.Set(srcH-1-y, x, src.At(x+bounds.Min.X, y+bounds.Min.Y))
		}
	}
	return dst
}
