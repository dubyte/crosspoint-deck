package bmp

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

// Encode writes img to w as an uncompressed 24-bit Windows BMP.
// Preserves grayscale levels for the XTEink X4's 4-level grayscale display.
func Encode(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	const headerSize = 54
	rowSize := width * 3
	padding := (4 - rowSize%4) % 4
	pixelDataSize := (rowSize + padding) * height
	fileSize := headerSize + pixelDataSize

	// BITMAPFILEHEADER (14 bytes)
	header := make([]byte, headerSize)
	header[0] = 'B'
	header[1] = 'M'
	binary.LittleEndian.PutUint32(header[2:6], uint32(fileSize))
	binary.LittleEndian.PutUint32(header[10:14], uint32(headerSize))

	// BITMAPINFOHEADER (40 bytes)
	binary.LittleEndian.PutUint32(header[14:18], 40) // header size
	binary.LittleEndian.PutUint32(header[18:22], uint32(width))
	binary.LittleEndian.PutUint32(header[22:26], uint32(height))
	binary.LittleEndian.PutUint16(header[26:28], 1)  // planes
	binary.LittleEndian.PutUint16(header[28:30], 24) // bit count
	binary.LittleEndian.PutUint32(header[30:34], 0)  // compression (BI_RGB)
	binary.LittleEndian.PutUint32(header[34:38], uint32(pixelDataSize))
	binary.LittleEndian.PutUint32(header[38:42], 2835) // X pixels/meter (~72 DPI)
	binary.LittleEndian.PutUint32(header[42:46], 2835) // Y pixels/meter
	binary.LittleEndian.PutUint32(header[46:50], 0)    // colors used
	binary.LittleEndian.PutUint32(header[50:54], 0)    // important colors

	if _, err := w.Write(header); err != nil {
		return fmt.Errorf("bmp: write header: %w", err)
	}

	row := make([]byte, rowSize+padding)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			// Preserve actual pixel values for the X4's 4-level grayscale display.
			// The SSD1677 controller dithers 24-bit input to 2-bit (4 levels) natively.
			row[x*3+0] = byte(b >> 8)
			row[x*3+1] = byte(g >> 8)
			row[x*3+2] = byte(r >> 8)
		}
		if _, err := w.Write(row); err != nil {
			return fmt.Errorf("bmp: write row %d: %w", y, err)
		}
	}
	return nil
}

// IsBlackWhite reports whether every pixel in img is already pure black or white.
func IsBlackWhite(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				continue // transparent pixels are ignored
			}
			if r != 0 && r != 0xFFFF {
				return false
			}
			if g != 0 && g != 0xFFFF {
				return false
			}
			if b != 0 && b != 0xFFFF {
				return false
			}
		}
	}
	return true
}

// FillWhite returns a new 800x480 RGBA image filled with white.
func FillWhite() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 800, 480))
}

// Black returns the color black as image/color.Color.
func Black() color.Color { return color.Black }

// White returns the color white as image/color.Color.
func White() color.Color { return color.White }
