package periodic

import (
	"flag"
	"image"
	"image/color"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

var elements = []struct {
	num, sym, name string
}{
	{"1", "H", "Hydrogen"}, {"2", "He", "Helium"}, {"3", "Li", "Lithium"}, {"4", "Be", "Beryllium"},
	{"5", "B", "Boron"}, {"6", "C", "Carbon"}, {"7", "N", "Nitrogen"}, {"8", "O", "Oxygen"},
	{"9", "F", "Fluorine"}, {"10", "Ne", "Neon"}, {"11", "Na", "Sodium"}, {"12", "Mg", "Magnesium"},
	{"13", "Al", "Aluminum"}, {"14", "Si", "Silicon"}, {"15", "P", "Phosphorus"}, {"16", "S", "Sulfur"},
	{"17", "Cl", "Chlorine"}, {"18", "Ar", "Argon"}, {"19", "K", "Potassium"}, {"20", "Ca", "Calcium"},
	{"21", "Sc", "Scandium"}, {"22", "Ti", "Titanium"}, {"23", "V", "Vanadium"}, {"24", "Cr", "Chromium"},
	{"25", "Mn", "Manganese"}, {"26", "Fe", "Iron"}, {"27", "Co", "Cobalt"}, {"28", "Ni", "Nickel"},
	{"29", "Cu", "Copper"}, {"30", "Zn", "Zinc"}, {"31", "Ga", "Gallium"}, {"32", "Ge", "Germanium"},
	{"33", "As", "Arsenic"}, {"34", "Se", "Selenium"}, {"35", "Br", "Bromine"}, {"36", "Kr", "Krypton"},
}

type Card struct {
	Portrait bool
	FontPath string
}

func (c *Card) Render() image.Image {
	var W, H int
	if c.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(color.White)
	dc.Clear()

	bodyY := layout.DrawReversedHeader(dc, "Periodic Table", W, 26, c.FontPath)

	cols := 4
	if c.Portrait {
		cols = 2
	}
	perCol := (len(elements) + cols - 1) / cols
	colW := float64(W) / float64(cols)
	startY := bodyY + 12
	lineH := 24.0

	for i, entry := range elements {
		col := i / perCol
		row := i % perCol
		x := float64(col)*colW + 20
		y := startY + float64(row)*lineH
		if y > float64(H)-20 {
			break
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 14)
		dc.SetColor(color.Black)
		dc.DrawString(entry.num+" "+entry.sym, x, y)
		
		lw, _ := dc.MeasureString(entry.num+" "+entry.sym+" ")
		_ = layout.LoadFontFace(dc, c.FontPath, 14)
		dc.DrawString(entry.name, x+lw, y)
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "periodic",
		Usage: "Generate a periodic table of elements card",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{}
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return c
		},
	}
}
