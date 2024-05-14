import subprocess
import time
import os
import tkinter as tk
from tkinter import ttk
import tkinter.messagebox as tkmessagebox

class App:
    def __init__(self, root):
        self.root = root
        self.root.title("vid-gen-go")
        self.root.geometry("720x770")
        # TODO: Play around with row and columns in .grid() to find a better param layout on the gui
        self.root.columnconfigure(0, weight=1)
        self.root.columnconfigure(1, weight=1)
        """for i in range(0, 33):
            self.root.rowconfigure(i, weight=1)"""
        self.vid_in_label = tk.Label(self.root, width=50, height=1, text="Video Input Name (from src/vid_in)", fg="black")
        #self.vid_in_label.grid(row=0)
        self.vid_in_label.grid(row=0, column=0)
        og_dir = os.getcwd()
        os.chdir(og_dir+"/src/vid_in")
        self.video_name_options = os.listdir(os.getcwd())
        os.chdir(og_dir)
        if len(self.video_name_options) == 0:
            self.video_name_options = ["No Videos To Select"]
        self.vid_in_var = tk.StringVar()
        self.vid_in_var.set(self.video_name_options[0])
        self.vid_in_drop = ttk.Combobox(self.root, textvariable=self.vid_in_var, values=self.video_name_options)
        #self.vid_in_drop.grid(row=1)
        self.vid_in_drop.grid(row=0, column=1)
        self.vid_in_drop.bind("<<ComboboxSelected>>", self.handle_vid_in_select)
        #self.vid_in = tk.Text(self.root, width=40, height=1)
        #self.vid_in.grid(row=1)
        #self.vid_in.grid(row=0, column=1)

        self.frames_dir_label = tk.Label(self.root, width=50, height=1, text="Output Directory Name for Output Video Frames", fg="black")
        #self.frames_dir_label.grid(row=2)
        self.frames_dir_label.grid(row=1, column=0)
        self.frames_dir = tk.Text(self.root, width=40, height=1)
        #self.frames_dir.grid(row=3)
        self.frames_dir.grid(row=1, column=1)

        self.vid_out_label = tk.Label(self.root, width=50, height=1, text="Video Output Name (outputs to src/vid_out)", fg="black")
        #self.vid_out_label.grid(row=4)
        self.vid_out_label.grid(row=2, column=0)
        self.vid_out = tk.Text(self.root, width=40, height=1)
        #self.vid_out.grid(row=5)
        self.vid_out.grid(row=2, column=1)

        self.r_exp_label = tk.Label(self.root, width=50, height=1, text="Pixel Red Expression", fg="black")
        #self.r_exp_label.grid(row=6)
        self.r_exp_label.grid(row=3, column=0)
        self.r_exp = tk.Text(self.root, width=40, height=1)
        #self.r_exp.grid(row=7)
        self.r_exp.grid(row=3, column=1)

        self.r_exp_mult_label = tk.Label(self.root, width=50, height=1, text="Pixel Red Expression Multiplier", fg="black")
        #self.r_exp_mult_label.grid(row=8)
        self.r_exp_mult_label.grid(row=4, column=0)
        self.r_exp_mult = tk.Text(self.root, width=40, height=1)
        #self.r_exp_mult.grid(row=9)
        self.r_exp_mult.grid(row=4, column=1)

        self.g_exp_label = tk.Label(self.root, width=50, height=1, text="Pixel Green Expression", fg="black")
        #self.g_exp_label.grid(row=10)
        self.g_exp_label.grid(row=5, column=0)
        self.g_exp = tk.Text(self.root, width=40, height=1)
        #self.g_exp.grid(row=11)
        self.g_exp.grid(row=5, column=1)

        self.g_exp_mult_label = tk.Label(self.root, width=50, height=1, text="Pixel Green Expression Multiplier", fg="black")
        #self.g_exp_mult_label.grid(row=12)
        self.g_exp_mult_label.grid(row=6, column=0)
        self.g_exp_mult = tk.Text(self.root, width=40, height=1)
        #self.g_exp_mult.grid(row=13)
        self.g_exp_mult.grid(row=6, column=1)

        self.b_exp_label = tk.Label(self.root, width=50, height=1, text="Pixel Blue Expression", fg="black")
        #self.b_exp_label.grid(row=14)
        self.b_exp_label.grid(row=7, column=0)
        self.b_exp = tk.Text(self.root, width=40, height=1)
        #self.b_exp.grid(row=15)
        self.b_exp.grid(row=7, column=1)

        self.b_exp_mult_label = tk.Label(self.root, width=50, height=1, text="Pixel Blue Expression Multiplier", fg="black")
        #self.b_exp_mult_label.grid(row=16)
        self.b_exp_mult_label.grid(row=8, column=0)
        self.b_exp_mult = tk.Text(self.root, width=40, height=1)
        #self.b_exp_mult.grid(row=17)
        self.b_exp_mult.grid(row=8, column=1)

        self.float_vcmd = (self.root.register(self.validate_float_value), '%P')

        self.r_scale_label = tk.Label(self.root, width=50, height=1, text="Pixel Red Scale Multiplier", fg="black")
        #self.r_scale_label.grid(row=18)
        self.r_scale_label.grid(row=9, column=0)
        self.r_scale = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.r_scale.grid(row=19)
        self.r_scale.grid(row=9, column=1)

        self.r_scaleadj_label = tk.Label(self.root, width=50, height=1, text="Pixel Red Scale Multiplier Variation", fg="black")
        #self.r_scaleadj_label.grid(row=20)
        self.r_scaleadj_label.grid(row=10, column=0)
        self.r_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.r_scaleadj.grid(row=21)
        self.r_scaleadj.grid(row=10, column=1)

        self.g_scale_label = tk.Label(self.root, width=50, height=1, text="Pixel Green Scale Multiplier", fg="black")
        #self.g_scale_label.grid(row=22)
        self.g_scale_label.grid(row=11, column=0)
        self.g_scale = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.g_scale.grid(row=23)
        self.g_scale.grid(row=11, column=1)

        self.g_scaleadj_label = tk.Label(self.root, width=50, height=1, text="Pixel Green Scale Multiplier Variation", fg="black")
        #self.g_scaleadj_label.grid(row=24)
        self.g_scaleadj_label.grid(row=12, column=0)
        self.g_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.g_scaleadj.grid(row=25)
        self.g_scaleadj.grid(row=12, column=1)

        self.b_scale_label = tk.Label(self.root, width=50, height=1, text="Pixel Blue Scale Multiplier", fg="black")
        #self.b_scale_label.grid(row=26)
        self.b_scale_label.grid(row=13, column=0)
        self.b_scale = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.b_scale.grid(row=27)
        self.b_scale.grid(row=13, column=1)

        self.b_scaleadj_label = tk.Label(self.root, width=50, height=1, text="Pixel Blue Scale Multiplier Variation", fg="black")
        #self.b_scaleadj_label.grid(row=28)
        self.b_scaleadj_label.grid(row=14, column=0)
        self.b_scaleadj = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.b_scaleadj.grid(row=29)
        self.b_scaleadj.grid(row=14, column=1)

        self.ir_vcmd = (self.root.register(self.validate_interp_ratio), '%P')

        self.interp_ratio_label = tk.Label(self.root, width=50, height=1, text="Original Video : FX Interpolation Ratio", fg="black")
        #self.interp_ratio_label.grid(row=30)
        self.interp_ratio_label.grid(row=15, column=0)
        self.interp_ratio = tk.Entry(self.root, validate="key", validatecommand=self.ir_vcmd)
        #self.interp_ratio.grid(row=31)
        self.interp_ratio.grid(row=15, column=1)

        self.interp_adj_label = tk.Label(self.root, width=50, height=1, text="Interpolation Ratio Variation", fg="black")
        #self.interp_adj_label.grid(row=32)
        self.interp_adj_label.grid(row=16, column=0)
        self.interp_adj = tk.Entry(self.root, validate="key", validatecommand=self.float_vcmd)
        #self.interp_adj.grid(row=33)
        self.interp_adj.grid(row=16, column=1)

        self.applyredux_var = tk.IntVar()
        self.applyredux_cbox = tk.Checkbutton(self.root, text="applyRedux", variable=self.applyredux_var)
        #self.applyredux_cbox.grid(row=34)
        self.applyredux_cbox.grid(row=17, column=0)

        self.reduxbefore_var = tk.IntVar()
        self.reduxbefore_cbox = tk.Checkbutton(self.root, text="reduxBefore", variable=self.reduxbefore_var)
        #self.reduxbefore_cbox.grid(row=35)
        self.reduxbefore_cbox.grid(row=17, column=1)

        self.applygfire_var = tk.IntVar()
        self.applygfire_cbox = tk.Checkbutton(self.root, text="applyGfire", variable=self.applygfire_var)
        #self.applygfire_cbox.grid(row=36)
        self.applygfire_cbox.grid(row=18, column=0)

        self.gfirebefore_var = tk.IntVar()
        self.gfirebefore_cbox = tk.Checkbutton(self.root, text="gfireBefore", variable=self.gfirebefore_var)
        #self.gfirebefore_cbox.grid(row=37)
        self.gfirebefore_cbox.grid(row=18, column=1)

        self.applyed_var = tk.IntVar()
        self.applyed_cbox = tk.Checkbutton(self.root, text="edgeDetect", variable=self.applyed_var)
        #self.applyed_cbox.grid(row=38)
        self.applyed_cbox.grid(row=19, column=0)

        self.edbefore_var = tk.IntVar()
        self.edbefore_cbox = tk.Checkbutton(self.root, text="edBefore", variable=self.edbefore_var)
        #self.edbefore_cbox.grid(row=39)
        self.edbefore_cbox.grid(row=19, column=1)

        self.applykmc_var = tk.IntVar()
        self.applykmc_cbox = tk.Checkbutton(self.root, text="applyKmc", variable=self.applykmc_var)
        #self.applykmc_cbox.grid(row=40)
        self.applykmc_cbox.grid(row=20, column=0)

        self.kmcbefore_var = tk.IntVar()
        self.kmcbefore_cbox = tk.Checkbutton(self.root, text="kmcBefore", variable=self.kmcbefore_var)
        #self.kmcbefore_cbox.grid(row=41)
        self.kmcbefore_cbox.grid(row=20, column=1)

        self.watershed_var = tk.IntVar()
        self.watershed_cbox = tk.Checkbutton(self.root, text="applyWater", variable=self.watershed_var)
        #self.watershed_cbox.grid(row=42)
        self.watershed_cbox.grid(row=21, column=0)

        self.waterbefore_var = tk.IntVar()
        self.waterbefore_cbox = tk.Checkbutton(self.root, text="waterBefore", variable=self.waterbefore_var)
        #self.waterbefore_cbox.grid(row=43)
        self.waterbefore_cbox.grid(row=21, column=1)

        self.applywave_var = tk.IntVar()
        self.applywave_cbox = tk.Checkbutton(self.root, text="applyWave", variable=self.applywave_var)
        #self.applywave_cbox.grid(row=44)
        self.applywave_cbox.grid(row=22, column=0)

        self.wavebefore_var = tk.IntVar()
        self.wavebefore_cbox = tk.Checkbutton(self.root, text="waveBefore", variable=self.wavebefore_var)
        #self.wavebefore_cbox.grid(row=45)
        self.wavebefore_cbox.grid(row=22, column=1)

        self.applysine_var = tk.IntVar()
        self.applysine_cbox = tk.Checkbutton(self.root, text="applySine", variable=self.applysine_var)
        #self.applysine_cbox.grid(row=46)
        self.applysine_cbox.grid(row=23, column=0)

        self.sinebefore_var = tk.IntVar()
        self.sinebefore_cbox = tk.Checkbutton(self.root, text="sineBefore", variable=self.sinebefore_var)
        #self.sinebefore_cbox.grid(row=47)
        self.sinebefore_cbox.grid(row=23, column=1)

        self.applycosine_var = tk.IntVar()
        self.applycosine_cbox = tk.Checkbutton(self.root, text="applyCosine", variable=self.applycosine_var)
        #self.applycosine_cbox.grid(row=48)
        self.applycosine_cbox.grid(row=24, column=0)

        self.cosinebefore_var = tk.IntVar()
        self.cosinebefore_cbox = tk.Checkbutton(self.root, text="cosineBefore", variable=self.cosinebefore_var)
        #self.cosinebefore_cbox.grid(row=49)
        self.cosinebefore_cbox.grid(row=24, column=1)

        self.applydither_var = tk.IntVar()
        self.applydither_cbox = tk.Checkbutton(self.root, text="applyDither", variable=self.applydither_var)
        #self.applydither_cbox.grid(row=50)
        self.applydither_cbox.grid(row=25, column=1)

        self.ditherbefore_var = tk.IntVar()
        self.ditherbefore_cbox = tk.Checkbutton(self.root, text="ditherBefore", variable=self.ditherbefore_var)
        #self.ditherbefore_cbox.grid(row=51)
        self.ditherbefore_cbox.grid(row=25, column=1)

        self.invertsrc_var = tk.IntVar()
        self.invertsrc_cbox = tk.Checkbutton(self.root, text="invertSrc", variable=self.invertsrc_var)
        #self.invertsrc_cbox.grid(row=52)
        self.invertsrc_cbox.grid(row=26, column=0)

        self.integer_vcmd = (self.root.register(self.validate_integer_value), '%P')

        self.bits_label = tk.Label(self.root, text="Bits Redux", fg="black")
        #self.bits_label.grid(row=53)
        self.bits_label.grid(row=27, column=0)
        self.bits_redux = tk.Entry(self.root, validate="key", validatecommand=self.integer_vcmd)
        #self.bits_redux.grid(row=54)
        self.bits_redux.grid(row=27, column=1)

        self.kmc_label = tk.Label(self.root, text="K-Means Clustering Factor", fg="black")
        #self.kmc_label.grid(row=55)
        self.kmc_label.grid(row=28, column=0)
        self.kmc_factor = tk.Entry(self.root, validate="key", validatecommand=self.integer_vcmd)
        #self.kmc_factor.grid(row=56)
        self.kmc_factor.grid(row=28, column=1)

        self.dstbsize_label = tk.Label(self.root, text="Discrete Sine Transform Block Size", fg="black")
        #self.dstbsize_label.grid(row=57)
        self.dstbsize_label.grid(row=29, column=0)
        self.dst_bsize = tk.Entry(self.root, validate="key", validatecommand=self.integer_vcmd)
        #self.dst_bsize.grid(row=58)
        self.dst_bsize.grid(row=29, column=1)

        self.dctbsize_label = tk.Label(self.root, text="Discrete Cosine Transform Block Size", fg="black")
        #self.dctbsize_label.grid(row=59)
        self.dctbsize_label.grid(row=30, column=0)
        self.dct_bsize = tk.Entry(self.root, validate="key", validatecommand=self.integer_vcmd)
        #self.dct_bsize.grid(row=60)
        self.dct_bsize.grid(row=30, column=1)

        self.gfire_vcmd = (self.root.register(self.validate_gfire), '%P')

        self.gfire_label = tk.Label(self.root, text="Grassfire Algorithm Tolerance", fg="black")
        #self.gfire_label.grid(row=61)
        self.gfire_label.grid(row=31, column=0)
        self.gfire_tol = tk.Entry(self.root, validate="key", validatecommand=self.gfire_vcmd)
        #self.gfire_tol.grid(row=62)
        self.gfire_tol.grid(row=31, column=1)
        
        self.gen_btn = tk.Button(self.root, text="Generate .mp4 Video", command=self.generate)
        #self.gen_btn.grid(row=63)
        self.gen_btn.grid(row=32, column=0)

    def handle_vid_in_select(self, event):
        self.vid_in_var.set(self.video_name_options[self.vid_in_drop.current()])

    def validate_float_value(self, value):
        try:
            if value.strip() == "":
                tkmessagebox.showerror("Invalid float Input", "Please enter a valid float f; you entered no value")
                return False
            float(value)
            return True
        except ValueError:
            tkmessagebox.showerror("Invalid float Input", "Please enter a valid float f; you may have not entered a valid numerical value as text input somewhere")
            return False

    def validate_interp_ratio(self, value):
        try:
            if value.strip() == "":
                tkmessagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered no value")
                return False
            float(value)
            if value > 1.0:
                tkmessagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered a float greater than 1.0")
                return False
            if value < 0.0:
                tkmessagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you entered a float less than 0.0")
                return False
            return True
        except ValueError:
            tkmessagebox.showerror("Invalid interp ratio Input", "For interp ratio, please enter a valid float f such that: 0.0 <= f <= 1.0; you may have not entered a valid numerical value as text input")
            return False
        
    def validate_integer_value(self, value):
        try:
            if value.strip() == "":
                tkmessagebox.showerror("Invalid integer Input", "Please enter a valid integer i such that: i >= 0; you entered no value")
                return False
            int_val = int(value)
            if int_val < 0:
                tkmessagebox.showerror("Invalid integer Input", "Please enter a valid integer i such that: i >= 0; you entered an integer less than 0")
                return False
            return True
        except ValueError:
            tkmessagebox.showerror("Invalid integer Input", "Please enter a valid integer i such that: i >= 0; you may have not entered a valid numerical value as text input somewhere")
            return False

    def validate_gfire(self, value):
        try:
            if value.strip() == "":
                tkmessagebox.showerror("Invalid Grassfire Algorithm integer Input", "Please enter a valid integer i such that: 0 <= i <= 255; you entered no value")
                return False
            int_val = int(value)
            if int_val < 0:
                tkmessagebox.showerror("Invalid Grassfire Algorithm integer Input", "Please enter a valid integer i such that: 0 <= i <= 255; you entered an integer less than 0")
                return False
            if int_val > 255:
                tkmessagebox.showerror("Invalid Grassfire Algorithm integer Input", "Please enter a valid integer i such that: 0 <= i <= 255; you entered an integer greater than 255")
                return False
            return True
        except ValueError:
            tkmessagebox.showerror("Invalid Grassfire Algorithm integer Input", "Please enter a valid integer i such that: 0 <= i <= 255; you may have not entered a valid numerical value as text input")
            return False

    def generate(self):
        cmd_str = "tkmain.exe %s %s %s %s %s %s %s %s %s %f %f %f %f %f %f %f %f %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"%(
            #"vid_in/"+self.vid_in.get("1.0", tk.END).strip(),
            "vid_in/"+self.vid_in_var.get(),
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
            int(self.applyredux_var.get()),
            int(self.reduxbefore_var.get()),
            int(self.applygfire_var.get()),
            int(self.gfirebefore_var.get()),
            int(self.applyed_var.get()),
            int(self.edbefore_var.get()),
            int(self.applykmc_var.get()),
            int(self.kmcbefore_var.get()),
            int(self.watershed_var.get()),
            int(self.waterbefore_var.get()),
            int(self.applywave_var.get()),
            int(self.wavebefore_var.get()),
            int(self.applysine_var.get()),
            int(self.sinebefore_var.get()),
            int(self.applycosine_var.get()),
            int(self.cosinebefore_var.get()),
            int(self.applydither_var.get()),
            int(self.ditherbefore_var.get()),
            int(self.invertsrc_var.get()),
            int(self.bits_redux.get()),
            int(self.kmc_factor.get()),
            int(self.dst_bsize.get()),
            int(self.dct_bsize.get()),
            int(self.gfire_tol.get()),
        )
        print(cmd_str)
        main_dir = os.getcwd()
        os.chdir(main_dir+"/src")
        start = time.time()
        cmd_result = subprocess.run(cmd_str, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        print(f'app.generate() ./src/main execution time: {round(time.time()-start, 2)} sec')
        os.chdir(main_dir)
        if cmd_result.returncode == 0:
            print('Go built and executed successfully:')
            print(cmd_result.stdout)
        else:
            print('Go building or execution failed with error:')
            print(cmd_result.stderr)
        
    def run(self):
        self.root.mainloop()

if __name__ == "__main__":
    pass
