package main

import (
    "fmt"
    "log"
    "math"
    "strconv"
    "image"
    "image/color"
    "image/png"
    "os"
    "os/exec"
    "sync"
)

const (
    WIDTH = 1600
    HEIGHT = 1600
    dist_amp = 0.01
    dist_freq = 0.1
    dist_phase = 0.001
)
/*
func savePng(fname string, newPng *image.RGBA) {
    out, err := os.Create(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()
    err = png.Encode(out, newPng)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Successfully created/rewritten", fname)
}
*/
func savePngWg(wg *sync.WaitGroup, pngdir,fname string, newPng *image.RGBA) {
    defer wg.Done()
    _, err := os.Stat("png_out/"+pngdir)
    if err != nil {
        if os.IsNotExist(err) {
            mkdir_err := os.Mkdir("png_out/"+pngdir, 0755)
            if mkdir_err != nil {
                fmt.Println(mkdir_err)
            }
            fmt.Println("Created directory png_out/"+pngdir)
        } else {
            log.Fatal(err)
            return
        }
    }
    fout, err := os.Create(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer fout.Close()
    err = png.Encode(fout, newPng)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Successfully created/rewritten", fname)
}

func clamp(value,min,max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func getPixelColorOne(x,y int, SCALE,COMPL,CLRFACTOR float64) (uint8, uint8, uint8) {
    // angle := (math.Pi * SCALE * math.Sin(float64(x*x*y/5))) // Change to +/-/* or divide/modulus by x+1 or y+1
    // angle := math.Pi * SCALE * (math.Cos(float64(x+y))*0.2)
    // angle := math.Pi * SCALE * math.Tan(float64(x+(y/((x*y)+1))) - math.Sin(float64(x*y)*0.1))
    // angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y)) // circular gradient
    angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) // sine wave ripple
    // angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))  // hyperbolic spiral
    // angle := math.Pi * SCALE * float64(x*x+y*y)  // square
    // angle := math.Pi * SCALE * math.Exp(-0.01 * math.Sqrt(float64(x*x+y*y)))  // exponential decay
    // angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4  // offset sine wave
    // angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y))  // lissajous curve
    // angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
    // angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2) // diamond
    // angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))  // hypotochoid
    // angle := math.Pi * SCALE * (math.Sin(float64(x)*0.05) * math.Exp(-float64(y)*0.1))
    // Remix Set 1 Generators
    //  centerX := float64(WIDTH) / 2
    //  centerY := float64(HEIGHT) / 2
    /*angle := math.Atan2(float64(y)-centerY, float64(x)-centerX)
    angle = angle + COMPL*math.Sin(angle)
    r := uint8((math.Sin(angle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)*/
    /*distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    distance = distance + COMPL*math.Sin(distance*0.1)
    r := uint8((math.Sin(distance*0.1) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(distance*0.1+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(distance*0.1+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)*/
    /*angleX := math.Pi * COMPL * (float64(x) - centerX) / centerX
    angleY := math.Pi * COMPL * (float64(y) - centerY) / centerY
    angle := math.Sin(angleX) + math.Cos(angleY)
    angle = angle + COMPL*math.Sin(angle*10)
    r := uint8((math.Sin(angle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)*/
    /*distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    angleX := math.Pi * COMPL * (float64(x) - centerX) / centerX
    angleY := math.Pi * COMPL * (float64(y) - centerY) / centerY
    combinedAngle := angleX + angleY + distance
    trippyAngle := combinedAngle + COMPL*math.Sin(combinedAngle)
    r := uint8((math.Sin(trippyAngle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(trippyAngle+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(trippyAngle+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)*/
    // Main RGB functions
    distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * COMPL + frequency) * CLRFACTOR + 128)
    g := uint8(math.Sin(angle * COMPL + frequency + 2*math.Pi/3) * CLRFACTOR + 128)
    b := uint8(math.Sin(angle * COMPL + frequency + 4*math.Pi/3) * CLRFACTOR + 128)
    return r, g, b
}

func generateSimplePng(pngDir,fnameInc,fnameDec string, AMP,COMPL,CLRFACTOR,FREQ,PHASE,SCALE float64) {
    var wg sync.WaitGroup
    newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx, dy := distort(x, y, AMP, FREQ, PHASE)
            r, g, b := getPixelColorOne(dx, dy, SCALE, COMPL, CLRFACTOR)
            newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    wg.Add(2)
    go savePngWg(&wg, pngDir, fnameInc, newPng)
    go savePngWg(&wg, pngDir, fnameDec, newPng)
    wg.Wait()
}

func runFfmpegOne(pngDir,pngName,vidName string) {
    ffmpegCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/"+pngDir+"/"+pngName+"_%d.png",
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "vid_out/"+vidName+".mp4",
    )
    cmdOutput, err := ffmpegCmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Output:")
    fmt.Println(string(cmdOutput))
}

func simpleCleanup(pngdir,pngname string, FRAMES int) {
    for i := 1; i < FRAMES+1; i++ {
        err := os.Remove("png_out/"+pngdir+"/"+pngname+"_"+strconv.FormatInt(int64(((FRAMES*2)-1)-(i-1)), 10)+".png")
        if err != nil {
            log.Fatal(err)
        }
    }
}

func routineSimple(pngDir,pngName,vidName string, FRAMES int, AMP,FREQ,MULTIPLIER,PHASE,SCALE float64) {
    var fnameInc,fnameDec string
    for i := 1; i < FRAMES+1; i++ {
        fnameInc = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(i-1), 10)+".png"
        fnameDec = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(((FRAMES*2)-1)-(i-1)), 10)+".png"
        generateSimplePng(pngDir, fnameInc, fnameDec, AMP, MULTIPLIER*float64(i), float64(4*i-1), FREQ, PHASE, SCALE)
    }
    runFfmpegOne(pngDir, pngName, vidName)
    simpleCleanup(pngDir, pngName, FRAMES)
}

func getPixelColorTwo(x,y int, SCALE,COMPL,CLRFACTOR float64) (uint8, uint8, uint8) {
	// angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))  // hypotochoid
	// angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y)) // lissajous curve
	// angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4  // offset sine wave
	// angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) // sine wave ripple
	// angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2) // diamond
	// angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))  // hyperbolic spiral
	// angle := math.Pi * SCALE * math.Exp(-0.01 * math.Sqrt(float64(x*x+y*y)))  // exponential decay
	// angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y)) // circular gradient
	// angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
	angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y))
	/*// Remix Set 1 Generators
    centerX := float64(WIDTH) / 2
    centerY := float64(HEIGHT) / 2
    angle := math.Atan2(float64(x*y)-centerY, float64(x*x*y*y)-centerX)
    // angle = angle + COMPL*math.Sin(angle)
    r := uint8((math.Sin(angle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)*/
	distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * COMPL + frequency) * CLRFACTOR + 128)
    g := uint8(math.Sin(angle * COMPL + frequency + 2*math.Pi/3) * CLRFACTOR + 128)
    b := uint8(math.Sin(angle * COMPL + frequency + 4*math.Pi/3) * CLRFACTOR + 128)
    return r, g, b
}

func distort(x,y int, AMP,FREQ,PHASE float64) (int, int) {
    dx := y + int(AMP * math.Sin(FREQ*float64(x)+PHASE))
    dy := x + int(AMP * math.Sin(FREQ * float64(y) + PHASE))
	dx = clamp(dx, 0, WIDTH-1)
	dy = clamp(dy, 0, HEIGHT-1)
    return dx, dy
}

func generatePngInterp(
    iterUp,iterDown int,
    pngDir,pngName string,
    AMP1,AMP2,COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2,FREQ1,FREQ2,PHASE1,PHASE2,SCALE1,SCALE2,INTERPFACTOR float64) {
    var wg sync.WaitGroup
    newPngA := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    newPngB := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    finalPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx1, dy1 := distort(x, y, AMP1, FREQ1, PHASE1)
            dx2, dy2 := distort(x, y, AMP2, FREQ2, PHASE2)
            r1, g1, b1 := getPixelColorOne(dx1, dy1, SCALE1, COMPL1, CLRFACTOR1)
            r2, g2, b2 := getPixelColorTwo(dx2, dy2, SCALE2, COMPL2, CLRFACTOR2)
            newPngA.Set(x, y, color.RGBA{r1, g1, b1, 255})
            newPngB.Set(x, y, color.RGBA{r2, g2, b2, 255})
            r := uint8(float64(r1)*(float64(1)-INTERPFACTOR) + float64(r2)*INTERPFACTOR)
            g := uint8(float64(g1)*(float64(1)-INTERPFACTOR) + float64(g2)*INTERPFACTOR)
            b := uint8(float64(b1)*(float64(1)-INTERPFACTOR) + float64(b2)*INTERPFACTOR)
            finalPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    fNameA1 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterUp),10)+"a.png"
    fNameA2 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterDown),10)+"a.png"
    fNameB1 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterUp),10)+"b.png"
    fNameB2 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterDown),10)+"b.png"
    fInterp1 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterUp),10)+"i.png"
    fInterp2 := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(iterDown),10)+"i.png"
    wg.Add(6)
    go savePngWg(&wg, pngDir, fNameA1, newPngA)
    go savePngWg(&wg, pngDir, fNameA2, newPngA)
    go savePngWg(&wg, pngDir, fNameB1, newPngB)
    go savePngWg(&wg, pngDir, fNameB2, newPngB)
    go savePngWg(&wg, pngDir, fInterp1, finalPng)
    go savePngWg(&wg, pngDir, fInterp2, finalPng)
    wg.Wait()
}

func runFfmpegInterp(pngDir,pngName,vidName string) {
    var cmdOutput1,cmdOutput2,cmdOutput3 []byte
    var err1,err2,err3 error
    ffmpegCmdA := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/"+pngDir+"/"+pngName+"_%da.png",
        "-c:v", "libx264",  
        "-pix_fmt", "yuv420p",
        "vid_out/"+vidName+"a.mp4",
    )
    ffmpegCmdB := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/"+pngDir+"/"+pngName+"_%db.png",
        "-c:v", "libx264",  
        "-pix_fmt", "yuv420p",
        "vid_out/"+vidName+"b.mp4",
    )
    ffmpegCmdFinal := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/"+pngDir+"/"+pngName+"_%di.png",
        "-c:v", "libx264",  
        "-pix_fmt", "yuv420p",
        "vid_out/"+vidName+".mp4",
    )
    var wg sync.WaitGroup
    wg.Add(3)
    go func(){
        defer wg.Done()
        cmdOutput1, err1 = ffmpegCmdA.CombinedOutput()
        fmt.Println("ffmpegCmdA Output:\n", string(cmdOutput1))
    }()
    go func() {
        defer wg.Done()
        cmdOutput2, err2 = ffmpegCmdB.CombinedOutput()
        fmt.Println("ffmpegCmdB Output:\n", string(cmdOutput2))
    }()
    go func() {
        defer wg.Done()
        cmdOutput3, err3 = ffmpegCmdFinal.CombinedOutput()
        fmt.Println("ffmpegCmdFinal Output:\n", string(cmdOutput3))
    }()
    wg.Wait()
    if err1 != nil {
        log.Fatal("Error running ffmpegCmdA:", err1)
    }
    if err2 != nil {
        log.Fatal("Error running ffmpegCmdB:", err2)
    }
    if err3 != nil {
        log.Fatal("Error running ffmpegCmdFinal:", err3)
    }
}

func interpCleanup(pngDir,pngName string, FRAMES int) {
    var wg sync.WaitGroup
    var err error
    for i := 1; i < FRAMES+1; i++ {
        idx := int64(((FRAMES*2)-1)-(i-1))
        wg.Add(3)
        go func() {
            defer wg.Done()
            err = os.Remove("png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(idx, 10)+"a.png")
            if err != nil {
                log.Fatal(err)
            }
        }()
        go func() {
            defer wg.Done()
            err = os.Remove("png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(idx, 10)+"b.png")
            if err != nil {
                log.Fatal(err)
            }
        }()
        go func(){
            defer wg.Done()
            err = os.Remove("png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(idx, 10)+"i.png")
            if err != nil {
                log.Fatal(err)
            }
        }()
    }
    wg.Wait()
}

func routineInterp(
    pngDir,pngName,vidName string,
    FRAMES int,
    AMP1,AMP2,FREQ1,FREQ2,MULTIPLIER1,MULTIPLIER2,PHASE1,PHASE2,SCALE1,SCALE2,INTERPFACTOR float64) {
    var COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2 float64
	for i := 1; i < FRAMES+1; i++ {
        COMPL1 = MULTIPLIER1*float64(i)
        COMPL2 = MULTIPLIER2*float64(i)
        CLRFACTOR1 = float64(4*i-1)
        CLRFACTOR2 = float64(2*i+1)
        generatePngInterp(
            i-1,((FRAMES*2)-1)-(i-1),
            pngDir,pngName,
            AMP1,AMP2,COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2,FREQ1,FREQ2,PHASE1,PHASE2,SCALE1,SCALE2,INTERPFACTOR,
        )
    }
	runFfmpegInterp(pngDir, pngName, vidName)
    interpCleanup(pngDir, pngName, FRAMES)
}

func routineOverlay(fInName,fOutName string, cropWidth int) {
    fIn, err := os.Open(fInName) 
    if err != nil {
        log.Fatal(err)
        return
    }
    defer fIn.Close()
    img, _, err := image.Decode(fIn)
    if err != nil {
        log.Fatal(err)
        return
    }
    bounds := img.Bounds()
    width := bounds.Max.Y
    height := bounds.Max.X
    if cropWidth > width || cropWidth > height {
        log.Fatal("Crop size is larger than input .png dimensions")
	return
    }
    startX := (height - cropWidth) / 2
    startY := (width - cropWidth) / 2
    cropRect := image.Rect(startX, startY, startX+cropWidth, startY+cropWidth)
    croppedImg := img.(interface {
        SubImage(r image.Rectangle) image.Image
    }).SubImage(cropRect)
    fOut, err := os.Create(fOutName)
    if err != nil {
        log.Fatal(err)
	return
    }
    defer fOut.Close()
    err = png.Encode(fOut, croppedImg)
    if err != nil {
        log.Fatal(err)
	return
    }
    fmt.Printf("\nImage cropped and saved successfully to %s\n", fOutName)
}

func main() {
    /*routineSimple(
        "trial2",          // pngDir
        "trial2_10232023", // pngName
        "trial2_10232023", // vidName
        60,                // FRAMES
        0.222,             // AMP
        0.001,             // FREQ
        0.777,             // MULTIPLIER
        0.0,               // PHASE
        0.5,               // SCALE
    )*/
    /*routineInterp(
        "trial4",          // pngDir
        "trial4_10232023", // pngName
        "trial4_10232023", // vidName
        60,                // FRAMES
        0.333,               // AMP1
        0.777,               // AMP2
        0.01,              // FREQ1
        0.05,               // FREQ2
        2.0,               // MULTIPLIER1
        1.5,            // MULTIPLIER2
        0.001,             // PHASE1
        0.05,               // PHASE2
        0.5,               // SCALE1
        0.333,               // SCALE2
        0.5,               // INTERPFACTOR
    )*/
    routineOverlay("png_in/IMG_0520.png", "png_in/new.png", 1400)
}
