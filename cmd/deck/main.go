package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dubyte/crosspoint-deck/pkg/render"
	"github.com/dubyte/crosspoint-deck/pkg/templates/calendar"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "calendar":
		calendarCmd(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`Usage: deck <command> [flags]

Commands:
  calendar    Generate a year-at-a-glance calendar card

Examples:
  deck calendar --year 2026 --output ./output/calendar-2026.bmp
  deck calendar --year 2026 --portrait --output ./output/calendar-2026-portrait.bmp
  deck calendar --year 2026 --font /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf --output ./output/calendar-2026.bmp
`)
}

func calendarCmd(args []string) {
	fs := flag.NewFlagSet("calendar", flag.ExitOnError)
	year := fs.Int("year", timeNowYear(), "Year to render")
	portrait := fs.Bool("portrait", false, "Render in portrait orientation (480x800)")
	out := fs.String("output", "calendar.bmp", "Output BMP file path")
	font := fs.String("font", "", "Path to TTF font (optional)")
	_ = fs.Parse(args)

	if err := os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: mkdir: %v\n", err)
		os.Exit(1)
	}

	card := &calendar.YearCard{
		Year:     *year,
		Portrait: *portrait,
		FontPath: *font,
	}

	if err := render.ToFile(card, *out); err != nil {
		fmt.Fprintf(os.Stderr, "error: render: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated: %s\n", *out)
}

func timeNowYear() int {
	return time.Now().Year()
}
