package chore

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
	Chores   []string
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

	bodyY := layout.DrawReversedHeader(dc, c.Title, W, 26, c.FontPath)

	_ = layout.LoadFontFace(dc, c.FontPath, 20)
	colW := float64(W) / 2
	startY := bodyY + 12
	lineH := 38.0
	perCol := 10
	if c.Portrait {
		perCol = 18
	}

	for i, chore := range c.Chores {
		y := startY + float64((i%perCol))*lineH
		x := 40.0
		if i >= perCol && !c.Portrait {
			x = colW + 40
			y = startY + float64((i-perCol))*lineH
		}
		if y > float64(H)-20 {
			break
		}
		dc.SetColor(color.Black)
		dc.DrawString("[ ] "+chore, x, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "chore",
		Usage: "Generate a chore chart checklist",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var choresRaw string
			fs.StringVar(&c.Title, "title", "Chores", "Chart title")
			fs.StringVar(&choresRaw, "chores", "Dishes,Laundry,Vacuum,Trash", "Comma-separated chores")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, choresRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Chores) == 0 {
		for _, s := range strings.Split(rc.raw, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				rc.Card.Chores = append(rc.Card.Chores, s)
			}
		}
	}
	return rc.Card.Render()
}
