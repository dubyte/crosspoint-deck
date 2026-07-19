package packing

import (
	"flag"
	"image"
	"image/color"
	"strings"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

// Card renders a packing checklist.
type Card struct {
	Title    string
	Items    []string
	Portrait bool
	FontPath string
}

// Render produces a packing checklist card.
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

	_ = layout.LoadFontFace(dc, p.FontPath, 22)

	// Title
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(p.Title, float64(W)/2, 35, 0.5, 0.5)

	// Items list
	_ = layout.LoadFontFace(dc, p.FontPath, 16)
	colW := float64(W) / 2
	startY := 70
	for i, item := range p.Items {
		y := float64(startY + (i%12)*30)
		x := 30.0
		if i >= 12 && !p.Portrait {
			x = colW + 30
			y = float64(startY + (i-12)*30)
		}
		if y > float64(H)-30 {
			break
		}
		dc.DrawString("[ ] "+item, x, y)
	}

	return dc.Image()
}

// Spec returns the card.Spec for packing.
func Spec() card.Spec {
	return card.Spec{
		Name:  "packing",
		Usage: "Generate a packing checklist card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var itemsRaw string
			fs.StringVar(&c.Title, "title", "Packing List", "List title")
			fs.StringVar(&itemsRaw, "items", "Passport,Tickets,Phone charger,Toiletries", "Comma-separated items")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &packingCard{c, itemsRaw}
		},
	}
}

type packingCard struct {
	*Card
	itemsRaw string
}

func (pc *packingCard) Render() image.Image {
	if pc.itemsRaw != "" && len(pc.Card.Items) == 0 {
		pc.Card.Items = strings.Split(pc.itemsRaw, ",")
		for i := range pc.Card.Items {
			pc.Card.Items[i] = strings.TrimSpace(pc.Card.Items[i])
		}
	}
	return pc.Card.Render()
}
