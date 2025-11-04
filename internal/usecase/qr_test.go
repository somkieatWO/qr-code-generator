package usecase

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestGenerate_NoIcon(t *testing.T) {
	g := NewQRGenerator(128)
	out, err := g.Generate("https://example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Fatalf("expected png bytes")
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode png: %v", err)
	}
}

func TestGenerate_WithIcon(t *testing.T) {
	iconPNG := createRedPNG(10, 10)
	g := NewQRGenerator(128)
	out, err := g.Generate("data with icon", iconPNG)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Fatalf("expected png bytes")
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode png: %v", err)
	}
}

func TestGenerate_EmptyText(t *testing.T) {
	g := NewQRGenerator(128)
	if _, err := g.Generate("", nil); err == nil {
		t.Fatalf("expected error for empty text")
	}
}

func createRedPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
