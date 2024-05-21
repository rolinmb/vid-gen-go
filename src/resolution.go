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

func min(a,b int) int {
    if a < b {
        return a
    }
    return b
}

func interpolate(v0,v1,r float64) float64 {
    return v0*(1-r) + v1*r
}

func interpolateColor(c00,c10,c01,c11 color.RGBA, dx,dy float64) color.RGBA {
    r := interpolate(interpolate(float64(c00.R), float64(c10.R), dx), interpolate(float64(c01.R), float64(c11.R), dx), dy)
    g := interpolate(interpolate(float64(c00.G), float64(c10.G), dx), interpolate(float64(c01.G), float64(c11.G), dx), dy)
    b := interpolate(interpolate(float64(c00.B), float64(c10.B), dx), interpolate(float64(c01.B), float64(c11.B), dx), dy)
    a := interpolate(interpolate(float64(c00.A), float64(c10.A), dx), interpolate(float64(c01.A), float64(c11.A), dx), dy)
    return color.RGBA{
        R: uint8(r),
        G: uint8(g),
        B: uint8(b),
        A: uint8(a),
    }
}

func upscaleResolution(src *image.RGBA, newWidth,newHeight int) *image.RGBA {
    srcBounds := src.Bounds()
    srcWidth := srcBounds.Dx()
    srcHeight := srcBounds.Dy()
    result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            gx := float64(x) * float64(srcWidth)/float64(newWidth)
            gy := float64(y) * float64(srcHeight)/float64(newHeight)
            gxi := int(gx)
            gyi := int(gy)
            c00 := src.At(gxi, gyi).(color.RGBA)
            c10 := src.At(min(gxi+1, srcWidth-1), gyi).(color.RGBA)
            c01 := src.At(gxi, min(gyi+1, srcHeight-1)).(color.RGBA)
            c11 := src.At(min(gxi+1, srcWidth-1), min(gyi+1, srcHeight-1)).(color.RGBA)
            result.set(x, y interpolateColor(c00, c10, c01, c11, gx-float64(gxi), gy-float64(gyi)))
        }
    }
    return result
}