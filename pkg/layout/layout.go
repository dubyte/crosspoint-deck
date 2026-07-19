package layout

import (
	"fmt"
	"os"

	"github.com/fogleman/gg"
)

// LoadFontFace loads a TTF font at the given size for gg.
// If path is empty, it attempts to find a system font.
func LoadFontFace(dc *gg.Context, path string, size float64) error {
	if path == "" {
		path = findSystemFont()
	}
	if path == "" {
		return fmt.Errorf("layout: no font found")
	}
	return dc.LoadFontFace(path, size)
}

// DrawCenteredText draws text centered at (cx, cy).
func DrawCenteredText(dc *gg.Context, text string, cx, cy float64) {
	dc.DrawStringAnchored(text, cx, cy, 0.5, 0.5)
}

// DrawWrappedText draws text wrapped to maxWidth, returning lines.
func DrawWrappedText(dc *gg.Context, text string, x, y, maxWidth, lineHeight float64) {
	words := splitWords(text)
	var line string
	var lines []string
	for _, word := range words {
		if line == "" {
			line = word
		} else {
			test := line + " " + word
			w, _ := dc.MeasureString(test)
			if w > maxWidth {
				lines = append(lines, line)
				line = word
			} else {
				line = test
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	for i, l := range lines {
		dc.DrawString(l, x, y+float64(i)*lineHeight)
	}
}

// Grid computes cell dimensions for a grid layout.
type Grid struct {
	Cols, Rows   int
	CellW, CellH float64
	MarginX      float64
	MarginY      float64
}

// NewGrid creates a grid fitting the given dimensions with optional margins.
func NewGrid(cols, rows int, width, height, marginX, marginY float64) Grid {
	return Grid{
		Cols:    cols,
		Rows:    rows,
		CellW:   (width - 2*marginX) / float64(cols),
		CellH:   (height - 2*marginY) / float64(rows),
		MarginX: marginX,
		MarginY: marginY,
	}
}

// Cell returns the top-left corner of cell (col, row).
func (g Grid) Cell(col, row int) (x0, y0 float64) {
	return g.MarginX + float64(col)*g.CellW, g.MarginY + float64(row)*g.CellH
}

func findSystemFont() string {
	candidates := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/Windows/Fonts/arial.ttf",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func splitWords(text string) []string {
	var words []string
	var word string
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(r)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}
