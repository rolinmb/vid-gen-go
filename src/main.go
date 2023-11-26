package main

import (
    "fmt"
    "log"
    "math"
    "go/parser"
    "go/token"
    "go/ast"
    "reflect"
    "strconv"
    "image"
    "image/color"
    "image/png"
    "os"
    "os/exec"
    "sync"
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

func distort(x,y,WIDTH,HEIGHT int, AMP,FREQ,PHASE float64) (int, int) {
    dx := y + int(AMP * math.Sin(FREQ * float64(x) + PHASE))
    dy := x + int(AMP * math.Sin(FREQ * float64(y) + PHASE))
	dx = clamp(dx, 0, WIDTH-1)
	dy = clamp(dy, 0, HEIGHT-1)
    return dx,dy
}

func evaluateASTNode(node interface{}, vars map[string]int) (float64, error) {
    switch n := node.(type) {
    case *ast.BasicLit:
        if n.Kind == token.INT {
            val, err := strconv.Atoi(n.Value)
            if err != nil {
                return 0, err
            }
            return float64(val), nil
        } else if n.Kind == token.FLOAT {
            val, err := strconv.ParseFloat(n.Value, 64)
            if err != nil {
                return 0, err
            }
            return val, nil
        }
    case *ast.Ident:
        varName := n.Name
        if val, ok := vars[varName]; ok {
            return float64(val), nil
        }
        return 0, fmt.Errorf("Undefined variable: %s", varName)
    case *ast.CallExpr:
        funcName := n.Fun.(*ast.Ident).Name
        args := n.Args
        if funcName == "sin" && len(args) == 1 {
            argVal, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Sin(argVal), nil
        } else if funcName == "cos" && len(args) == 1 {
            argVal, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Cos(argVal), nil
        } else if funcName == "tan" && len(args) == 1 {
            argVal, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Tan(argVal), nil
        } else if funcName == "exp" && len(args) == 1 {
            argVal, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Exp(argVal), nil
        } else if funcName == "sqrt" && len(args) == 1 {
            argVal1, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Sqrt(argVal1), nil
        } else if funcName == "abs" && len(args) == 1 {
            argVal1, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            return math.Abs(argVal1), nil
        } else if funcName == "pow" && len(args) == 2 {
            argVal1, err := evaluateASTNode(args[0], vars)
            if err != nil {
                return 0, err
            }
            argVal2, err := evaluateASTNode(args[1], vars)
            if err != nil {
                return 0, err
            }
            return math.Pow(argVal1, argVal2), nil
        }
    case *ast.BinaryExpr:
        left, err := evaluateASTNode(n.X, vars)
        if err != nil {
            return 0, err
        }
        right, err := evaluateASTNode(n.Y, vars)
        if err != nil {
        return 0, err
        }
        switch n.Op {
        case token.ADD:
            return left + right, nil
        case token.SUB:
            return left - right, nil
        case token.MUL:
            return left * right, nil
        case token.QUO:
            return left / right, nil
        }
    case *ast.ParenExpr:
        return evaluateASTNode(n.X, vars)
    case *ast.UnaryExpr:
        operand, err := evaluateASTNode(n.X, vars)
        if err != nil {
        return 0, err
        }
        switch n.Op {
        case token.ADD:
            return operand, nil
        case token.SUB:
            return -operand, nil
        }
    }
    return 0, fmt.Errorf("Unsupported expression: %s", reflect.TypeOf(node))
}

func getPixelColorCustom(
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
 /*
func getPixelColorOne(x,y,WIDTH,HEIGHT int, SCALE,COMPL,CLRFACTOR float64) (uint8, uint8, uint8) {
    angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y))
    distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * COMPL + frequency) * CLRFACTOR + 128)
    g := uint8(math.Sin(angle * COMPL + frequency + 2*math.Pi/3) * CLRFACTOR + 128)
    b := uint8(math.Sin(angle * COMPL + frequency + 4*math.Pi/3) * CLRFACTOR + 128)
    return r,g,b
}
func getPixelColorTwo(x,y,WIDTH,HEIGHT int, SCALE,COMPL,CLRFACTOR float64) (uint8, uint8, uint8) {
    angle := math.Pi * SCALE * ((math.Sin(float64(x)*0.1) + math.Sin(float64(y)*0.1)) / (1 + math.Sqrt(float64(x*x+y*y))))
	distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * COMPL + frequency) * CLRFACTOR + 128)
    g := uint8(math.Sin(angle * COMPL + frequency + 2*math.Pi/3) * CLRFACTOR + 128)
    b := uint8(math.Sin(angle * COMPL + frequency + 4*math.Pi/3) * CLRFACTOR + 128)
    return r,g,b
}
*/
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
	        r,g,b := getPixelColorCustom(dx, dy, WIDTH, HEIGHT, SCALE, COMPL, CLRFACTOR, expr)
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
    AMP,FREQ,MULTIPLIER,PHASE,SCALE float64) {
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
            r1,g1,b1 := getPixelColorCustom(dx1, dy1, WIDTH, HEIGHT, SCALE1, COMPL1, CLRFACTOR1, expr1)
            r2,g2,b2 := getPixelColorCustom(dx2, dy2, WIDTH, HEIGHT, SCALE2, COMPL2, CLRFACTOR2, expr2)
            newPngA.Set(x, y, color.RGBA{r1, g1, b1, 255})
            newPngB.Set(x, y, color.RGBA{r2, g2, b2, 255})
            r := uint8(float64(r1)*(float64(1)-INTERPFACTOR) + float64(r2)*INTERPFACTOR)
            g := uint8(float64(g1)*(float64(1)-INTERPFACTOR) + float64(g2)*INTERPFACTOR)
            b := uint8(float64(b1)*(float64(1)-INTERPFACTOR) + float64(b2)*INTERPFACTOR)
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
    AMP1,AMP2,FREQ1,FREQ2,MULTIPLIER1,MULTIPLIER2,PHASE1,PHASE2,SCALE1,SCALE2,INTERPFACTOR float64) {
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
    AMP1,AMP2,FREQ1,FREQ2,MULTIPLIER1,MULTIPLIER2,PHASE1,PHASE2,SCALE1,SCALE2 float64) {
    if cropWidth < 1 || cropHeight < 1 {
        log.Fatalf("cropWidth (= %d) or cropHeight (= %d) cannot be negative or zero", cropWidth, cropHeight)
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
        return
    }
    if err2 != nil {
        log.Fatal(err2)
        return
    }
    fIn, err := os.Open(fInName) 
    if err != nil {
        log.Fatal(err)
        return
    }
    defer fIn.Close()
    srcPng, _, err := image.Decode(fIn)
    if err != nil {
        log.Fatal(err)
        return
    }
    srcBounds := srcPng.Bounds()
    srcWidth := srcBounds.Max.X
    srcHeight := srcBounds.Max.Y
    if cropWidth > srcWidth || cropHeight > srcHeight {
        log.Fatalf("cropWidth (= %d) or cropHeight (= %d) is larger than input .png dimensions (%dx%d)", cropWidth, cropHeight, srcWidth, srcHeight)
	    return
    }
    startX := (srcWidth - cropWidth) / 2
    startY := (srcHeight - cropHeight) / 2
    cropRect := image.Rect(startX, startY, startX+cropWidth, startY+cropHeight)
    croppedPng := srcPng.(interface {
        SubImage(r image.Rectangle) image.Image
    }).SubImage(cropRect)
    fOut, err := os.Create(fOutName)
    if err != nil {
        log.Fatal(err)
	    return
    }
    defer fOut.Close()
    err = png.Encode(fOut, croppedPng)
    if err != nil {
        log.Fatal(err)
	    return
    }
    fmt.Printf("\nImage cropped and saved successfully to %s\n\n", fOutName)
    var fNameInc,fNameDec string
    var COMPL1,COMPL2,CLRFACTOR1,CLRFACTOR2 float64
    for i := 1; i < FRAMES+1; i++ {
        COMPL1 = MULTIPLIER1*float64(i)
        COMPL2 = MULTIPLIER2*float64(i)
        CLRFACTOR1 = float64(3*i-1)
        CLRFACTOR2 = float64(2*i-1)
	    pngResult := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))
        for x := 0; x < cropWidth; x++ {
            for y := 0; y < cropHeight; y++ {
                dx1,dy1 := distort(x, y, cropWidth, cropHeight, AMP1, FREQ1, PHASE1)
                dx2,dy2 := distort(x, y, cropWidth, cropHeight, AMP2, FREQ2, PHASE2)
                r1,g1,b1 := getPixelColorCustom(dx1, dy1, cropWidth, cropHeight, SCALE1, COMPL1, CLRFACTOR1, expr1)
                r2,g2,b2 := getPixelColorCustom(dx2, dy2, cropWidth, cropHeight, SCALE2, COMPL2, CLRFACTOR2, expr2)
                rI := uint8(float64(r1)/2 + float64(r2)/2)
                gI := uint8(float64(g1)/2 + float64(g2)/2)
                bI := uint8(float64(b1)/2 + float64(b2)/2)   
                rC,gC,bC, _ := croppedPng.At(x, y).RGBA()
                r := uint8(float64(rI)/2 + float64(rC)/2)
                g := uint8(float64(gI)/2 + float64(gC)/2)
                b := uint8(float64(bI)/2 + float64(bC)/2)
                pngResult.Set(x, y, color.RGBA{r, g, b, 255})
            }
        }
        fNameInc = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(i - 1),10)+".png"
        fNameDec = "png_out/"+pngDir+"/"+pngName+"_"+strconv.FormatInt(int64(((FRAMES*2)-1)-(i-1)),10)+".png"
        wg.Add(2)
        go savePngWg(&wg, pngDir, fNameInc, pngResult)
        go savePngWg(&wg, pngDir, fNameDec, pngResult) 
        wg.Wait()
    }
    runFfmpegOverlay(pngDir, pngName, vidName)
    overlayCleanup(pngDir, pngName, FRAMES)
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
        0.777, // MULTIPLIER
        0.0,   // PHASE
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
        2.0,   // MULTIPLIER1
        1.5,   // MULTIPLIER2
        0.001, // PHASE1
        0.05,  // PHASE2
        0.5,   // SCALE1
        0.333, // SCALE2
        0.5,   // INTERPFACTOR
    )*/
    fmt.Println("[main.go : routineOverlay() started]")
    routineOverlay(
        "png_in/IMG_0520.png",  // fInName
        "png_in/new.png",       // fOutName
        "overlay0520_3",          // pngDir
        "overlay0520_3",          // pngName
        "overlay0520_3",          // vidName
        "abs(y-x) + abs(x-y)",     // EXPRSSION1
	    "x*x + y*y", // EXPRESSION2
        1512,  // cropWidth
	    1512,  // cropHeight
        60,    // FRAMES
        2.0,   // AMP1
        2.0,   // AMP2
        0.01,  // FREQ1
        0.01,  // FREQ2
        1.777,   // MULTIPLIER1
        10.0,   // MULTIPLIER2
        0.001,  // PHASE1
        0.001,  // PHASE2
        10.0,  // SCALE1
        10.0,  // SCALE2
    )
}
