package calendar

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/fogleman/gg"
)

// YearCard renders a year-at-a-glance calendar.
type YearCard struct {
	Year int
	// Portrait renders at 480x800 instead of 800x480.
	Portrait bool
	// FontPath is optional; if empty, common system fonts are tried.
	FontPath string
}

// Render produces the calendar image.
// Landscape: 800x480 with a 4x3 month grid.
// Portrait:  480x800 with a 3x4 month grid.
func (y *YearCard) Render() image.Image {
	var W, H int
	var cols, rows int
	if y.Portrait {
		W, H = 480, 800
		cols, rows = 3, 4
	} else {
		W, H = 800, 480
		cols, rows = 4, 3
	}

	dc := gg.NewContext(W, H)

	// White background
	dc.SetColor(color.White)
	dc.Clear()

	// Try to load a font
	fontPath := y.FontPath
	if fontPath == "" {
		fontPath = findFont()
	}
	if fontPath != "" {
		_ = dc.LoadFontFace(fontPath, 14)
	}

	// Title: year
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(fmt.Sprintf("%d", y.Year), float64(W)/2, 24, 0.5, 0.5)

	cellW := float64(W) / float64(cols)
	cellH := (float64(H) - 48) / float64(rows) // leave top 48px for title + margin

	startMonth := time.Date(y.Year, time.January, 1, 0, 0, 0, 0, time.UTC)

	for m := 0; m < 12; m++ {
		col := m % cols
		row := m / cols
		x0 := float64(col) * cellW
		y0 := 48 + float64(row)*cellH

		monthDate := startMonth.AddDate(0, m, 0)
		drawMonth(dc, monthDate, x0, y0, cellW, cellH)
	}

	return dc.Image()
}

func drawMonth(dc *gg.Context, t time.Time, x0, y0, cw, ch float64) {
	// Month name
	dc.SetColor(color.Black)
	monthName := t.Format("Jan")
	dc.DrawStringAnchored(monthName, x0+cw/2, y0+14, 0.5, 0.5)

	// Day headers: S M T W T F S
	days := []string{"S", "M", "T", "W", "T", "F", "S"}
	headerY := y0 + 30
	cellW := cw / 7
	for i, d := range days {
		dc.DrawStringAnchored(d, x0+cellW*float64(i)+cellW/2, headerY, 0.5, 0.5)
	}

	// Days grid
	firstDay := t
	lastDay := firstDay.AddDate(0, 1, -1)
	startWeekday := int(firstDay.Weekday()) // 0=Sunday

	dayY := headerY + 18
	dayH := (ch - 48) / 6 // up to 6 weeks

	day := 1
	for week := 0; week < 6; week++ {
		for wd := 0; wd < 7; wd++ {
			if week == 0 && wd < startWeekday {
				continue
			}
			if day > lastDay.Day() {
				break
			}
			cx := x0 + cellW*float64(wd) + cellW/2
			cy := dayY + dayH*float64(week)
			dc.DrawStringAnchored(fmt.Sprintf("%d", day), cx, cy, 0.5, 0.5)
			day++
		}
	}
}

func findFont() string {
	candidates := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",           // macOS
		"/Windows/Fonts/arial.ttf",                       // Windows
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
