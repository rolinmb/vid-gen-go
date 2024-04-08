package main

import (
    "image"
    "image/color"
    "math"
    "container/list"
)

var GFIRENEIGHBORS = [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

func makeBinaryImage(srcPng image.Image, tol uint8) [][]int {
    bounds := srcPng.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    binaryImg := make([][]int, height)
    for y := 0; y < height; y++ {
        binaryImg[y] = make([]int, width)
        for x := 0; x < width; x++ {
            r, g, b, _ := srcPng.At(x, y).RGBA()
            gray := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
            if gray >= tol {
                binaryImg[y][x] = 1
            } else {
                binaryImg[y][x] = 0
            }
        }
    }
    return binaryImg
}

func distMapToGrayscale(distMap [][]float64) image.Image {
    rows := len(distMap)
    cols := len(distMap[0])
    grayImg := image.NewGray(image.Rect(0, 0, cols, rows))
    for y := 0; y < rows; y++ {
        for x := 0; x < cols; x++ {
            grayImg.SetGray(x, y, color.Gray{ Y: uint8(distMap[y][x] * 255 / math.MaxFloat64) })
        }
    }
    return grayImg
}

func applyGrassfire(srcPng image.Image, tol uint8) image.Image {
    binaryImg := makeBinaryImage(srcPng, tol)
    rows := len(binaryImg)
    cols := len(binaryImg[0])
    distMap := make([][]float64, rows)
    for i := range distMap {
        distMap[i] = make([]float64, cols)
    }
    q := list.New()
    for i:= 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            if distMap[i][j] == 1 {
                distMap[i][j] = 0
                q.PushBack([2]int{i, j})
            } else {
                distMap[i][j] = math.Inf(1)
            }
        }
    }
    for q.Len() > 0 {
        front := q.Front()
        q.Remove(front)
        cur := front.Value.([2]int)
        curRow, curCol := cur[0], cur[1]
        for _, n := range GFIRENEIGHBORS {
            newRow, newCol := curRow+n[0], curCol+n[1]
            if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols && distMap[newRow][newCol] > distMap[curRow][curCol]+1 {
                distMap[newRow][newCol] = distMap[curRow][curCol] + 1
                q.PushBack([2]int{newRow, newCol})
            }
        }
    }
    return distMapToGrayscale(distMap)
}