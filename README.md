# CrossPoint Deck

Turn your XTEink X4 e-reader into a glanceable, high-utility information hub.

CrossPoint Deck is a content-generation layer for the [CrossPoint Reader](https://github.com/crosspoint-reader/crosspoint-reader) ecosystem. Instead of long-form books, it produces **cards** — single-page, self-contained utilities such as business cards, year-at-a-glance calendars, QR code sheets, and cheat sheets — formatted for the device's 800×480 e-ink display and synced via the existing [CrossPoint Sync](https://github.com/zabirauf/crosspoint-sync) pipeline.

## The Idea

E-ink excels at static, high-contrast information that you reference in seconds, not minutes. CrossPoint Deck treats the reader as a **deck of cards**: each card is an independent, instantly readable single-page asset. You browse folders on the device, tap a `.bmp` file, and the native image viewer displays it full-screen with no rendering overhead, no page turns, and no loading delay.

## What It Does

- **Generates cards** — Go packages that render `.bmp` images with anti-aliased edges for the X4's 4-level grayscale display (SSD1677 controller).
- **Organizes collections** — cards are grouped into folders (e.g., `/Cards/Work/`, `/Cards/Travel/`) that the firmware's existing file browser navigates natively.
- **Syncs upstream** — cards are pushed to the device through the same WebSocket/HTTP upload pipeline that `crosspoint-sync` uses for EPUBs. No custom protocol, no firmware patch.

## What It Does *Not* Do

- **It is not a firmware fork.** CrossPoint Deck does not modify the CrossPoint Reader firmware, add new activities, or consume RAM on the ESP32-C3.
- **It is not a sync client.** It does not reimplement discovery, WebSocket chunking, or queue management; it delegates all transport to `crosspoint-sync` or manual SD-card copy.
- **It does not render PDFs or EPUBs.** Cards are flat bitmaps. The firmware's existing BMP viewer is the rendering engine.

## Target Hardware

- **Device:** XTEink X4 (ESP32-C3, 800×480 monochrome e-ink, ~380 KB usable RAM)
- **Display:** 4-level grayscale via SSD1677 controller (white, light gray, dark gray, black). BMP encoder preserves anti-aliased edges; device dithers natively.
- **Storage:** SD card. Cards live as ordinary files in ordinary folders.
- **Navigation:** Firmware file browser + native BMP viewer (prev/next through sibling `.bmp` files in the same folder).

## How a Card Reaches the Device

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Card Template  │ ──► │  Render Script   │ ──► │   800×480 BMP   │
│  (Go code)      │     │  (fogleman/gg)   │     │  (24-bit uncomp)│
└─────────────────┘     └──────────────────┘     └─────────────────┘
                                                        │
┌─────────────────┐     ┌──────────────────┐            │
│  XTEink X4      │ ◄── │ crosspoint-sync  │ ◄─────────┘
│  (BMP viewer)   │     │  (upload queue)  │   WebSocket/HTTP
└─────────────────┘     └──────────────────┘
```

1. A template renders a card to an **uncompressed 24-bit BMP** at exactly 800×480 pixels (or 480×800 portrait). Anti-aliased edges and subtle shading are preserved for the X4's grayscale display.
2. The BMP is placed in a collection folder (e.g., `/Cards/Business/`).
3. `crosspoint-sync` (or a manual copy) pushes the folder to the SD card.
4. On the device, you browse to the folder and select the `.bmp`. The native viewer opens it full-screen.

## Card Types (Examples)

| Collection | Example Cards |
|---|---|
| **Business** | Contact card with QR vCard, meeting-room schedules, one-page pitch sheets |
| **Calendar** | Year-at-a-glance, quarterly planning grids, holiday lists |
| **Reference** | Keyboard shortcuts, language cheat sheets, unit-conversion tables |
| **Travel** | Packing checklists, itinerary summaries, translation cards |
| **Health** | Medication schedules, gym routines, measurement logs |

Each card is a single `.bmp` file. A folder of `.bmp` files becomes a swipeable deck thanks to the firmware's built-in sibling-image navigation.

## Relationship to the Ecosystem

| Project | Role |
|---|---|
| **crosspoint-reader** | The firmware. Provides the BMP viewer, file browser, and sleep-screen BMP support that Deck relies on. |
| **crosspoint-sync** | The companion app. Provides the upload queue, device discovery, and folder management that gets Deck assets onto the device. |
| **crosspoint-deck** | **This repo.** Generates the card assets and defines collection structures. It sits *above* the sync layer and feeds it standard files. |

## Quick Start

CrossPoint Deck uses [Mage](https://magefile.org/) as its task runner. Tasks are defined in `magefile.go` as Go functions, providing cross-platform shell completion and discoverability.

```bash
# List available tasks
$ mage -l

# Run the default task (build)
$ mage

# Generate calendars
$ mage calendar          # landscape 800×480
$ mage calendarPortrait  # portrait 480×800
$ mage all               # both orientations

# Verify output format and run linting
$ mage verify

# Clean build artifacts
$ mage clean
```

You can also use the compiled CLI directly:

```bash
# Build
$ go build ./cmd/deck

# Generate a landscape calendar (default: 800×480)
$ ./deck calendar --year 2026 --output ./output/calendar-2026.bmp

# Generate a portrait calendar (480×800)
$ ./deck calendar --year 2026 --portrait --output ./output/calendar-2026-portrait.bmp

# Use a custom font
$ ./deck calendar --year 2026 --font /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf --output ./output/calendar-2026.bmp
```

## Available Card Types

Run `deck --help` to see all commands.

| Command | Description |
|---|---|
| `calendar` | Year-at-a-glance calendar with month gaps (landscape/portrait) |
| `wifi` | WiFi access card with QR code and bold labels |
| `business` | Business card with QR vCard and reversed header |
| `cheatsheet` | Keyboard shortcuts cheat sheet with bold key bindings |
| `meeting` | Meeting room schedule with reversed room header |
| `packing` | Packing checklist with checkbox items |
| `emergency` | Emergency contact card with bold labels |
| `habit` | Habit tracker grid with bold habit names |
| `chore` | Chore chart checklist with checkbox items |
| `coffee` | Coffee brew guide with ratios, temp, and steps |
| `convert` | Common unit conversions reference |
| `library` | Library card with card number and branch |
| `loyalty` | Loyalty cards list with store numbers |
| `maintenance` | Home maintenance log checklist |
| `morse` | Morse code reference chart |
| `nato` | NATO phonetic alphabet reference |
| `owner` | Owner identification card (name, email, optional phone) |
| `plant` | Plant care guide with water, light, humidity |
| `recipe` | Recipe card with ingredients and steps |
| `resistor` | Resistor color code reference |
| `shopping` | Shopping list checklist |
| `stretch` | Stretching routine guide |
| `timezones` | World time zones reference |
| `workout` | Bodyweight workout card with exercises and rounds |

## Starter Packs

CrossPoint Deck can bundle cards into `.deckpack.zip` files for easy distribution.

```bash
# Generate starter packs
$ mage PackCalendar     # 2026 Calendar Pack
$ mage PackDeveloper    # Developer Reference Pack
$ mage PackTravel       # Travel Essentials Pack
$ mage Packs            # All packs
```

A `.deckpack.zip` contains:

- `bmps/` — pre-rendered BMP files
- `manifest.json` — pack metadata (name, author, tags)
- `preview.png` — thumbnail for gallery display

## Project Status

CrossPoint Deck is in active development. Current capabilities:

- Pure-Go BMP encoder (24-bit uncompressed, preserves grayscale)
- Registry-based template system (add a template in 2 steps)
- 24 card types across 8 categories with unified design system
- Reversed black header bar + bold/regular typographic hierarchy on all cards
- 4-level grayscale support for the X4's SSD1677 display controller
- Pack generation and distribution format (`.deckpack.zip`)
- Mage-based build automation

See [ROADMAP.md](./ROADMAP.md) for planned phases.

## License

MIT — use, modify, distribute, and build on this freely. See [LICENSE](./LICENSE) for the full text.
