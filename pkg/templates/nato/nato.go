package nato

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

var alphabet = []struct{ letter, word string }{
	{"A", "Alpha"}, {"B", "Bravo"}, {"C", "Charlie"}, {"D", "Delta"},
	{"E", "Echo"}, {"F", "Foxtrot"}, {"G", "Golf"}, {"H", "Hotel"},
	{"I", "India"}, {"J", "Juliet"}, {"K", "Kilo"}, {"L", "Lima"},
	{"M", "Mike"}, {"N", "November"}, {"O", "Oscar"}, {"P", "Papa"},
	{"Q", "Quebec"}, {"R", "Romeo"}, {"S", "Sierra"}, {"T", "Tango"},
	{"U", "Uniform"}, {"V", "Victor"}, {"W", "Whiskey"}, {"X", "X-ray"},
	{"Y", "Yankee"}, {"Z", "Zulu"},
}

type Card struct {
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

	bodyY := layout.DrawReversedHeader(dc, "NATO Phonetic", W, 22, c.FontPath)

	half := 13
	colW := float64(W) / 2
	startY := bodyY + 12
	lineH := 28.0

	for i, entry := range alphabet {
		x := 30.0
		col := i
		y := startY + float64(col)*lineH
		if i >= half {
			x = colW + 30
			y = startY + float64(col-half)*lineH
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 16)
		dc.SetColor(color.Black)
		lw, _ := dc.MeasureString(entry.letter)
		dc.DrawString(entry.letter, x, y)

		_ = layout.LoadFontFace(dc, c.FontPath, 16)
		dc.DrawString("  "+entry.word, x+lw, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "nato",
		Usage: "Generate a NATO phonetic alphabet reference card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
