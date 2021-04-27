package fractal

// EscapeIterFn returns the iteration at which the absolute value of the
// iterates of a complex function (parametrized by `c`) exceeds some escape
// radius.
type EscapeIterFn func(c complex128) int

// JuliaSetFn defines a function f(z) of the form f(x) = z**exp + c.
// The parameter `c` is variable.
type JuliaSetFn struct {
	Exp         int
	MaxIters    int
	EscapeRad   float64
	InitIterate complex128
}

// The EscapeIter method returns the iteration at which the absolute value of
// the iterates of the JuliaSetFn function (with parameter `c`) exceeds the
// `escapeRadius`.
func (f *JuliaSetFn) EscapeIter(c complex128) int {
	v := f.InitIterate
	escapeRadius := f.EscapeRad * f.EscapeRad
	for n := 0; n < f.MaxIters; n++ {
		for i := 1; i < f.Exp; i++ {
			v *= v
		}
		v += c
		abs := real(v)*real(v) + imag(v)*imag(v)
		if abs > escapeRadius {
			return n
		}
	}
	return f.MaxIters
}
