package emergency

import (
	"flag"
	"image"
	"image/color"
	"strings"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

// Card renders an emergency information card.
type Card struct {
	Country   string
	Contacts  []EmergencyContact
	BloodType string
	Allergies string
	Portrait  bool
	FontPath  string
}

// EmergencyContact holds a label and number.
type EmergencyContact struct {
	Label  string
	Number string
}

// Render produces an emergency info card.
func (e *Card) Render() image.Image {
	var W, H int
	if e.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	title := "Emergency Info"
	if e.Country != "" {
		title = "Emergency · " + e.Country
	}
	bodyY := layout.DrawReversedHeader(dc, title, W, 22, e.FontPath)

	// Contacts with bold labels
	startY := bodyY + 16
	lineH := 32.0
	_ = layout.LoadFontFaceBold(dc, e.FontPath, 18)

	for i, contact := range e.Contacts {
		y := float64(startY) + float64(i)*lineH
		if y > float64(H)-80 {
			break
		}
		_ = layout.LoadFontFaceBold(dc, e.FontPath, 18)
		dc.SetColor(color.Black)
		w, _ := dc.MeasureString(contact.Label)
		dc.DrawString(contact.Label, float64(W)/2-w-8, y)

		_ = layout.LoadFontFace(dc, e.FontPath, 18)
		dc.DrawString(contact.Number, float64(W)/2+8, y)
	}

	// Medical info
	if e.BloodType != "" || e.Allergies != "" {
		medY := float64(H) - 60
		_ = layout.LoadFontFace(dc, e.FontPath, 16)
		dc.SetColor(color.Black)
		if e.BloodType != "" {
			_ = layout.LoadFontFaceBold(dc, e.FontPath, 16)
			w, _ := dc.MeasureString("Blood")
			dc.DrawString("Blood", float64(W)/2-w-4, medY)

			_ = layout.LoadFontFace(dc, e.FontPath, 16)
			dc.DrawString(e.BloodType, float64(W)/2+4, medY)
		}
		if e.Allergies != "" {
			_ = layout.LoadFontFaceBold(dc, e.FontPath, 16)
			w, _ := dc.MeasureString("Allergies")
			dc.DrawString("Allergies", float64(W)/2-w-4, medY+25)

			_ = layout.LoadFontFace(dc, e.FontPath, 16)
			dc.DrawString(e.Allergies, float64(W)/2+4, medY+25)
		}
	}

	return dc.Image()
}

// Spec returns the card.Spec for emergency.
func Spec() card.Spec {
	return card.Spec{
		Name:  "emergency",
		Usage: "Generate an emergency contact card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			var contactsRaw string
			fs.StringVar(&c.Country, "country", "", "Country name")
			fs.StringVar(&contactsRaw, "contacts", "Police:911,Ambulance:911,Fire:911", "Comma-separated label:number pairs")
			fs.StringVar(&c.BloodType, "blood", "", "Blood type")
			fs.StringVar(&c.Allergies, "allergies", "", "Known allergies")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &emergencyCard{c, contactsRaw}
		},
	}
}

type emergencyCard struct {
	*Card
	contactsRaw string
}

func (ec *emergencyCard) Render() image.Image {
	if ec.contactsRaw != "" && len(ec.Card.Contacts) == 0 {
		for _, pair := range strings.Split(ec.contactsRaw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				ec.Card.Contacts = append(ec.Card.Contacts, EmergencyContact{
					Label:  strings.TrimSpace(parts[0]),
					Number: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	return ec.Card.Render()
}
