package pomodoro

import (
	"flag"
	"image"
	"fmt"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Title    string
	Sessions int
	Portrait bool
	FontPath string
}

func (c *Card) Render() image.Image {
	W, H := 800, 480
	if c.Portrait {
		W, H = 480, 800
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(layout.ColorWhite)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, c.Title, W, 26, c.FontPath)
	
	// Layout grid of checkboxes
	cols := 4
	if c.Portrait {
		cols = 2
	}
	
	rows := (c.Sessions + cols - 1) / cols
	
	gridW := float64(W) - 80
	gridH := float64(H) - bodyY - 60
	
	if rows == 0 {
		rows = 1
	}
	
	g := layout.NewGrid(cols, rows, gridW, gridH, 40, bodyY+30)
	
	count := 1
	for r := 0; r < rows; r++ {
		for col := 0; col < cols; col++ {
			if count > c.Sessions {
				break
			}
			
			x, y := g.Cell(col, r)
			// padding inside cell
			cx := x + 10
			cy := y + 10
			cw := g.CellW - 20
			ch := g.CellH - 20
			
			// Draw hard drop shadow box
			layout.DrawHardDropShadowPanel(dc, cx, cy, cw, ch, 6, layout.ColorWhite, layout.ColorLightGray)
			
			// Draw text
			_ = layout.LoadFontFaceBold(dc, c.FontPath, 24)
			dc.SetColor(layout.ColorDarkGray)
			layout.DrawCenteredText(dc, fmt.Sprintf("Session %d", count), cx+cw/2, cy+ch/2)
			
			count++
		}
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "pomodoro",
		Usage: "Generate a Pomodoro / deep work tracker card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Title, "title", "Deep Work Tracker", "Title")
			fs.IntVar(&c.Sessions, "sessions", 8, "Number of sessions to track")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font")
			return c
		},
	}
}
