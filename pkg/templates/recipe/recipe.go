package recipe

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
	Title       string
	Ingredients []string
	Steps       []string
	Time        string
	Servings    string
	Portrait    bool
	FontPath    string
}

func (r *Card) Render() image.Image {
	var W, H int
	if r.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, r.Title, W, 22, r.FontPath)

	// Meta line
	metaY := bodyY + 14
	meta := ""
	if r.Time != "" {
		meta = "⏱ " + r.Time
	}
	if r.Servings != "" {
		if meta != "" {
			meta += "  ·  "
		}
		meta += "🍽 " + r.Servings
	}
	if meta != "" {
		_ = layout.LoadFontFace(dc, r.FontPath, 14)
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(meta, float64(W)/2, metaY, 0.5, 0.5)
		metaY += 26
	}

	sectionY := metaY + 8

	if r.Portrait {
		// Stacked
		sectionY = r.drawSection(dc, r.FontPath, "Ingredients", r.Ingredients, 30, sectionY, float64(W)-60, false)
		sectionY += 16
		_ = r.drawSection(dc, r.FontPath, "Steps", r.Steps, 30, sectionY, float64(W)-60, true)
	} else {
		// Side by side
		midX := float64(W) / 2
		r.drawSection(dc, r.FontPath, "Ingredients", r.Ingredients, 20, sectionY, midX-30, false)
		r.drawSection(dc, r.FontPath, "Steps", r.Steps, midX+10, sectionY, midX-30, true)
	}

	return dc.Image()
}

func (r *Card) drawSection(dc *gg.Context, fontPath, title string, items []string, x, y, maxW float64, numbered bool) float64 {
	_ = layout.LoadFontFaceBold(dc, fontPath, 16)
	dc.SetColor(color.Black)
	dc.DrawString(title, x, y)
	curY := y + 24

	_ = layout.LoadFontFace(dc, fontPath, 14)
	for i, item := range items {
		if curY > float64(dc.Height())-20 {
			break
		}
		prefix := "• "
		if numbered {
			prefix = string(rune('1'+i)) + ". "
		}
		dc.DrawString(prefix+item, x+10, curY)
		curY += 22
	}
	return curY
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "recipe",
		Usage: "Generate a recipe card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var ingRaw, stepsRaw string
			fs.StringVar(&c.Title, "title", "Recipe", "Recipe name")
			fs.StringVar(&ingRaw, "ingredients", "", "Comma-separated ingredients")
			fs.StringVar(&stepsRaw, "steps", "", "Comma-separated steps")
			fs.StringVar(&c.Time, "time", "", "Prep/cook time")
			fs.StringVar(&c.Servings, "servings", "", "Number of servings")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, ingRaw, stepsRaw}
		},
	}
}

type rawCard struct {
	*Card
	ingRaw   string
	stepsRaw string
}

func (rc *rawCard) Render() image.Image {
	if rc.ingRaw != "" && len(rc.Card.Ingredients) == 0 {
		rc.Card.Ingredients = split(rc.ingRaw)
	}
	if rc.stepsRaw != "" && len(rc.Card.Steps) == 0 {
		rc.Card.Steps = split(rc.stepsRaw)
	}
	return rc.Card.Render()
}

func split(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
