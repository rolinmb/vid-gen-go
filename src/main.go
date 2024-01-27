package main

import (
    "fmt"
    "log"
    "math"
    "go/parser"
    //"go/token"
    //"go/ast"
    //"reflect"
    "strings"
    "strconv"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
    "os/exec"
    "io/ioutil"
    "sync"
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
        y = (x + n/x) / 2
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
            mkdir_err := os.Mkdir("png_out/"+pngdir, 0700)
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

func distort(x,y,WIDTH,HEIGHT int, AMP,FREQ,PHASE float64) (int, int) {
    dx := y + int(AMP * math.Sin(FREQ * float64(x) + PHASE))
    dy := x + int(AMP * math.Sin(FREQ * float64(y) + PHASE))
	dx = clamp(dx, 0, WIDTH-1)
	dy = clamp(dy, 0, HEIGHT-1)
    return dx,dy
}

func getPixelColor(
    x,y,WIDTH,HEIGHT int,
    SCALE,COMPL,CLRFACTOR float64,
    expr interface{}) (uint8, uint8, uint8) {
    vars := map[string]int{"x": x, "y": y}
    result, err := evaluateASTNode(expr, vars)
    if err != nil {
        log.Fatal(err)
	    return uint8(0),uint8(0),uint8(0)
    }
    angle := math.Pi * SCALE * result
    distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
    frequency := distance * SCALE
    r := uint8(math.Sin(angle * COMPL + frequency) * CLRFACTOR + 128)
    g := uint8(math.Sin(angle * COMPL + frequency + 2*math.Pi/2) * CLRFACTOR + 128)
    b := uint8(math.Sin(angle * COMPL + frequency + 4*math.Pi/3 * CLRFACTOR + 128)) 
    return r,g,b
}

func generateSimplePng(
    pngDir,pngName string,
    i,WIDTH,HEIGHT,FRAMES int,
    AMP,COMPL,CLRFACTOR,FREQ,PHASE,SCALE float64,
    expr interface{}) {
    var wg sync.WaitGroup
    newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx,dy := distort(x, y, WIDTH, HEIGHT, AMP, FREQ, PHASE)
	        r,g,b := getPixelColor(dx, dy, WIDTH, HEIGHT, SCALE, COMPL, CLRFACTOR, expr)
	        newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    fnameInc := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(i-1), 10)+".png"
    fnameDec := "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(((FRAMES*2)-1)-(i-1)), 10)+".png"
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

func routineSimple(
    pngDir,pngName,vidName,EXPRESSION string,
    WIDTH,HEIGHT,FRAMES int,
    AMP,FREQ,PHASE,MULTIPLIER,SCALE float64) {
    expr, err := parser.ParseExpr(EXPRESSION)
    if err != nil {
        log.Fatal(err)
	return
    }
    for i := 1; i < FRAMES+1; i++ {
        generateSimplePng(
            pngDir, pngName,
            i, WIDTH, HEIGHT, FRAMES,
            AMP, MULTIPLIER*float64(i), float64(4*i-1), FREQ, PHASE, SCALE,
            expr,
        )
    }
    runFfmpegOne(pngDir, pngName, vidName)
    simpleCleanup(pngDir, pngName, FRAMES)
}

func generatePngInterp(
    pngDir,pngName string,
    i,WIDTH,HEIGHT,FRAMES int,
    AMP1,AMP2,COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2,FREQ1,FREQ2,PHASE1,PHASE2,SCALE1,SCALE2,INTERPFACTOR float64,
    expr1,expr2 interface{}) {
    var wg sync.WaitGroup
    newPngA := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    newPngB := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    finalPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    var iterUp,iterDown int
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx1,dy1 := distort(x, y, WIDTH, HEIGHT, AMP1, FREQ1, PHASE1)
            dx2,dy2 := distort(x, y, WIDTH, HEIGHT, AMP2, FREQ2, PHASE2)
            r1,g1,b1 := getPixelColor(dx1, dy1, WIDTH, HEIGHT, SCALE1, COMPL1, CLRFACTOR1, expr1)
            r2,g2,b2 := getPixelColor(dx2, dy2, WIDTH, HEIGHT, SCALE2, COMPL2, CLRFACTOR2, expr2)
            newPngA.Set(x, y, color.RGBA{r1, g1, b1, 255})
            newPngB.Set(x, y, color.RGBA{r2, g2, b2, 255})
            r := uint8(float64(r1)*(INTERPFACTOR) + float64(r2)*(1.0-INTERPFACTOR))
            g := uint8(float64(g1)*(INTERPFACTOR) + float64(g2)*(1.0-INTERPFACTOR))
            b := uint8(float64(b1)*(INTERPFACTOR) + float64(b2)*(1.0-INTERPFACTOR))
            finalPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    iterUp = i-1
    iterDown = ((FRAMES*2)-1)-(i-1)
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
	    wg.Wait()
    }
}

func routineInterp(
    pngDir,pngName,vidName,EXPRESSION1,EXPRESSION2 string,
    WIDTH,HEIGHT,FRAMES int,
    AMP1,AMP2,FREQ1,FREQ2,PHASE1,PHASE2,MULTIPLIER1,MULTIPLIER2,SCALE1,SCALE2,INTERPFACTOR float64) {
    var expr1,expr2 interface{}
    var err1,err2 error
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
    	expr1, err1 = parser.ParseExpr(EXPRESSION1)
    }()
    go func() {
        defer wg.Done()
	    expr2, err2 = parser.ParseExpr(EXPRESSION2)
    }()   
    wg.Wait()
    if err1 != nil {
        log.Fatal(err1)
	    return
    }
    if err2 != nil {
        log.Fatal(err2)
	    return
    }
    var COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2 float64
	for i := 1; i < FRAMES+1; i++ {
        COMPL1 = MULTIPLIER1*float64(i)
        COMPL2 = MULTIPLIER2*float64(i)
        CLRFACTOR1 = float64(4*i-1)
        CLRFACTOR2 = float64(2*i+1)
        generatePngInterp(
            pngDir, pngName,
            i, WIDTH, HEIGHT, FRAMES,
            AMP1, AMP2, COMPL1, COMPL2, CLRFACTOR1, CLRFACTOR2, FREQ1, FREQ2, PHASE1, PHASE2, SCALE1, SCALE2, INTERPFACTOR,
	        expr1, expr2,
    	)
    }
	runFfmpegInterp(pngDir, pngName, vidName)
    interpCleanup(pngDir, pngName, FRAMES)
}

func runFfmpegOverlay(pngDir,pngName,vidName string) {
    ffmpegCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/"+pngDir+"/"+pngName+"_%d.png",
        "-c:v", "libx264",  
        "-pix_fmt", "yuv420p",
        "vid_out/"+vidName+".mp4",
    )
    cmdOut, err := ffmpegCmd.CombinedOutput()
    fmt.Println("runFfmpegOverlay Output:\n", string(cmdOut))
    if err != nil {
        log.Fatal(err)
        return
    }
}

func overlayCleanup(pngDir,pngName string, FRAMES int) {
    for i := 1; i < (FRAMES+1); i++ {
        idx := int64(((FRAMES*2)-1)-(i-1))
        err := os.Remove("png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(idx, 10)+".png")
        if err != nil {
            log.Fatal(err)
            return
        }
    }
}

func routineOverlay(
    fInName,fOutName,pngDir,pngName,vidName,EXPRESSION1,EXPRESSION2 string,
    cropWidth,cropHeight,FRAMES int,
    AMP1,AMP1FACTOR,AMP2,AMP2FACTOR,FREQ1,FREQ2,PHASE1,PHASE2,MULTIPLIER1,MULTIPLIER2,SCALE1,SCALE1FACTOR,SCALE2,SCALE2FACTOR,INTERPFACTOR1,IF1AMP,IF1FREQ,INTERPFACTOR2,IF2AMP,IF2FREQ float64,
    IF1CONST,IF2CONST,edgeDetect bool) {
    if cropWidth < 1 || cropHeight < 1 {
        log.Fatalf("cropWidth (= %d) or cropHeight (= %d) cannot be negative or zero", cropWidth, cropHeight)
        return
    }
    if INTERPFACTOR1 < 0 || INTERPFACTOR2 < 0 {
        log.Fatalf("INTERPFACTOR1 (= %f) or INTERPFACTOR2 (= %f) cannot be negative", INTERPFACTOR1, INTERPFACTOR2)
        return
    }
    if FRAMES == 0 {
        log.Fatalf("FRAMES (= %d) cannot be zero", FRAMES)
        return
    }
    var expr1,expr2 interface{}
    var err1,err2 error
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()   
        expr1, err1 = parser.ParseExpr(EXPRESSION1)
    }()
    go func() {
        defer wg.Done()
        expr2, err2 = parser.ParseExpr(EXPRESSION2)
    }()
    wg.Wait()
    if err1 != nil {
        log.Fatal(err1)
    }
    if err2 != nil {
        log.Fatal(err2)
    }
    fIn, err := os.Open(fInName) 
    if err != nil {
        log.Fatal(err)
    }
    defer fIn.Close()
    srcPng, _, err := image.Decode(fIn)
    if err != nil {
        log.Fatal(err)
    }
    srcBounds := srcPng.Bounds()
    srcWidth := srcBounds.Max.X
    srcHeight := srcBounds.Max.Y
    if cropWidth > srcWidth || cropHeight > srcHeight {
        log.Fatalf("cropWidth (= %d) or cropHeight (= %d) is larger than input .png dimensions (%dx%d)", cropWidth, cropHeight, srcWidth, srcHeight)
    }
    /*startX := (srcWidth - cropWidth) / 2 // crop from center v1
    startY := (srcHeight - cropHeight) / 2
    cropRect := image.Rect(startX, startY, (startX+cropWidth) / 2, (startY+cropHeight) / 2)*/
    cropRect := image.Rect(0, 0, cropWidth, cropHeight)
    croppedPng := srcPng.(interface {
        SubImage(r image.Rectangle) image.Image
    }).SubImage(cropRect)
    if edgeDetect {
        croppedRGBA := image.NewRGBA(croppedPng.Bounds())
        draw.Draw(croppedRGBA, croppedPng.Bounds(), croppedPng, croppedPng.Bounds().Min, draw.Over)
        croppedPng = getEdges(croppedRGBA)
    }
    fOut, err := os.Create(fOutName)
    if err != nil {
        log.Fatal(err)
    }
    defer fOut.Close()
    err = png.Encode(fOut, croppedPng)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\n[Image cropped and saved successfully to %s]\n\n", fOutName)
    var fNameInc,fNameDec string
    var COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2,IF1,IF2,SCL1,SCL2 float64
    IF1 = INTERPFACTOR1
    IF2 = INTERPFACTOR2
    SCL1 = SCALE1
    SCL2 = SCALE2
    for i := 1; i < FRAMES + 1; i++ {
        pngResult := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))
        COMPL1 = MULTIPLIER1*float64(i)
        COMPL2 = MULTIPLIER2*float64(i)
        CLRFACTOR1 = float64(3*i-1)// + (2*math.Sin(float64(i)))
        CLRFACTOR2 = float64(2*i-1)// + (1*math.Cos(float64(i)))
        if !IF1CONST {
            IF1 += (0.5*IF1AMP*math.Sin(IF1FREQ*float64(i)))
        }
        if !IF2CONST {
            IF2 += (0.5*IF2AMP*math.Sin(IF2FREQ*float64(i)))
        }
        for x := 0; x < cropWidth; x++ {
            for y := 0; y < cropHeight; y++ {
                dx1,dy1 := distort(x, y, cropWidth, cropHeight, AMP1, FREQ1, PHASE1)
                dx2,dy2 := distort(x, y, cropWidth, cropHeight, AMP2, FREQ2, PHASE2)
                r1,g1,b1 := getPixelColor(dx1, dy1, cropWidth, cropHeight, SCL1, COMPL1, CLRFACTOR1, expr1)
                r2,g2,b2 := getPixelColor(dx2, dy2, cropWidth, cropHeight, SCL2, COMPL2, CLRFACTOR2, expr2)
                rI := uint8(float64(r1)*IF1 + float64(r2)*(1.0-IF1))
                gI := uint8(float64(g1)*IF1 + float64(g2)*(1.0-IF1))
                bI := uint8(float64(b1)*IF1 + float64(b2)*(1.0-IF1))   
                rC,gC,bC, _ := croppedPng.At(x, y).RGBA()
                r := uint8(float64(rC)*(1.0-IF2) + float64(rI)*IF2)
                g := uint8(float64(gC)*(1.0-IF2) + float64(gI)*IF2)
                b := uint8(float64(bC)*(1.0-IF2) + float64(bI)*IF2)
                pngResult.Set(x, y, color.RGBA{r, g, b, 255})
            }
        }
        fNameInc = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(i - 1),10)+".png"
        fNameDec = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(((FRAMES*2)-1)-(i-1)),10)+".png"
        wg.Add(2)
        go savePngWg(&wg, pngDir, fNameInc, pngResult)
        go savePngWg(&wg, pngDir, fNameDec, pngResult) 
        wg.Wait()
        AMP1 += AMP1FACTOR
        AMP2 += AMP2FACTOR
        SCALE1 *= SCALE1FACTOR
        SCALE2 *= SCALE2FACTOR
    }
    runFfmpegOverlay(pngDir, pngName, vidName)
    overlayCleanup(pngDir, pngName, FRAMES)
}

func routineVideoFx(
    inVidName,framesDir,outVidName,expressionR,multFnR,expressionG,multFnG,expressionB,multFnB string,
    scaleR,scaleAdjR,scaleG,scaleAdjG,scaleB,scaleAdjB,interpRatio,interpAdj float64,
    edgeDetect bool) {
    _, err := os.Stat(framesDir)
    if err != nil {
        if os.IsNotExist(err) {
        mkdir_err := os.Mkdir(framesDir, 0700)
        if mkdir_err != nil {
            fmt.Printf("\nroutineVideoFx(): Error occurred while creating 'src/%s'\n", framesDir, mkdir_err)
        }
        } else {
        log.Fatalf("routineVideoFx(): Error occured while checking if 'src/%s' exists: %v", framesDir, err)
        }
        fmt.Printf("\nroutineVideoFx(): Created directory 'src/%s'\n", framesDir)
    } else {
        os.RemoveAll(framesDir)
        os.MkdirAll(framesDir, 0700)
        fmt.Printf("\nroutineVideoFx(): Cleared contents of directory 'src/%s'\n", framesDir)
    }
    if _, err := os.Stat(inVidName); err != nil {
        log.Fatalf("routineVideoFx(): Error locating .mp4 video input 'src/%s': %v", inVidName, err)
    }
    teardownCmd := exec.Command(
        "ffmpeg", "-i", inVidName,
        "-vf", "fps=30", framesDir+"/"+outVidName+"_%03d.png",
    )
    teardownOut, err := teardownCmd.CombinedOutput()
    if err != nil {
        log.Fatalf("routineVideoFx(): Error occured while running teardownCmd: %v", err)
    }
    fmt.Printf("\nroutineVideoFx(): teardownCmd Output:\n\n%s\n(Successfully took apart 'src/%s' into individual frames)\n", string(teardownOut), inVidName)
    frameFiles, err := ioutil.ReadDir(framesDir)
    if err != nil {
        log.Fatalf("routineVideoFx(): Error occured while trying to read names of frame .pngs in 'src/%s': %v", framesDir, err)
    }
    var EXPR,EXPG,EXPB,MFNR,MFNG,MFNB interface {}
    var errR,errG,errB,errmR,errmG,errmB error
    var wg sync.WaitGroup
    wg.Add(6)
    go func() {
        defer wg.Done()
        EXPR, errR = parser.ParseExpr(expressionR)
    }()
    go func() {
        defer wg.Done()
        EXPG, errG = parser.ParseExpr(expressionG)
    }()
    go func() {
        defer wg.Done()
        EXPB, errB = parser.ParseExpr(expressionB)
    }()
    go func() {
        defer wg.Done()
        MFNR, errmR = parser.ParseExpr(multFnR)
    }()
    go func() {
        defer wg.Done()
        MFNG, errmG = parser.ParseExpr(multFnG)
    }()
    go func() {
        defer wg.Done()
        MFNB, errmB = parser.ParseExpr(multFnB)
    }()
    wg.Wait()
    if errR != nil {
        log.Fatal(errR)
    }
    if errG != nil {
        log.Fatal(errG)
    }
    if errB != nil {
        log.Fatal(errB)
    }
    if errmR != nil {
        log.Fatal(errR)
    }
    if errmG != nil {
        log.Fatal(errG)
    }
    if errmB != nil {
        log.Fatal(errB)
    }
    var ir float64
    if interpRatio < 0.0 {
        ir = 0.0
    } else if interpRatio > 1.0 {
        ir = 1.0
    } else {
        ir = interpRatio
    }
    adjFactor := interpAdj / float64(len(frameFiles))
    for _, pngFile := range frameFiles {
        rawFrame, err := os.Open(framesDir+"/"+pngFile.Name()) 
        if err != nil {
        log.Fatalf("routineVideoFx(): Error loading frame 'src/%s/%s' raw: %v", framesDir, pngFile.Name(), err)
        }
        framePng, _, err := image.Decode(rawFrame)
        if err != nil {
        log.Fatalf("routineVideoFx(): Error decoding raw frame 'src/%s/%s' with go/image: %v", framesDir, pngFile.Name(), err)
        }
        if edgeDetect {
            frameRGBA := image.NewRGBA(framePng.Bounds())
            draw.Draw(frameRGBA, framePng.Bounds(), framePng, framePng.Bounds().Min, draw.Over)
            framePng = getEdges(frameRGBA)
        }
        newPng := image.NewRGBA(image.Rect(0, 0, framePng.Bounds().Max.X, framePng.Bounds().Max.Y))
        for x := 0; x < framePng.Bounds().Max.X; x++ {
            for y := 0; y < framePng.Bounds().Max.Y; y++ {
                vars := map[string]int{ "x": x, "y": y }
                rt, err := evaluateASTNode(EXPR, vars)
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed Rgb expression: %v", err)
                }
                rval := uint8(rt)
                gt, err := evaluateASTNode(EXPG, vars) 
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed rGb expression: %v", err)
                }
                gval := uint8(gt)
                bt, err := evaluateASTNode(EXPB, vars)
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed rgB expression: %v", err)
                }
                bval := uint8(bt)
                multr, err := evaluateASTNode(MFNR, vars)
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed Rgb multiplier expression: %v", err)
                }
                multg, err := evaluateASTNode(MFNG, vars)
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed rGb multiplier expression: %v", err)
                }
                multb, err := evaluateASTNode(MFNB, vars)
                if err != nil {
                    log.Fatalf("routineVideoFx(): Error evaulating parsed rgB multiplier expression: %v", err)
                }
                rs, gs, bs, _ := framePng.At(x, y).RGBA()
                newPng.Set(x, y, color.RGBA{
                    uint8((ir*float64(rs) + (1.0-ir)*(scaleR*multr*float64(rval)))), 
                    uint8((ir*float64(gs) + (1.0-ir)*(scaleG*multg*float64(gval)))), 
                    uint8((ir*float64(bs) + (1.0-ir)*(scaleB*multb*float64(bval)))),
                    255,
                })
            }
        }
        segments := strings.Split(pngFile.Name(), "_")
        idxStr := strings.Replace(segments[len(segments)-1], ".png", "", -1)
        newFname := fmt.Sprintf("%s/%s_fx_%s.png", framesDir, outVidName, idxStr)
        newFrame, err := os.Create(newFname)
        if err != nil {
            log.Fatal("routineVideoFx(): Error creating 'src/%s/%s': %v", framesDir, newFname, err)
        }
        err = png.Encode(newFrame, newPng)
        if err != nil {
            log.Fatal("routineVideoFx(): Error encoding raw .png data to save to 'src/%s/%s': %v", framesDir, newFname, err)
        }
        newFrame.Close()
        //fmt.Printf("\nroutineVideoFx(): Successfully created FX'd frame 'src/%s/%s'\n", framesDir, newFname)
        rawFrame.Close()
        err = os.Remove(framesDir+"/"+pngFile.Name())
        if err != nil {
            log.Fatalf("routineVideoFx(): Failed to remove source frame 'src/%s/%s': %v", framesDir, pngFile.Name(), err)
        }
        ir += adjFactor
        if ir < 0.0 {
            ir = 0.0
        } else if  ir > (interpRatio + interpAdj) {
            ir = interpRatio + interpAdj
        }
        scaleR *= scaleAdjR
        scaleG *= scaleAdjG
        scaleB *= scaleAdjB
    }
    // Maybe need to get the framerate of vidInName so we can pass it to recombineCmd and thus the resulting outVidName.mp4 is of the same FPS.
    recombineCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", framesDir+"/"+outVidName+"_fx_%03d.png",
        "-c:v", "libx264",  
        "-pix_fmt", "yuv420p",
        "vid_out/"+outVidName+".mp4",
    )
    recombineOut, err := recombineCmd.CombinedOutput()
    if err != nil {
        log.Fatalf("routineVideoFx(): Error occured while running recombineCmd: %v", err)
    }
    fmt.Printf("\nroutineVideoFx(): recombineCmd Output:\n\n%s\n(Successfully created 'src/vid_out/%s.mp4')\n", string(recombineOut), outVidName)
}

func main() {
    /*fmt.Println("[main.go : routineSimple() started]")
    routineSimple(
        "trial1",          // pngDir
        "trial1_11262023", // pngName
        "trial1_11262023", // vidName
	    "x*x+y*y",    // EXPRESSION
        1000,  // WIDTH
        1000,  // HEIGHT
        60,    // FRAMES
        0.222, // AMP
        0.001, // FREQ
        0.0,   // PHASE
        0.777, // MULTIPLIER
        0.5,   // SCALE
    )
    fmt.Println("[main.go : routineInterp() started]")
    routineInterp(
        "trial2",          // pngDir
        "trial2_11262023", // pngName
        "trial2_11262023", // vidName
	    "cos(x)*sin(y)", // EXPRESSION1
	    "x*x+y*y", // EXPRESSION2
        1000,  // WIDTH
        1000,  // HEIGHT
        60,    // FRAMES
        0.333, // AMP1
        0.777, // AMP2
        0.01,  // FREQ1
        0.05,  // FREQ2
        0.001, // PHASE1
        0.05,  // PHASE2
        2.0,   // MULTIPLIER1
        1.5,   // MULTIPLIER2
        0.5,   // SCALE1
        0.333, // SCALE2
        0.5,   // INTERPFACTOR // (factor < 0.5 => less of EXPRESSION1; factor > 0.5 => more of EXPRESSION1)
    )*/
    /*fmt.Println("[main.go : routineOverlay() started]")
    routineOverlay(
        "png_in/IMG_0432.png",  // fInName
        "png_in/temp.png",       // fOutName
        "IMG_0432_1",          // pngDir
        "IMG_0432_1",          // pngName
        "IMG_0432_1",          // vidName
        "(x + y)*(x - y)",     // EXPRSSION1
	    "-1*((x*x*y + y*y)*(x*x - y*y*y))", // EXPRESSION2
        1512,  // cropWidth
	    1512,  // cropHeight
        100,    // FRAMES
        0.0,   // AMP1
        0.0,   // AMP1FACTOR
        0.0,   // AMP2
        0.0,   // AMP2FACTOR
        0.0,  // FREQ1
        0.0,  // FREQ2
        0.0,  // PHASE1
        0.0,  // PHASE2
        10.0,  // MULTIPLIER1
        10.0,  // MULTIPLIER2
        1000.0,   // SCALE1
        1.0,   // SCALE1FACTOR
        1000.0,   // SCALE2
        1.0,   // SCALE2FACTOR
        0.5,   // INTERPFACTOR1 (factor < 0.5 => less of EXPRESSION1; factor > 0.5 => more of EXPRESSION1)
        0.0,   // IF1AMP
        0.0,   // IF1FREQ
        0.1, // INTERPFACTOR2 (factor < 0.5 => more of fInName; factor > 0.5 => less of fInName)
        0.0,   // IF2AMP
        0.0,   // IF2FREQ
        true,  // IF1CONST
        true,  // IF2CONST
        false, // edgeDetect
    )
    */fmt.Println("[main.go : routineVideoFx() started]")
    routineVideoFx(
        "vid_in/driving_fog.mp4", // inVidName 
        "png_out/driving_fog2", // framesDir
        "driving_fog2", // outVidName
        "(x*x + y*y)", // expressionR
        "1.05 + 0.2*x", // multFnR 
        "(x*y)", // expressionG
        "1.05 + 0.2*y", // multFnG
        "(x*x - y*y)", // expressionB
        "1.05 + 0.1*x + 0.1*y", // multFnB
        1.0, // scaleR
        1.0, // scaleAdjR
        1.0, // scaleG
        1.0, // scaleAdjG
        1.0, // scaleB
        1.0, // scaleAdjB
        0.98, // interpRatio (ratio < 0.5 => less of inVidName; ratio > 0.5 => more of inVidName)
        -0.03, // interpAdj (value represents difference in interp ratio by final frame)
        false, // edgeDetect
    )
}
