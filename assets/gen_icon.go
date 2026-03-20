//go:build ignore

// gen_icon.go generates a 64x64 pixel-art cat face icon as icon.png.
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

	// Transparent background
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.Transparent)
		}
	}

	// Each "pixel" is 2x2 actual pixels (32x32 grid → 64x64 image)
	px := 2
	set := func(gx, gy int, c color.Color) {
		for dy := 0; dy < px; dy++ {
			for dx := 0; dx < px; dx++ {
				xx := gx*px + dx
				yy := gy*px + dy
				if xx >= 0 && xx < size && yy >= 0 && yy < size {
					img.Set(xx, yy, c)
				}
			}
		}
	}

	// Colors
	outline := color.RGBA{50, 50, 50, 255}     // dark outline
	fur := color.RGBA{255, 180, 80, 255}       // orange fur
	furDark := color.RGBA{220, 150, 60, 255}   // darker fur / stripes
	furLight := color.RGBA{255, 210, 140, 255} // lighter fur
	white := color.RGBA{255, 255, 255, 255}    // eye whites / muzzle
	eyeGreen := color.RGBA{80, 200, 100, 255}  // iris
	pupil := color.RGBA{30, 30, 30, 255}       // pupil
	nosePink := color.RGBA{255, 140, 150, 255} // nose
	innerEar := color.RGBA{255, 160, 160, 255} // inner ear pink
	mouth := color.RGBA{80, 60, 60, 255}       // mouth line

	// ---- Left ear (rows 2-8, cols 3-8) ----
	// Outline
	for _, p := range [][2]int{
		{4, 2}, {3, 3}, {3, 4}, {3, 5}, {3, 6}, {4, 7},
		{5, 2}, {5, 7},
		{6, 3}, {6, 7},
		{7, 4}, {7, 5}, {7, 6}, {7, 7},
	} {
		set(p[0], p[1], outline)
	}
	// Fill
	for _, p := range [][2]int{
		{4, 3}, {4, 4}, {4, 5}, {4, 6},
		{5, 3}, {5, 4}, {5, 5}, {5, 6},
		{6, 4}, {6, 5}, {6, 6},
	} {
		set(p[0], p[1], fur)
	}
	// Inner ear pink
	for _, p := range [][2]int{
		{4, 4}, {4, 5},
		{5, 4}, {5, 5},
	} {
		set(p[0], p[1], innerEar)
	}

	// ---- Right ear (rows 2-8, cols 23-28) ----
	for _, p := range [][2]int{
		{27, 2}, {28, 3}, {28, 4}, {28, 5}, {28, 6}, {27, 7},
		{26, 2}, {26, 7},
		{25, 3}, {25, 7},
		{24, 4}, {24, 5}, {24, 6}, {24, 7},
	} {
		set(p[0], p[1], outline)
	}
	for _, p := range [][2]int{
		{27, 3}, {27, 4}, {27, 5}, {27, 6},
		{26, 3}, {26, 4}, {26, 5}, {26, 6},
		{25, 4}, {25, 5}, {25, 6},
	} {
		set(p[0], p[1], fur)
	}
	for _, p := range [][2]int{
		{27, 4}, {27, 5},
		{26, 4}, {26, 5},
	} {
		set(p[0], p[1], innerEar)
	}

	// ---- Head outline (rows 7-26, cols 5-26) ----
	// Top
	for x := 8; x <= 23; x++ {
		set(x, 7, outline)
	}
	// Sides
	for y := 8; y <= 25; y++ {
		set(6, y, outline)
		set(25, y, outline)
	}
	// Corners
	set(7, 7, outline)
	set(7, 8, outline)
	set(24, 7, outline)
	set(24, 8, outline)
	// Bottom
	for x := 7; x <= 24; x++ {
		set(x, 26, outline)
	}
	// Bottom corners
	set(7, 25, outline)
	set(24, 25, outline)

	// ---- Head fill ----
	for y := 8; y <= 25; y++ {
		for x := 7; x <= 24; x++ {
			set(x, y, fur)
		}
	}
	// Top row fill
	for x := 8; x <= 23; x++ {
		set(x, 8, fur)
	}

	// ---- Forehead stripes ----
	for _, p := range [][2]int{
		{14, 8}, {15, 8}, {16, 8}, {17, 8},
		{15, 9}, {16, 9},
		{11, 9}, {12, 9},
		{19, 9}, {20, 9},
	} {
		set(p[0], p[1], furDark)
	}

	// ---- Eyes (row 13-16) ----
	// Left eye white
	for y := 13; y <= 16; y++ {
		for x := 9; x <= 13; x++ {
			set(x, y, white)
		}
	}
	// Right eye white
	for y := 13; y <= 16; y++ {
		for x := 18; x <= 22; x++ {
			set(x, y, white)
		}
	}
	// Left iris
	for _, p := range [][2]int{
		{10, 14}, {11, 14}, {12, 14},
		{10, 15}, {11, 15}, {12, 15},
	} {
		set(p[0], p[1], eyeGreen)
	}
	// Right iris
	for _, p := range [][2]int{
		{19, 14}, {20, 14}, {21, 14},
		{19, 15}, {20, 15}, {21, 15},
	} {
		set(p[0], p[1], eyeGreen)
	}
	// Left pupil
	for _, p := range [][2]int{{11, 14}, {11, 15}} {
		set(p[0], p[1], pupil)
	}
	// Right pupil
	for _, p := range [][2]int{{20, 14}, {20, 15}} {
		set(p[0], p[1], pupil)
	}
	// Eye highlights
	set(12, 14, white)
	set(21, 14, white)

	// ---- Nose (row 18-19) ----
	for _, p := range [][2]int{
		{15, 18}, {16, 18},
		{15, 19}, {16, 19},
	} {
		set(p[0], p[1], nosePink)
	}

	// ---- Muzzle / cheeks (row 17-22) ----
	// White muzzle area
	for y := 19; y <= 23; y++ {
		for x := 12; x <= 19; x++ {
			set(x, y, furLight)
		}
	}
	for _, p := range [][2]int{
		{13, 20}, {14, 20}, {15, 20}, {16, 20}, {17, 20}, {18, 20},
		{14, 21}, {15, 21}, {16, 21}, {17, 21},
	} {
		set(p[0], p[1], white)
	}

	// ---- Mouth ----
	set(15, 21, mouth)
	set(16, 21, mouth)
	set(14, 22, mouth)
	set(17, 22, mouth)

	// ---- Whiskers ----
	whiskerClr := color.RGBA{100, 80, 70, 200}
	// Left whiskers
	for _, p := range [][2]int{
		{7, 17}, {8, 18}, {7, 19}, {8, 20},
	} {
		set(p[0], p[1], whiskerClr)
	}
	// Right whiskers
	for _, p := range [][2]int{
		{24, 17}, {23, 18}, {24, 19}, {23, 20},
	} {
		set(p[0], p[1], whiskerClr)
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
