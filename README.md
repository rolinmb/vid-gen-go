Generating animations from .png images using golang and ffmpeg

TODO:
    - src/main.go: In routineVideoFx() may need to detect source FPS of inVideoName so when running recombineCmd we rebuild the .mp4 at the original framerate
    - src/main.go && gui/app.py: Need to add new parameters from routineVideoFx() into routineVideoFxTk() and the python GUI app 

Build and run go program executable (from src): go build -o main && ./main

Build and run tkinter GUI program (from src): go build -o tkmain.exe && cd .. && python gui/main.py 

Example functions (so I can remove comments from getPixelColorOne and getPixelColorTwo):
    - angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))
    - angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y))
    - angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4
    - angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) 
    - angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2)
    - angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))
    - angle := math.Pi * SCALE * math.Exp(-0.01 * math.Sqrt(float64(x*x+y*y)))
    - angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y))
    - angle := math.Pi * SCALE * float64(x*x+y*y)
    - angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
    - angle := math.Pi * SCALE * ((math.Sin(float64(x)*0.1) + math.Sin(float64(y)*0.1)) / (1 + math.Sqrt(float64(x*x+y*y))))
    - angle := math.Pi * 2.0 * SCALE * math.Sqrt(float64(x*x+y*y))
    - angle := math.Pi * SCALE * (math.Cos(float64(x+y))*0.2)
    - angle := math.Pi * SCALE * math.Tan(float64(x+(y/((x*y)+1))) - math.Sin(float64(x*y)*0.1))
    - angle := math.Pi * SCALE * (math.Sin(float64(x)*0.05) * math.Exp(-float64(y)*0.1))
    - angle := (math.Pi * SCALE * math.Sin(float64(x*x*y/5)))
    - angle := math.Pi * SCALE * sin(sqrt(pow(x,2) + pow(y,2)))
    - angle := math.Pi * SCALE *sin(0.1 * x) + cos(0.1 * y)
