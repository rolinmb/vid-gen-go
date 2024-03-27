package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "math/rand"
    "log"
)

type Point []float64

type Cluster struct {
    Centroid Point
    Points []Point
}

func eucDist(a,b Point) float64 {
    sum := 0.0
    for i := range a {
        sum += math.Pow(a[i]-b[i], 2)
    }
    return math.Sqrt(sum)
}

func initCentroids(data []Point, k int) []Point {
    centroids := make([]Point, k)
    for i := range centroids {
        centroids[i] = data[rand.Intn(len(data))]
    }
    return centroids
}

func assignPoints(data []Point, centroids []Point) map[int][]Point {
    assignments := make(map[int][]Point)
    for _, point := range data {
        minDist := math.MaxFloat64
        var closestCluster int
        for i, centroid := range centroids {
            dist := eucDist(point, centroid)
            if dist < minDist {
                minDist = dist
                closestCluster = i
            }
        }
        assignments[closestCluster] = append(assignments[closestCluster], point)
    }
    return assignments
}

func updateCentroids(assignments map[int][]Point) []Point {
    centroids := make([]Point, len(assignments))
    for clusterIdx, points := range assignments {
        centroid := make(Point, len(points[0]))
        for _, point := range points {
            for i := range point {
                centroid[i] += point[i]
            }
        }
        for i := range centroid {
            centroid[i] /= float64(len(points))
        }
        centroids[clusterIdx] = centroid
    }
    return centroids
}

func kMeans(data []Point, k int) ([]Cluster, error) {
    if k <= 0 || k > len(data) {
        return nil, fmt.Errorf("kMeans(): Error, invalid number of clusters 'k' [k <= 0 or k > len(data)]")
    }
    centroids := initCentroids(data, k)
    for {
        assignments := assignPoints(data, centroids)
        newCentroids := updateCentroids(assignments)
        converged := true
        for i, centroid := range centroids {
            if eucDist(centroid, newCentroids[i]) > 0.0001 {
                converged = false
                break
            }
        }
        if converged {
            break
        }
        centroids = newCentroids
    }
    clusters := make([]Cluster, k)
    for i, centroid := range centroids {
        clusters[i].Centroid = centroid
        clusters[i].Points = assignPoints(data, centroids)[i]
    }
    return clusters, nil
}

func nearestNeigbor(point Point, centroids []Point) int {
    minDist := math.MaxFloat64
    var closestCluster int
    for i, centroid := range centroids {
        dist := eucDist(point, centroid)
        if dist < minDist {
            minDist = dist
            closestCluster = i
        }
    }
    return closestCluster
}

func applyKmeans(img image.Image, k int) *image.RGBA {
    bounds := img.Bounds()
    data := make([]Point, 0, bounds.Dx()*bounds.Dy())
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r,g,b,_ := img.At(x, y).RGBA()
            data = append(data, Point{float64(r), float64(g), float64(b)})
        }
    }
    clusters, err := kMeans(data, k)
    if err != nil {
        log.Fatalf("applyKmeans(): Error %v", err)
    }
    clusteredImg := image.NewRGBA(bounds)
    for _, cluster := range clusters {
        for _, point := range cluster.Points {
            for _, p := range cluster.Points {
                clusteredImg.Set(int(p[0]), int(p[1]), color.RGBA{uint8(point[0]), uint8(point[1]), uint8(point[2]), 255})
            }
        }
    }
    return clusteredImg
}

func imgToPixels(img image.Image) [][]float64 {
    bounds := img.Bounds()
    pixels := make([][]float64, 0, bounds.Dx()*bounds.Dy())
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            pixels = append(pixels, []float64{float64(r), float64(g), float64(b)})
        }
    }
    return pixels
}

func pixelsToImage(assignments []int, centers [][]float64, bounds image.Rectangle) *image.RGBA {
    img := image.NewRGBA(bounds)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            idx := assignments[y*bounds.Dx()+x]
            center := centers[idx]
            img.Set(x, y, color.RGBA{uint8(center[0]), uint8(center[1]), uint8(center[2]), 255})
        }
    }
    return img
}