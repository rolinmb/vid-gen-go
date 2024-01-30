Generating animations from .png images using golang and ffmpeg

TODO: In routineVideoFx(); need to detect source FPS of inVideoName so when running recombineCmd we rebuild the .mp4 at the original framerate

Run (from src): go build -o main && ./main

Example functions (so I can remove comments from getPixelColorOne and getPixelColorTwo):
'''
    // Hypotochoid
    angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))
    // Lissajous Curve
	angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y))
    // Offset Sine
	angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4
    // Sine Ripple
	angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) 
    // Diamond
	angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2)
    // Hyperbolic Spiral
	angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))
    // Exponential Decay
	angle := math.Pi * SCALE * math.Exp(-0.01 * math.Sqrt(float64(x*x+y*y)))
    // Circular Gradient
	angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y))
    // Square
    angle := math.Pi * SCALE * float64(x*x+y*y)
    // Others (custom w/o descriptions)
	angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
    angle := math.Pi * SCALE * ((math.Sin(float64(x)*0.1) + math.Sin(float64(y)*0.1)) / (1 + math.Sqrt(float64(x*x+y*y))))
	angle := math.Pi * 2.0 * SCALE * math.Sqrt(float64(x*x+y*y))
    angle := math.Pi * SCALE * (math.Cos(float64(x+y))*0.2)
    angle := math.Pi * SCALE * math.Tan(float64(x+(y/((x*y)+1))) - math.Sin(float64(x*y)*0.1))
    angle := math.Pi * SCALE * (math.Sin(float64(x)*0.05) * math.Exp(-float64(y)*0.1))
    angle := (math.Pi * SCALE * math.Sin(float64(x*x*y/5)))
    angle := math.Pi * SCALE * sin(sqrt(pow(x,2) + pow(y,2)))
    angle := math.Pi * SCALE *sin(0.1 * x) + cos(0.1 * y)
    /* Used for "redetach remixes" */
    /* 1 */
    centerX := float64(WIDTH) / 2
    centerY := float64(HEIGHT) / 2
    angle := math.Atan2(float64(y)-centerY, float64(x)-centerX)
    angle = angle + COMPL*math.Sin(angle)
    r := uint8((math.Sin(angle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    /* 2 */
    distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    distance = distance + COMPL*math.Sin(distance*0.1)
    r := uint8((math.Sin(distance*0.1) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(distance*0.1 + 2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(distance*0.1 + 4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    /* 3 */
    angleX := math.Pi * COMPL * (float64(x) - centerX) / centerX
    angleY := math.Pi * COMPL * (float64(y) - centerY) / centerY
    angle := math.Sin(angleX) + math.Cos(angleY)
    angle = angle + COMPL*math.Sin(angle*10)
    r := uint8((math.Sin(angle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(angle + 2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(angle + 4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    /* 4 */
    distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    angleX := math.Pi * COMPL * (float64(x) - centerX) / centerX
    angleY := math.Pi * COMPL * (float64(y) - centerY) / centerY
    combinedAngle := angleX + angleY + distance
    trippyAngle := combinedAngle + COMPL*math.Sin(combinedAngle)
    r := uint8((math.Sin(trippyAngle) + 1) * 127.5 * CLRFACTOR)
    g := uint8((math.Sin(trippyAngle + 2*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
    b := uint8((math.Sin(trippyAngle + 4*math.Pi/3) + 1) * 127.5 * CLRFACTOR)
'''
