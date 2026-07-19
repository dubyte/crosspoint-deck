# Phase 3 Implementation Prompt — The Laminated Life

Copy this entire document into a new AI session to implement Phase 3 cards.

## Context

You are working on `crosspoint-deck`, a Go CLI that generates 800×480 BMP cards for the XTEink X4 e-ink reader (SSD1677 controller, 4-level grayscale). The project renders cards from Go templates using the `gg` graphics library.

**Project root:** `/home/dubyte/Documents/Workspace/vibe/crosspoint-deck`

**Build:** `go build -o deck ./cmd/deck`

**Generate all:** `mage All`

**Key files:**

- `pkg/card/card.go` — `Card` interface, `Spec`, `Factory`
- `pkg/layout/layout.go` — `LoadFontFace`, `LoadFontFaceBold`, `DrawReversedHeader`, `Grid`
- `pkg/bmp/encoder.go` — 24-bit BMP encoder (preserves grayscale)
- `pkg/render/render.go` — `ToFile()` with dimensions validation
- `cmd/deck/main.go` — CLI entry, `registry` slice
- `magefile.go` — Build tasks

## Design System

Every card follows this unified design:

```
┌──────────────────────────────────┐
│▓▓▓▓▓▓▓▓▓▓ TITLE ▓▓▓▓▓▓▓▓▓▓▓▓▓▓│  ← reversed header: black bar, white text
│──────────────────────────────────│  ← 2px divider line
│                                  │
│  Label    Value                  │  ← bold label, regular value
│  Label    Value                  │
│                                  │
└──────────────────────────────────┘
```

- **`DrawReversedHeader(dc, title, W, fontSize, fontPath)`** returns `bodyY` — use it.
- **`LoadFontFaceBold(dc, fontPath, size)`** for labels/keys.
- **`LoadFontFace(dc, fontPath, size)`** for values/descriptions.
- Font sizes: title 22pt, body 14–16pt, small 12pt.
- All cards have `--portrait` (480×800) and `--font` flags.
- Pattern: landscape = 800×480, portrait = 480×800.

## How to Add a Card

### 1. Create template file

`pkg/templates/<name>/<name>.go`

Minimal template:

```go
package <name>

import (
 "flag"
 "image"
 "image/color"

 "github.com/dubyte/crosspoint-deck/pkg/card"
 "github.com/dubyte/crosspoint-deck/pkg/layout"
 "github.com/fogleman/gg"
)

type Card struct {
 Portrait bool
 FontPath string
 // Add template-specific fields
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

 bodyY := layout.DrawReversedHeader(dc, "Card Title", W, 22, c.FontPath)

 // Your layout here, starting at bodyY

 return dc.Image()
}

func Spec() card.Spec {
 return card.Spec{
  Name:  "name",
  Usage: "Short description",
  New: func(fs *flag.FlagSet) card.Card {
   c := &Card{}
   fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
   fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
   // Add your flags here
   return c
  },
 }
}
```

If the card needs a comma-separated input that gets parsed lazily, use the wrapper pattern (see `packing.go` or `cheatsheet.go`):

```go
type rawCard struct {
 *Card
 rawData string
}

func (rc *rawCard) Render() image.Image {
 if rc.rawData != "" && len(rc.Card.Items) == 0 {
  // parse rc.rawData into rc.Card.Items
 }
 return rc.Card.Render()
}
```

Return the wrapper from `Spec().New`.

### 2. Register in main.go

Add the import and append to `registry` in `cmd/deck/main.go`:

```go
import "github.com/dubyte/crosspoint-deck/pkg/templates/<name>"

var registry = []card.Spec{
 // ...existing...
 <name>.Spec(),
}
```

### 3. Add mage target

Add a function to `magefile.go`:

```go
func <Name>() error {
 mg.Deps(Build)
 fmt.Println("Generating <name> card...")
 return sh.Run("./deck", "<name>", "--output", "./output/<name>.bmp")
}
```

Then update the `All()` function to include `<Name>` in `mg.Deps`.

### 4. Build and test

```bash
go build -o deck ./cmd/deck
./deck <name> --output ./output/<name>.bmp
```

---

## Phase 3 Cards — Implementation Order

Implement in this order. Each card should be tested individually before moving on.

### Batch 1: Simple text/list cards (use packing/emergency patterns)

**1. Chore Chart** (`chore`)

- Modeled on habit tracker but simpler — just a checklist
- Flags: `--title`, `--chores` (comma-separated), `--portrait`, `--font`
- Layout: reversed header with title, then `[ ] chore` items in 2 columns
- Each item is `[ ] <chore name>`

**2. Shopping List** (`shopping`)

- Same pattern as packing but different title
- Flags: `--title` (default "Shopping List"), `--items` (comma-separated), `--portrait`, `--font`
- Layout: reversed header, `[ ] item` in 2 columns

**3. Loyalty Cards** (`loyalty`)

- Store membership numbers — simple label/value list
- Flags: `--title`, `--stores` (comma-separated `StoreName:MemberID` pairs), `--portrait`, `--font`
- Layout: reversed header, bold store name, regular member ID
- Parse like emergency contacts

**4. Library Card** (`library`)

- Single card with library info
- Flags: `--name`, `--card-number`, `--branch`, `--phone`, `--portrait`, `--font`
- Layout: reversed header "Library Card", bold labels (Card #, Branch, Phone), regular values
- Optionally render a barcode if card number is numeric (use QR package with a barcode-format string, or just display the number large)

**5. NATO Phonetic Alphabet** (`nato`)

- No flags needed except `--portrait`, `--font`
- Layout: reversed header "NATO Phonetic", 2 columns: `A  Alpha`, `B  Bravo`, etc.
- Letter in bold, word in regular. 13 rows per column.

**6. Morse Code Chart** (`morse`)

- No flags except `--portrait`, `--font`
- Layout: reversed header "Morse Code", 2 columns: `A  ·−`, `B  −···`, etc.
- Letter in bold, morse in regular. 13 rows per column.

**7. Resistor Color Codes** (`resistor`)

- No flags except `--portrait`, `--font`
- Layout: reversed header "Resistor Codes", 2 columns
- Each row: digit in bold, color bands (drawn as small filled rectangles) + color name in regular
- Since the X4 is B&W, draw the bands as labeled patterns: a black square for black, outlined square for white, etc. Or just list the color name — band drawing is a stretch goal.

**8. Common Conversions** (`convert`)

- No flags except `--portrait`, `--font`
- Layout: reversed header "Conversions", 2 columns
- Each row: conversion pair, e.g. `1 in = 2.54 cm`, `1 mi = 1.61 km`, `0°C = 32°F`, `1 kg = 2.2 lb`, `1 gal = 3.79 L`

### Batch 2: Medium complexity (use cheatsheet/business patterns)

**9. Recipe Card** (`recipe`)

- Flags: `--title`, `--ingredients` (comma-separated), `--steps` (comma-separated), `--time`, `--servings`, `--portrait`, `--font`
- Layout: reversed header with recipe name, then "⏱ TIME  ·  🍽 SERVINGS" row, then "Ingredients" bold label + ingredient list, then "Steps" bold label + numbered steps
- Landscape for ingredients/steps side-by-side; portrait for long recipes

**10. Coffee Brew Guide** (`coffee`)

- Flags: `--method`, `--ratio`, `--temp`, `--time`, `--steps` (comma-separated), `--portrait`, `--font`
- Layout: reversed header with method name, then ratio/temp/time in bold, then numbered steps

**11. Plant Care Guide** (`plant`)

- Flags: `--plant`, `--water`, `--light`, `--humidity`, `--food`, `--notes`, `--portrait`, `--font`
- Layout: reversed header with plant name, bold labels (Water, Light, Humidity, Food), regular values, notes at bottom

**12. Bodyweight Workout** (`workout`)

- Flags: `--title`, `--exercises` (comma-separated `Exercise:Reps` pairs), `--rounds`, `--rest`, `--portrait`, `--font`
- Layout: reversed header with title, then "Rounds: N · Rest: Xs" line, then exercise list with bold exercise names and regular rep counts, numbered

**13. Stretching Routine** (`stretch`)

- Same pattern as workout
- Flags: `--title`, `--stretches` (comma-separated `Name:Duration` pairs), `--portrait`, `--font`
- Layout: reversed header, numbered stretches with bold name and regular "Hold: Xs"

### Batch 3: Reference cards with more layout

**14. World Time Zones** (`timezones`)

- Flags: `--local` (your city/timezone, e.g. "New York"), `--cities` (comma-separated `City:Offset` pairs), `--portrait`, `--font`
- Layout: reversed header "World Time", then "Local: New York (EST)" line, then city list with offset: `Tokyo     +14h`, `London     +5h`, etc.
- Bold city name, regular offset

**15. Home Maintenance Log** (`maintenance`)

- Flags: `--year`, `--tasks` (comma-separated `Task:Month` pairs), `--portrait`, `--font`
- Layout: reversed header "Maintenance YYYY", then simple grid or list: task name + due month/year
- Bold task, regular schedule

---

## Registration Checklist

After implementing each card, verify:

```bash
cd /home/dubyte/Documents/Workspace/vibe/crosspoint-deck
go build -o deck ./cmd/deck         # must compile
./deck help                          # new command appears
./deck <name> --output ./output/<name>.bmp  # generates without warnings
ls -lh output/<name>.bmp             # ~1.1 MB
```

## Files You'll Edit

- `pkg/templates/<name>/<name>.go` — new file per card (15 files)
- `cmd/deck/main.go` — add import + registry entry per card
- `magefile.go` — add mage target + update `All()` per card

## Completion Check

```bash
cd /home/dubyte/Documents/Workspace/vibe/crosspoint-deck
mage All   # should generate all 24 cards (9 existing + 15 new)
```
