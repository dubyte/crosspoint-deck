# CrossPoint Deck — Roadmap

This roadmap synthesizes findings from three expert analyses (Product & UX, Technical Architecture, Community & Distribution, Xerography & Print-to-E-Ink Transition) into actionable phases.

---

## Phase 1: Foundation ✅ COMPLETE

**Goal:** Refactor the architecture so adding a new card type takes <1 hour of AI-agent time.

**Delivered:**

1. ✅ **Registry pattern** — `card.Spec` with factory function; `main.go` has a single `registry` slice (one line per template)
2. ✅ **Shared utilities** — `pkg/layout/` (fonts, text, grids, tables, reversed headers, bold/regular helpers), `pkg/qr/` (QR generation)
3. ✅ **Twenty-four card types** covering five core rendering primitives:
   - WiFi QR, Business vCard — QR rasterization + scanning
   - Cheatsheet, Meeting, Packing, Emergency, Owner, Shopping, Chore — text layout + tables
   - Calendar — grid layout + variable-density content
   - Habit Tracker, Maintenance — checkbox grids
   - Coffee, Recipe, Plant, Workout, Stretch — structured reference cards
   - Library, Loyalty — identity cards
   - Nato, Morse, Periodic, Convert, Timezones — quick reference charts
4. ✅ **Design system** — Reversed black header bar + 2px divider + bold/regular typographic hierarchy on all cards
5. ✅ **Grayscale pipeline** — BMP encoder preserves anti-aliased edges; X4's SSD1677 controller dithers 24-bit to 4-level natively
6. ✅ **Build automation** — Mage tasks for each card + `mage All`

---

## Phase 2: Content Library (Active — In Progress)

**Goal:** Prove value by producing a diverse library of cards and establishing distribution formats.

**Delivered:**

1. ✅ **`.deckpack` format** — ZIP of BMPs + `manifest.json` + `preview.png`
2. ✅ **3 starter packs** — Calendar, Developer Reference, Travel Essentials
3. ✅ **24 cards done** covering all planned Phase 2 and most Phase 3 cards

**Remaining:**

1. **Gallery integration proposal** — spec + preview assets for `crosspoint-sync` maintainer
2. **2–4 niche card types** from Phase 4 catalog (crossword, sudoku, guitar chords, chess)

**Effort:** ~1–2 days remaining for core library

---

## Phase 3: The Laminated Life — Print-to-E-Ink Cards ✅ DONE

**Goal:** Cards that replace things people already print, laminate, and pin up.

**Xerography insight:** The best e-ink cards answer "what would I laminate?" — references that are always visible, never interactive, and replace the paper you already carry. These aren't "apps on a tiny screen"; they're the index cards, cheat sheets, and wall references that work because they're instant and mindless.

**Design principles for this phase:**

- **Always-on > interactive.** If it needs buttons, it's a phone app, not a card.
- **One glance = full answer.** No scrolling, no menus, no "next page."
- **Replace the laminated original.** If nobody laminates it today, question the card.
- **Personal data > generic content.** A recipe with *your* measurements beats a generic one.

**Deliverables:**

1. ✅ **Home & kitchen pack** — recipe card, chore chart, plant care guide, shopping list, coffee brew guide
2. ✅ **Identity pack** — loyalty cards, library card, owner card
3. ✅ **Reference pack** — NATO phonetic alphabet, morse code, periodic table, world time zones, common conversions
4. ✅ **Fitness pack** — bodyweight workout, stretching routine

**Remaining Phase 3 / Phase 4 ideas:**

- Run log, gym membership card
- First aid quick reference, baby feeding log, blood pressure log
- Crossword, sudoku, guitar chords, chess openings, constellation chart
- Daily quote card

---

## Phase 4: Community On-Ramp (Lower Priority)

**Goal:** Open the tool to non-technical users who will never install Go.

**Deliverables:**

1. **WASM web configurator** on GitHub Pages — compile Go renderer to WASM; zero backend
2. **Static HTML form** for 2–3 most popular card types (calendar, business card, WiFi QR)
3. **Template exchange** via GitHub Discussions — `.decktemplate` spec, community showcase
4. **Promotion** — "Show HN" post, Reddit r/eink, Mastodon/Bluesky with demo GIFs

**Effort:** 1–2 weeks of AI-agent work

---

## Full Card Type Catalog

### Phase 1 — Done ✅

| # | Card | Category | Why it works on e-ink |
|---|---|---|---|
| 1 | Year-at-a-glance Calendar | Calendar | Replaces wall calendar; always visible |
| 2 | WiFi Guest Access QR | Access | Guests scan without asking; QR is perfect for e-ink |
| 3 | Business Card with QR vCard | Identity | Share contact instantly; QR scans from e-ink |
| 4 | Keyboard Shortcuts Cheat Sheet | Reference | Replaces printed cheat sheet taped to monitor |
| 5 | Meeting Room Schedule | Office | Replaces paper schedule outside conference rooms |
| 6 | Packing Checklist | Travel | Replaces printed packing list; reuse trip after trip |
| 7 | Emergency Info Card | Safety | Laminated card replacement; always accessible |
| 8 | Habit Tracker Grid | Productivity | Replaces paper habit tracker on fridge/bulletin |

### Phase 2 — Done ✅

| # | Card | Category | Complexity | Why it works |
|---|---|---|---|---|
| 9 | Owner Identification | Identity | Simple | "If found" card; privacy-first with optional phone |
| 10 | Chore Chart | Home | Simple | Weekly rotation; checkboxes like habit tracker |
| 11 | Coffee Brew Guide | Kitchen | Simple | Ratios, temps, timers for each method |
| 12 | Plant Care Guide | Home | Simple | Water/light per plant; replaces sticky notes in pots |
| 13 | Recipe Card | Kitchen | Medium | Single recipe, big type; beats phone with wet hands |
| 14 | Shopping List Template | Kitchen | Simple | Reusable checklist; update before market runs |
| 15 | Loyalty Card Numbers | Identity | Simple | Barcode + member number; replaces keychain clutter |
| 16 | Library Card | Identity | Simple | Card number + branch info; always in your "wallet" |
| 17 | NATO Phonetic Alphabet | Reference | Simple | Alpha Bravo Charlie; classic laminated reference |
| 18 | Morse Code Chart | Reference | Simple | Dots and dashes; fits one card perfectly |
| 19 | Periodic Table | Reference | Simple | First 36 elements; science reference |
| 20 | World Time Zones | Reference | Medium | Major cities mapped to your local time |
| 21 | Common Conversions | Reference | Simple | Metric ↔ imperial; °F ↔ °C; laminated kitchen card |
| 22 | Bodyweight Workout | Fitness | Simple | No-equipment circuit; replaces gym poster |
| 23 | Stretching Routine | Fitness | Medium | Stick-figure diagrams + hold times |
| 24 | Home Maintenance Log | Home | Simple | Filter changes, battery swaps, seasonal tasks |

### Phase 3 — The Laminated Life ✅ DONE

All planned laminated-life cards are implemented. See Phase 2 table above.

### Phase 4 — Fun & Niche (Next)

| # | Card | Category | Complexity | Why it works |
|---|---|---|---|---|
| 30 | Daily Quote | Fun | Simple | Rotating inspiration; e-ink's always-on makes it work |
| 31 | Crossword Puzzle | Fun | Medium | Generated grid; solve with a pen on a screen protector |
| 32 | Sudoku | Fun | Medium | Always-ready puzzle; no app, no ads |
| 33 | Guitar Chord Chart | Music | Medium | Common chords in a grid; replaces printed chord sheet |
| 34 | Chess Openings | Games | Medium | Reference during study; like a printed opening book page |
| 35 | Constellation Chart | Reference | Medium | Seasonal sky map; dark-adapted eyes hate phone screens |
| 36 | First Aid Quick Reference | Safety | Simple | CPR steps, Heimlich; replaces Red Cross foldout |
| 37 | Baby Feeding Log | Health | Simple | Time + amount tracker; new parents live on paper logs |
| 38 | Blood Pressure Log | Health | Medium | Date + reading grid; doctor-visit ready |
| 39 | Price List / Menu | Business | Simple | Replaces laminated menu; update seasonally |
| 40 | Boarding Pass Summary | Travel | Medium | Flight details at a glance; no app needed |
| 41 | Itinerary Summary | Travel | Simple | Day-by-day plan; beats scrolling a PDF |
| 42 | Medication Schedule | Health | Simple | Replaces paper meds chart; critical info always visible |
| 43 | Weekly Planner Grid | Productivity | Simple | Time-blocking on paper, now on e-ink |



---

## Key Decisions Locked In

| Decision | Choice | Rationale |
|---|---|---|
| **Template system** | Code-driven (Go packages) | AI agents write Go better than custom DSLs; full expressiveness for calendars, QR, tables |
| **Registry** | Manual slice in `main.go` | Explicit, debuggable, one line per template; no `init()` magic |
| **Shared utilities** | Free functions in `pkg/layout/` | Go-idiomatic, composable, no forced coupling via embedded structs |
| **Design system** | Reversed header + bold labels | Unified deck identity across all cards; proven on e-ink contrast |
| **Grayscale strategy** | Hardware dithering (SSD1677) | Preserve anti-aliased edges; let device handle 4-level mapping |
| **Font strategy** | System font fallback chain | DejaVu → Liberation → Noto → Helvetica → Arial; bold variants supported |
| **Web configurator** | WASM compilation of Go renderer | Identical output to CLI; zero backend; native to existing codebase |
| **Distribution** | `.deckpack` ZIP + manifest.json | Simple for consumers; gallery-ready with PNG preview |

---

## Design Principles (from Xerography)

These principles guide which cards to build and how to design them:

1. **Always-on > interactive.** If it requires buttons, menus, or scrolling, it's a phone app, not a card. Cards work because they're mindless to use.
2. **One glance = full answer.** The user should absorb the card's entire value in under 3 seconds. If they need to study it, it's a document, not a card.
3. **Replace the laminated original.** Every good e-ink card has a physical ancestor: the gym card, the recipe card, the emergency contact sheet. If nobody laminates it today, question whether it needs to exist.
4. **Personal data over generic content.** A recipe with *your* measurements beats a generic one. A workout with *your* routine beats someone else's. Templates enable personalization; packs deliver it.
5. **B&W is a feature, not a limitation.** E-ink's constraint forces clarity. If a design needs color to work, it won't survive the medium. High-contrast B&W with bold/regular hierarchy is the palette.

---

## Decisions You Still Need to Make

1. **Gallery integration?** Reach out to `crosspoint-sync` maintainer with spec + preview assets for a "Card Packs" gallery section.
2. **Which Phase 4 cards first?** Daily quote + first aid are highest utility; crossword/sudoku are highest fun factor.
3. **Monetization path?** GitHub Sponsors is the low-friction default; labeled sponsored packs are optional.
4. **`.decktemplate` spec?** Needed before Phase 4 community exchange; simple JSON schema + SVG with `{{placeholders}}`.

---

*Updated: 2026-07-19*
*Status: 24 cards implemented. Phase 1–3 complete. Phase 4 (community on-ramp) and gallery integration remain.*
*Expert panel: Product & UX, Technical Architecture, Community & Distribution, Xerography & Print-to-E-Ink*
