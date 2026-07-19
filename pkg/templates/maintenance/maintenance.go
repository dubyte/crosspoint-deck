package maintenance

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
	Year     string
	Tasks    []Task
	Portrait bool
	FontPath string
}

type Task struct {
	Name  string
	Month string
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

	title := "Maintenance"
	if m.Year != "" {
		title += " " + m.Year
	}
	bodyY := layout.DrawReversedHeader(dc, title, W, 22, m.FontPath)

	startY := bodyY + 16
	lineH := 28.0
	_ = layout.LoadFontFace(dc, m.FontPath, 16)
	for i, t := range m.Tasks {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-30 {
			break
		}

		_ = layout.LoadFontFaceBold(dc, m.FontPath, 16)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(t.Name)
		dc.DrawString(t.Name, float64(W)/2-w-8, y)

		_ = layout.LoadFontFace(dc, m.FontPath, 16)
		dc.DrawString(t.Month, float64(W)/2+8, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "maintenance",
		Usage: "Generate a home maintenance log card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var tasksRaw string
			fs.StringVar(&c.Year, "year", "", "Year (e.g. 2026)")
			fs.StringVar(&tasksRaw, "tasks", "HVAC filter:Jan,Smoke detector:Mar,AC check:Jun", "Comma-separated Task:Month pairs")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, tasksRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Tasks) == 0 {
		for _, pair := range strings.Split(rc.raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				rc.Card.Tasks = append(rc.Card.Tasks, Task{
					Name:  strings.TrimSpace(parts[0]),
					Month: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return rc.Card.Render()
}
