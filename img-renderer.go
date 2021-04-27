package fractal

import (
	"sync"
)

// ColorFunction returns the index of the color to plot for the pixel
// corresponding to the given complex point.
type ColorFunction func(complex128) uint8

type Img interface {
	//Set(x, y int, c color.Color)
	SetColorIndex(x, y int, idx uint8)
}

// ImgRenderer renders images of a 2-dimensional complex plot, using a
// ColorFunction, which maps complex points to color.Color values.
type ImgRenderer struct {
	// image being rendered (animation frame/standalone image)
	Img Img
	// image dimensions (width & height)
	ImgWidth, ImgHeight uint
	// plot dimensions (with & height)
	PlotWidth, PlotHeight float64
	// plot center
	CX, CY float64
}

// renderSlice processes a vertical slice (slice of the x-axis) of the image.
// It can (should) be used to process different vertical slices of an image
// concurrently and therefore speed up the render process.
//
// Allowing the colorFn to be a per-slice parameter enables one to render
// different slices in the same image with different color palettes.
//
// Warning: increasing the number of slices will improve performance, but the
// performance gains will begin to diminish after some point.
// This depends on factors sich as the number of logical cores available on the
// machine.
func (r *ImgRenderer) RenderSlice(
	wg *sync.WaitGroup,
	sliceLb, sliceUb uint,
	colorFn ColorFunction,
) {
	defer wg.Done()
	for py := uint(0); py < r.ImgHeight; py++ {
		y := r.CY + (float64(py)/float64(r.ImgHeight)-0.5)*r.PlotHeight
		for px := sliceLb; px < sliceUb; px++ {
			x := r.CX + (float64(px)/float64(r.ImgWidth)-0.5)*r.PlotWidth
			r.Img.SetColorIndex(int(px), int(py), colorFn(complex(x, y)))
		}
	}
}

// renderImg renders the underlying image using the given colorFn for all slices
// of the image.
// This method spawns go-routines to process the image in vertical slices and
// the sliceWidth parameter determines the width of the vertical slice that each
// rendering go-routine will work on.
func (r *ImgRenderer) RenderImg(colorFn ColorFunction, sliceWidth uint) {
	var wg sync.WaitGroup
	var i uint
	for {
		var sliceLb, sliceUb uint = i, i + sliceWidth
		if sliceUb > r.ImgWidth {
			sliceUb = r.ImgWidth
		}
		wg.Add(1)
		go r.RenderSlice(&wg, sliceLb, sliceUb, colorFn)
		if sliceUb == r.ImgWidth {
			break
		}
		i = sliceUb
	}
	wg.Wait()
}
