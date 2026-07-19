package owner

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Name     string
	Email    string
	Portrait bool
	FontPath string
}

func (o *Card) Render() image.Image {
	var W, H int
	if o.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "Owner", W, 26, o.FontPath)

	// Label
	_ = layout.LoadFontFace(dc, o.FontPath, 20)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("This e-reader belongs to", float64(W)/2, bodyY+30, 0.5, 0.5)

	// Name - large and bold
	_ = layout.LoadFontFaceBold(dc, o.FontPath, 36)
	dc.DrawStringAnchored(o.Name, float64(W)/2, bodyY+90, 0.5, 0.5)

	// Email
	if o.Email != "" {
		_ = layout.LoadFontFace(dc, o.FontPath, 22)
		dc.DrawStringAnchored(o.Email, float64(W)/2, bodyY+140, 0.5, 0.5)
	}

	// Optional: subtle "If found, please contact" at bottom
	_ = layout.LoadFontFace(dc, o.FontPath, 14)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("If found, please contact owner at above email",
		float64(W)/2, float64(H)-30, 0.5, 0.5)

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "owner",
		Usage: "Generate an owner identification card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Name, "name", "", "Owner name")
			fs.StringVar(&c.Email, "email", "", "Owner email")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
