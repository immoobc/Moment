//go:build ignore

// gen_icon.go generates a 64x64 pixel-art peach icon as icon.png.
// Run: go run assets/gen_icon.go
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	const size = 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// transparent background
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.Transparent)
		}
	}

	// Colors
	stem := color.RGBA{90, 140, 60, 255}     // green stem
	leaf := color.RGBA{100, 180, 70, 255}    // green leaf
	leafD := color.RGBA{70, 150, 50, 255}    // darker leaf
	peachL := color.RGBA{255, 180, 130, 255} // light peach
	peachM := color.RGBA{255, 140, 100, 255} // mid peach
	peachD := color.RGBA{230, 110, 80, 255}  // dark peach
	blush := color.RGBA{255, 100, 100, 255}  // blush
	outline := color.RGBA{140, 70, 50, 255}  // brown outline
	highlight := color.RGBA{255, 220, 200, 255}

	// Helper: draw a filled rect of "pixels" (each pixel = 2x2 actual pixels for crispness at 64x64)
	px := 2
	set := func(gx, gy int, c color.Color) {
		for dy := 0; dy < px; dy++ {
			for dx := 0; dx < px; dx++ {
				xx := gx*px + dx
				yy := gy*px + dy
				if xx < size && yy < size {
					img.Set(xx, yy, c)
				}
			}
		}
	}

	// Stem (top center)
	for _, p := range [][2]int{{15, 4}, {15, 5}, {15, 6}, {16, 3}, {16, 4}, {16, 5}} {
		set(p[0], p[1], stem)
	}

	// Leaf (right of stem)
	for _, p := range [][2]int{
		{17, 4}, {18, 3}, {18, 4}, {19, 3}, {19, 4}, {20, 4},
		{17, 5}, {18, 5}, {19, 5},
	} {
		set(p[0], p[1], leaf)
	}
	for _, p := range [][2]int{{18, 4}, {19, 4}} {
		set(p[0], p[1], leafD)
	}

	// Peach body outline (roughly circular, rows 7-27, cols 8-24)
	// Top outline
	for x := 12; x <= 20; x++ {
		set(x, 7, outline)
	}
	// Left outline
	for y := 8; y <= 10; y++ {
		set(10, y, outline)
		set(11, y, outline)
	}
	for y := 11; y <= 24; y++ {
		set(8, y, outline)
		set(9, y, outline)
	}
	// Right outline
	for y := 8; y <= 10; y++ {
		set(21, y, outline)
		set(22, y, outline)
	}
	for y := 11; y <= 24; y++ {
		set(23, y, outline)
		set(24, y, outline)
	}
	// Bottom outline
	for x := 10; x <= 22; x++ {
		set(x, 25, outline)
		set(x, 26, outline)
	}
	// Bottom corners
	set(9, 24, outline)
	set(9, 25, outline)
	set(23, 24, outline)
	set(23, 25, outline)

	// Fill peach body — light
	for y := 8; y <= 24; y++ {
		xStart, xEnd := 10, 22
		if y <= 10 {
			xStart, xEnd = 12, 20
		}
		if y >= 24 {
			xStart, xEnd = 10, 22
		}
		for x := xStart; x <= xEnd; x++ {
			set(x, y, peachL)
		}
	}

	// Mid tone (lower half)
	for y := 17; y <= 24; y++ {
		for x := 10; x <= 22; x++ {
			set(x, y, peachM)
		}
	}

	// Dark shadow (bottom-left)
	for y := 20; y <= 24; y++ {
		for x := 10; x <= 14; x++ {
			set(x, y, peachD)
		}
	}

	// Blush (center)
	for _, p := range [][2]int{
		{14, 15}, {15, 15}, {16, 15},
		{14, 16}, {15, 16}, {16, 16},
		{15, 17},
	} {
		set(p[0], p[1], blush)
	}

	// Highlight (top-right)
	for _, p := range [][2]int{
		{18, 9}, {19, 9},
		{18, 10}, {19, 10},
		{19, 11},
	} {
		set(p[0], p[1], highlight)
	}

	// Peach crease (vertical center line)
	for y := 8; y <= 24; y++ {
		set(16, y, peachD)
	}

	f, err := os.Create("assets/icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
