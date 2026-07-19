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
	Country     string
	Contacts    []EmergencyContact
	BloodType   string
	Allergies   string
	Portrait    bool
	FontPath    string
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

	_ = layout.LoadFontFace(dc, e.FontPath, 22)

	// Title
	dc.SetColor(color.Black)
	title := "Emergency Info"
	if e.Country != "" {
		title = "Emergency: " + e.Country
	}
	dc.DrawStringAnchored(title, float64(W)/2, 35, 0.5, 0.5)

	// Contacts
	_ = layout.LoadFontFace(dc, e.FontPath, 18)
	startY := 75
	for i, contact := range e.Contacts {
		y := float64(startY + i*35)
		if y > float64(H)-80 {
			break
		}
		dc.DrawString(contact.Label+": "+contact.Number, 30, y)
	}

	// Medical info
	if e.BloodType != "" || e.Allergies != "" {
		medY := float64(H) - 60
		_ = layout.LoadFontFace(dc, e.FontPath, 16)
		if e.BloodType != "" {
			dc.DrawString("Blood: "+e.BloodType, 30, medY)
		}
		if e.Allergies != "" {
			dc.DrawString("Allergies: "+e.Allergies, 30, medY+25)
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
