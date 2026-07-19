# CrossPoint Deck

Turn your XTEink X4 e-reader into a glanceable, high-utility information hub.

CrossPoint Deck is a content-generation layer for the [CrossPoint Reader](https://github.com/crosspoint-reader/crosspoint-reader) ecosystem. Instead of long-form books, it produces **cards** — single-page, self-contained utilities such as business cards, year-at-a-glance calendars, QR code sheets, and cheat sheets — formatted for the device's 800×480 e-ink display and synced via the existing [CrossPoint Sync](https://github.com/zabirauf/crosspoint-sync) pipeline.

## The Idea

E-ink excels at static, high-contrast information that you reference in seconds, not minutes. CrossPoint Deck treats the reader as a **deck of cards**: each card is an independent, instantly readable single-page asset. You browse folders on the device, tap a `.bmp` file, and the native image viewer displays it full-screen with no rendering overhead, no page turns, and no loading delay.

## What It Does

- **Generates cards** — scripts and templates that produce `.bmp` images optimized for the X4's monochrome display.
- **Organizes collections** — cards are grouped into folders (e.g., `/Cards/Work/`, `/Cards/Travel/`) that the firmware's existing file browser navigates natively.
- **Syncs upstream** — cards are pushed to the device through the same WebSocket/HTTP upload pipeline that `crosspoint-sync` uses for EPUBs. No custom protocol, no firmware patch.

## What It Does *Not* Do

- **It is not a firmware fork.** CrossPoint Deck does not modify the CrossPoint Reader firmware, add new activities, or consume RAM on the ESP32-C3.
- **It is not a sync client.** It does not reimplement discovery, WebSocket chunking, or queue management; it delegates all transport to `crosspoint-sync` or manual SD-card copy.
- **It does not render PDFs or EPUBs.** Cards are flat bitmaps. The firmware's existing BMP viewer is the rendering engine.

## Target Hardware

- **Device:** XTEink X4 (ESP32-C3, 800×480 monochrome e-ink, ~380 KB usable RAM)
- **Display:** Single 48 KB framebuffer. Monochrome (1-bit) with optional grayscale support.
- **Storage:** SD card. Cards live as ordinary files in ordinary folders.
- **Navigation:** Firmware file browser + native BMP viewer (prev/next through sibling `.bmp` files in the same folder).

## How a Card Reaches the Device

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Card Template  │ ──► │  Render Script   │ ──► │   800×480 BMP   │
│  (SVG/HTML/...) │     │  (rasterizer)    │     │  (24-bit uncomp)│
└─────────────────┘     └──────────────────┘     └─────────────────┘
                                                        │
┌─────────────────┐     ┌──────────────────┐            │
│  XTEink X4      │ ◄── │ crosspoint-sync  │ ◄─────────┘
│  (BMP viewer)   │     │  (upload queue)  │   WebSocket/HTTP
└─────────────────┘     └──────────────────┘
```

1. A template or script renders a card to an **uncompressed 24-bit BMP** at exactly 800×480 pixels.
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

## Project Status

CrossPoint Deck is in early conceptual development. The repo currently contains:

- Documentation and constraints for future AI agents and contributors.
- Placeholder structure for card-generation scripts.

No implementation details are final. The goal is to establish scope, boundaries, and rendering constraints before writing generators.

## License

MIT — use, modify, distribute, and build on this freely. See [LICENSE](./LICENSE) for the full text.
