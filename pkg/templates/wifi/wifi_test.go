package wifi

import (
	"flag"
	"testing"

	"github.com/dubyte/crosspoint-deck/pkg/card"
)

func TestWiFiCard_Render(t *testing.T) {
	c := &Card{
		SSID:     "TestNet",
		Password: "secret",
	}
	card.AssertDimensions(t, c, 800, 480)
}

func TestWiFiCard_Portrait(t *testing.T) {
	c := &Card{
		SSID:     "TestNet",
		Password: "secret",
		Portrait: true,
	}
	card.AssertDimensions(t, c, 480, 800)
}

func TestSpec(t *testing.T) {
	spec := Spec()
	if spec.Name != "wifi" {
		t.Errorf("name = %q, want wifi", spec.Name)
	}
	if spec.Usage == "" {
		t.Error("usage should not be empty")
	}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := spec.New(fs)
	if c == nil {
		t.Fatal("spec.New returned nil")
	}
	card.AssertDimensions(t, c, 800, 480)
}
