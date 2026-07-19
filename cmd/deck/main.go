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
	"github.com/dubyte/crosspoint-deck/pkg/templates/wifi"
)

var registry = []card.Spec{
	calendar.Spec(),
	wifi.Spec(),
	business.Spec(),
	cheatsheet.Spec(),
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
