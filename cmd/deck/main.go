package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dubyte/crosspoint-deck/pkg/card"
	"github.com/dubyte/crosspoint-deck/pkg/render"
	"github.com/dubyte/crosspoint-deck/pkg/templates/business"
	"github.com/dubyte/crosspoint-deck/pkg/templates/calendar"
	"github.com/dubyte/crosspoint-deck/pkg/templates/cheatsheet"
	"github.com/dubyte/crosspoint-deck/pkg/templates/chore"
	"github.com/dubyte/crosspoint-deck/pkg/templates/coffee"
	"github.com/dubyte/crosspoint-deck/pkg/templates/convert"
	"github.com/dubyte/crosspoint-deck/pkg/templates/emergency"
	"github.com/dubyte/crosspoint-deck/pkg/templates/habit"
	"github.com/dubyte/crosspoint-deck/pkg/templates/library"
	"github.com/dubyte/crosspoint-deck/pkg/templates/loyalty"
	"github.com/dubyte/crosspoint-deck/pkg/templates/maintenance"
	"github.com/dubyte/crosspoint-deck/pkg/templates/meeting"
	"github.com/dubyte/crosspoint-deck/pkg/templates/morse"
	"github.com/dubyte/crosspoint-deck/pkg/templates/nato"
	"github.com/dubyte/crosspoint-deck/pkg/templates/packing"
	"github.com/dubyte/crosspoint-deck/pkg/templates/plant"
	"github.com/dubyte/crosspoint-deck/pkg/templates/recipe"
	"github.com/dubyte/crosspoint-deck/pkg/templates/resistor"
	"github.com/dubyte/crosspoint-deck/pkg/templates/shopping"
	"github.com/dubyte/crosspoint-deck/pkg/templates/stretch"
	"github.com/dubyte/crosspoint-deck/pkg/templates/timezones"
	"github.com/dubyte/crosspoint-deck/pkg/templates/wifi"
	"github.com/dubyte/crosspoint-deck/pkg/templates/workout"
)

var registry = []card.Spec{
	calendar.Spec(),
	wifi.Spec(),
	business.Spec(),
	cheatsheet.Spec(),
	chore.Spec(),
	coffee.Spec(),
	convert.Spec(),
	emergency.Spec(),
	habit.Spec(),
	library.Spec(),
	loyalty.Spec(),
	maintenance.Spec(),
	meeting.Spec(),
	morse.Spec(),
	nato.Spec(),
	packing.Spec(),
	plant.Spec(),
	recipe.Spec(),
	resistor.Spec(),
	shopping.Spec(),
	stretch.Spec(),
	timezones.Spec(),
	workout.Spec(),
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "help" || cmd == "-h" || cmd == "--help" {
		usage()
		os.Exit(0)
	}

	for _, spec := range registry {
		if spec.Name == cmd {
			runSpec(spec, os.Args[2:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
	usage()
	os.Exit(1)
}

func usage() {
	fmt.Print("Usage: deck <command> [flags]\n\nCommands:\n")
	for _, spec := range registry {
		fmt.Printf("  %-12s %s\n", spec.Name, spec.Usage)
	}
	fmt.Print(`
Examples:
  deck calendar --year 2026 --output ./output/calendar-2026.bmp
  deck wifi --ssid MyNetwork --password secret --output ./output/wifi.bmp
  deck business --name "John Doe" --phone "+1-555-0100" --output ./output/card.bmp
  deck cheatsheet --title "Vim" --output ./output/vim.bmp
`)
}

func runSpec(spec card.Spec, args []string) {
	fs := flag.NewFlagSet(spec.Name, flag.ExitOnError)
	out := fs.String("output", spec.Name+".bmp", "Output BMP file path")

	c := spec.New(fs)
	_ = fs.Parse(args)

	if err := os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: mkdir: %v\n", err)
		os.Exit(1)
	}

	if err := render.ToFile(c, *out); err != nil {
		fmt.Fprintf(os.Stderr, "error: render: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated: %s\n", *out)
}
