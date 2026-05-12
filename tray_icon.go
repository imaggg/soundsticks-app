//go:build linux || windows

package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

// makeTrayIcon draws a 32×32 dark-background speaker icon as PNG bytes.
func makeTrayIcon() []byte {
	const sz = 32
	img := image.NewNRGBA(image.Rect(0, 0, sz, sz))
	bg := color.NRGBA{0x11, 0x11, 0x13, 0xff}
	fg := color.NRGBA{0xf2, 0xf2, 0xf7, 0xff}

	for y := 0; y < sz; y++ {
		for x := 0; x < sz; x++ {
			img.SetNRGBA(x, y, bg)
		}
	}
	r := 5.0
	for y := 0; y < sz; y++ {
		for x := 0; x < sz; x++ {
			var dx, dy float64
			if x < int(r) {
				dx = r - float64(x)
			} else if x >= sz-int(r) {
				dx = float64(x) - float64(sz-1) + r
			}
			if y < int(r) {
				dy = r - float64(y)
			} else if y >= sz-int(r) {
				dy = float64(y) - float64(sz-1) + r
			}
			if dx > 0 && dy > 0 && math.Sqrt(dx*dx+dy*dy) > r {
				img.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0})
			}
		}
	}

	for y := 11; y <= 21; y++ {
		for x := 5; x <= 9; x++ {
			img.SetNRGBA(x, y, fg)
		}
	}
	for y := 9; y <= 23; y++ {
		t := float64(y-9) / 14.0
		xEnd := int(math.Round(10 + t*5))
		for x := 10; x <= xEnd; x++ {
			img.SetNRGBA(x, y, fg)
		}
	}

	arcAt(img, 16.0, 16.0, 3.5, fg)
	arcAt(img, 16.0, 16.0, 6.5, fg)

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func arcAt(img *image.NRGBA, cx, cy, r float64, c color.NRGBA) {
	const sz = 32
	for deg := -55.0; deg <= 55.0; deg += 0.4 {
		rad := deg * math.Pi / 180
		for dr := -0.6; dr <= 0.6; dr += 0.3 {
			x := int(math.Round(cx + (r+dr)*math.Cos(rad)))
			y := int(math.Round(cy + (r+dr)*math.Sin(rad)))
			if x >= 0 && x < sz && y >= 0 && y < sz {
				img.SetNRGBA(x, y, c)
			}
		}
	}
}
