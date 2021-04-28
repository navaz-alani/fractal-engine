// This file generates the animation similar to the one at
// https://en.wikipedia.org/wiki/Julia_set#/media/File:JSr07885.gif

package main

import (
	"fmt"
	"image/gif"
	"math"
	"math/cmplx"
	"os"
	"runtime"

	"github.com/navaz-alani/fractal-engine"
)

const (
	numFrames = 100
	maxIters  = 1000
)

var iterateAnglularDiff float64 = 2 * math.Pi / numFrames

func main() {
	f, err := os.Create("renders/circle-anim.gif")
	if err != nil {
		fmt.Printf("error creating render file: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// set up color palette
	pal := fractal.GreenBlackPalette{
		MaxIters: maxIters,
	}
	pal.Precompute()

	// set up the julia set fn
	mandelbrot := fractal.JuliaSetFn{
		Exp:       2,
		EscapeRad: 2,
		MaxIters:  maxIters,
	}

	// set up animation renderer
	gifRenderer := fractal.GIFRenderer{
		NFrames:      numFrames,
		ImgWidth:     512,
		ImgHeight:    512,
		PlotWidth:    3,
		PlotHeight:   3,
		ZoomFactor:   1,
		FrameDelayFn: func(u uint) int { return 8 },
		ColorFnGen: func(u uint) fractal.ColorFunction {
			return func(c complex128) uint8 {
				mandelbrot := mandelbrot
				mandelbrot.InitIterate = c
				p := cmplx.Rect(0.7885, float64(u)*iterateAnglularDiff)
				return pal.PixelColorIdx(mandelbrot.EscapeIter(p))
			}
		},
		Palette:  pal.Palette(),
		Progress: true,
	}

	// render the animation
	anim := gifRenderer.RenderParallel(2 * uint(runtime.NumCPU()))
	// encode the animation to file
	if err := gif.EncodeAll(f, anim); err != nil {
		fmt.Printf("error encoding gif: %s\n", err.Error())
		os.Exit(1)
	}
}
