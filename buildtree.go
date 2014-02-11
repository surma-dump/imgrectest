package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		Skip      int           `goptions:"-s, description='Skip a certain amount of images'"`
		NumImages int           `goptions:"-n, description='Number of images to process'"`
		Path      string        `goptions:"-p, description='Path to data folder'"`
		Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		Path: "./data",
	}
)

func main() {
	goptions.ParseAndFail(&options)
	imgs := readImages(options.Path)
	imgs = imgs[options.Skip:]
	if options.NumImages != 0 {
		imgs = imgs[0:options.NumImages]
	}

	f, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Could not create result file: %s", err)
	}
	defer f.Close()

	b := imgs[0].Bounds()
	canvas := image.NewGray16(image.Rect(0, 0, 4*b.Dx(), len(imgs)*b.Dy()))
	for i, img := range imgs {
		b = img.Bounds().Canon()
		draw.Draw(canvas, b.Add(image.Point{0, i * b.Dy()}), img, image.Point{0, 0}, draw.Over)

		xGradKernel := [][]float64{
			{-1, 0, 1},
		}
		yGradKernel := [][]float64{
			{1},
			{0},
			{-1},
		}
		xGrad, yGrad := Convolute(img, xGradKernel), Convolute(img, yGradKernel)
		draw.Draw(canvas, b.Add(image.Point{1 * b.Dx(), i * b.Dy()}), xGrad, image.Point{0, 0}, draw.Over)
		draw.Draw(canvas, b.Add(image.Point{2 * b.Dx(), i * b.Dy()}), yGrad, image.Point{0, 0}, draw.Over)
		draw.Draw(canvas, b.Add(image.Point{3 * b.Dx(), i * b.Dy()}), &DistanceImage{xGrad, yGrad}, image.Point{0, 0}, draw.Over)
	}
	if err := png.Encode(f, canvas); err != nil {
		log.Fatalf("Could not encode result: %s", err)
	}
}

func readImages(dir string) []image.Image {
	imgc := make(chan image.Image)
	images := make([]image.Image, 0)
	wg := &sync.WaitGroup{}

	go func() {
		for img := range imgc {
			images = append(images, img)
		}
	}()

	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			f, err := os.Open(path)
			if err != nil {
				log.Printf("Could not open %s: %s", path, err)
				return
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				log.Printf("Could not decode %s: %s", path, err)
				return
			}
			imgc <- img
		}(path)
		return nil
	})
	wg.Wait()
	close(imgc)
	return images
}

type DistanceImage struct {
	A, B image.Image
}

func (di *DistanceImage) At(x, y int) color.Color {
	cA, cB := di.A.At(x, y), di.B.At(x, y)
	rA, _, _, _ := cA.RGBA()
	rB, _, _, _ := cB.RGBA()

	return color.Gray16{uint16(math.Sqrt(float64(rA*rA) + float64(rB*rB)))}
}

func (di *DistanceImage) ColorModel() color.Model {
	return di.A.ColorModel()
}

func (di *DistanceImage) Bounds() image.Rectangle {
	return di.A.Bounds()
}
