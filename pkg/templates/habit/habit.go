package habit

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Title    string
	Habits   []string
	Days     int
	Portrait bool
	FontPath string
}

func (h *Card) Render() image.Image {
	var W, H int
	if h.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, h.Title, W, 26, h.FontPath)

	if len(h.Habits) == 0 {
		return dc.Image()
	}

	days := h.Days
	if days <= 0 {
		days = 7
	}
	if days > 31 {
		days = 31
	}

	cols := days
	rows := len(h.Habits)
	marginX := 90.0
	marginY := bodyY + 24
	cellW := (float64(W) - marginX - 20) / float64(cols)
	cellH := (float64(H) - marginY - 20) / float64(rows)
	if cellH > 44 {
		cellH = 44
	}

	// Day headers
	_ = layout.LoadFontFace(dc, h.FontPath, 16)
	for d := 0; d < days; d++ {
		x := marginX + float64(d)*cellW + cellW/2
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(fmtDay(d+1), x, marginY-12, 0.5, 0.5)
	}

	// Habit rows
	_ = layout.LoadFontFaceBold(dc, h.FontPath, 16)
	for r, habit := range h.Habits {
		y := marginY + float64(r)*cellH + cellH/2
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(habit, marginX-12, y, 1.0, 0.5)

		for d := 0; d < days; d++ {
			x := marginX + float64(d)*cellW + cellW/2
			boxSize := min(cellW, cellH) * 0.6
			dc.DrawRectangle(x-boxSize/2, y-boxSize/2, boxSize, boxSize)
			dc.SetColor(color.Black)
			dc.SetLineWidth(1.5)
			dc.Stroke()
		}
	}

	return dc.Image()
}

func fmtDay(n int) string {
	return fmt.Sprintf("%2d", n)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "habit",
		Usage: "Generate a habit tracker card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var habitsRaw string
			fs.StringVar(&c.Title, "title", "Habit Tracker", "Tracker title")
			fs.StringVar(&habitsRaw, "habits", "Read,Exercise,Meditate", "Comma-separated habit names")
			fs.IntVar(&c.Days, "days", 7, "Number of days (max 31)")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &habitCard{c, habitsRaw}
		},
	}
}

type habitCard struct {
	*Card
	habitsRaw string
}

func (hc *habitCard) Render() image.Image {
	if hc.habitsRaw != "" && len(hc.Card.Habits) == 0 {
		hc.Card.Habits = strings.Split(hc.habitsRaw, ",")
		for i := range hc.Card.Habits {
			hc.Card.Habits[i] = strings.TrimSpace(hc.Card.Habits[i])
		}
	}
	return hc.Card.Render()
}
