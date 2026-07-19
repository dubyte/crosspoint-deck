//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// Build compiles the deck CLI binary.
func Build() error {
	mg.Deps(Clean)
	fmt.Println("Building deck...")
	return sh.Run("go", "build", "-o", "deck", "./cmd/deck")
}

// Calendar generates a landscape year-at-a-glance calendar (800x480).
func Calendar() error {
	mg.Deps(Build)
	fmt.Println("Generating landscape calendar...")
	return sh.Run("./deck", "calendar", "--year", "2026", "--output", "./output/calendar-2026.bmp")
}

// CalendarPortrait generates a portrait calendar (480x800).
func CalendarPortrait() error {
	mg.Deps(Build)
	fmt.Println("Generating portrait calendar...")
	return sh.Run("./deck", "calendar", "--year", "2026", "--portrait", "--output", "./output/calendar-2026-portrait.bmp")
}

// WiFi generates a WiFi access card.
func WiFi() error {
	mg.Deps(Build)
	fmt.Println("Generating WiFi card...")
	return sh.Run("./deck", "wifi", "--ssid", "MyNetwork", "--password", "secret123", "--output", "./output/wifi.bmp")
}

// Business generates a business card.
func Business() error {
	mg.Deps(Build)
	fmt.Println("Generating business card...")
	return sh.Run("./deck", "business", "--name", "John Doe", "--title", "Developer", "--phone", "+1-555-0100", "--email", "john@example.com", "--output", "./output/business.bmp")
}

// Cheatsheet generates a shortcuts cheat sheet.
func Cheatsheet() error {
	mg.Deps(Build)
	fmt.Println("Generating cheat sheet...")
	return sh.Run("./deck", "cheatsheet", "--title", "Vim", "--shortcuts", "i:insert,Esc:normal,:w:save,:q:quit", "--output", "./output/cheatsheet.bmp")
}

// Meeting generates a meeting room schedule card.
func Meeting() error {
	mg.Deps(Build)
	fmt.Println("Generating meeting schedule...")
	return sh.Run("./deck", "meeting", "--room", "Boardroom", "--output", "./output/meeting.bmp")
}

// Packing generates a packing checklist.
func Packing() error {
	mg.Deps(Build)
	fmt.Println("Generating packing list...")
	return sh.Run("./deck", "packing", "--title", "Trip", "--items", "Passport,Phone,Charger,Camera,Snacks", "--output", "./output/packing.bmp")
}

// Emergency generates an emergency contact card.
func Emergency() error {
	mg.Deps(Build)
	fmt.Println("Generating emergency card...")
	return sh.Run("./deck", "emergency", "--country", "USA", "--contacts", "Police:911,Fire:911,Ambulance:911", "--blood", "O+", "--output", "./output/emergency.bmp")
}

// Habit generates a habit tracker.
func Habit() error {
	mg.Deps(Build)
	fmt.Println("Generating habit tracker...")
	return sh.Run("./deck", "habit", "--title", "Daily Habits", "--habits", "Read,Exercise,Meditate", "--days", "7", "--output", "./output/habit.bmp")
}

// All generates all card types.
func All() error {
	mg.Deps(Calendar, CalendarPortrait, WiFi, Business, Cheatsheet, Meeting, Packing, Emergency, Habit)
	fmt.Println("All cards generated.")
	return nil
}

// Verify runs go vet, tests, and checks the BMP output format.
func Verify() error {
	fmt.Println("Running go vet...")
	if err := sh.Run("go", "vet", "./..."); err != nil {
		return err
	}

	fmt.Println("Running tests...")
	if err := sh.Run("go", "test", "./..."); err != nil {
		return err
	}

	for _, f := range []string{
		"./output/calendar-2026.bmp",
		"./output/wifi.bmp",
		"./output/business.bmp",
		"./output/cheatsheet.bmp",
	} {
		if _, err := os.Stat(f); err == nil {
			fmt.Printf("Checking %s...\n", f)
			if err := sh.Run("file", f); err != nil {
				return err
			}
		}
	}

	return nil
}

// Test runs the Go test suite.
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "-v", "./...")
}

// Clean removes build artifacts and output files.
func Clean() error {
	fmt.Println("Cleaning...")
	_ = sh.Rm("deck")
	_ = sh.Rm("deck.exe")
	_ = os.RemoveAll("./output")
	return nil
}


