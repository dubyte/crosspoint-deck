# Agent Instructions: CrossPoint Deck

This document is a reference for any AI agent or human contributor working on the `crosspoint-deck` repository. It records the discovered constraints of the upstream CrossPoint ecosystem and sets architectural boundaries for this project.

## 1. Ecosystem Constraints (Discovered from Source)

### 1.1 Target Hardware — XTEink X4

- **MCU:** ESP32-C3 (single-core RISC-V @ 160 MHz)
- **RAM:** ~380 KB usable. **No PSRAM.** This is a hard ceiling.
- **Display:** 800 × 480 pixels, monochrome e-ink, single 48 KB framebuffer.
- **Storage:** SD card (FAT32). All card assets must live on SD card as ordinary files.
- **Refresh:** Slow. Full updates take 1–2 seconds. Partial/fast refreshes are possible but the BMP viewer uses `HALF_REFRESH` for initial display.

### 1.2 Firmware Rendering Engine — CrossPoint Reader

The firmware's native image pipeline is the only rendering engine Deck may use. Key behaviors:

- **Supported image format:** `.bmp` only. The firmware has a native `Bitmap` parser and `BmpViewerActivity`.
- **BMP specs that work:** Uncompressed BMP with 24-bit color depth. Other bit depths may not be handled by the parser.
- **Resolution target:** 800 × 480 (portrait). The viewer scales images that are larger, but scales down (not up) with centering. For pixel-perfect clarity, match the panel exactly.
- **Orientation:** The firmware supports four orientations (Portrait, Inverted, Landscape CW, Landscape CCW). The BMP viewer does not rotate images automatically; if you want landscape cards, generate them at 480 × 800 or accept portrait cropping.
- **Grayscale:** The display is 1-bit black/white. The firmware can simulate grayscale via dithering (`drawBitmap` + `displayGrayBuffer` with MSB/LSB passes), but this is **slower** and visually softer. For card content that must be crisp (text, QR codes, fine lines), use **pure black and white** (no grayscale).
- **Contrast:** E-ink has limited contrast ratio. Thin gray lines or low-contrast color pairs disappear. Always maximize contrast: black on white or white on black.

### 1.3 File Discovery & Browsing

- The firmware's **Browse Files** screen lists all files and folders on the SD card. Hidden dotfiles (starting with `.`) are hidden unless the user enables `showHiddenFiles`.
- Selecting a `.bmp` file opens `BmpViewerActivity`, which renders it full-screen.
- If multiple `.bmp` files exist in the same folder, the viewer provides **prev/next navigation** using the side buttons (Left/Right or Volume Up/Down). This means a folder of cards is naturally a swipeable deck.
- The viewer also offers a "Set as sleep cover" action (Confirm button), so any card can become the sleep screen.

### 1.4 Sync Pipeline — CrossPoint Sync

`crosspoint-sync` is the companion iOS app. Deck must interoperate with it, not replace it.

- **Discovery:** UDP broadcast on port 8134, or manual IP entry.
- **HTTP API (port 80):**
  - `GET /api/files?path=/Folder` — list files
  - `POST /mkdir` — create folders
  - `POST /upload` — HTTP multipart upload (fallback)
  - `POST /delete` — delete files
- **WebSocket Upload (port 81):**
  - Protocol: `START:filename:size:path` → `READY` → binary chunks (64 KB window) → `PROGRESS:received:total` → `DONE`
  - Chunk size: 4 KB; window: 16 chunks (64 KB) before backpressure wait.
  - The app uses `react-native-udp` and requires a native dev build (Expo Go is insufficient).
- **Default paths:** The sync app uploads EPUBs to `/` by default and clipped articles to `/Articles`. Deck should use its own folder (e.g., `/Cards/`) to avoid cluttering the book root.
- **Folder creation:** `ensureRemotePath` recursively creates missing folders before upload.

### 1.5 Memory & Performance Boundaries

- **Firmware RAM budget:** Any payload that forces the firmware to allocate large temporary buffers risks an OOM reboot. BMPs are streamed from SD card during rendering, but the parser may allocate line buffers. Keep BMP row widths reasonable (800 px × 3 bytes = 2.4 KB per row — trivial).
- **File size:** Uncompressed 800×480×24-bit BMP = ~1.15 MB per card. This is acceptable for SD card storage but should not grow larger unnecessarily. Do not embed high-resolution photos; use vector-based designs.
- **No server-side rendering on device:** Do not propose running Python, Node, or any interpreter on the ESP32. All generation happens on the host (phone or computer) before sync.

## 2. Architectural Boundaries

### 2.1 What CrossPoint Deck Is

- A **content generator.** It produces static `.bmp` files and folder structures.
- A **consumer of the sync API.** It hands files to `crosspoint-sync` (or instructs the user to copy them manually).
- **Display-agnostic in design, display-specific in output.** Templates can be authored in SVG/HTML/CSS at arbitrary resolution, but the render pipeline must export exact 800×480 BMPs.

### 2.2 What CrossPoint Deck Is Not

- **Not a firmware patch.** No C++, no PlatformIO, no `lib/`, `src/`, or `freeink-sdk/` changes.
- **Not a mobile app.** No React Native, no Expo, no iOS signing, no app store submission.
- **Not a new protocol.** No custom UDP ports, no new WebSocket message types, no REST endpoints.
- **Not an EPUB generator.** Cards are bitmaps, not reflowable documents. Do not introduce EPUB, HTML-in-ZIP, or CSS layout engines for card content.

### 2.3 Layer Diagram

```
┌─────────────────────────────────────────┐
│         CrossPoint Deck (this repo)      │
│  Templates → Render → 800×480 BMP files  │
│  Collection metadata (JSON/YAML indexes)   │
└─────────────────────────────────────────┘
                    │
                    ▼ (push via sync API or SD copy)
┌─────────────────────────────────────────┐
│      CrossPoint Sync (upstream iOS app) │
│  Discovery / Upload Queue / WebSocket   │
└─────────────────────────────────────────┘
                    │
                    ▼ (WebSocket / HTTP / SD card)
┌─────────────────────────────────────────┐
│      CrossPoint Reader (upstream firmware)│
│  SD card → File Browser → BMP Viewer     │
│  800×480 e-ink display                  │
└─────────────────────────────────────────┘
```

**Rule:** Deck code must never import from `crosspoint-sync` source paths or `crosspoint-reader` firmware headers. It may document the API contract, but it does not embed it.

## 3. Code Generation Guidelines

### 3.1 Asset Pipeline: Vector SVG → Monochrome BMP

The recommended pipeline for card generation:

1. **Design in SVG** (or programmatic SVG generation).
   - Use precise coordinates; text should be converted to paths or use known font metrics.
   - Color must be reduced to black (`#000000`) and white (`#FFFFFF`) only. No anti-aliased grays.
2. **Rasterize to exact resolution.**
   - Target: 800 × 480 px.
   - If the design is landscape-first, render at 480 × 800 and accept that the viewer will display it in the current orientation without rotation.
3. **Export as uncompressed 24-bit BMP.**
   - The firmware parser expects standard BMP headers with 24-bit RGB and no compression.
   - Tools: ImageMagick (`convert -depth 24 -compress none`), Pillow (Python), or sharp (Node.js).

### 3.2 Why Not PNG/JPG?

- The firmware does **not** have a PNG or JPEG decoder in the BMP viewer path.
- EPUB-embedded JPG/PNG are decoded by the EPUB engine, which is heavyweight and unnecessary for a single static card.
- BMP is parsed as a stream: headers → row data. Minimal RAM, no decompression codec.

### 3.3 Contrast and Legibility Requirements

- **Minimum line thickness:** 2 px for black lines on white. 1 px lines can disappear in e-ink ghosting.
- **Text size:** No smaller than 12 px height for body text; 16 px preferred. E-ink dot gain makes small text bleed.
- **QR codes:** Use a minimum quiet zone of 4 modules. Test at actual display size; high-density QR codes may fail to scan if dithered.
- **Dithering:** Do not use Floyd-Steinberg or ordered dither for text/line art. Use thresholding (black if luminance < 128, else white). Reserve grayscale simulation for photographs only, and only if the card design explicitly needs them.

### 3.4 Folder Naming Conventions

Cards are organized on the SD card under a root collection folder. Suggested structure:

```
/Cards/
  Business/
    contact-qrcode.bmp
    meeting-schedule.bmp
  Calendar/
    2026-year-at-a-glance.bmp
    Q3-planning.bmp
  Reference/
    keyboard-shortcuts.bmp
    spanish-phrases.bmp
```

- Root folder: `/Cards/` (or user-configurable, but default to this).
- Collection subfolders: PascalCase, no spaces if possible (firmware file browser handles spaces, but URLs and terminal commands are simpler without them).
- Filenames: kebab-case, descriptive, `.bmp` extension only.
- Do not use hidden folders (starting with `.`) unless the user explicitly wants them hidden from the browser.

### 3.5 Metadata Sidecars (Optional)

If a card needs a title, description, or update timestamp that the firmware cannot display (the BMP viewer has no overlay text), store a sidecar JSON file next to the BMP:

```
/Cards/Business/contact-qrcode.bmp
/Cards/Business/contact-qrcode.json
```

The sidecar is for the **host-side** generator or a future mobile UI; the firmware ignores it. Keep sidecars small (< 1 KB).

## 4. Testing and Verification Rules

- **Never claim a card "should work" without inspecting the BMP header.**
  - Verify: `file mycard.bmp` should report `PC bitmap, Windows 3.x format, 800 x 480 x 24`.
  - Verify: ImageMagick `identify -verbose mycard.bmp | grep "Print size"` or pixel dimensions.
- **Never generate BMPs with compression.** The firmware parser does not handle RLE or Huffman-compressed BMPs.
- **Test contrast empirically:** If you cannot threshold the design to pure black/white without losing meaning, redesign it. Do not rely on grayscale as a crutch.
- **Respect upstream scope:** CrossPoint Reader's `SCOPE.md` rejects interactive apps, active connectivity, and complex annotation. Deck must not propose features that violate this scope (e.g., animated cards, live data, touch interaction).

## 5. Commit and Code Style Conventions

- Follow the host language's standard tooling (Prettier, Black, gofmt, etc.) once a language is chosen.
- Do not commit generated `.bmp` files to git. Use `.gitignore` for `/output/` or `/dist/` directories.
- Commit templates and source files (SVG, JSON, scripts). Treat rendered BMPs as build artifacts.
- If committing metadata about the upstream API, cite the source file path in the upstream repo (e.g., `crosspoint-reader/src/activities/util/BmpViewerActivity.cpp`) so future agents can verify.

## 6. Summary Checklist for Every Change

Before generating or modifying card assets, confirm:

- [ ] Output format is uncompressed 24-bit BMP.
- [ ] Resolution matches 800 × 480 (X4 portrait) or 480 × 800 (X4 landscape, if accepting no auto-rotation).
- [ ] Design is pure black and white for text/line-art cards; grayscale only if explicitly needed.
- [ ] Folder structure is under `/Cards/` or a user-defined root, not mixed with `/Books/` or `/Articles/`.
- [ ] No firmware code, no mobile app code, no new network protocols were introduced.
- [ ] File sizes are reasonable (~1 MB per card, not 10 MB).
- [ ] Contrast is verified by thresholding the design to 1-bit.
