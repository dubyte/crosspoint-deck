package card

import (
	"flag"
	"image"
)

// Card is anything that can render itself to an 800x480 (or 480x800) bitmap.
type Card interface {
	// Render returns an image.Image. Implementations must ensure the result
	// is exactly 800x480 for portrait cards or 480x800 for landscape cards.
	// All content should use pure black (#000000) and white (#FFFFFF) only.
	Render() image.Image
}

// AssertDimensions verifies a card renders at the expected dimensions.
func AssertDimensions(t interface {
	Helper()
	Errorf(format string, args ...interface{})
}, c Card, wantW, wantH int) {
	t.Helper()
	img := c.Render()
	bounds := img.Bounds()
	gotW, gotH := bounds.Dx(), bounds.Dy()
	if gotW != wantW || gotH != wantH {
		t.Errorf("dimensions = %dx%d, want %dx%d", gotW, gotH, wantW, wantH)
	}
}

// Factory creates a Card, registering template-specific flags on fs.
// The returned Card should close over flag pointers.
type Factory func(fs *flag.FlagSet) Card

// Spec describes a card template for CLI registration.
type Spec struct {
	Name    string // CLI subcommand name, e.g. "calendar"
	Usage   string // one-line description for help text
	New     Factory
}

// Meta holds optional sidecar metadata for host-side tooling.
type Meta struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
