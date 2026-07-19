package qr

import (
	"image"
	"image/color"

	"github.com/skip2/go-qrcode"
)

// Generate creates a QR code image for the given text.
// The size parameter specifies the desired width/height in pixels.
// Returns a black-on-white image.
func Generate(text string, size int) (image.Image, error) {
	code, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// go-qrcode returns an image.Image; we scale it to the target size
	img := code.Image(size)
	return img, nil
}

// GenerateWiFi creates a WiFi network QR code.
func GenerateWiFi(ssid, password, encryption string, size int) (image.Image, error) {
	if encryption == "" {
		encryption = "WPA"
	}
	text := "WIFI:S:" + ssid + ";T:" + encryption + ";P:" + password + ";;"
	return Generate(text, size)
}

// IsDark reports whether a pixel in a QR image is a dark module.
func IsDark(img image.Image, x, y int) bool {
	if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
		return false
	}
	c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
	return c.Y < 128
}
