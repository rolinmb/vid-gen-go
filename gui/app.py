import subprocess
import time
import os
import tkinter as tk

class App:
    def __init__(self, root):
        self.root = root
        self.root.title("vid-gen-go")
        self.root.geometry("720x720")
        self.vid_in_label = tk.Label(self.root, width=100, height=1, text="video input name (from src/vid_in)", fg="black")
        self.vid_in_label.grid(row=0)
        # TODO: make a selection list of video names in src/vid_in instead of having to enter exact vid name
        self.vid_in = tk.Text(self.root, width=100, height=1)
        self.vid_in.grid(row=1)
        self.vid_out_label = tk.Label(self.root, width=100, height=1, text="video output name (outputs to src/vid_out)", fg="black")
        self.vid_out_label.grid(row=2)
        self.vid_out = tk.Text(self.root, width=100, height=1)
        self.vid_out.grid(row=3)

    def run(self):
        self.root.mainloop()

if __name__ == "__main__":
    pass
