package fractal

import (
	"image/color"
	//	"math"
)

// ColorPalette translates the number of iterations to a display color.
type ColorPalette interface {
	PixelColor(n int) color.Color
	Precompute()
	Palette() []color.Color
}

// BWPalette is a black and white ColorPalette.
type BWPalette struct {
	Inverse      bool
	MaxIters     int
	Contrast     uint8
	ColorPalette []color.Color
}

func (p *BWPalette) Precompute() {
	p.ColorPalette = make([]color.Color, 0xff)
	for i := uint8(0); i < 0xff; i++ {
		p.ColorPalette[i] = color.Gray{0xff - (p.Contrast * i)}
	}
}

func (p *BWPalette) PixelColor(iters int) color.Color {
	if iters < p.MaxIters {
		return p.ColorPalette[iters%0xff]
	}
	return color.Black
}

func (p *BWPalette) Palette() []color.Color {
	return p.ColorPalette
}

// UltraFractalPalette is a color palette based on UltraFractal's
// implementaiton.
type UltraFractalPalette struct {
	MaxIters     int
	ColorPalette []color.Color
}

func (p *UltraFractalPalette) Precompute() {
	p.ColorPalette = []color.Color{
		color.RGBA{66, 30, 15, 0xff},
		color.RGBA{25, 7, 26, 0xff},
		color.RGBA{9, 1, 47, 0xff},
		color.RGBA{4, 4, 73, 0xff},
		color.RGBA{0, 7, 100, 0xff},
		color.RGBA{12, 44, 138, 0xff},
		color.RGBA{24, 82, 177, 0xff},
		color.RGBA{57, 125, 209, 0xff},
		color.RGBA{134, 181, 229, 0xff},
		color.RGBA{211, 236, 248, 0xff},
		color.RGBA{241, 233, 191, 0xff},
		color.RGBA{248, 201, 95, 0xff},
		color.RGBA{255, 170, 0, 0xff},
		color.RGBA{204, 128, 0, 0xff},
		color.RGBA{153, 87, 0, 0xff},
		color.RGBA{106, 52, 3, 0xff},
	}
}

func (p *UltraFractalPalette) PixelColor(iters int) color.Color {
	if iters < p.MaxIters {
		return p.ColorPalette[iters%16]
	}
	return color.RGBA{0, 0, 0, 0xff}
}

func (p *UltraFractalPalette) Palette() []color.Color {
	return p.ColorPalette
}
