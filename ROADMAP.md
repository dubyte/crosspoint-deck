# CrossPoint Deck — Roadmap

This roadmap synthesizes findings from three expert analyses (Product & UX, Technical Architecture, Community & Distribution) into actionable phases. It is designed for a director who approves direction, not implementation details.

---

## Phase 1: Foundation (High Priority — Start Now)

**Goal:** Refactor the architecture so adding a new card type takes <1 hour of AI-agent time.

**Why first:** Every subsequent phase depends on this. Without a clean template system, Phase 2 becomes a mess of copy-paste code, and Phase 3 (WASM) requires a stable interface to compile against.

**Deliverables:**
1. **Registry pattern** — `card.Spec` with factory function; `main.go` has a single `registry` slice (one line per template)
2. **Shared utilities** — `pkg/layout/` (fonts, text, grids, tables), `pkg/qr/` (QR generation), `pkg/fonts/` (embedded default font via `go:embed`)
3. **Three card types** that validate the three core rendering primitives:
   - **WiFi QR card** — validates QR rasterization and scanning on real hardware
   - **Contact QR vCard** — validates complex QR payloads and demonstrates shareable value
   - **Keyboard Shortcuts Cheat Sheet** — validates text-table layout and proves the "reference deck" pattern
4. **Test scaffolding** — `card.AssertDimensions`, `card.AssertBlackWhite`
5. **Build automation** — Mage tasks for each new card

**Effort:** 3–4 days of AI-agent work

---

## Phase 2: Content Library (Medium Priority — Start After Phase 1)

**Goal:** Prove value by producing a diverse library of cards and establishing distribution formats.

**Why second:** You need working cards before you can share them or ask the `crosspoint-sync` maintainer for gallery integration.

**Deliverables:**
1. **8–12 additional card types** from the catalog (see below)
2. **`.deckpack` format spec** — ZIP of BMPs + `manifest.json` + `preview.png`
3. **3 starter packs** published as GitHub releases:
   - "2026 Calendar Pack" (landscape + portrait)
   - "Developer Reference Pack" (shortcuts for 3–4 tools)
   - "Travel Essentials Pack" (phrases, emergency numbers, boarding pass template)
4. **Gallery integration proposal** — a finished spec + preview assets to hand to the `crosspoint-sync` maintainer for a "Card Packs" tab

**Effort:** ~1 week of AI-agent work

---

## Phase 3: Community On-Ramp (Lower Priority — Start After Phase 2)

**Goal:** Open the tool to non-technical users who will never install Go.

**Why third:** A web configurator is high-impact but only valuable once there are enough card types and a proven pack format to configure. Premature investment here means configuring a tool that only makes calendars.

**Deliverables:**
1. **WASM web configurator** on GitHub Pages — compile the Go renderer to `GOOS=js GOARCH=wasm`; zero backend
2. **Static HTML form** for 2–3 most popular card types (calendar, business card, WiFi QR)
3. **Template exchange** via GitHub Discussions — submission guidelines, `.decktemplate` spec, community showcase
4. **Promotion** — "Show HN" post, Reddit r/eink, Mastodon/Bluesky with demo GIFs

**Effort:** 1–2 weeks of AI-agent work

---

## Full Card Type Catalog

| # | Card | Category | Complexity | Phase |
|---|---|---|---|---|
| 1 | Year-at-a-glance Calendar | Calendar | Simple | ✅ Done |
| 2 | **WiFi Guest Access QR** | Business | Simple | 1 |
| 3 | **Contact Card with QR vCard** | Business | Simple | 1 |
| 4 | **Keyboard Shortcuts Cheat Sheet** | Reference | Simple | 1 |
| 5 | Meeting Room Schedule | Business | Simple | 2 |
| 6 | Price List / Menu | Business | Simple | 2 |
| 7 | Boarding Pass Summary | Travel | Medium | 2 |
| 8 | Packing Checklist | Travel | Simple | 2 |
| 9 | Itinerary Summary | Travel | Simple | 2 |
| 10 | Emergency Info Card | Travel | Simple | 2 |
| 11 | Medication Schedule | Health | Simple | 2 |
| 12 | Workout / Gym Routine | Health | Simple | 2 |
| 13 | Measurement Log (with sparkline) | Health | Medium | 2 |
| 14 | Weekly Planner Grid | Productivity | Simple | 2 |
| 15 | Habit Tracker | Productivity | Simple | 2 |
| 16 | Kanban / Sprint Snapshot | Productivity | Medium | 2 |
| 17 | First Aid Quick Reference | Reference | Simple | 2 |
| 18 | Daily Quote / Word of the Day | Fun | Simple | 2 |
| 19 | Sudoku / Logic Puzzle | Fun | Medium | 2 |

---

## Key Decisions Locked In

| Decision | Choice | Rationale |
|---|---|---|
| **Template system** | Code-driven (Go packages) | AI agents write Go better than custom DSLs; full expressiveness for calendars, QR, tables |
| **Registry** | Manual slice in `main.go` | Explicit, debuggable, one line per template; no `init()` magic |
| **Shared utilities** | Free functions in `pkg/layout/` | Go-idiomatic, composable, no forced coupling via embedded structs |
| **Font strategy** | Embedded default via `go:embed` | Eliminates platform-specific font discovery; binary is self-contained |
| **Web configurator** | WASM compilation of Go renderer | Identical output to CLI; zero backend; native to existing codebase |
| **Distribution** | `.deckpack` ZIP + manifest.json | Simple for consumers; gallery-ready with PNG preview |

---

## Decisions You Still Need to Make

1. **Start Phase 1 now?** All later phases depend on it.
2. **Which card first in Phase 1?**
   - *WiFi QR* (simplest, safest QR test)
   - *Contact vCard* (highest demo value)
   - *Shortcuts* (no QR, tests tables)
3. **Target date for Phase 2 packs?** Needed before gallery outreach.
4. **Monetization path?** GitHub Sponsors is the low-friction default; labeled sponsored packs are optional.

---

*Generated: 2026-07-18*
*Expert panel: Product & UX (Oracle), Technical Architecture (Oracle), Community & Distribution (General)*
