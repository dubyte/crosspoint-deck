package coffee

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
	Method   string
	Ratio    string
	Temp     string
	Time     string
	Steps    []string
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

	title := "Brew Guide"
	if c.Method != "" {
		title = c.Method
	}
	bodyY := layout.DrawReversedHeader(dc, title, W, 22, c.FontPath)

	// Brew params
	paramY := bodyY + 16
	params := []struct{ label, value string }{
		{"Ratio", c.Ratio},
		{"Temp", c.Temp},
		{"Time", c.Time},
	}
	for _, p := range params {
		if p.value == "" {
			continue
		}
		_ = layout.LoadFontFaceBold(dc, c.FontPath, 16)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(p.label)
		dc.DrawString(p.label, float64(W)/2-w-6, paramY)

		_ = layout.LoadFontFace(dc, c.FontPath, 16)
		dc.DrawString(p.value, float64(W)/2+6, paramY)
		paramY += 28
	}

	// Steps
	if len(c.Steps) > 0 {
		paramY += 12
		_ = layout.LoadFontFaceBold(dc, c.FontPath, 16)
		dc.DrawString("Steps", 30, paramY)
		paramY += 26

		_ = layout.LoadFontFace(dc, c.FontPath, 14)
		for i, step := range c.Steps {
			if paramY > float64(H)-20 {
				break
			}
			dc.DrawString(string(rune('1'+i))+". "+step, 40, paramY)
			paramY += 22
		}
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "coffee",
		Usage: "Generate a coffee brew guide card",
		New: func(fs *flag.FlagSet) card.Card {
			card := &Card{}
			var stepsRaw string
			fs.StringVar(&card.Method, "method", "Pour Over", "Brew method")
			fs.StringVar(&card.Ratio, "ratio", "1:16", "Coffee to water ratio")
			fs.StringVar(&card.Temp, "temp", "93°C", "Water temperature")
			fs.StringVar(&card.Time, "time", "3 min", "Brew time")
			fs.StringVar(&stepsRaw, "steps", "Rinse filter,Add grounds,Bloom 30s,Pour slowly,Enjoy", "Comma-separated steps")
			fs.BoolVar(&card.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&card.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{card, stepsRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Steps) == 0 {
		for _, s := range strings.Split(rc.raw, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				rc.Card.Steps = append(rc.Card.Steps, s)
			}
		}
	}
	return rc.Card.Render()
}
