package convert

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

var conversions = []string{
	"1 in  =  2.54 cm",
	"1 ft  =  0.305 m",
	"1 mi  =  1.61 km",
	"1 oz  =  28.35 g",
	"1 lb  =  0.454 kg",
	"1 gal  =  3.79 L",
	"0°C  =  32°F",
	"100°C  =  212°F",
	"1 kg  =  2.20 lb",
	"1 km  =  0.62 mi",
	"1 L  =  0.26 gal",
	"1 m  =  3.28 ft",
	"1 cup  =  237 mL",
	"1 tbsp  =  15 mL",
	"1 tsp  =  5 mL",
	"1 mph  =  1.61 km/h",
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

	bodyY := layout.DrawReversedHeader(dc, "Conversions", W, 26, c.FontPath)

	half := len(conversions) / 2
	colW := float64(W) / 2
	startY := bodyY + 12
	lineH := 32.0

	_ = layout.LoadFontFace(dc, c.FontPath, 18)
	for i, conv := range conversions {
		x := 30.0
		col := i
		y := startY + float64(col)*lineH
		if i >= half {
			x = colW + 30
			y = startY + float64(col-half)*lineH
		}
		dc.SetColor(color.Black)
		dc.DrawString(conv, x, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "convert",
		Usage: "Generate a common conversions reference card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
