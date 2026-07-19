package meeting

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Room     string
	Events   []string
	Portrait bool
	FontPath string
}

func (m *Card) Render() image.Image {
	var W, H int
	if m.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, m.Room, W, 26, m.FontPath)

	_ = layout.LoadFontFace(dc, m.FontPath, 22)
	startY := bodyY + 20
	for i, event := range m.Events {
		y := startY + float64(i)*42
		if y > float64(H)-30 {
			break
		}
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(event, float64(W)/2, y, 0.5, 0.5)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "meeting",
		Usage: "Generate a meeting room schedule card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Room, "room", "Conference A", "Room name")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
