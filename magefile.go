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

// CalendarWithFont generates a calendar using a custom font path.
func CalendarWithFont() error {
	mg.Deps(Build)
	font := findSystemFont()
	if font == "" {
		return fmt.Errorf("no system font found; use --font flag manually")
	}
	fmt.Printf("Generating calendar with font: %s\n", font)
	return sh.Run("./deck", "calendar", "--year", "2026", "--font", font, "--output", "./output/calendar-2026.bmp")
}

// All generates both landscape and portrait calendars.
func All() error {
	mg.Deps(Calendar, CalendarPortrait)
	fmt.Println("All calendars generated.")
	return nil
}

// Verify runs go vet and checks the BMP output format.
func Verify() error {
	fmt.Println("Running go vet...")
	if err := sh.Run("go", "vet", "./..."); err != nil {
		return err
	}

	if _, err := os.Stat("./output/calendar-2026.bmp"); err == nil {
		fmt.Println("Checking BMP format...")
		return sh.Run("file", "./output/calendar-2026.bmp")
	}

	fmt.Println("No BMP found to verify. Run 'mage calendar' first.")
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

// findSystemFont attempts to locate a suitable TTF on the host.
func findSystemFont() string {
	candidates := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/Windows/Fonts/arial.ttf",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
