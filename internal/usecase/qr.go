package usecase

import (
	"bytes"
	"errors"
	"image"
	"image/png"

	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

type QRGenerator struct {
	size int
}

func NewQRGenerator(size int) *QRGenerator {
	if size <= 0 {
		size = 256
	}
	return &QRGenerator{size: size}
}

func (g *QRGenerator) Size() int { return g.size }

// Generate creates a QR code PNG for the provided text. If iconData is provided it will be
// composited at the center of the QR code (scaled to ~20% of QR size).
func (g *QRGenerator) Generate(text string, iconData []byte) ([]byte, error) {
	if text == "" {
		return nil, errors.New("text is required")
	}

	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	img := qr.Image(g.size)

	if len(iconData) > 0 {
		iconImg, _, err := image.Decode(bytes.NewReader(iconData))
		if err != nil {
			return nil, err
		}
		// scale icon to 20% of QR size
		iconTargetSize := g.size / 5
		resized := image.NewRGBA(image.Rect(0, 0, iconTargetSize, iconTargetSize))
		draw.NearestNeighbor.Scale(resized, resized.Bounds(), iconImg, iconImg.Bounds(), draw.Over, nil)

		// composite centered
		out := image.NewRGBA(img.Bounds())
		draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
		center := image.Pt((g.size-iconTargetSize)/2, (g.size-iconTargetSize)/2)
		draw.Draw(out, image.Rectangle{Min: center, Max: center.Add(resized.Bounds().Size())}, resized, image.Point{}, draw.Over)
		img = out
	}

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
