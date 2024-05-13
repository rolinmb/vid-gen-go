import subprocess
import time
import os
import tkinter as tk

class App:
    def __init__(self, root):
        self.root = root
        self.root.title("vid-gen-go")
        self.root.geometry("900x1000")

        self.vid_in_label = tk.Label(self.root, width=100, height=1, text="Video Input Name (from src/vid_in)", fg="black")
        self.vid_in_label.grid(row=0)
        # TODO: make a selection list of video names in src/vid_in instead of having to enter exact vid name
        self.vid_in = tk.Text(self.root, width=100, height=1)
        self.vid_in.grid(row=1)

        self.frames_dir_label = tk.Label(self.root, width=100, height=1, text="Output Directory Name for Output Video Frames", fg="black")
        self.frames_dir_label.grid(row=2)
        self.frames_dir = tk.Text(self.root, width=100, height=1)
        self.frames_dir.grid(row=3)

        self.vid_out_label = tk.Label(self.root, width=100, height=1, text="Video Output Name (outputs to src/vid_out)", fg="black")
        self.vid_out_label.grid(row=4)
        self.vid_out = tk.Text(self.root, width=100, height=1)
        self.vid_out.grid(row=5)

        self.r_exp_label = tk.Label(self.root, width=100, height=1, text="Pixel Red Expression")
        self.r_exp_label.grid(row=6)
        self.r_exp = tk.Text(self.root, width=100, height=1)
        self.r_exp.grid(row=7)

        self.r_exp_mult_label = tk.Label(self.root, width=100, height=1, text="Pixel Red Expression Multiplier")
        self.r_exp_mult_label.grid(row=8)
        self.r_exp_mult = tk.Text(self.root, width=100, height=1)
        self.r_exp_mult.grid(row=9)

        self.g_exp_label = tk.Label(self.root, width=100, height=1, text="Pixel Green Expression")
        self.g_exp_label.grid(row=10)
        self.g_exp = tk.Text(self.root, width=100, height=1)
        self.g_exp.grid(row=11)

        self.g_exp_mult_label = tk.Label(self.root, width=100, height=1, text="Pixel Green Expression Multiplier")
        self.g_exp_mult_label.grid(row=12)
        self.g_exp_mult = tk.Text(self.root, width=100, height=1)
        self.g_exp_mult.grid(row=13)

        self.b_exp_label = tk.Label(self.root, width=100, height=1, text="Pixel Blue Expression")
        self.b_exp_label.grid(row=14)
        self.b_exp = tk.Text(self.root, width=100, height=1)
        self.b_exp.grid(row=15)

        self.b_exp_mult_label = tk.Label(self.root, width=100, height=1, text="Pixel Blue Expression Multiplier")
        self.b_exp_mult_label.grid(row=16)
        self.b_exp_mult = tk.Text(self.root, width=100, height=1)
        self.b_exp_mult.grid(row=17)

        self.scale_vcmd = (self.root.register(self.validate_scale_value), '%P')

        self.r_scale_label = tk.Label(self.root, width=100, height=1, text="Pixel Red Scale Multiplier")
        self.r_scale_label.grid(row=18)
        self.r_scale = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.r_scale.grid(row=19)

        self.r_scaleadj_label = tk.Label(self.root, width=100, height=1, text="Pixel Red Scale Multiplier Variation")
        self.r_scaleadj_label.grid(row=20)
        self.r_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.r_scaleadj.grid(row=21)

        self.g_scale_label = tk.Label(self.root, width=100, height=1, text="Pixel Green Scale Multiplier")
        self.g_scale_label.grid(row=22)
        self.g_scale = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.g_scale.grid(row=23)

        self.g_scaleadj_label = tk.Label(self.root, width=100, height=1, text="Pixel Green Scale Multiplier Variation")
        self.g_scaleadj_label.grid(row=24)
        self.g_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.g_scaleadj.grid(row=25)

        self.b_scale_label = tk.Label(self.root, width=100, height=1, text="Pixel Blue Scale Multiplier")
        self.b_scale_label.grid(row=26)
        self.b_scale = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.b_scale.grid(row=27)

        self.b_scaleadj_label = tk.Label(self.root, width=100, height=1, text="Pixel Blue Scale Multiplier Variation")
        self.b_scaleadj_label.grid(row=28)
        self.b_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.b_scaleadj.grid(row=29)

        self.ir_vcmd = (self.root.register(self.validate_interp_ratio), '%P')

        self.interp_ratio_label = tk.Label(self.root, width=100, height=1, text="Original Video : FX Interpolation Ratio")
        self.interp_ratio_label.grid(row=30)
        self.interp_ratio = tk.Entry(self.root, validate="key", validatecommand=self.ir_vcmd)
        self.interp_ratio.grid(row=31)

        self.interp_adj_label = tk.Label(self.root, width=100, height=1, text="Interpolation Ratio Variation")
        self.interp_adj_label.grid(row=32)
        self.interp_adj = tk.Entry(self.root, validate="key", validatecommand=self.scale_vcmd)
        self.interp_adj.grid(row=33)

        self.applyredux_var = tk.IntVar()
        self.applyredux_cbox = tk.Checkbutton(self.root, text="applyRedux", variable=self.applyredux_var)
        self.applyredux_cbox.grid(row=34)

        self.reduxbefore_var = tk.IntVar()
        self.reduxbefore_cbox = tk.Checkbutton(self.root, text="reduxBefore", variable=self.reduxbefore_var)
        self.reduxbefore_cbox.grid(row=35)

        self.applygfire_var = tk.IntVar()
        self.applygfire_cbox = tk.Checkbutton(self.root, text="applyGfire", variable=self.applyredux_var)
        self.applygfire_cbox.grid(row=34)

        self.gfirebefore_var = tk.IntVar()
        self.gfirebefore_cbox = tk.Checkbutton(self.root, text="gfireBefore", variable=self.reduxbefore_var)
        self.gfirebefore_cbox.grid(row=35)

        self.applyed_var = tk.IntVar()
        self.applyed_cbox = tk.Checkbutton(self.root, text="edgeDetect", variable=self.applyed_var)
        self.applyed_cbox.grid(row=36)

        self.edbefore_var = tk.IntVar()
        self.edbefore_cbox = tk.Checkbutton(self.root, text="edBefore", variable=self.edbefore_var)
        self.edbefore_cbox.grid(row=37)

        self.applykmc_var = tk.IntVar()
        self.applykmc_cbox = tk.Checkbutton(self.root, text="applyKmc", variable=self.applykmc_var)
        self.applykmc_cbox.grid(row=36)

        self.kmcbefore_var = tk.IntVar()
        self.kmcbefore_cbox = tk.Checkbutton(self.root, text="kmcBefore", variable=self.kmcbefore_var)
        self.kmcbefore_cbox.grid(row=37)

        self.watershed_var = tk.IntVar()
        self.watershed_cbox = tk.Checkbutton(self.root, text="applyWater", variable=self.watershed_var)
        self.watershed_cbox.grid(row=38)

        self.waterbefore_var = tk.IntVar()
        self.waterbefore_cbox = tk.Checkbutton(self.root, text="waterBefore", variable=self.waterbefore_var)
        self.waterbefore_cbox.grid(row=39)

        self.applywave_var = tk.IntVar()
        self.applywave_cbox = tk.Checkbutton(self.root, text="applyWave", variable=self.applywave_var)
        self.applywave_cbox.grid(row=40)

        self.wavebefore_var = tk.IntVar()
        self.wavebefore_cbox = tk.Checkbutton(self.root, text="waveBefore", variable=self.wavebefore_var)
        self.wavebefore_cbox.grid(row=41)

        self.applysine_var = tk.IntVar()
        self.applysine_cbox = tk.Checkbutton(self.root, text="applySine", variable=self.applysine_var)
        self.applysine_cbox.grid(row=42)

        self.sinebefore_var = tk.IntVar()
        self.sinebefore_cbox = tk.Checkbutton(self.root, text="sineBefore", variable=self.sinebefore_var)
        self.sinebefore_cbox.grid(row=43)

        self.applycosine_var = tk.IntVar()
        self.applycosine_cbox = tk.Checkbutton(self.root, text="applyCosine", variable=self.applycosine_var)
        self.applycosine_cbox.grid(row=44)

        self.cosinebefore_var = tk.IntVar()
        self.cosinebefore_cbox = tk.Checkbutton(self.root, text="cosineBefore", variable=self.cosinebefore_var)
        self.cosinebefore_cbox.grid(row=45)

        self.applydither_var = tk.IntVar()
        self.applydither_cbox = tk.Checkbutton(self.root, text="applyDither", variable=self.applydither_var)
        self.applydither_cbox.grid(row=46)

        self.ditherbefore_var = tk.IntVar()
        self.ditherbefore_cbox = tk.Checkbutton(self.root, text="ditherBefore", variable=self.ditherbefore_var)
        self.ditherbefore_cbox.grid(row=47)

    def validate_scale_value(self, value):
        try:
            if value.strip() == "":
                tk.messagebox.showerror("Invalid float Input", "Please enter a valid float f; you entered no value")
                return False
            float(value)
            return True
        except ValueError:
            tk.messagebox.showerror("Invalid float Input", "Please enter a valid float f; you may have not entered a numerical value as text input somewhere")
            return False

    def validate_interp_ratio(self, value):
        try:
            if value.strip() == "":
                tk.messagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered no value")
                return False
            float(value)
            if value > 1.0:
                tk.messagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered a float greater than 1.0")
                return False
            if value < 0.0:
                tk.messagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered a float less than 0.0")
                return False
            return True
        except ValueError:
            tk.messagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you may have not entered a numerical value as text input")
            return False

    def generate(self):
        """
        Example parameters to pass to routineVideoFx()
        "vid_in/work_tv.mp4", // inVidName 
        "png_out/worktv_6", // framesDir
        "worktv_6", // outVidName
        "x+y", // expressionR
        "1.0001", // multFnR 
        "y+x", // expressionG
        "1.0001", // multFnG
        "x-y", // expressionB
        "1.0001", // multFnB
        1.001, // scaleR
        1.005, // scaleAdjR
        1.001, // scaleG
        1.005, // scaleAdjG
        1.001, // scaleB
        1.005, // scaleAdjB
        0.9999, // interpRatio (ratio < 0.5 => less of inVidName; ratio > 0.5 => more of inVidName)
        0.007, // interpAdj (value represents difference in interp ratio by final frame)
        false, true, // applyRedux, reduxBefore
        false, true, // applyGfire, gfireBefore
        false, true, // edgeDetect, edBefore
        false, true, // applyKmc, kmcBefore
        false, true, // applyWater, wtrBefore
        false, true, // applyWave, waveBefore
        false, true, // applySine, sinBefore
        false, true, // applyCosine, cosBefore
        false, true, // applyDither, ditherBefore
        false, // invertSrc
        2, // bitsRedux
        5, // kmcFactor
        8, // dstBlockSize
        8, // dctBlockSize
        uint8(128), // gfireTol
        """
        start = time.time()
        main_dir = os.getcwd()
        os.chdir("../src")
        cmd = "go build -o main && ./main %s %s %s %s %s %s %s %s %s %f %f %f %f %f %f %f %f %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"%(
            "vid_in/"+self.vid_in.get("1.0", tk.END).strip(), # will eventually replace this by a tkinter select option of files in src/vid_in
            "png_out/"+self.frames_dir.get("1.0", tk.END).strip(),
            self.vid_out.get("1.0", tk.END).strip(),
            self.r_exp.get("1.0", tk.END).strip(),
            self.r_exp_mult.get("1.0", tk.END).strip(),
            self.g_exp.get("1.0", tk.END).strip(),
            self.g_exp_mult.get("1.0", tk.END).strip(),
            self.b_exp.get("1.0", tk.END).strip(),
            self.b_exp_mult.get("1.0", tk.END).strip(),
            float(self.r_scale.get()),
            float(self.r_scaleadj.get()),
            float(self.g_scale.get()),
            float(self.g_scaleadj.get()),
            float(self.b_scale.get()),
            float(self.b_scaleadj.get()),
            float(self.interp_ratio.get()),
            float(self.interp_adj.get()),
          
        )
        cmd_result = subprocess.run(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        end = time.time()
        os.chdir(main_dir)
        if cmd_result.returncode == 0:
            print('Go executed successfully')
            print('Go Output:')
            print(cmd_result.stdout)
        else:
            print('Go failed with error:')
            print(cmd_result.stderr)
        print(f'app.generate() ./src/main execution time: {round(end-start, 2)} sec')

    def run(self):
        self.root.mainloop()

if __name__ == "__main__":
    pass
