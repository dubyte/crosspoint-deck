package card

import (
	"image"
	"testing"
)

// testCard is a minimal Card implementation for testing.
type testCard struct {
	w, h int
}

func (t *testCard) Render() image.Image {
	return image.NewRGBA(image.Rect(0, 0, t.w, t.h))
}

func TestAssertDimensions(t *testing.T) {
	c := &testCard{w: 800, h: 480}
	AssertDimensions(t, c, 800, 480)
}

func TestAssertDimensionsFail(t *testing.T) {
	c := &testCard{w: 100, h: 100}
	// Use a mock testing.T to capture failure
	mt := &mockT{}
	AssertDimensions(mt, c, 800, 480)
	if !mt.failed {
		t.Error("expected AssertDimensions to fail for wrong dimensions")
	}
}

type mockT struct {
	failed bool
}

func (m *mockT) Helper() {}
func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
}
func (m *mockT) Error(args ...interface{}) {
	m.failed = true
}
