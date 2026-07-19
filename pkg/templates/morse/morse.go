package morse

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

var codes = []struct{ letter, code string }{
	{"A", "·−"}, {"B", "−···"}, {"C", "−·−·"}, {"D", "−··"},
	{"E", "·"}, {"F", "··−·"}, {"G", "−−·"}, {"H", "····"},
	{"I", "··"}, {"J", "·−−−"}, {"K", "−·−"}, {"L", "·−··"},
	{"M", "−−"}, {"N", "−·"}, {"O", "−−−"}, {"P", "·−−·"},
	{"Q", "−−·−"}, {"R", "·−·"}, {"S", "···"}, {"T", "−"},
	{"U", "··−"}, {"V", "···−"}, {"W", "·−−"}, {"X", "−··−"},
	{"Y", "−·−−"}, {"Z", "−−··"},
	{"0", "−−−−−"}, {"1", "·−−−−"}, {"2", "··−−−"}, {"3", "···−−"},
	{"4", "····−"}, {"5", "·····"}, {"6", "−····"}, {"7", "−−···"},
	{"8", "−−−··"}, {"9", "−−−−·"},
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

	bodyY := layout.DrawReversedHeader(dc, "Morse Code", W, 22, c.FontPath)

	half := len(codes) / 2
	colW := float64(W) / 2
	startY := bodyY + 12
	lineH := 28.0

	for i, entry := range codes {
		x := 30.0
		col := i
		y := startY + float64(col)*lineH
		if i >= half {
			x = colW + 30
			y = startY + float64(col-half)*lineH
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 16)
		dc.SetColor(color.Black)
		lw, _ := dc.MeasureString(entry.letter)
		dc.DrawString(entry.letter, x, y)

		_ = layout.LoadFontFace(dc, c.FontPath, 16)
		dc.DrawString("  "+entry.code, x+lw, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "morse",
		Usage: "Generate a Morse code reference card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
