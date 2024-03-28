package main

import (
    "image"
    "image/color"
)

type Queue struct { //  FIFO queue
    items []image.Point
}

func (q *Queue) enqueue(item image.Point) { // adds an item to the end of the queue
    q.items = append(q.items, item)
}

func (q *Queue) dequeue() image.Point { // removes and returns the item at the front of the queue
    item := q.items[0]
    q.items = q.items[1:]
    return item
}

func (q *Queue) isEmpty() bool {
    return len(q.items) == 0
}

func applyWatershed(img image.Image) *image.RGBA {
    bounds := img.Bounds()
    result := image.NewRGBA(bounds)
    marker := image.NewGray(bounds)
    var q Queue
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            marker.SetGray(x, y, color.Gray{Y: 255})
            if x == bounds.Min.X || x == bounds.Max.X-1 || y == bounds.Min.Y || y == bounds.Max.Y-1 {
                q.enqueue(image.Point{X: x, Y: y})
                marker.SetGray(x, y, color.Gray{Y: 0})
            }
        }
    }
    for !q.isEmpty() {
        p := q.dequeue()
        for dy := -1; dy <= 1; dy++ {
            for dx := -1; dx <= 1; dx++ {
                x := p.X + dx
                y := p.Y + dy
                if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
                    if marker.GrayAt(x, y).Y == 255 {
                        q.enqueue(image.Point{X: x, Y: y})
                        marker.SetGray(x, y, color.Gray{Y: 0})
                    }
                }
            }
        }
    }
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            if marker.GrayAt(x, y).Y == 0 {
                result.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
            } else {
                result.Set(x, y, img.At(x, y))
            }
        }
    }
    return result
}