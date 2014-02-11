package main

import (
	"image"
	"image/color"
	"math"
)

type OOBImage struct {
	image.Image
	OOBColor color.Color
}

func (i *OOBImage) At(x, y int) color.Color {
	if !(image.Point{x, y}).In(i.Bounds()) {
		return i.OOBColor
	}
	return i.Image.At(x, y)
}

func Convolute(img image.Image, matrix [][]float64) image.Image {
	if len(matrix)%2 == 0 || len(matrix[0])%2 == 0 {
		panic("Matrix needs odd dimensions")
	}

	img = &OOBImage{
		Image:    img,
		OOBColor: color.Black,
	}

	b := img.Bounds().Canon()
	cImg := image.NewGray16(img.Bounds())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			cImg.Set(x, y, color.Gray16{convolute(img, matrix, image.Point{x, y})})
		}
	}
	return cImg
}

func convolute(img image.Image, matrix [][]float64, p image.Point) uint16 {
	s := float64(0)
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[0]); x++ {
			r, g, b, _ := img.At(p.X-(len(matrix)-1)/2+x, p.Y-(len(matrix)-1)/2+y).RGBA()
			s += float64(r+g+b) / 3 * matrix[y][x]
		}
	}
	if s < 0 {
		return 0
	} else if s > math.MaxUint16 {
		return math.MaxUint16
	}
	return uint16(s)
}
