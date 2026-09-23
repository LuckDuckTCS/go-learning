package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 4096, 4096
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	for py := 0; py < height; py++ {
		wg.Add(1)
		y := float64(py)/height*(ymax-ymin) + ymin
		go func(y float64, py int) {
			defer wg.Done()
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				img.Set(px, py, mandelbrot(z))
			}
		}(y, py)
	}
	wg.Wait()
	f, err := os.Create("out.png")
	if err != nil {
		fmt.Printf("create output file error: %v", err)
		os.Exit(2)
	}
	defer f.Close()
	fmt.Fprintln(os.Stderr, time.Since(start))
	png.Encode(f, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
