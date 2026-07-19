package timezones

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
	Local    string
	Cities   []City
	Portrait bool
	FontPath string
}

type City struct {
	Name   string
	Offset string
}

func (t *Card) Render() image.Image {
	var W, H int
	if t.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "World Time", W, 22, t.FontPath)

	// Local timezone
	infoY := bodyY + 14
	if t.Local != "" {
		_ = layout.LoadFontFace(dc, t.FontPath, 14)
		dc.SetColor(color.Black)
		dc.DrawStringAnchored("Local: "+t.Local, float64(W)/2, infoY, 0.5, 0.5)
		infoY += 32
	}

	// Cities
	startY := infoY + 8
	lineH := 30.0
	for i, c := range t.Cities {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-30 {
			break
		}

		_ = layout.LoadFontFaceBold(dc, t.FontPath, 16)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(c.Name)
		dc.DrawString(c.Name, float64(W)/2-w-8, y)

		_ = layout.LoadFontFace(dc, t.FontPath, 16)
		dc.DrawString(c.Offset, float64(W)/2+8, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "timezones",
		Usage: "Generate a world time zones reference card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var citiesRaw string
			fs.StringVar(&c.Local, "local", "", "Your city/timezone (e.g. New York EST)")
			fs.StringVar(&citiesRaw, "cities", "Tokyo:+14h,London:+5h,Sydney:+18h", "Comma-separated City:Offset pairs")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, citiesRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Cities) == 0 {
		for _, pair := range strings.Split(rc.raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				rc.Card.Cities = append(rc.Card.Cities, City{
					Name:   strings.TrimSpace(parts[0]),
					Offset: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return rc.Card.Render()
}
