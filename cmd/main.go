package main

import (
  "github.com/navaz-alani/fractal-engine"

	"flag"
	"fmt"
	"image"
	"image/color"
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
	colorPatelle   = flag.String("palette", "uf", "color palette ('bw', 'bw-inv', 'uf')")
	contrast       = flag.Uint("bw-contrast", 15, "bw-palette constrast")
	juliaExp       = flag.Uint("julia-exp", 2, "exponent of the function generating the set f(z)=z**exp+c")
	juliaEscapeRad = flag.Float64("julia-escape-rad", 2, "iterate absolute value escape radius")
  initIterateX   = flag.Float64("init-iteratex", 0, "real part of the initial iterate")
  initIterateY   = flag.Float64("init-iteratey", 0, "imag part of the initial iterate")
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
	} else if *colorPatelle != "bw" && *colorPatelle != "bw-inv" && *colorPatelle != "uf" {
		exitOnErr("invalid palette - expected 'bw', 'bw-inv' or 'uf'", nil, 1)
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
	} else {
		ufPalette := &fractal.UltraFractalPalette{
			MaxIters: int(*iterations),
		}
		pal = ufPalette
	}

	pal.Precompute()

	// prepare the color fn
	juliaSetFn := &fractal.JuliaSetFn{
		Exp:       int(*juliaExp),
		MaxIters:  int(*iterations),
		EscapeRad: *juliaEscapeRad,
    InitIterate: complex(*initIterateX, *initIterateY),
	}
	colorFn := func(c complex128) color.Color {
		return pal.PixelColor(juliaSetFn.EscapeIter(c))
	}

	if *mode == "img" {
		img := image.NewRGBA(
			image.Rect(0, 0, int(*imgWidth), int(*imgHeight)),
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
		anim := renderer.Render()
		if err := gif.EncodeAll(file, anim); err != nil {
			exitOnErr("failed to encode render to gif", err, 1)
		}
	}
}
