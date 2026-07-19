package wifi

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/dubyte/crosspoint-deck/pkg/qr"
	"github.com/fogleman/gg"
)

type Card struct {
	SSID     string
	Password string
	Portrait bool
	FontPath string
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

	bodyY := layout.DrawReversedHeader(dc, "WiFi", W, 26, w.FontPath)

	infoY := bodyY + 30
	lineH := 42.0

	_ = layout.LoadFontFaceBold(dc, w.FontPath, 20)
	dc.SetColor(color.Black)
	nw, _ := dc.MeasureString("SSID")
	dc.DrawString("SSID", float64(W)/2-nw-10, infoY)

	_ = layout.LoadFontFace(dc, w.FontPath, 26)
	dc.DrawString(w.SSID, float64(W)/2+10, infoY)
	infoY += lineH

	_ = layout.LoadFontFaceBold(dc, w.FontPath, 20)
	dc.SetColor(color.Black)
	pw, _ := dc.MeasureString("PWD")
	dc.DrawString("PWD", float64(W)/2-pw-10, infoY)

	_ = layout.LoadFontFace(dc, w.FontPath, 26)
	dc.DrawString(w.Password, float64(W)/2+10, infoY)

	// QR
	qrText := fmt.Sprintf("WIFI:T:WPA;S:%s;P:%s;;", w.SSID, w.Password)
	qrSize := 240
	if w.Portrait {
		qrSize = 280
	}
	qrImg, err := qr.Generate(qrText, qrSize)
	if err == nil {
		x := (W - qrSize) / 2
		y := H - qrSize - 30
		dc.DrawImage(qrImg, x, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "wifi",
		Usage: "Generate a WiFi access card with QR code",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.SSID, "ssid", "", "WiFi network name")
			fs.StringVar(&c.Password, "password", "", "WiFi password")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
