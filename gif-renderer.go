package fractal

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"runtime"
	"sync"
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

type frameRenderJob struct {
	frameID    uint
	img        *image.Paletted
	renderer   *ImgRenderer
	jobStream  chan<- *frameRenderJob
	sliceWidth uint
}

func (gr *GIFRenderer) renderJob(job *frameRenderJob) {
	job.renderer.Img = job.img
	job.renderer.RenderImg(gr.ColorFnGen(job.frameID), job.sliceWidth)
	job.jobStream <- job
}

func (gr *GIFRenderer) RenderParallel(maxJobs uint) *gif.GIF {
	anim := &gif.GIF{}

	// set up main render loop state variables
	// map of frames which have been rendered, but not added to anim
	renderJobs := make(map[uint]*image.Paletted)
	// channel over which a frameRenderJob reports its completed job
	jobStream := make(chan *frameRenderJob)
	var (
		runningJobs     uint
		nextFrameID     uint // id of next job to dispatch
		nextFrameNeeded uint // id of next frame to add to anim
	)
	// initial plot dimension - will need to be updated each frame
	plotWidth, plotHeight := gr.PlotWidth, gr.PlotHeight
	// slice width for each render - const
	sliceWidth := gr.ImgWidth / uint(runtime.NumCPU())

	// pools - to prevent a lot of heap-reallocations, we pool frequently
	// allocated objects

	// image renderers
	var imgRendererPool sync.Pool
	imgRendererPool.New = func() interface{} {
		return &ImgRenderer{
			ImgWidth:  gr.ImgWidth,
			ImgHeight: gr.ImgHeight,
			CX:        gr.CX,
			CY:        gr.CY,
		}
	}
	// pool of render jobs
	var renderJobPool sync.Pool
	renderJobPool.New = func() interface{} {
		return &frameRenderJob{
			jobStream:  jobStream,
			sliceWidth: sliceWidth,
		}
	}

	// render loop - dispatch jobs and compile anim
	for {
		select {
		case render := <-jobStream:
			{
				runningJobs--
				renderJobs[render.frameID] = render.img
				// return the allocated objects for reuse
				imgRendererPool.Put(render.renderer)
				renderJobPool.Put(render)
			}
		default:
			{
				if nextFrameID%20 == 0 {
					// run garbage collection every 20 iterations
					runtime.GC()
				}

				// check if a new job can be dispatched
				if runningJobs < maxJobs && nextFrameID < gr.NFrames {
					if nextFrameID != 0 {
						plotWidth *= gr.ZoomFactor
						plotHeight *= gr.ZoomFactor
					}
					// set up new renderer
					renderer := imgRendererPool.Get().(*ImgRenderer)
					renderer.PlotWidth = plotWidth
					renderer.PlotHeight = plotHeight
					// set up new job
					newJob := renderJobPool.Get().(*frameRenderJob)
					newJob.frameID = nextFrameID
					newJob.renderer = renderer
					newJob.img = image.NewPaletted(
						image.Rect(0, 0, int(gr.ImgWidth), int(gr.ImgHeight)),
						gr.Palette,
					)
					// dispatch new routine to work on this render
					go gr.renderJob(newJob)
					// update render state counters
					nextFrameID++
					runningJobs++
				}

				// check if an image can be added to the anim
				if img := renderJobs[nextFrameNeeded]; img != nil {
					anim.Image = append(anim.Image, img)
					anim.Delay = append(anim.Delay, gr.FrameDelayFn(uint(nextFrameNeeded)))
					delete(renderJobs, nextFrameNeeded)
					nextFrameNeeded++
					// log progress if required
					if gr.Progress {
						fmt.Fprintf(os.Stdout, "Rendered %d of %d frames\n", nextFrameNeeded, gr.NFrames)
					}
				}

				// check if all images have been added to the animation - done rendering
				if nextFrameNeeded == gr.NFrames {
					return anim
				}
			}
		}
	}
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
	sliceWidth := gr.ImgWidth / uint(runtime.NumCPU())
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
