package main

import (
    "image"
    "image/color"
)

func reduxResolution(pngSrc image.Image, nbits int) *image.RGBA {
    bounds := pngSrc.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    shift := 8 - nbits
    newPng := image.NewRGBA(bounds)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r, g, b, a := pngSrc.At(x, y).RGBA()
            r = (r >> shift) << shift
            g = (g >> shift) << shift
            b = (b >> shift) << shift
            a = (a >> shift) << shift
            newPng.Set(
                x,
                y,
                color.RGBA{
                    R: uint8(r >> 8),
                    G: uint8(g >> 8),
                    B: uint8(b >> 8),
                    A: uint8(a >> 8),
                },
            )
        }
    }
    return newPng
}