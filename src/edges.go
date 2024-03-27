package main

import (
    "image"
    "image/color"
    "image/draw"
)

var sobelHoriz = [3][3]int{
    {-1, 0, 1},
    {-2, 0, 2},
    {-1, 0, 1},
}

var sobelVerti = [3][3]int{
    {-1, -2, -1},
    {0, 0, 0},
    {1, 2, 1},
}

func intSqrt(n int) int {
    var x,y int
    for x = n; x > 0; x = y {
        y = (x + n/x)
        if y >= x {
            return x
        }
    }
    return 0
}

func getEdges(inputPng *image.RGBA) *image.RGBA {
    bounds := inputPng.Bounds()
    grayPng := image.NewGray(bounds)
    draw.Draw(grayPng, grayPng.Bounds(), inputPng, image.Point{}, draw.Over)
    width,height := bounds.Dx(), bounds.Dy()
    gradX := image.NewGray(bounds)
    gradY := image.NewGray(bounds)
    for y := 1; y < height-1; y++ {
        for x := 1; x < width-1 ; x++ {
            var gx,gy int
            for i := 0; i < 3; i++ {
                for j := 0; j < 3; j++ {
                    px := int(grayPng.GrayAt(x+i-1, y+j-1).Y)
                    gx += px*sobelHoriz[i][j]
                    gy += px*sobelVerti[i][j]
                }
            }
            gradX.SetGray(x, y, color.Gray{uint8(gx)})
            gradY.SetGray(x, y, color.Gray{uint8(gy)})
        }
    }
    edges := image.NewRGBA(bounds)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            magX := int(gradX.GrayAt(x, y).Y)
            magY := int(gradY.GrayAt(x, y).Y)
            mag := uint8(intSqrt(magX*magX + magY*magY))
            edges.Set(x, y, color.RGBA{mag, mag, mag, 255})
        }
    }
    return edges
}