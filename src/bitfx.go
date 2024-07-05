package main

import (
	"image"
    "image/draw"
    "image/color"
)

func bitFx(pngSrc image.Image, delays []int, atts []float64) *image.RGBA {
    pngBounds := pngSrc.Bounds()
    newRgba := image.NewRGBA(pngBounds)
    draw.Draw(newRgba, pngBounds, pngSrc, pngBounds.Min, draw.Src)
    for y := pngBounds.Min.Y; y < pngBounds.Max.Y; y++ {
        for x := pngBounds.Min.X; x < pngBounds.Max.X ; x++ {
            var r,g,b,a float64
            for i := 0; i < len(delays); i++ {
                nx := x - delays[i]
                ny := y - delays[i]
                if nx >= pngBounds.Min.X && ny >= pngBounds.Min.Y {
                    c := newRgba.At(x, y)
                    cr,cg,cb,ca := c.RGBA()
                    aten := atts[i]
                    r += aten * float64(cr)
                    g += aten * float64(cg)
                    b += aten * float64(cb)
                    a += aten * float64(ca) 
                }
            }
            newRgba.Set(x, y, color.RGBA{
                R: uint8(r / 256),
                G: uint8(g / 256),
                B: uint8(b / 256),
                A: uint8(a / 256),
            })
        }
    }
    return newRgba
}
