package stretch

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
	Stretches []Stretch
	Portrait  bool
	FontPath  string
}

type Stretch struct {
	Name     string
	Duration string
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

	startY := bodyY + 24
	lineH := 40.0
	for i, st := range s.Stretches {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-30 {
			break
		}

		_ = layout.LoadFontFace(dc, s.FontPath, 20)
		dc.SetColor(color.Black)
		num := string(rune('1'+i)) + ". "
		nw, _ := dc.MeasureString(num)
		dc.DrawString(num, 30, y)

		_ = layout.LoadFontFaceBold(dc, s.FontPath, 20)
		dc.DrawString(st.Name, 30+nw, y)

		_ = layout.LoadFontFace(dc, s.FontPath, 20)
		ew, _ := dc.MeasureString(st.Name)
		dc.DrawString(" · Hold: "+st.Duration, 30+nw+ew, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "stretch",
		Usage: "Generate a stretching routine card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var raw string
			fs.StringVar(&c.Title, "title", "Stretching", "Routine title")
			fs.StringVar(&raw, "stretches", "Neck roll:30s,Shoulder stretch:30s", "Comma-separated Name:Duration pairs")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, raw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Stretches) == 0 {
		for _, pair := range strings.Split(rc.raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				rc.Card.Stretches = append(rc.Card.Stretches, Stretch{
					Name:     strings.TrimSpace(parts[0]),
					Duration: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return rc.Card.Render()
}
