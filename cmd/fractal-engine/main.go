package main

import (
	"github.com/navaz-alani/fractal-engine"

	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	mode           = flag.String("mode", "img", "render mode ('img' or 'gif')")
	imgWidth       = flag.Uint("img-width", 4096, "width of the image")
	imgHeight      = flag.Uint("img-height", 2048, "height of the image")
	plotWidth      = flag.Float64("plot-width", 4, "width of the complex plot (x_max-x_min)")
	plotHeight     = flag.Float64("plot-height", 2, "height of the complex plot (y_max-y_min)")
	plotCenterX    = flag.Float64("plot-cx", 0, "x-coordinate of the plot center")
	plotCenterY    = flag.Float64("plot-cy", 0, "y-coordinate of the plot center")
	zoomFactor     = flag.Float64("zoom-fact", 0.9, "zoom factor")
	numFrames      = flag.Uint("num-frames", 100, "number of frames in the animation")
	frameDelay     = flag.Uint("frame-delay", 8, "delay between animation frames (100s of ms)")
	iterations     = flag.Uint("iterations", 1000, "number of iterations")
	ofName         = flag.String("of-name", "render.png", "output file name")
	colorPatelle   = flag.String("palette", "uf", "color palette ('bw', 'bw-inv', 'uf', 'gb')")
	contrast       = flag.Uint("bw-contrast", 15, "bw-palette constrast")
	juliaExp       = flag.Uint("julia-exp", 2, "exponent of the function generating the set f(z)=z**exp+c")
	juliaEscapeRad = flag.Float64("julia-escape-rad", 2, "iterate absolute value escape radius")
	initIterateX   = flag.Float64("init-iteratex", 0, "real part of the initial iterate")
	initIterateY   = flag.Float64("init-iteratey", 0, "imag part of the initial iterate")
	parallel       = flag.Bool("parallel", true, "whether to render animations in parallel")
	progress       = flag.Bool("progress", false, "display render progress (frame count)")
)

func exitOnErr(ctx string, err error, code int) {
	if err != nil {
		fmt.Printf("%s: %s\n", ctx, err.Error())
	} else {
		fmt.Printf("%s", ctx)
	}
	os.Exit(code)
}

func main() {
	flag.Parse()

	profile, err := os.Create("profile.prof")
	if err != nil {
		exitOnErr("failed to create profile file", err, 1)
	}
	pprof.StartCPUProfile(profile)
	defer pprof.StopCPUProfile()

	if *ofName == "" {
		exitOnErr("output filename not specified; use flag -h for help", nil, 1)
	}
	file, err := os.Create(*ofName)
	if err != nil {
		exitOnErr("failed to create output file: %s", err, 1)
	}
	defer file.Close()

	if *mode != "img" && *mode != "gif" {
		exitOnErr("invalid mode - expected 'img' or 'gif'", nil, 1)
	} else if *colorPatelle != "bw" &&
		*colorPatelle != "bw-inv" &&
		*colorPatelle != "uf" &&
		*colorPatelle != "gb" {
		exitOnErr("invalid palette - expected 'bw', 'bw-inv', 'uf' or 'gb'", nil, 1)
	}

	// setup the color palette
	var pal fractal.ColorPalette

	if *colorPatelle == "bw" || *colorPatelle == "bw-inv" {
		bwPalette := &fractal.BWPalette{
			Contrast: uint8(*contrast),
			MaxIters: int(*iterations),
		}
		if *colorPatelle == "bw-inv" {
			bwPalette.Inverse = true
		}
		pal = bwPalette
	} else if *colorPatelle == "uf" {
		ufPalette := &fractal.UltraFractalPalette{
			MaxIters: int(*iterations),
		}
		pal = ufPalette
	} else {
		gbPalette := &fractal.GreenBlackPalette{
			MaxIters: int(*iterations),
		}
		pal = gbPalette
	}

	pal.Precompute()

	// prepare the color fn
	juliaSetFn := &fractal.JuliaSetFn{
		Exp:         int(*juliaExp),
		MaxIters:    int(*iterations),
		EscapeRad:   *juliaEscapeRad,
		InitIterate: complex(*initIterateX, *initIterateY),
	}
	colorFn := func(c complex128) uint8 {
		return pal.PixelColorIdx(juliaSetFn.EscapeIter(c))
	}

	if *mode == "img" {
		img := image.NewPaletted(
			image.Rect(0, 0, int(*imgWidth), int(*imgHeight)),
			pal.Palette(),
		)
		renderer := fractal.ImgRenderer{
			Img:        img,
			ImgWidth:   *imgWidth,
			ImgHeight:  *imgHeight,
			PlotWidth:  *plotWidth,
			PlotHeight: *plotHeight,
			CX:         *plotCenterX,
			CY:         *plotCenterY,
		}
		renderer.RenderImg(colorFn, *imgWidth/uint(runtime.NumCPU()))
		if err := png.Encode(file, img); err != nil {
			exitOnErr("failed to encode render to png", err, 1)
		}
	} else {
		renderer := fractal.GIFRenderer{
			NFrames:    *numFrames,
			ImgWidth:   *imgWidth,
			ImgHeight:  *imgHeight,
			PlotWidth:  *plotWidth,
			PlotHeight: *plotHeight,
			CX:         *plotCenterX,
			CY:         *plotCenterY,
			ZoomFactor: *zoomFactor,
			ColorFnGen: func(u uint) fractal.ColorFunction {
				return colorFn
			},
			FrameDelayFn: func(u uint) int {
				return int(*frameDelay)
			},
			Palette:  pal.Palette(),
			Progress: *progress,
		}
		var anim *gif.GIF
		if *parallel {
			anim = renderer.RenderParallel(uint(runtime.NumCPU()))
		} else {
			anim = renderer.Render()
		}
		if err := gif.EncodeAll(file, anim); err != nil {
			exitOnErr("failed to encode render to gif", err, 1)
		}
	}
}
