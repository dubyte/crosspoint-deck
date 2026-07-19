package card

import "image"

// Card is anything that can render itself to an 800x480 (or 480x800) bitmap.
type Card interface {
	// Render returns an image.Image. Implementations must ensure the result
	// is exactly 800x480 for portrait cards or 480x800 for landscape cards.
	// All content should use pure black (#000000) and white (#FFFFFF) only.
	Render() image.Image
}

// Meta holds optional sidecar metadata for host-side tooling.
type Meta struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
