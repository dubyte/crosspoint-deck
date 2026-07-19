package shopping

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

func (s *Card) Render() image.Image {
	var W, H int
	if s.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, s.Title, W, 26, s.FontPath)

	_ = layout.LoadFontFace(dc, s.FontPath, 20)
	colW := float64(W) / 2
	startY := bodyY + 12
	perCol := 10
	lineH := 36.0
	if s.Portrait {
		perCol = 18
		lineH = 42
	}

	for i, item := range s.Items {
		y := startY + float64((i%perCol))*lineH
		x := 50.0
		if i >= perCol && !s.Portrait {
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
		Name:  "shopping",
		Usage: "Generate a shopping list card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var itemsRaw string
			fs.StringVar(&c.Title, "title", "Shopping List", "List title")
			fs.StringVar(&itemsRaw, "items", "Milk,Eggs,Bread,Butter", "Comma-separated items")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, itemsRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Items) == 0 {
		for _, s := range strings.Split(rc.raw, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				rc.Card.Items = append(rc.Card.Items, s)
			}
		}
	}
	return rc.Card.Render()
}
