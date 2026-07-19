package pack

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dubyte/crosspoint-deck/pkg/bmp"
	"github.com/dubyte/crosspoint-deck/pkg/card"
)

// Manifest describes a deck pack.
type Manifest struct {
	Name        string     `json:"name"`
	Author      string     `json:"author"`
	Version     string     `json:"version"`
	License     string     `json:"license,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Cards       []CardMeta `json:"cards"`
}

// CardMeta describes a single card in a pack.
type CardMeta struct {
	File  string   `json:"file"`
	Title string   `json:"title"`
	Tags  []string `json:"tags,omitempty"`
}

// PackBuilder constructs a deck pack.
type PackBuilder struct {
	manifest Manifest
	cards    []cardEntry
}

type cardEntry struct {
	meta CardMeta
	c    card.Card
}

// NewBuilder creates a new pack builder.
func NewBuilder(name, author, version string) *PackBuilder {
	return &PackBuilder{
		manifest: Manifest{
			Name:      name,
			Author:    author,
			Version:   version,
			CreatedAt: time.Now().UTC(),
			Cards:     []CardMeta{},
		},
	}
}

// SetDescription sets the pack description.
func (b *PackBuilder) SetDescription(d string) *PackBuilder {
	b.manifest.Description = d
	return b
}

// SetLicense sets the pack license.
func (b *PackBuilder) SetLicense(l string) *PackBuilder {
	b.manifest.License = l
	return b
}

// AddTag adds a tag to the pack.
func (b *PackBuilder) AddTag(t string) *PackBuilder {
	b.manifest.Tags = append(b.manifest.Tags, t)
	return b
}

// AddCard adds a card to the pack.
func (b *PackBuilder) AddCard(title string, c card.Card) *PackBuilder {
	b.cards = append(b.cards, cardEntry{
		meta: CardMeta{Title: title},
		c:    c,
	})
	return b
}

// Build renders all cards and writes a .deckpack.zip to the given path.
// It also writes a preview.png (first card rendered at half size).
func (b *PackBuilder) Build(outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("pack: mkdir: %w", err)
	}

	// Render all cards to BMPs
	bmpDir := filepath.Join(outputPath, "bmps")
	if err := os.MkdirAll(bmpDir, 0755); err != nil {
		return err
	}

	for i, entry := range b.cards {
		filename := fmt.Sprintf("card_%03d.bmp", i+1)
		b.manifest.Cards = append(b.manifest.Cards, CardMeta{
			File:  "bmps/" + filename,
			Title: entry.meta.Title,
		})

		path := filepath.Join(bmpDir, filename)
		if err := b.renderCard(entry.c, path); err != nil {
			return fmt.Errorf("pack: render %s: %w", entry.meta.Title, err)
		}
	}

	// Write manifest
	manifestPath := filepath.Join(outputPath, "manifest.json")
	mf, err := os.Create(manifestPath)
	if err != nil {
		return err
	}
	defer mf.Close()
	enc := json.NewEncoder(mf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b.manifest); err != nil {
		return err
	}
	mf.Close()

	// Generate preview from first card
	if len(b.cards) > 0 {
		previewPath := filepath.Join(outputPath, "preview.png")
		if err := b.generatePreview(b.cards[0].c, previewPath); err != nil {
			return fmt.Errorf("pack: preview: %w", err)
		}
	}

	return nil
}

// BuildZip creates a .deckpack.zip file.
func (b *PackBuilder) BuildZip(zipPath string) error {
	tmpDir, err := os.MkdirTemp("", "deckpack-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := b.Build(tmpDir); err != nil {
		return err
	}

	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	// Walk tmpDir and add files to zip
	return filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(tmpDir, path)
		if err != nil {
			return err
		}

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}

		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		_, err = io.Copy(w, r)
		return err
	})
}

func (b *PackBuilder) renderCard(c card.Card, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return bmp.Encode(f, c.Render())
}

func (b *PackBuilder) generatePreview(c card.Card, path string) error {
	img := c.Render()
	bounds := img.Bounds()

	// Scale down to 400x240 for preview
	preview := image.NewRGBA(image.Rect(0, 0, 400, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 400; x++ {
			srcX := bounds.Min.X + x*bounds.Dx()/400
			srcY := bounds.Min.Y + y*bounds.Dy()/240
			preview.Set(x, y, img.At(srcX, srcY))
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, preview)
}
