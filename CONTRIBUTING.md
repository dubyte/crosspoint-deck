# Contributing to CrossPoint Deck

Thank you for considering contributing! This guide covers how to add new card types, fork the project, and submit changes.

## Ways to Contribute

1. **Add a new card type** — The most impactful contribution. Got an idea for a card you'd laminate in real life? Turn it into a template.
2. **Improve an existing card** — Better layout, more options, or portrait-mode support.
3. **Report bugs** — If a card looks wrong on your XTEink X4, open an issue with a photo.
4. **Share your cards** — Post photos of your deck on social media or [GitHub Discussions](https://github.com/dubyte/crosspoint-deck/discussions).

## Quick Start for Contributors

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or newer
- A system font (DejaVu, Liberation, Noto, Helvetica, or Arial — the tool auto-detects)

### Fork and Clone

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR-USERNAME/crosspoint-deck.git
cd crosspoint-deck
```

### Build and Verify

```bash
go build ./cmd/deck
mage verify    # or: go test ./... && go vet ./...
```

## Adding a New Card Type

Adding a card takes about 20–30 minutes if you copy an existing simple card as a starting point. Here's the full process:

### Step 1: Pick a Starting Point

The simplest cards to copy are `owner` or `nato` — they are pure text with a reversed header and minimal layout. More complex examples: `recipe` (two-column layout), `calendar` (grid math), `wifi` or `business` (QR codes).

```bash
# Copy a simple template as your starting point
cp -r pkg/templates/owner pkg/templates/my-card
```

### Step 2: Implement Your Card

Edit `pkg/templates/my-card/my-card.go`:

```go
package mycard

import (
    "flag"
    "image"
    "image/color"

    "github.com/dubyte/crosspoint-deck/pkg/card"
    "github.com/dubyte/crosspoint-deck/pkg/layout"
    "github.com/fogleman/gg"
)

type Card struct {
    Title    string
    Items    []string
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

    // Reversed black header bar — required by the design system
    bodyY := layout.DrawReversedHeader(dc, c.Title, W, 26, c.FontPath)

    // Your content here
    _ = layout.LoadFontFace(dc, c.FontPath, 20)
    dc.SetColor(color.Black)
    for i, item := range c.Items {
        dc.DrawString(item, 30, bodyY+30+float64(i)*36)
    }

    return dc.Image()
}

func Spec() card.Spec {
    return card.Spec{
        Name:  "my-card",
        Usage: "Generate a my-card description",
        New: func(fs *flag.FlagSet) card.Card {
            c := &Card{}
            fs.StringVar(&c.Title, "title", "My Card", "Card title")
            // Add more flags for your data fields
            fs.BoolVar(&c.Portrait, "portrait", false, "Render in portrait orientation")
            fs.StringVar(&c.FontPath, "font", "", "Path to TTF font (optional)")
            return c
        },
    }
}
```

**Design system rules you must follow:**
- Use `layout.DrawReversedHeader` for the title bar
- Use bold labels + regular values for text hierarchy
- Pure black (#000000) and white (#FFFFFF) only
- Minimum line thickness: 2 px
- Minimum text size: 12 px (14–16 px preferred)
- Do NOT rotate in the template — `render.ToFile` handles orientation

### Step 3: Add Tests

Create `pkg/templates/my-card/my-card_test.go`:

```go
package mycard

import (
    "flag"
    "testing"

    "github.com/dubyte/crosspoint-deck/pkg/card"
)

func TestMyCard_Render(t *testing.T) {
    c := &Card{
        Title: "Test",
        Items: []string{"One", "Two", "Three"},
    }
    card.AssertDimensions(t, c, 800, 480)
}

func TestMyCard_Portrait(t *testing.T) {
    c := &Card{
        Title:    "Test",
        Items:    []string{"One", "Two"},
        Portrait: true,
    }
    card.AssertDimensions(t, c, 480, 800)
}

func TestSpec(t *testing.T) {
    spec := Spec()
    if spec.Name != "my-card" {
        t.Errorf("name = %q, want my-card", spec.Name)
    }
    fs := flag.NewFlagSet("test", flag.ContinueOnError)
    c := spec.New(fs)
    if c == nil {
        t.Fatal("spec.New returned nil")
    }
    card.AssertDimensions(t, c, 800, 480)
}
```

### Step 4: Register the Card

Add your import and one line to `cmd/deck/main.go`:

```go
import (
    // ... existing imports ...
    "github.com/dubyte/crosspoint-deck/pkg/templates/my-card"
)

var registry = []card.Spec{
    // ... existing specs ...
    mycard.Spec(),
}
```

### Step 5: Add a Mage Target (Optional)

In `magefile.go`, add:

```go
// MyCard generates a my-card example.
func MyCard() error {
    mg.Deps(Build)
    fmt.Println("Generating my-card...")
    return sh.Run("./deck", "my-card", "--title", "Example", "--output", "./output/my-card.bmp")
}
```

### Step 6: Verify Everything

```bash
go test ./pkg/templates/my-card/...
go vet ./...
./deck my-card --title "Test" --output /tmp/test.bmp
file /tmp/test.bmp    # Should say 800 x 480 x 24
```

### Step 7: Update Documentation

1. Add your card to the table in `EXAMPLES.md` (copy an existing entry)
2. Generate a PNG preview: `./deck my-card ... --output /tmp/my-card.bmp && convert /tmp/my-card.bmp docs/images/example-my-card.png`
3. Update the card count in `README.md` if you feel like it

## Submitting Changes

1. **Open an issue first** for large changes (new card categories, breaking layout changes). Small fixes and new card types can go straight to PR.
2. **One card per PR** — keeps review focused and fast.
3. **Include a generated PNG** in the PR description so reviewers can see what it looks like without building.
4. **Follow the design system** — all cards must use the reversed header, bold/regular hierarchy, and pure B&W palette.

## Design Philosophy

The best cards answer: **"Would I laminate this and pin it up?"**

- **Always-on > interactive.** If it needs buttons or scrolling, it's a phone app.
- **One glance = full answer.** No menus, no "next page."
- **Replace the laminated original.** Every good card has a physical ancestor.
- **Personal data over generic content.** A recipe with *your* measurements beats a generic one.

If your card idea fits these principles, it's probably a great addition.

## Code of Conduct

Be kind, be constructive, and assume good intent. This is a small project by one person; thoughtful contributions are deeply appreciated.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
