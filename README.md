Generating animations and videos using golang and ffmpeg

TODO:

    - src/main.go: In routineVideoFx() may need to detect source FPS of inVideoName so when running recombineCmd we rebuild the .mp4 at the original framerate
    
    - src/main.go && gui/app.py: Need to add new parameters from routineVideoFx() into routineVideoFxTk() and the python GUI app 

Build and run go program executable (from src): go build -o main && ./main

Build and run tkinter GUI program (from src): go build -o tkmain.exe && cd .. && python gui/main.py
