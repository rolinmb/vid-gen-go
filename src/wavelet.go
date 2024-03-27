package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func waveletTransform(img image.Image) ([][]float64, [][]float64, [][]float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	approximations := make([][]float64, 3)
	details := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		approximations[i] = make([]float64, width*height)
		details[i] = make([]float64, width*height)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			red := float64(r) / 65535.0
			green := float64(g) / 65535.0
			blue := float64(b) / 65535.0
			approximations[0][y*width+x] = red
			approximations[1][y*width+x] = green
			approximations[2][y*width+x] = blue
		}
	}
	for i := 0; i < 3; i++ {
		approximations[i], details[i] = waveletStep(approximations[i])
	}
	return approximations, details
}

func waveletStep(signal []float64) ([]float64, []float64) {
	n := len(signal)
	approximation := make([]float64, n)
	detail := make([]float64, n)
	for i := 0; i < n/2; i++ {
		approximation[i] = (signal[2*i] + signal[2*i+1]) / math.Sqrt2
		detail[i] = (signal[2*i] - signal[2*i+1]) / math.Sqrt2
	}
	return approximation, detail
}

func applyWaveletInverse(approximations, details [][]float64, width, height int) *image.RGBA {
	reconstructedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for c := 0; c < 3; c++ {
		reconstructedChannel := make([]float64, width*height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				reconstructedChannel[y*width+x] = approximations[c][y*width+x]
			}
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				reconstructedChannel[y*width+x] += details[c][y*width+x]
			}
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				value := reconstructedChannel[y*width+x]
				if value < 0 {
					value = 0
				} else if value > 1 {
					value = 1
				}
				reconstructedImg.Set(x, y, color.RGBA{
					R: uint8(255 * value),
					G: uint8(255 * value),
					B: uint8(255 * value),
					A: 255,
				})
			}
		}
	}
	return reconstructedImg
}