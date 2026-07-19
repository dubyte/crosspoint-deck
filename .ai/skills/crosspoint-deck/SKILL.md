---
name: crosspoint-deck
description: >
  Generate example cards for the CrossPoint Deck e-ink card system.
  Use when creating, testing, or demonstrating CrossPoint Deck card templates
  (calendar, wifi, business, cheatsheet, meeting, packing, emergency, habit).
  Triggers: crosspoint deck, crosspoint-deck, deck card, e-ink card, XTEink X4,
  BMP card, deck template, generate deck card.
---

# CrossPoint Deck Skill

Generate single-page 800×480 BMP cards for the XTEink X4 e-ink reader.
Cards are produced by `./deck <command>` in the project root.

## When This Skill MUST Be Used

- Generating example cards for any of the 8 card types
- Testing card rendering after code changes
- Creating starter packs (`.deckpack.zip` files)
- Debugging card layout, font, or dimension issues

## Prerequisites

```bash
cd /home/dubyte/Documents/Workspace/vibe/crosspoint-deck
go build -o deck ./cmd/deck   # or `mage build`
```

Cards are written to `./output/` as 1.1 MB uncompressed 24-bit BMPs (800×480 or 480×800).
All card types support `--portrait` (480×800) and `--font <path>` for a custom TTF.

## Card Types and Example Commands

### calendar — Year-at-a-glance

4×3 month grid (landscape) or 3×4 (portrait). Days start on Sunday.

```bash
./deck calendar --year 2026 --output ./output/calendar.bmp
./deck calendar --year 2026 --portrait --output ./output/calendar-portrait.bmp
```

Flags: `--year` (default: current year), `--portrait`, `--font`

### wifi — WiFi access card with QR code

Network name, password, and a scannable QR code for `WIFI:S:...;;`.

```bash
./deck wifi --ssid "HomeBase" --password "p4ssw0rd!" --encryption WPA \
  --output ./output/wifi.bmp
```

Flags: `--ssid`, `--password`, `--encryption` (WPA/WEP/nopass), `--portrait`, `--font`

### business — Business card with QR vCard

Contact info + scannable vCard QR code.

```bash
./deck business \
  --name "Jane Smith" \
  --title "Software Engineer" \
  --phone "+1-555-0199" \
  --email "jane@example.com" \
  --website "janesmith.dev" \
  --output ./output/business.bmp
```

Flags: `--name`, `--title`, `--phone`, `--email`, `--website`, `--portrait`, `--font`

### cheatsheet — Keyboard shortcuts

Two-column layout with `key:description` pairs.

```bash
./deck cheatsheet \
  --title "Vim Essentials" \
  --shortcuts "h:left,j:down,k:up,l:right,i:insert,Esc:normal,:w:save,:q:quit,dd:delete line,yy:yank line,p:paste,u:undo,Ctrl-r:redo,/:search,n:next match,%:jump to bracket,gg:top of file,G:bottom of file" \
  --output ./output/cheatsheet.bmp
```

Flags: `--title`, `--shortcuts` (comma-separated `key:desc` pairs), `--portrait`, `--font`

### meeting — Meeting room schedule

Room name with event list. Events are hardcoded in the template source; for custom
events, edit `pkg/templates/meeting/meeting.go` or pipe data via a wrapper script.

```bash
./deck meeting --room "Boardroom A" --output ./output/meeting.bmp
```

Flags: `--room`, `--portrait`, `--font`

### packing — Packing checklist

Title + comma-separated items, each prefixed with `[ ]`.

```bash
./deck packing \
  --title "Weekend Trip" \
  --items "Passport,Phone charger,Toothbrush,Socks x3,T-shirt x2,Hoodie,Sunglasses,Sunscreen,Book,Snacks,Water bottle,Power bank,Headphones" \
  --output ./output/packing.bmp
```

Flags: `--title`, `--items` (comma-separated), `--portrait`, `--font`

### emergency — Emergency contact card

Country, emergency numbers, blood type, allergies.

```bash
./deck emergency \
  --country "USA" \
  --contacts "Police:911,Fire:911,Ambulance:911,Poison Control:800-222-1222" \
  --blood "O+" \
  --allergies "Penicillin,Peanuts" \
  --output ./output/emergency.bmp
```

Flags: `--country`, `--contacts` (comma-separated `label:number` pairs), `--blood`, `--allergies`, `--portrait`, `--font`

### habit — Habit tracker grid

Habits as rows, days as columns. Each cell is an empty checkbox outline.

```bash
./deck habit \
  --title "Daily Habits" \
  --habits "Read 30min,Exercise,Meditate,Water 8 cups,Sleep 7h,No sugar,Stretch" \
  --days 14 \
  --output ./output/habit.bmp
```

Flags: `--title`, `--habits` (comma-separated), `--days` (max 31, default 7), `--portrait`, `--font`

## Generate All Examples at Once

```bash
cd /home/dubyte/Documents/Workspace/vibe/crosspoint-deck
rm -rf output && mkdir -p output && \
./deck calendar --year 2026 --output ./output/calendar.bmp && \
./deck wifi --ssid "HomeBase" --password "p4ssw0rd!" --output ./output/wifi.bmp && \
./deck business --name "Jane Smith" --title "Engineer" --phone "+1-555-0199" --email "jane@example.com" --website "janesmith.dev" --output ./output/business.bmp && \
./deck cheatsheet --title "Vim" --shortcuts "h:left,j:down,k:up,l:right,i:insert,Esc:normal,:w:save,:q:quit" --output ./output/cheatsheet.bmp && \
./deck meeting --room "Boardroom A" --output ./output/meeting.bmp && \
./deck packing --title "Weekend Trip" --items "Passport,Phone charger,Toothbrush,Socks,Hoodie" --output ./output/packing.bmp && \
./deck emergency --country "USA" --contacts "Police:911,Fire:911,Ambulance:911" --blood "O+" --allergies "Penicillin" --output ./output/emergency.bmp && \
./deck habit --title "Daily" --habits "Read,Exercise,Meditate" --days 7 --output ./output/habit.bmp && \
ls -lh output/
```

## Mage Shortcuts

The `magefile.go` has individual targets that wrap the CLI:

```bash
mage Build
mage Calendar       # landscape calendar
mage CalendarPortrait
mage WiFi
mage Business
mage Cheatsheet
mage Meeting
mage Packing
mage Emergency
mage Habit
mage All            # every card type at once
mage Packs          # all starter packs (.deckpack.zip)
```

## Starter Packs

Three `.deckpack.zip` files in `./packs/`:

```bash
mage PackCalendar     # 2026-calendar.deckpack.zip
mage PackDeveloper    # developer-reference.deckpack.zip
mage PackTravel       # travel-essentials.deckpack.zip
```

Each pack contains pre-rendered BMPs, a `manifest.json`, and a `preview.png`.

## Font Fallback

If `--font` is omitted, `layout.LoadFontFace` tries these paths in order:

- `/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf`
- `/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf`
- `/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf`
- `/usr/share/fonts/truetype/freefont/FreeSans.ttf`
- macOS: `/System/Library/Fonts/Helvetica.ttc`
- Windows: `/Windows/Fonts/arial.ttf`

## Adding a New Card Type

1. Create `pkg/templates/<name>/<name>.go` implementing `card.Card`
2. Export a `Spec() card.Spec` function
3. Register in `cmd/deck/main.go`'s `registry` slice
4. Add a mage target in `magefile.go`

## Project Layout

```
pkg/
├── bmp/encoder.go       # 24-bit uncompressed BMP encoder
├── card/card.go         # Card interface + Spec + Factory
├── layout/layout.go     # Font loading, grid, text helpers
├── qr/                  # QR code generation (WiFi + vCard)
├── render/render.go     # Render + dimension validation
└── templates/           # One subdirectory per card type
    ├── business/
    ├── calendar/
    ├── cheatsheet/
    ├── emergency/
    ├── habit/
    ├── meeting/
    ├── packing/
    └── wifi/
```

## Quick Reference

- **Design system:** Reversed black header bar + 2px divider + bold labels/regular values. Use `layout.DrawReversedHeader` and `layout.LoadFontFaceBold`.
- **Adding a card:** Create `pkg/templates/<name>/<name>.go`, register in `cmd/deck/main.go`'s `registry`, add mage target in `magefile.go`.
- **Phase 3 prompt:** `.ai/prompts/phase3-implementation.md` — detailed instructions for implementing 15 new card types.
- **Roadmap:** `ROADMAP.md` — full card catalog and phase planning.
- **Agent rules:** `agents.md` — hardware constraints, design principles, commit conventions.
