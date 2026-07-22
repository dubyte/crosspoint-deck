package transit

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
	Outbound  string
	Inbound   string
	Portrait  bool
	FontPath  string
}

func (c *Card) Render() image.Image {
	W, H := 800, 480
	if c.Portrait {
		W, H = 480, 800
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(layout.ColorWhite)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, c.Title, W, 26, c.FontPath)
	startY := bodyY + 20

	if c.Portrait {
		drawSchedule(dc, "Outbound", c.Outbound, c.FontPath, 40, startY, 400)
		drawSchedule(dc, "Inbound", c.Inbound, c.FontPath, 40, startY+350, 400)
	} else {
		drawSchedule(dc, "Outbound", c.Outbound, c.FontPath, 40, startY, 340)
		drawSchedule(dc, "Inbound", c.Inbound, c.FontPath, 420, startY, 340)
		// divider
		dc.SetColor(layout.ColorDarkGray)
		dc.SetLineWidth(2)
		dc.DrawLine(400, startY, 400, float64(H)-40)
		dc.Stroke()
	}

	return dc.Image()
}

func drawSchedule(dc *gg.Context, title, raw, font string, x, y, w float64) {
	_ = layout.LoadFontFaceBold(dc, font, 22)
	
	// Draw category header with drop shadow effect
	layout.DrawHardDropShadowPanel(dc, x, y, w, 36, 4, layout.ColorLightGray, layout.ColorDarkGray)
	
	dc.SetColor(layout.ColorBlack)
	dc.DrawStringAnchored(title, x+10, y+18, 0, 0.35)
	
	y += 50
	_ = layout.LoadFontFace(dc, font, 20)
	
	times := strings.Split(raw, ",")
	lineH := 36.0
	for _, t := range times {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		
		// Split into time and AM/PM or destination if separated by space
		parts := strings.SplitN(t, " ", 2)
		mainText := parts[0]
		subText := ""
		if len(parts) > 1 {
			subText = parts[1]
		}
		
		dc.SetColor(layout.ColorBlack)
		dc.DrawString(mainText, x+10, y+lineH/2)
		
		if subText != "" {
			dc.SetColor(layout.ColorDarkGray)
			tw, _ := dc.MeasureString(mainText)
			dc.DrawString(subText, x+10+tw+10, y+lineH/2)
		}
		
		// subtle bottom border
		dc.SetColor(color.RGBA{230, 230, 230, 255})
		dc.SetLineWidth(1)
		dc.DrawLine(x, y+lineH-4, x+w, y+lineH-4)
		dc.Stroke()
		
		y += lineH
	}
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "transit",
		Usage: "Generate a transit schedule card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Title, "title", "Transit Schedule", "Title")
			fs.StringVar(&c.Outbound, "outbound", "", "Comma-separated outbound times (e.g., '07:15 AM,07:45 AM,08:20 AM')")
			fs.StringVar(&c.Inbound, "inbound", "", "Comma-separated inbound times")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font")
			return c
		},
	}
}
