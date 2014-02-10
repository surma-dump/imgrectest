package main

import (
	"image"
	"image/color"
	"testing"
)

func TestOOBImage(t *testing.T) {
	img := image.NewGray16(image.Rect(0, 0, 3, 3))
	for i := 0; i < 9; i++ {
		img.Set(i%3, i/3, color.White)
	}
	oob := &OOBImage{
		Image:    img,
		OOBColor: color.Black,
	}

	if r, _, _, _ := oob.At(0, 0).RGBA(); r != 0xFFFF {
		t.Fatalf("Expected white, got 0x%x", r)
	}

	if r, _, _, _ := oob.At(-1, 0).RGBA(); r != 0x0000 {
		t.Fatalf("Expected black, got 0x%x", r)
	}
}

func TestInternalConvolute(t *testing.T) {
	img := image.NewGray16(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		img.Set(i%4, i/4, color.Gray16{uint16(i * 1000)})
	}
	matrix := [][]float64{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}

	oob := &OOBImage{
		Image:    img,
		OOBColor: color.Black,
	}

	if v := convolute(oob, matrix, image.Point{0, 0}); v != 5000/9 {
		t.Fatalf("Unexpected convolution result at (0, 0): %d", v)
	}

	if v := convolute(oob, matrix, image.Point{1, 1}); v != (0+2000+8000+10000)/9 {
		t.Fatalf("Unexpected convolution result at (1, 1): %d", v)
	}
}

func TestConvolute(t *testing.T) {
	img := image.NewGray16(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		img.Set(i%4, i/4, color.Gray16{uint16(i * 1000)})
	}
	matrix := [][]float64{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}

	cImg := Convolute(img, matrix)
	if r, _, _, _ := cImg.At(0, 0).RGBA(); r != 5000/9 {

		t.Fatalf("Unexpected convolution result at (0, 0): %d", r)
	}

	if r, _, _, _ := cImg.At(1, 1).RGBA(); r != (0+2000+8000+10000)/9 {
		t.Fatalf("Unexpected convolution result at (1, 1): %d", r)
	}
}
