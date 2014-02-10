package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		NumImages int           `goptions:"-n, description='Number of images to process'"`
		Path      string        `goptions:"-p, description='Path to data folder'"`
		Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		NumImages: 0,
		Path:      "./data",
	}
)

func main() {
	goptions.ParseAndFail(&options)
	imgs := readImages(options.Path)
	if options.NumImages != 0 {
		imgs = imgs[0:options.NumImages]
	}

	f, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Could not create result file: %s", err)
	}
	defer f.Close()

	b := imgs[0].Bounds()
	canvas := image.NewGray16(image.Rect(0, 0, 3*b.Dx(), len(imgs)*b.Dy()))
	for i, img := range imgs {
		draw.Draw(canvas, img.Bounds().Canon().Add(image.Point{0, i * b.Dy()}), img, image.Point{0, 0}, draw.Over)
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
