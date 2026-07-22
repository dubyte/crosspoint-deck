package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"image/png"
	_ "image/jpeg"
	"golang.org/x/image/bmp"

	"github.com/fogleman/gg"
)

func main() {
	inPath := flag.String("in", "", "Input BMP file")
	outPath := flag.String("out", "", "Output PNG file")
	flag.Parse()

	if *inPath == "" || *outPath == "" {
		fmt.Println("Usage: frame -in <input.bmp> -out <output.png>")
		os.Exit(1)
	}

	f, err := os.Open(*inPath)
	if err != nil {
		fmt.Println("Error opening input:", err)
		os.Exit(1)
	}
	defer f.Close()

	src, err := bmp.Decode(f)
	if err != nil {
		fmt.Println("Error decoding input (ensure it's BMP):", err)
		os.Exit(1)
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	bezel := 40.0
	outerR := 20.0
	shadow := 10.0

	canvasW := float64(w) + bezel*2 + shadow*2
	canvasH := float64(h) + bezel*2 + shadow*2

	dc := gg.NewContext(int(canvasW), int(canvasH))

	// transparent background
	dc.SetColor(color.Transparent)
	dc.Clear()

	// draw shadow
	dc.DrawRoundedRectangle(shadow, shadow, float64(w)+bezel*2, float64(h)+bezel*2, outerR)
	dc.SetColor(color.RGBA{0, 0, 0, 50})
	dc.Fill()

	// draw device bezel
	dc.DrawRoundedRectangle(shadow/2, shadow/2, float64(w)+bezel*2, float64(h)+bezel*2, outerR)
	dc.SetColor(color.RGBA{30, 30, 30, 255})
	dc.Fill()

	// draw screen
	dc.DrawImage(src, int(shadow/2+bezel), int(shadow/2+bezel))

	// draw screen border/shadow for depth
	dc.DrawRectangle(shadow/2+bezel, shadow/2+bezel, float64(w), float64(h))
	dc.SetColor(color.RGBA{0, 0, 0, 100})
	dc.SetLineWidth(2)
	dc.Stroke()

	outF, err := os.Create(*outPath)
	if err != nil {
		fmt.Println("Error creating output:", err)
		os.Exit(1)
	}
	defer outF.Close()

	if err := png.Encode(outF, dc.Image()); err != nil {
		fmt.Println("Error encoding PNG:", err)
		os.Exit(1)
	}

	fmt.Printf("Framed image saved to %s\n", *outPath)
}
