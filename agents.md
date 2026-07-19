# Agent Instructions: CrossPoint Deck

## Project

Go CLI + library that generates 800×480 monochrome BMP "cards" for the XTEink X4 e-ink reader. Cards are single-page utilities (calendars, QR codes, cheat sheets) synced to the device via the upstream `crosspoint-sync` pipeline.

## Build & Dev Commands

- **Task runner:** [Mage](https://magefile.org/). Targets are Go functions in `magefile.go` and `packs.go`.
  - `mage -l` — list tasks
  - `mage` — default is `Build`
  - `mage verify` — `go vet` + `go test ./...` + `file` check on output BMPs
  - `mage all` — generate every card type
  - `mage clean` — remove `deck` binary and `./output/`
- **Direct CLI:**
  - `go build ./cmd/deck` — produces `./deck` binary
  - `./deck <command> --output ./output/foo.bmp` — generate a single card
  - `./deck --display portrait|landscape <command>` — override orientation
- **No `make`, no Docker, no CI config in repo.**

## Architecture (Non-Obvious)

### Entrypoint & Registry

- `cmd/deck/main.go` holds the **only** registry: a `[]card.Spec` slice. Adding a card type requires **two steps**:
  1. Create `pkg/templates/<name>/` with `Spec()` + `Render()`
  2. Add one import + one `spec.Name()` call in `main.go`
- **No `init()` magic, no reflection, no auto-discovery.** The registry is intentionally explicit and debuggable.

### Package Layout

| Package | Purpose |
|---|---|
| `pkg/card/` | `Card` interface (`Render() image.Image`) + `Spec` struct + `AssertDimensions` test helper |
| `pkg/render/` | `ToFile(c, display, path)` — validates dimensions, rotates if needed, calls `bmp.Encode` |
| `pkg/bmp/` | Pure-Go encoder: uncompressed 24-bit BMP, no compression, preserves grayscale values |
| `pkg/layout/` | Shared drawing primitives: `DrawReversedHeader`, `LoadFontFace`/`LoadFontFaceBold`, `Grid`, `DrawWrappedText` |
| `pkg/qr/` | QR code generation wrapper |
| `pkg/pack/` | `.deckpack.zip` builder (BMPs + `manifest.json` + `preview.png`) |
| `pkg/templates/<name>/` | One package per card type. Each exports `Spec() card.Spec` |

### Design System (Enforced Convention)

Every card **must** use:
- **Reversed black header bar** (`layout.DrawReversedHeader`) with white bold title
- **2px black divider** below the header
- **Bold label + regular value** typographic hierarchy
- **Pure black (#000000) and white (#FFFFFF)** only. Grayscale is handled by the X4's hardware dithering; do NOT threshold or apply software dithering.
- **Minimum line thickness:** 2 px. 1 px lines disappear on e-ink.
- **Minimum text size:** 12 px height; 14–16 px preferred.

### Orientation Rules

- Cards render at **800×480** (landscape, default) or **480×800** (portrait).
- The `Card.Render()` implementation decides content orientation via a `Portrait bool` field (or similar).
- `render.ToFile` auto-rotates 90° clockwise if the `display` arg (from `--display` flag) differs from content orientation. **Do not rotate in the template.**

### Font Loading

- `layout.LoadFontFace` / `LoadFontFaceBold` accept an optional `--font` path.
- If empty, they walk a **hardcoded system font fallback chain**: DejaVu → Liberation → Noto → Helvetica → Arial (bold variants for bold).
- If no font is found, rendering fails. On headless CI, install `fonts-dejavu-core` or similar.

## Adding a New Card Type

1. `mkdir pkg/templates/<name>`
2. Implement:
   - `type Card struct { ... }` with fields for your data
   - `func (c *Card) Render() image.Image` — use `gg.NewContext(W, H)`, call `layout.DrawReversedHeader`, draw content
   - `func Spec() card.Spec` — return `Name`, `Usage`, and a `Factory` that registers flags and returns the card
3. Add `pkg/templates/<name>` import and `name.Spec()` to `cmd/deck/main.go` registry slice
4. Add `mage <Name>` target in `magefile.go` (optional but preferred)
5. Add test in `pkg/templates/<name>/<name>_test.go` using `card.AssertDimensions(t, c, 800, 480)`

## Testing

- Run: `go test ./...` or `mage test`
- Each template should have a `_test.go` with at least:
  - `card.AssertDimensions(t, c, 800, 480)` for landscape
  - `card.AssertDimensions(t, c, 480, 800)` for portrait if supported
  - A test that `Spec().New(flag.NewFlagSet(...))` returns non-nil
- `go vet ./...` is part of `mage verify`

## Output Verification

- Generated BMPs must be **uncompressed 24-bit**.
  - Verify: `file foo.bmp` should report `PC bitmap, Windows 3.x format, 800 x 480 x 24` (or 480×800).
- `./output/` and `./packs/*.zip` are **gitignored**. Never commit generated assets.
- Commit templates, source, and manifest JSON. Treat BMPs as build artifacts.

## Upstream Constraints (Hard Boundaries)

- **Target device:** XTEink X4, ESP32-C3, ~380 KB RAM, no PSRAM.
- **Only `.bmp` is viewable** via the firmware's `BmpViewerActivity`. No PNG/JPG support.
- **File browser navigation:** Placing multiple `.bmp` files in the same SD card folder gives prev/next button navigation. Design folder structures as "decks."
- **Do NOT** write firmware code (C++, PlatformIO), mobile app code (React Native), or new network protocols. This repo is strictly a **content generator** that produces static files.

## Common Mistakes to Avoid

- **Thresholding grayscale:** The BMP encoder writes actual 24-bit pixel values. Let the X4's SSD1677 dither to 4 levels. Do not pre-threshold to black/white.
- **Embedding photos:** Uncompressed 800×480 BMP is ~1.15 MB. Photo cards bloat quickly; prefer vector/text designs.
- **Forgetting the registry:** A template package without an entry in `main.go` is unreachable from the CLI.
- **Custom `--font` without fallback:** If you change font loading logic, ensure the fallback chain still works or CI will break.
- **Software dithering:** Do not apply Floyd-Steinberg or ordered dithering. The hardware does it better.
