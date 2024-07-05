package main

import (
	"image"
	"image/png"
	"image/bmp"
	"io/ioutil"
	"log"
	"os"
    "strings"
)

const(
    HEADERSIZE = 54
)

func applyParams(img *image.RGBA, delays []int, atts []float64) {
    bounds := img.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X ; x++ {
            var r,g,b,a float64
            for i := 0; i < len(delays); i++ {
                nx := x - delays[i]
                ny := y - delays[i]
                if nx >= bounds.Min.X && ny >= bounds.Min.Y {
                    c := img.At(x, y)
                    cr,ch,cb,ca := c.RGBA()
                    aten := atts[i]
                    r += aten * float64(cr)
                    g += aten * float64(cg)
                    b += aten * float64(cb)
                    a += aten * float64(ca) 
                }
            }
            img.Set(x, y, color.RBGA{
                R: uint8(r / 256),
                G: uint8(g / 256),
                B: uint8(b / 256),
                A: uint8(a / 256),
            })
        }
    }
}

func applyBmpFx(pngSrc image.Image, pngFileName string, delays []int, atts []float64) *image.RBGA {
    bmpFileName := strings.Replace(".png", ".bpm")
    bmpSrc, err := os.Create(bpmFileName)
    if err != nil {
        log.Fatalf("applyBmpFx(): Error creating 'src/%s': %v", bmpFileName, err)
    }
    err = bmp.Encode(bmpSrc, pngSrc)
    if err != nil {
        log.Fatalf("applyBmpFx(): Error encoding 'src/%s' from 'src/%s': %v", bmpFileName, pngFileName, err)
    }
    bmpSrc.close()
    bmpData, err := ioutil.ReadFile(bmpFileName)
    if len(bmpData) < HEADERSIZE {
        log.Fatalf("applyBmpFx(): 'src/%s' is too small to contain a valid header", bmpFileName)
    }
    bmpDataBytes := bmpData[HEADERSIZE:]
    log.Printf("applyBmpFx(): First few bytes of 'src/%s': %v", bmpFileName, bmpDataBytes)
    pngBounds := pngSrc.Bounds()
    newBmpRgba := image.NewRGBA(pngBounds)
    copy(rgba.Pix, bmpDataBytes)
    applyParams(newBmpRgba, delays, atts)
    err = os.Remove(bmpFileName)
    if err != nil {
        log.Fatalf("applyBmpFx(): Error while removing 'src/%s': %v", bmpFileName, err)
    }
    return newBmpRgba
}
