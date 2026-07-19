package render

import (
	"fmt"
	"os"

	"github.com/dubyte/crosspoint-deck/pkg/bmp"
	"github.com/dubyte/crosspoint-deck/pkg/card"
)

// Target dimensions for the XTEink X4 display.
const (
	PortraitW  = 800
	PortraitH  = 480
	LandscapeW = 480
	LandscapeH = 800
)

// ToFile renders a card and writes an uncompressed 24-bit BMP to path.
// It validates exact dimensions and warns if the design contains gray pixels.
func ToFile(c card.Card, path string) error {
	img := c.Render()

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if !isValidDimension(w, h) {
		return fmt.Errorf("render: expected %dx%d or %dx%d, got %dx%d",
			PortraitW, PortraitH, LandscapeW, LandscapeH, w, h)
	}

	if !bmp.IsBlackWhite(img) {
		// We still proceed because the encoder will threshold, but warn.
		// In a future version this could be a strict error.
		fmt.Fprintf(os.Stderr, "render: warning: design contains non-black/white pixels; they will be thresholded\n")
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
	return (w == PortraitW && h == PortraitH) || (w == LandscapeW && h == LandscapeH)
}
