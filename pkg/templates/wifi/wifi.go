package wifi

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/dubyte/crosspoint-deck/pkg/qr"
	"github.com/fogleman/gg"
)

// WiFiCard renders a WiFi access card with QR code.
type WiFiCard struct {
	SSID       string
	Password   string
	Encryption string
	Portrait   bool
	FontPath   string
}

// Render produces a WiFi card with network info and QR code.
func (w *WiFiCard) Render() image.Image {
	var W, H int
	if w.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "WiFi Access", W, 22, w.FontPath)

	// Network info with bold labels
	infoY := bodyY + 16

	_ = layout.LoadFontFaceBold(dc, w.FontPath, 16)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("Network", float64(W)/2, infoY, 0.5, 0.5)

	_ = layout.LoadFontFace(dc, w.FontPath, 22)
	dc.DrawStringAnchored(w.SSID, float64(W)/2, infoY+28, 0.5, 0.5)

	infoY += 72
	_ = layout.LoadFontFaceBold(dc, w.FontPath, 16)
	dc.DrawStringAnchored("Password", float64(W)/2, infoY, 0.5, 0.5)

	_ = layout.LoadFontFace(dc, w.FontPath, 22)
	dc.DrawStringAnchored(w.Password, float64(W)/2, infoY+28, 0.5, 0.5)

	// QR code
	qrSize := 240
	if w.Portrait {
		qrSize = 280
	}
	qrImg, err := qr.GenerateWiFi(w.SSID, w.Password, w.Encryption, qrSize)
	if err == nil {
		x := (W - qrSize) / 2
		y := H - qrSize - 20
		dc.DrawImage(qrImg, x, y)
	}

	return dc.Image()
}

// Spec returns the card.Spec for wifi.
func Spec() card.Spec {
	return card.Spec{
		Name:  "wifi",
		Usage: "Generate a WiFi access card with QR code",
		New: func(fs *flag.FlagSet) card.Card {
			c := &WiFiCard{}
			fs.StringVar(&c.SSID, "ssid", "MyNetwork", "WiFi network name")
			fs.StringVar(&c.Password, "password", "", "WiFi password")
			fs.StringVar(&c.Encryption, "encryption", "WPA", "Encryption type (WPA, WEP, nopass)")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
