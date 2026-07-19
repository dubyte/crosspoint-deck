package workout

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
	Exercises []Exercise
	Rounds    string
	Rest      string
	Portrait  bool
	FontPath  string
}

type Exercise struct {
	Name string
	Reps string
}

func (w *Card) Render() image.Image {
	var W, H int
	if w.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, w.Title, W, 26, w.FontPath)

	// Rounds/rest
	metaY := bodyY + 16
	meta := ""
	if w.Rounds != "" {
		meta = "Rounds: " + w.Rounds
	}
	if w.Rest != "" {
		if meta != "" {
			meta += "  ·  "
		}
		meta += "Rest: " + w.Rest + "s"
	}
	if meta != "" {
		_ = layout.LoadFontFace(dc, w.FontPath, 18)
		dc.SetColor(color.Black)
		dc.DrawStringAnchored(meta, float64(W)/2, metaY, 0.5, 0.5)
		metaY += 36
	}

	// Exercises
	startY := metaY + 10
	lineH := 40.0
	for i, ex := range w.Exercises {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-30 {
			break
		}

		_ = layout.LoadFontFace(dc, w.FontPath, 20)
		dc.SetColor(color.Black)
		num := string(rune('1'+i)) + ". "
		nw, _ := dc.MeasureString(num)
		dc.DrawString(num, 50, y)

		_ = layout.LoadFontFaceBold(dc, w.FontPath, 20)
		dc.DrawString(ex.Name, 50+nw, y)

		_ = layout.LoadFontFace(dc, w.FontPath, 20)
		ew, _ := dc.MeasureString(ex.Name)
		dc.DrawString(" · "+ex.Reps, 50+nw+ew, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "workout",
		Usage: "Generate a bodyweight workout card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var exRaw string
			fs.StringVar(&c.Title, "title", "Workout", "Workout title")
			fs.StringVar(&exRaw, "exercises", "Push-ups:10,Squats:15,Plank:30s", "Comma-separated Exercise:Reps pairs")
			fs.StringVar(&c.Rounds, "rounds", "3", "Number of rounds")
			fs.StringVar(&c.Rest, "rest", "60", "Rest between rounds (seconds)")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &rawCard{c, exRaw}
		},
	}
}

type rawCard struct {
	*Card
	raw string
}

func (rc *rawCard) Render() image.Image {
	if rc.raw != "" && len(rc.Card.Exercises) == 0 {
		for _, pair := range strings.Split(rc.raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				rc.Card.Exercises = append(rc.Card.Exercises, Exercise{
					Name: strings.TrimSpace(parts[0]),
					Reps: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return rc.Card.Render()
}
