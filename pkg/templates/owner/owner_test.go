package owner

import (
	"flag"
	"testing"

	"github.com/dubyte/crosspoint-deck/pkg/card"
)

func TestOwnerCard_Render(t *testing.T) {
	c := &Card{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "+1-555-0100",
	}
	card.AssertDimensions(t, c, 800, 480)
}

func TestOwnerCard_Portrait(t *testing.T) {
	c := &Card{
		Name:     "Test User",
		Email:    "test@example.com",
		Portrait: true,
	}
	card.AssertDimensions(t, c, 480, 800)
}

func TestSpec(t *testing.T) {
	spec := Spec()
	if spec.Name != "owner" {
		t.Errorf("name = %q, want owner", spec.Name)
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
