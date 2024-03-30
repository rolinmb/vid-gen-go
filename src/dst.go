package main

import (
    "image"
    "math"
)

func DST(block [][]float64) [][]float64 {
    n := len(block)
    dst := make([][]float64, n)
    for u := 0; u < n; u++ {
        dst[u] = make([]float64, n)
        for v := 0; v < n; v++ {
            var sum float64
            for i := 0; i < n; i++ {
                for j := 0; j < n; j++ {
                    sum += block[i][j] *
                        math.Sin((float64(i)+0.5)*float64(u)*math.Pi/float64(n)) *
                        math.Sin((float64(j)+0.5)*float64(v)*math.Pi/float64(n))
                }
            }
            dst[u][v] = sum * (2.0 / math.Sqrt(float64(n)))
        }
    }
    return dst
}

func IDST(dst [][]float64) [][]float64 {
    n := len(dst)
    block := make([][]float64, n)
    for i := 0; i < n; i++ {
        block[i] = make([]float64, n)
        for j := 0; j < n; j++ {
            var sum float64
            for u := 0; u < n; u++ {
                for v := 0; v < n; v++ {
                    sum += dst[u][v] *
                        math.Sin((float64(i)+0.5)*float64(u)*math.Pi/float64(n)) *
                        math.Sin((float64(j)+0.5)*float64(v)*math.Pi/float64(n))
                }
            }
            block[i][j] = sum * (2.0 / math.Sqrt(float64(n)))
        }
    }
    return block
}

func applyDst(img image.Image, blockSize int) image.Image {
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    dstImg := image.NewRGBA(bounds)
    for y := 0; y < height; y += blockSize {
        for x := 0; x < width; x += blockSize {
            block := extractBlock(img, x, y, blockSize)
            dstBlock := DST(block)
            storeBlock(dstImg, dstBlock, x, y, blockSize)
        }
    }
    return dstImg
}