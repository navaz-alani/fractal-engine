package fractal

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"runtime"
)

// ColorFunctionGenerator is a function which takes the iteration that the
// animation is currently at, and returns a function to color the plot for that
// iteration.
// This can be useful for animations where the plot changes/evolves over time
// (e.g. varying the point `c` parameter of a function whose Julia set is being
// plotted).
type ColorFunctionGenerator func(uint) ColorFunction

// FrameDelayFunction is a function which takes the frame that the animation
// has just rendered and returns the delay between the previous frame and the
// current frame in the animation.
// This can be useful to making animations appear to speed up or slow down.
// If each frame of the animation has the same delay, then this function ought
// to be constant.
type FrameDelayFunction func(uint) int

// GIFRenderer renders GIF animations of the complex plane.
// It supports zooming into a point on the plane, or just rendering a
// changing-plot on the complex plane (both of which can be achieved using the
// ColorFunctionGenerator).
type GIFRenderer struct {
	// frame count of the gif animation
	NFrames uint
	// image dimensions
	ImgWidth, ImgHeight uint
	// plot dimensions
	PlotWidth, PlotHeight float64
	// center of the plot (do not modify for origin)
	CX, CY float64
	// factor by which the plot width and height are shrunk each frame of the
	// animation.
	// If the animation does not require zoom, then set this to 1.
	ZoomFactor float64
	// ColorFnGen is the ColorFunctionGenerator used to render each frame of the
	// animation.
	ColorFnGen ColorFunctionGenerator
	// FrameDelayFn is the FrameDelayFunction which is used to control the speed
	// of the animation.
	FrameDelayFn FrameDelayFunction
	// Palette is the color palette for the images.
	Palette []color.Color
	// Display render progress or not
	Progress bool
}

func (gr *GIFRenderer) Render() *gif.GIF {
	anim := &gif.GIF{}
	// set up the base image renderer
	baseImgRenderer := ImgRenderer{
		ImgWidth:   gr.ImgWidth,
		ImgHeight:  gr.ImgHeight,
		PlotWidth:  gr.PlotWidth,
		PlotHeight: gr.PlotHeight,
		CX:         gr.CX,
		CY:         gr.CY,
	}
	// determine the renderer sliceWidth (using at most runtime.NumCPU() slices)
	sliceWidth := gr.ImgWidth / (2 * uint(runtime.NumCPU()))
	// render the frames
	for f := uint(0); f < gr.NFrames; f++ {
		// modify image renderer's plot dimensions after the first frame
		if f != 0 {
			baseImgRenderer.PlotWidth *= gr.ZoomFactor
			baseImgRenderer.PlotHeight *= gr.ZoomFactor
		}
		// set frame to render
		img := image.NewPaletted(
			image.Rect(0, 0, int(gr.ImgWidth), int(gr.ImgHeight)),
			gr.Palette,
		)
		baseImgRenderer.Img = img
		baseImgRenderer.RenderImg(gr.ColorFnGen(f), sliceWidth)
		// add rendered image to animation
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, gr.FrameDelayFn(f))
		if gr.Progress {
			fmt.Fprintf(os.Stdout, "Rendered %d of %d frames\n", f+1, gr.NFrames)
		}
	}
	return anim
}
