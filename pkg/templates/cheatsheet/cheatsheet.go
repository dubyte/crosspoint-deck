package cheatsheet

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
	Title     string
	Shortcuts []Shortcut
	Portrait  bool
	FontPath  string
}

type Shortcut struct {
	Keys        string
	Description string
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

	startY := bodyY + 14
	lineH := 34.0
	colW := float64(W) / 2
	for i, s := range c.Shortcuts {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-20 {
			break
		}
		x := 50.0
		col := i
		if !c.Portrait && i >= (len(c.Shortcuts)+1)/2 {
			x = colW + 50
			col = i - (len(c.Shortcuts)+1)/2
			y = float64(startY) + float64(col)*lineH
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 18)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(s.Keys)
		dc.DrawString(s.Keys, x, y)

		_ = layout.LoadFontFace(dc, c.FontPath, 18)
		dc.DrawString(": "+s.Description, x+w, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "cheatsheet",
		Usage: "Generate a keyboard shortcuts cheat sheet",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{
				Shortcuts: []Shortcut{},
			}
			var shortcutsRaw string
			fs.StringVar(&c.Title, "title", "Shortcuts", "Cheat sheet title")
			fs.StringVar(&shortcutsRaw, "shortcuts", "", "Comma-separated key:description pairs")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &cheatsheetCard{c, shortcutsRaw}
		},
	}
}

type cheatsheetCard struct {
	*Card
	shortcutsRaw string
}

func (cc *cheatsheetCard) Render() image.Image {
	if cc.shortcutsRaw != "" && len(cc.Card.Shortcuts) == 0 {
		cc.Card.Shortcuts = ParseShortcuts(cc.shortcutsRaw)
	}
	return cc.Card.Render()
}

func ParseShortcuts(raw string) []Shortcut {
	var result []Shortcut
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, ":", 2)
		if len(parts) == 2 {
			result = append(result, Shortcut{
				Keys:        strings.TrimSpace(parts[0]),
				Description: strings.TrimSpace(parts[1]),
			})
		}
	}
	return result
}
