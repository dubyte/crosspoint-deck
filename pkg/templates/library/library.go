package library

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Name       string
	CardNumber string
	Branch     string
	Phone      string
	Portrait   bool
	FontPath   string
}

func (l *Card) Render() image.Image {
	var W, H int
	if l.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "Library Card", W, 22, l.FontPath)

	// Name
	_ = layout.LoadFontFace(dc, l.FontPath, 18)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(l.Name, float64(W)/2, bodyY+20, 0.5, 0.5)

	// Card number - large
	_ = layout.LoadFontFaceBold(dc, l.FontPath, 22)
	dc.DrawStringAnchored(l.CardNumber, float64(W)/2, bodyY+70, 0.5, 0.5)

	// Details with bold labels
	infoY := bodyY + 120
	lineH := 32.0
	for _, entry := range []struct{ label, value string }{
		{"Card #", l.CardNumber},
		{"Branch", l.Branch},
		{"Phone", l.Phone},
	} {
		if entry.value == "" {
			continue
		}
		_ = layout.LoadFontFaceBold(dc, l.FontPath, 16)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(entry.label)
		dc.DrawString(entry.label, float64(W)/2-w-6, infoY)

		_ = layout.LoadFontFace(dc, l.FontPath, 16)
		dc.DrawString(entry.value, float64(W)/2+6, infoY)
		infoY += lineH
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "library",
		Usage: "Generate a library card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Name, "name", "", "Patron name")
			fs.StringVar(&c.CardNumber, "card-number", "", "Library card number")
			fs.StringVar(&c.Branch, "branch", "", "Library branch")
			fs.StringVar(&c.Phone, "phone", "", "Library phone")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
