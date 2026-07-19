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

type Card struct {
	Title    string
	Items    []string
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

	bodyY := layout.DrawReversedHeader(dc, p.Title, W, 26, p.FontPath)

	_ = layout.LoadFontFace(dc, p.FontPath, 20)
	colW := float64(W) / 2
	startY := bodyY + 12
	perCol := 10
	lineH := 36.0
	if p.Portrait {
		perCol = 18
		lineH = 42
	}

	for i, item := range p.Items {
		y := startY + float64((i%perCol))*lineH
		x := 50.0
		if i >= perCol && !p.Portrait {
			x = colW + 50
			y = startY + float64((i-perCol))*lineH
		}
		if y > float64(H)-20 {
			break
		}
		dc.SetColor(color.Black)
		dc.DrawString("[ ] "+item, x, y)
	}

	return dc.Image()
}

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
