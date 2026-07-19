package plant

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Plant    string
	Water    string
	Light    string
	Humidity string
	Food     string
	Notes    string
	Portrait bool
	FontPath string
}

func (p *Card) Render() image.Image {
	var W, H int
	if p.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, p.Plant, W, 22, p.FontPath)

	infoY := bodyY + 16
	lineH := 32.0
	for _, entry := range []struct{ label, value string }{
		{"Water", p.Water},
		{"Light", p.Light},
		{"Humidity", p.Humidity},
		{"Food", p.Food},
	} {
		if entry.value == "" {
			continue
		}
		_ = layout.LoadFontFaceBold(dc, p.FontPath, 18)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(entry.label)
		dc.DrawString(entry.label, float64(W)/2-w-6, infoY)

		_ = layout.LoadFontFace(dc, p.FontPath, 18)
		dc.DrawString(entry.value, float64(W)/2+6, infoY)
		infoY += lineH
	}

	if p.Notes != "" {
		infoY += 10
		_ = layout.LoadFontFace(dc, p.FontPath, 14)
		dc.SetColor(color.Black)
		layout.DrawWrappedText(dc, p.Notes, 30, infoY, float64(W)-60, 20)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "plant",
		Usage: "Generate a plant care guide card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Plant, "plant", "", "Plant name")
			fs.StringVar(&c.Water, "water", "", "Watering instructions")
			fs.StringVar(&c.Light, "light", "", "Light requirements")
			fs.StringVar(&c.Humidity, "humidity", "", "Humidity needs")
			fs.StringVar(&c.Food, "food", "", "Fertilizer/feeding instructions")
			fs.StringVar(&c.Notes, "notes", "", "Additional notes")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
