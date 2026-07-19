package loyalty

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
	Stores   []Store
	Portrait bool
	FontPath string
}

type Store struct {
	Name   string
	Member string
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

	bodyY := layout.DrawReversedHeader(dc, l.Title, W, 22, l.FontPath)

	startY := bodyY + 16
	lineH := 32.0
	for i, s := range l.Stores {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-30 {
			break
		}

		_ = layout.LoadFontFaceBold(dc, l.FontPath, 18)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(s.Name)
		dc.DrawString(s.Name, float64(W)/2-w-8, y)

		_ = layout.LoadFontFace(dc, l.FontPath, 18)
		dc.DrawString(s.Member, float64(W)/2+8, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "loyalty",
		Usage: "Generate a loyalty/membership card list",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var storesRaw string
			fs.StringVar(&c.Title, "title", "Loyalty Cards", "Card title")
			fs.StringVar(&storesRaw, "stores", "", "Comma-separated StoreName:MemberID pairs")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, storesRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Stores) == 0 {
		for _, pair := range strings.Split(rc.raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				rc.Card.Stores = append(rc.Card.Stores, Store{
					Name:   strings.TrimSpace(parts[0]),
					Member: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return rc.Card.Render()
}
