package kitchenmath

import (
	"flag"
	"image"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
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

	bodyY := layout.DrawReversedHeader(dc, "Kitchen Math & Safe Temps", W, 26, c.FontPath)
	
	startY := bodyY + 20
	
	// Layout
	if c.Portrait {
		drawConversions(dc, c.FontPath, 30, startY, 420)
		drawTemps(dc, c.FontPath, 30, startY+350, 420)
	} else {
		drawConversions(dc, c.FontPath, 40, startY, 340)
		drawTemps(dc, c.FontPath, 420, startY, 340)
		// divider
		dc.SetColor(layout.ColorDarkGray)
		dc.SetLineWidth(2)
		dc.DrawLine(400, startY, 400, float64(H)-40)
		dc.Stroke()
	}

	return dc.Image()
}

func drawConversions(dc *gg.Context, font string, x, y, w float64) {
	_ = layout.LoadFontFaceBold(dc, font, 20)
	dc.SetColor(layout.ColorBlack)
	dc.DrawString("Volume Conversions", x, y)
	y += 30

	_ = layout.LoadFontFace(dc, font, 18)
	items := []struct{ k, v string }{
		{"1 Tablespoon", "3 teaspoons"},
		{"1/4 Cup", "4 Tablespoons"},
		{"1/3 Cup", "5 Tbsp + 1 tsp"},
		{"1/2 Cup", "8 Tablespoons"},
		{"1 Cup", "16 Tbsp | 237 mL"},
		{"1 Pint", "2 Cups | 473 mL"},
		{"1 Quart", "2 Pints | 946 mL"},
		{"1 Gallon", "4 Quarts | 3.8 L"},
	}
	for _, item := range items {
		layout.DrawLeaderDotsText(dc, item.k, item.v, x, y, w)
		y += 30
	}
}

func drawTemps(dc *gg.Context, font string, x, y, w float64) {
	_ = layout.LoadFontFaceBold(dc, font, 20)
	dc.SetColor(layout.ColorBlack)
	dc.DrawString("Safe Internal Temps", x, y)
	y += 30

	_ = layout.LoadFontFace(dc, font, 18)
	items := []struct{ k, v string }{
		{"Poultry (Chicken, Turkey)", "165°F (74°C)"},
		{"Ground Meats (Beef, Pork)", "160°F (71°C)"},
		{"Whole Meats (Steak, Pork)", "145°F (63°C)"},
		{"Fish & Shellfish", "145°F (63°C)"},
		{"Leftovers & Casseroles", "165°F (74°C)"},
		{"Eggs", "160°F (71°C)"},
	}
	for _, item := range items {
		layout.DrawLeaderDotsText(dc, item.k, item.v, x, y, w)
		y += 30
	}
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "kitchen-math",
		Usage: "Generate a kitchen conversions and safe temperatures card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font")
			return c
		},
	}
}
