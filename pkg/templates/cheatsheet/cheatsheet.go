package cheatsheet

import (
	"flag"
	"image"
	"strings"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/layout"
	"github.com/fogleman/gg"
)

type Card struct {
	Title      string
	Categories []Category
	Portrait   bool
	FontPath   string
}

type Category struct {
	Name      string
	Shortcuts []Shortcut
}

type Shortcut struct {
	Keys        string
	Description string
}

func (c *Card) Render() image.Image {
	var W, H int
	if c.Portrait {
		W, H = 480, 800
	} else {
		W, H = 800, 480
	}

	dc := gg.NewContext(W, H)
	dc.SetColor(layout.ColorWhite)
	dc.Clear()

	// Asymmetric "Hero" Zone (Top Banner)
	dc.SetColor(layout.ColorBlack)
	dc.DrawRectangle(0, 0, float64(W), 80)
	dc.Fill()
	
	_ = layout.LoadFontFaceBold(dc, c.FontPath, 36)
	dc.SetColor(layout.ColorWhite)
	dc.DrawStringAnchored(c.Title, 40, 40, 0, 0.35)

	// Columns setup
	var cols int
	var colW float64
	if c.Portrait {
		cols = 1
		colW = float64(W) - 60
	} else {
		cols = 3
		colW = (float64(W) - 100) / float64(cols)
	}

	startX := 30.0
	startY := 110.0
	lineH := 30.0
	gutter := 20.0

	currX := startX
	currY := startY
	currCol := 0

	for _, cat := range c.Categories {
		// Category Header
		if currY+lineH*2 > float64(H)-20 {
			// Move to next column
			currCol++
			if currCol >= cols {
				break // out of space
			}
			currX = startX + float64(currCol)*(colW+gutter)
			currY = startY
		}

		_ = layout.LoadFontFaceBold(dc, c.FontPath, 20)
		dc.SetColor(layout.ColorBlack)
		dc.DrawString(cat.Name, currX, currY)
		
		// underline category
		dc.SetLineWidth(2)
		dc.DrawLine(currX, currY+8, currX+colW, currY+8)
		dc.Stroke()
		
		currY += lineH

		_ = layout.LoadFontFace(dc, c.FontPath, 16)
		for _, s := range cat.Shortcuts {
			if currY+lineH > float64(H)-20 {
				currCol++
				if currCol >= cols {
					break
				}
				currX = startX + float64(currCol)*(colW+gutter)
				currY = startY
			}
			
			// Draw Action -> Shortcut
			// Action is Description, Shortcut is Keys
			layout.DrawLeaderDotsText(dc, s.Description, s.Keys, currX, currY, colW)
			currY += lineH
		}
		currY += 10 // extra space after category
	}

	return dc.Image()
}

func Spec() card.Spec {
	return card.Spec{
		Name:  "cheatsheet",
		Usage: "Generate a categorized keyboard shortcuts cheat sheet",
		New: func(fs *flag.FlagSet) card.Card {
			c := &Card{
				Categories: []Category{},
			}
			var shortcutsRaw string
			fs.StringVar(&c.Title, "title", "Shortcuts", "Cheat sheet title")
			fs.StringVar(&shortcutsRaw, "shortcuts", "", "Format: Category|desc:key,desc:key;Category2|desc:key")
			fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
			fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
			return &cheatsheetCard{c, shortcutsRaw}
		},
	}
}

type cheatsheetCard struct {
	*Card
	shortcutsRaw string
}

func (cc *cheatsheetCard) Render() image.Image {
	if cc.shortcutsRaw != "" && len(cc.Card.Categories) == 0 {
		cc.Card.Categories = ParseCategories(cc.shortcutsRaw)
	}
	return cc.Card.Render()
}

func ParseCategories(raw string) []Category {
	var result []Category
	// Fallback to flat comma-separated if no categories are detected
	if !strings.Contains(raw, "|") && !strings.Contains(raw, ";") {
		cat := Category{Name: "General"}
		for _, item := range strings.Split(raw, ",") {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) == 2 {
				// Old format was key:desc, new is desc:key but let's stick to key:desc if we just split by comma
				cat.Shortcuts = append(cat.Shortcuts, Shortcut{
					Keys:        strings.TrimSpace(parts[0]),
					Description: strings.TrimSpace(parts[1]),
				})
			}
		}
		result = append(result, cat)
		return result
	}

	// Format: Cat|desc:key,desc:key;Cat2|desc:key
	for _, catRaw := range strings.Split(raw, ";") {
		catRaw = strings.TrimSpace(catRaw)
		if catRaw == "" {
			continue
		}
		parts := strings.SplitN(catRaw, "|", 2)
		if len(parts) != 2 {
			continue
		}
		cat := Category{Name: strings.TrimSpace(parts[0])}
		for _, item := range strings.Split(parts[1], ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			kv := strings.SplitN(item, ":", 2)
			if len(kv) == 2 {
				cat.Shortcuts = append(cat.Shortcuts, Shortcut{
					Description: strings.TrimSpace(kv[0]),
					Keys:        strings.TrimSpace(kv[1]),
				})
			}
		}
		result = append(result, cat)
	}
	return result
}
