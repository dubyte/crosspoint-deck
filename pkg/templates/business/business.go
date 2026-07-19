package business

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
	Name     string
	Title    string
	Phone    string
	Email    string
	Website  string
	Portrait bool
	FontPath string
}

func (b *Card) Render() image.Image {
	var W, H int
	if b.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, b.Name, W, 26, b.FontPath)

	// Title
	_ = layout.LoadFontFace(dc, b.FontPath, 20)
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(b.Title, float64(W)/2, bodyY+14, 0.5, 0.5)

	// Contact info
	infoY := bodyY + 52
	lineH := 34.0
	for _, entry := range []struct{ label, value string }{
		{"Phone", b.Phone},
		{"Email", b.Email},
		{"Web", b.Website},
	} {
		if entry.value == "" {
			continue
		}
		_ = layout.LoadFontFaceBold(dc, b.FontPath, 18)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(entry.label)
		dc.DrawString(entry.label, float64(W)/2-w-8, infoY)
		_ = layout.LoadFontFace(dc, b.FontPath, 18)
		dc.DrawString(entry.value, float64(W)/2+8, infoY)
		infoY += lineH
	}

	// QR
	vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTITLE:%s\nTEL:%s\nEMAIL:%s\nURL:%s\nEND:VCARD",
		b.Name, b.Title, b.Phone, b.Email, b.Website)
	qrSize := 200
	if b.Portrait {
		qrSize = 260
	}
	qrImg, err := qr.Generate(vCard, qrSize)
	if err == nil {
		x := (W - qrSize) / 2
		y := H - qrSize - 30
		dc.DrawImage(qrImg, x, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "business",
		Usage: "Generate a business card with QR vCard",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.StringVar(&c.Name, "name", "", "Full name")
			fs.StringVar(&c.Title, "title", "", "Job title")
			fs.StringVar(&c.Phone, "phone", "", "Phone number")
			fs.StringVar(&c.Email, "email", "", "Email address")
			fs.StringVar(&c.Website, "website", "", "Website URL")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
