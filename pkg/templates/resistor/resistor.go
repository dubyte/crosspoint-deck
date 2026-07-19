package resistor

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

// pony tail: color swatches dropped — B&W e-ink can't represent colors meaningfully.
// Add pattern-based swatches if a 4-gray palette mapping each resistor color is defined.

var bands = []struct{ digit, name string }{
	{"0", "Black"}, {"1", "Brown"}, {"2", "Red"}, {"3", "Orange"},
	{"4", "Yellow"}, {"5", "Green"}, {"6", "Blue"}, {"7", "Violet"},
	{"8", "Grey"}, {"9", "White"},
}

type Card struct {
	Portrait bool
	FontPath string
}

func (c *Card) Render() image.Image {
	var W, H int
	if c.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "Resistor Codes", W, 22, c.FontPath)

	half := 5
	colW := float64(W) / 2
	startY := bodyY + 16
	lineH := 34.0

	for i, b := range bands {
		x := 60.0
		col := i
		y := startY + float64(col)*lineH
		if i >= half {
			x = colW + 60
			y = startY + float64(col-half)*lineH
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 18)
		dc.SetColor(color.Black)
		lw, _ := dc.MeasureString(b.digit)
		dc.DrawString(b.digit, x, y)

		_ = layout.LoadFontFace(dc, c.FontPath, 18)
		dc.DrawString("  "+b.name, x+lw, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "resistor",
		Usage: "Generate a resistor color code reference card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
