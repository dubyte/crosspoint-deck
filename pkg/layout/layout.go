package layout

import (
	"fmt"
	"image/color"
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

// LoadFontFaceBold loads a bold TTF font at the given size.
// If path is empty, it attempts to find a system bold font.
func LoadFontFaceBold(dc *gg.Context, path string, size float64) error {
	if path == "" {
		// Derive bold from regular path if one was previously used, else find system bold.
		path = findSystemFontBold()
	}
	if path == "" {
		// Fall back to regular.
		return LoadFontFace(dc, "", size)
	}
	return dc.LoadFontFace(path, size)
}

// DrawReversedHeader draws a black bar with white text and a 2px divider below it.
// Returns the y coordinate where body content should start.
func DrawReversedHeader(dc *gg.Context, title string, W int, fontSize float64, fontPath string) float64 {
	headerH := 56.0

	// Black bar
	dc.SetColor(color.Black)
	dc.DrawRectangle(0, 0, float64(W), headerH)
	dc.Fill()

	// White title text
	_ = LoadFontFaceBold(dc, fontPath, fontSize)
	dc.SetColor(color.White)
	dc.DrawStringAnchored(title, float64(W)/2, headerH/2, 0.5, 0.35)

	// Divider line
	dc.SetColor(color.Black)
	dc.SetLineWidth(2)
	dc.DrawLine(20, headerH+10, float64(W)-20, headerH+10)
	dc.Stroke()

	return headerH + 24 // body starts below the divider
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
	return findFont(candidatesRegular)
}

func findSystemFontBold() string {
	return findFont(candidatesBold)
}

var candidatesRegular = []string{
	"/usr/share/fonts/TTF/DejaVuSans.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	"/usr/share/fonts/liberation/LiberationSans-Regular.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
	"/usr/share/fonts/noto/NotoSans-Regular.ttf",
	"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
	"/System/Library/Fonts/Helvetica.ttc",
	"/Windows/Fonts/arial.ttf",
}

var candidatesBold = []string{
	"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	"/usr/share/fonts/liberation/LiberationSans-Bold.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	"/usr/share/fonts/noto/NotoSans-Bold.ttf",
	"/usr/share/fonts/truetype/noto/NotoSans-Bold.ttf",
	"/System/Library/Fonts/Helvetica.ttc",
	"/Windows/Fonts/arialbd.ttf",
}

func findFont(candidates []string) string {
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
