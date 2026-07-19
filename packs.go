//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/dubyte/crosspoint-deck/pkg/pack"
	"github.com/dubyte/crosspoint-deck/pkg/templates/calendar"
	"github.com/dubyte/crosspoint-deck/pkg/templates/cheatsheet"
	"github.com/dubyte/crosspoint-deck/pkg/templates/emergency"
	"github.com/dubyte/crosspoint-deck/pkg/templates/wifi"
	"github.com/magefile/mage/mg"
)

// PackCalendar generates the 2026 calendar starter pack.
func PackCalendar() error {
	fmt.Println("Building calendar starter pack...")
	os.MkdirAll("./packs", 0755)
	b := pack.NewBuilder("2026 Calendar Pack", "crosspoint-deck", "1.0.0").
		SetDescription("Year-at-a-glance calendar for 2026 in landscape and portrait orientations").
		SetLicense("MIT").
		AddTag("calendar").
		AddTag("2026")

	b.AddCard("2026 Calendar (Landscape)", &calendar.YearCard{Year: 2026, Portrait: false})
	b.AddCard("2026 Calendar (Portrait)", &calendar.YearCard{Year: 2026, Portrait: true})

	return b.BuildZip("./packs/2026-calendar.deckpack.zip")
}

// PackDeveloper generates the developer reference starter pack.
func PackDeveloper() error {
	fmt.Println("Building developer reference pack...")
	os.MkdirAll("./packs", 0755)
	b := pack.NewBuilder("Developer Reference Pack", "crosspoint-deck", "1.0.0").
		SetDescription("Essential keyboard shortcuts and cheat sheets for developers").
		SetLicense("MIT").
		AddTag("reference").
		AddTag("developer")

	b.AddCard("Vim Shortcuts", &cheatsheet.Card{
		Title: "Vim",
		Shortcuts: []cheatsheet.Shortcut{
			{Keys: "i", Description: "Insert mode"},
			{Keys: "Esc", Description: "Normal mode"},
			{Keys: ":w", Description: "Save"},
			{Keys: ":q", Description: "Quit"},
			{Keys: "dd", Description: "Delete line"},
			{Keys: "yy", Description: "Yank line"},
			{Keys: "p", Description: "Paste"},
			{Keys: "u", Description: "Undo"},
		},
	})

	b.AddCard("Git Shortcuts", &cheatsheet.Card{
		Title: "Git",
		Shortcuts: []cheatsheet.Shortcut{
			{Keys: "git status", Description: "Check status"},
			{Keys: "git add .", Description: "Stage all"},
			{Keys: "git commit", Description: "Commit"},
			{Keys: "git push", Description: "Push"},
			{Keys: "git pull", Description: "Pull"},
			{Keys: "git log", Description: "View log"},
			{Keys: "git diff", Description: "Show diff"},
			{Keys: "git branch", Description: "List branches"},
		},
	})

	return b.BuildZip("./packs/developer-reference.deckpack.zip")
}

// PackTravel generates the travel essentials starter pack.
func PackTravel() error {
	fmt.Println("Building travel essentials pack...")
	os.MkdirAll("./packs", 0755)
	b := pack.NewBuilder("Travel Essentials Pack", "crosspoint-deck", "1.0.0").
		SetDescription("Essential cards for travel: WiFi, emergency contacts, and packing checklists").
		SetLicense("MIT").
		AddTag("travel").
		AddTag("essentials")

	b.AddCard("WiFi Access", &wifi.Card{
		SSID:     "HotelWiFi",
		Password: "travel2026",
	})

	b.AddCard("Emergency USA", &emergency.Card{
		Country: "USA",
		Contacts: []emergency.EmergencyContact{
			{Label: "Police", Number: "911"},
			{Label: "Ambulance", Number: "911"},
			{Label: "Fire", Number: "911"},
		},
		BloodType: "O+",
	})

	return b.BuildZip("./packs/travel-essentials.deckpack.zip")
}

// Packs generates all starter packs.
func Packs() error {
	mg.Deps(PackCalendar, PackDeveloper, PackTravel)
	fmt.Println("All starter packs generated in ./packs/")
	return nil
}
