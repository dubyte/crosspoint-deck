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

// Chore generates a chore chart card.
func Chore() error {
	mg.Deps(Build)
	fmt.Println("Generating chore chart...")
	return sh.Run("./deck", "chore", "--output", "./output/chore.bmp")
}

// Shopping generates a shopping list card.
func Shopping() error {
	mg.Deps(Build)
	fmt.Println("Generating shopping list...")
	return sh.Run("./deck", "shopping", "--output", "./output/shopping.bmp")
}

// Loyalty generates a loyalty cards list.
func Loyalty() error {
	mg.Deps(Build)
	fmt.Println("Generating loyalty cards...")
	return sh.Run("./deck", "loyalty", "--stores", "Airline:FF123456,Gym:MEM789,Library:LIB001", "--output", "./output/loyalty.bmp")
}

// Library generates a library card.
func Library() error {
	mg.Deps(Build)
	fmt.Println("Generating library card...")
	return sh.Run("./deck", "library", "--name", "Jane Reader", "--card-number", "29103000123456", "--branch", "Downtown", "--phone", "+1-555-0200", "--output", "./output/library.bmp")
}

// Nato generates a NATO phonetic alphabet reference card.
func Nato() error {
	mg.Deps(Build)
	fmt.Println("Generating NATO phonetic alphabet...")
	return sh.Run("./deck", "nato", "--output", "./output/nato.bmp")
}

// Owner generates an owner identification card.
func Owner() error {
	mg.Deps(Build)
	fmt.Println("Generating owner card...")
	return sh.Run("./deck", "owner", "--name", "John Doe", "--email", "john@example.com", "--output", "./output/owner.bmp")
}

// Morse generates a Morse code reference card.
func Morse() error {
	mg.Deps(Build)
	fmt.Println("Generating Morse code chart...")
	return sh.Run("./deck", "morse", "--output", "./output/morse.bmp")
}

// Resistor generates a resistor color code reference card.
func Resistor() error {
	mg.Deps(Build)
	fmt.Println("Generating resistor codes...")
	return sh.Run("./deck", "resistor", "--output", "./output/resistor.bmp")
}

// Convert generates a common conversions reference card.
func Convert() error {
	mg.Deps(Build)
	fmt.Println("Generating conversions chart...")
	return sh.Run("./deck", "convert", "--output", "./output/convert.bmp")
}

// Recipe generates a recipe card.
func Recipe() error {
	mg.Deps(Build)
	fmt.Println("Generating recipe card...")
	return sh.Run("./deck", "recipe", "--title", "Pasta Carbonara", "--ingredients", "Spaghetti,Eggs,Pancetta,Parmesan,Black pepper", "--steps", "Cook pasta,Fry pancetta,Mix eggs & cheese,Combine all,Serve hot", "--time", "20 min", "--servings", "2", "--output", "./output/recipe.bmp")
}

// Coffee generates a coffee brew guide card.
func Coffee() error {
	mg.Deps(Build)
	fmt.Println("Generating coffee brew guide...")
	return sh.Run("./deck", "coffee", "--output", "./output/coffee.bmp")
}

// Plant generates a plant care guide card.
func Plant() error {
	mg.Deps(Build)
	fmt.Println("Generating plant care guide...")
	return sh.Run("./deck", "plant", "--plant", "Monstera", "--water", "Weekly", "--light", "Indirect bright", "--humidity", "Moderate", "--food", "Monthly", "--output", "./output/plant.bmp")
}

// Workout generates a bodyweight workout card.
func Workout() error {
	mg.Deps(Build)
	fmt.Println("Generating workout card...")
	return sh.Run("./deck", "workout", "--output", "./output/workout.bmp")
}

// Stretch generates a stretching routine card.
func Stretch() error {
	mg.Deps(Build)
	fmt.Println("Generating stretching routine...")
	return sh.Run("./deck", "stretch", "--output", "./output/stretch.bmp")
}

// Timezones generates a world time zones reference card.
func Timezones() error {
	mg.Deps(Build)
	fmt.Println("Generating time zones card...")
	return sh.Run("./deck", "timezones", "--local", "New York EST", "--cities", "Tokyo:+14h,London:+5h,Paris:+6h,LA:-3h", "--output", "./output/timezones.bmp")
}

// Maintenance generates a home maintenance log card.
func Maintenance() error {
	mg.Deps(Build)
	fmt.Println("Generating maintenance log...")
	return sh.Run("./deck", "maintenance", "--output", "./output/maintenance.bmp")
}

// All generates all card types.
func All() error {
	mg.Deps(Calendar, CalendarPortrait, WiFi, Business, Cheatsheet, Meeting, Packing, Emergency, Habit,
		Chore, Shopping, Loyalty, Library, Nato, Morse, Resistor, Convert,
		Recipe, Coffee, Plant, Workout, Stretch, Timezones, Maintenance)
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
