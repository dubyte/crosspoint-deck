package calendar

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
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

	dc.SetColor(color.White)
	dc.Clear()

	// Reversed header with year
	_ = layout.LoadFontFaceBold(dc, y.FontPath, 26)
	headerH := 50.0
	dc.SetColor(color.Black)
	dc.DrawRectangle(0, 0, float64(W), headerH)
	dc.Fill()
	dc.SetColor(color.White)
	dc.DrawStringAnchored(fmt.Sprintf("%d", y.Year), float64(W)/2, headerH/2, 0.5, 0.35)

	// 2px divider
	dc.SetColor(color.Black)
	dc.SetLineWidth(2)
	dc.DrawLine(20, headerH+10, float64(W)-20, headerH+10)
	dc.Stroke()

	gap := 16.0
	topY := headerH + 22
	cellW := (float64(W) - gap*float64(cols-1)) / float64(cols)
	cellH := (float64(H) - topY - gap*float64(rows-1)) / float64(rows)

	startMonth := time.Date(y.Year, time.January, 1, 0, 0, 0, 0, time.UTC)

	for m := 0; m < 12; m++ {
		col := m % cols
		row := m / cols
		x0 := float64(col) * (cellW + gap)
		y0 := topY + float64(row)*(cellH+gap)

		monthDate := startMonth.AddDate(0, m, 0)
		drawMonth(dc, monthDate, x0, y0, cellW, cellH, y.FontPath)
	}

	return dc.Image()
}

func drawMonth(dc *gg.Context, t time.Time, x0, y0, cw, ch float64, fontPath string) {
	// Use smaller fonts when cell is tight (landscape), larger when roomy (portrait).
	tight := ch < 140
	monthSize, headerSize, daySize := 18.0, 14.0, 16.0
	monthOff := 16.0
	headerOff := 34.0
	dayOff := 18.0
	dayBase := 52.0
	if tight {
		monthSize, headerSize, daySize = 14, 12, 12
		monthOff = 12
		headerOff = 26
		dayOff = 14
		dayBase = 40
	}

	// Month name in bold
	_ = layout.LoadFontFaceBold(dc, fontPath, monthSize)
	dc.SetColor(color.Black)
	monthName := t.Format("Jan")
	dc.DrawStringAnchored(monthName, x0+cw/2, y0+monthOff, 0.5, 0.5)

	// Day headers: S M T W T F S
	_ = layout.LoadFontFace(dc, fontPath, headerSize)
	days := []string{"S", "M", "T", "W", "T", "F", "S"}
	headerY := y0 + headerOff
	cellW := cw / 7
	for i, d := range days {
		dc.DrawStringAnchored(d, x0+cellW*float64(i)+cellW/2, headerY, 0.5, 0.5)
	}

	// Days grid
	firstDay := t
	lastDay := firstDay.AddDate(0, 1, -1)
	startWeekday := int(firstDay.Weekday()) // 0=Sunday

	dayY := headerY + dayOff
	dayH := (ch - dayBase) / 6

	_ = layout.LoadFontFace(dc, fontPath, daySize)
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

// Spec returns the card.Spec for calendar.
func Spec() card.Spec {
	return card.Spec{
		Name:  "calendar",
		Usage: "Generate a year-at-a-glance calendar card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &YearCard{}
			fs.IntVar(&c.Year, "year", time.Now().Year(), "Year to render")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation (480x800)")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
