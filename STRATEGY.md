# CrossPoint Deck — Community & Distribution Strategy

> Role: Community & Distribution strategist  
> Date: 2026-07-18  
> Updated: 2026-07-19  
> Context: Early-stage Go CLI, MIT license, one maintainer (dubyte), targeting XTEink X4 via existing `crosspoint-sync` pipeline.

---

## 1. Sharing & Discovering Card Templates

Today, `crosspoint-deck` renders opaque `.bmp` files. For a community to form, we must separate **consumption** from **creation**:

| Artifact | Audience | Content | Editable? |
|---|---|---|---|
| **Card Pack** (`.deckpack`) | Consumers | Pre-rendered BMPs + manifest | No — use as-is |
| **Card Template** (`.decktemplate`) | Creators | SVG layout + JSON parameter schema | Yes — remix & regenerate |

**Template format (`.decktemplate`)**
A ZIP or directory containing:
- `template.json` — title, author, parameter schema (e.g., `name: string, year: number`), tags
- `template.svg` — visual layout with Mustache-style placeholders (`{{name}}`)
- `preview.png` — 400×240 gallery thumbnail
- `README.md` — usage instructions and parameter descriptions

The CLI renders it:
```bash
deck render --template business-card \
  --data '{"name":"Alice Chen","phone":"+1-555-0100"}' \
  --output ./output/
```

**Why templates matter:** A shared BMP of a business card is useless to anyone except the original owner. A shared template lets a hundred people generate their own card in seconds.

---

## 2. Card Pack Format (`.deckpack`)

For users who just want ready-made cards (calendars, cheat sheets, phrase books), a pack bundles pre-rendered assets.

**Pack structure:**
```
2026-minimal-calendar.deckpack.zip
├── manifest.json
├── preview.png
├── README.md
└── bmps/
    ├── january.bmp
    ├── february.bmp
    └── ...
```

**`manifest.json` schema:**
```json
{
  "name": "2026 Minimal Calendar",
  "author": "dubyte",
  "version": "1.0.0",
  "license": "MIT",
  "tags": ["calendar", "minimal", "landscape", "date-critical"],
  "cards": [
    {"file": "bmps/january.bmp", "title": "January 2026", "tags": ["january"]}
  ]
}
```

- Packs are ordinary ZIPs. No custom unpacker needed.
- Users download, unzip, and push the `/bmps/` folder via `crosspoint-sync` or SD-card copy.
- The manifest is for gallery display and future CLI helpers (`deck install <url>`).

---

## 3. Integration with crosspointsync.com Gallery

The existing gallery at `crosspointsync.com/gallery.html` hosts sleep backgrounds (PNG images submitted by users). It has categories: *All Devices*, *Popular*, *Newest*, *Trending*.

**Opportunity:** Propose a new **"Card Packs"** (or **"Utility Decks"**) category.

**What the gallery would display:**
- `preview.png` from the pack (rendered on-host, not on-device)
- Title, author, and tags from `manifest.json`
- Download link to the `.deckpack.zip`

**What does *not* need to change in `crosspoint-sync`:**
- The iOS app itself does not need to parse ZIPs or manifests. The gallery is a web surface; users download packs manually and upload BMP folders through the existing sync flow.
- If the maintainer wants deeper integration later (e.g., "Install Pack" → auto-unzip → queue upload), the manifest format is ready.

**Proposed collaboration path:**
1. Draft a one-page spec (`docs/deckpack-spec.md`) and share it with the `crosspoint-sync` maintainer.
2. Offer to provide gallery-ready preview assets for the first 3 packs.
3. If accepted, card packs appear on the same gallery page with a filter tab; if declined, host a simple static gallery on GitHub Pages and link from the README.

---

## 4. Web-Based Card Editor / Configurator

**Verdict: Yes. A static web configurator should be the primary on-ramp for non-technical users.**

The Go CLI is perfect for power users and CI, but it excludes ~90% of potential users who will never install Go or run `mage`. A web editor turns "useful for others" into **"usable by others."**

**Recommended architecture:**
- **GitHub Pages** hosts a static site at `crosspoint-reader.github.io/crosspoint-deck` (zero hosting cost, zero backend).
- **WebAssembly** compiles the existing Go render pipeline (`GOOS=js GOARCH=wasm`) so the browser generates BMPs using the *exact same code* as the CLI. Identical output guaranteed.
- **HTML/JS form** wraps the WASM binary: pick a template → fill parameters → click "Download BMP".

**MVP scope:**
- Two card types max: **Year Calendar** and **Business Card**.
- No server. No account. No database. The page is a ~2 MB WASM blob + a form.
- Include a "Download Pack" option for multi-card sets (e.g., all 12 months).

**Why WASM over Canvas:**
- Canvas `toDataURL()` produces PNG; we need uncompressed 24-bit BMP.
- More importantly, WASM reuses the existing Go rendering logic (fonts, QR codes, layout math). No second implementation to maintain.

---

## 5. Building Community

### GitHub (the engine)
- **Repo = registry.** `templates/` and `packs/` directories are the canonical index.
- **GitHub Discussions** for "Show and tell", "Template requests", and "Q&A". Early-stage projects feel alive when the maintainer responds quickly in Discussions.
- **Releases** page hosts compiled `.deckpack` files with attached preview images.
- **Issue templates:** "New template submission", "Pack request", "Bug in renderer".

### crosspointsync.com Gallery (the storefront)
- The gallery is the *discovery* layer. Users browse beautiful images and click download.
- Aim for a "Card Packs" tab. Until then, create a `gallery.md` in the repo that embeds pack previews and links to releases.

### Social & niche forums
- **Reddit:** r/eink, r/selfhosted — post short GIFs of card swiping on the X4. The e-ink community is small but passionate; a novel use case gets attention.
- **Hacker News:** "Show HN" once the web configurator is live. Lead with the problem: "My e-reader sits unused 23 hours a day. I turned it into a desk dashboard."
- **Mastodon/Bluesky:** Share #CrossPointDeck cards with screenshots. Tag e-ink and minimalist accounts.
- **Discord:** e-ink enthusiast servers (e.g., Remarkable, Boox communities) often have #projects channels.

### Content strategy
- **"Card of the week"** — a rendered BMP shared on social with a link to the template.
- **Demo GIFs** showing prev/next navigation through a folder of cards on the X4. Visual proof is everything for hardware projects.
- **Real use cases, not features:** "How I replaced my paper gym log with a card deck."

---

## 6. Monetization & Sustainability

**Core principle:** Do not paywall the tool or the packs. MIT license means the project itself stays free and open.

**Viable, low-friction options:**

| Approach | Friction | Notes |
|---|---|---|
| **GitHub Sponsors** | Very low | One-time or monthly donations from users who want to support development. The most honest model for a single-maintainer utility. |
| **Sponsored packs** | Low | A company or project sponsors a reference card (e.g., "Neovim Shortcuts, sponsored by Neovim" or "AWS CLI Cheatsheet"). Pack is still free; sponsor gets a small credit line in the manifest and on the gallery page. Must be clearly labeled. |
| **Affiliate / hardware referral** | Medium | If XTEink or a reseller offers an affiliate program, link from the README and web configurator. Users already own the device, so this only captures new buyers. |

**Avoid:** Subscription tiers, "Pro" templates, API keys, or anything that fragments the open ecosystem. The total addressable market (XTEink X4 owners) is small; monetization should fund maintenance, not build a business.

---

## Top 3 Actionable Recommendations (Ranked by Impact vs. Effort)

### 1. Define the `.deckpack` Format & Publish 3 Starter Packs
**Impact: High | Effort: Low (~2–3 days)**

Before building any tooling, prove that cards are worth sharing. Create the ZIP+BMP+manifest spec, generate three high-quality starter packs (e.g., **2026 Calendar**, **Developer Shortcuts**, **Travel Phrases**), and publish them as GitHub releases with PNG previews. Simultaneously, draft a one-page proposal for the `crosspoint-sync` maintainer to add a "Card Packs" section to the gallery.

*Why first:* This creates immediate, concrete value. It validates the format on real hardware. And it gives you something to promote on Reddit/HN instead of saying "here's a CLI you could use."

**Immediate action items:**
- [x] `.deckpack` format implemented in `pkg/pack/` (ZIP + manifest.json + preview.png)
- [x] 3 starter packs defined in `packs.go`: Calendar, Developer Reference, Travel Essentials
- [x] 24 card types implemented across all categories
- [ ] Create GitHub releases with attached `.deckpack.zip` files
- [ ] Open a friendly issue/PR on `zabirauf/crosspoint-sync` proposing gallery integration

---

### 2. Ship a WebAssembly Card Configurator on GitHub Pages
**Impact: Very High | Effort: Medium (~1–2 weeks)**

This is the single biggest barrier-removal for community growth. Compile the Go render pipeline to WASM and wrap it in a static HTML form. Start with just two card types: a **year calendar** and a **business card**.

*Why second:* It depends on having at least one solid template (from #1) to expose in the UI. It also takes longer, so parallelizing it after the pack format is defined avoids blocking.

**Immediate action items:**
- [ ] Add a `mage wasm` build target to `magefile.go`
- [ ] Create `web/` directory with a vanilla-JS form wrapper (no build step, no React)
- [ ] Add a `deck render --wasm` mode or separate `cmd/wasm` entry point
- [ ] Set up GitHub Actions to build and deploy `web/` to `gh-pages` on every tag
- [ ] Add a prominent "Create Card Online" badge to README.md

---

### 3. Establish the Template Exchange via GitHub Discussions
**Impact: Medium | Effort: Low (~1–2 days)**

Create a `templates/` directory, define the `.decktemplate` spec (JSON schema + SVG), and turn on GitHub Discussions with a "Show your cards" category. This creates a social loop: users submit templates via PR, you (or CI) render preview PNGs, and the best ones get promoted to official packs.

*Why third:* Templates are for creators, not consumers. They matter for long-term ecosystem health, but the immediate audience is smaller than pack consumers. Launch this after #1 and #2 are working so new template authors have a gallery and a web editor to point to.

**Immediate action items:**
- [ ] Create `templates/README.md` with submission guidelines and JSON schema
- [ ] Enable GitHub Discussions; set categories: "Show and tell", "Template requests", "Q&A"
- [ ] Add a PR template for new template submissions
- [ ] Write a `CONTRIBUTING.md` section explaining the difference between packs and templates

---

## Summary

The path from "personal CLI tool" to "community resource" is:
1. **Make shareable artifacts** (packs) that work with the existing sync pipeline.
2. **Lower the creation barrier** (web configurator) so non-technical users can participate.
3. **Build social infrastructure** (Discussions, template exchange) so creators have a venue and an incentive to contribute.

All three recommendations respect the hard constraints: no firmware changes, no competing mobile app, no over-engineering. They leverage what already exists (`crosspoint-sync` gallery, GitHub, Go's WASM target) and add only lightweight standards (ZIP + JSON manifest) on top.
