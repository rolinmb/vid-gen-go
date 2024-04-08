package main

import (
    "image"
    "image/color"
    "math"
)

func nearestPixel(c color.Color) color.Color {
    r, g, b, _ := c.RGBA()
    return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
}

func pixelDelta(c1, c2 color.Color) [3]int32 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return [3]int32{int32(r1) - int32(r2), int32(g1) - int32(g2), int32(b1) - int32(b2)}
}

func distributeErr(img *image.RGBA, x,y int, quantError [3]int32, factor float64) {
    r, g, b, _ := img.At(x, y).RGBA()
    r += uint32(math.Round(float64(quantError[0]) * factor))
    g += uint32(math.Round(float64(quantError[1]) * factor))
    b += uint32(math.Round(float64(quantError[2]) * factor))
    img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
}

func fsDither(srcPng image.Image) *image.RGBA {
    bounds := srcPng.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    newPng := image.NewRGBA(bounds)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            newPng.Set(x, y, srcPng.At(x, y))
        }
    }
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            oldPxl := newPng.At(x, y)
            newPxl := nearestPixel(oldPxl)
            newPng.Set(x, y, newPxl)
            qErr := pixelDelta(oldPxl, newPxl)
            if x+1 < width {
                distributeErr(newPng, x+1, y, qErr, 7.0/16.0)
            }
            if x-1 >= 0 && y+1 < height {
                distributeErr(newPng, x-1, y+1, qErr, 3.0/16.0)
            }
            if y+1 < height {
                distributeErr(newPng, x, y+1, qErr, 5.0/16.0)
            }
            if x+1 < width && y+1 < height {
                distributeErr(newPng, x+1, y+1, qErr, 1.0/16.0)
            }
        }
    }
    return newPng
}