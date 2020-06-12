package Fractal

type Complex struct {
	real, imag float64
}

func (i *Complex) Multiply(c *Complex) (r *Complex) {
	r = &Complex{}

	r.real = i.real*c.real - i.imag*c.imag
	r.imag = i.real*c.imag + i.imag*c.real

	return
}

func (i *Complex) Add(c *Complex) (r *Complex) {
	r = &Complex{}

	r.real = i.real + c.real
	r.imag = i.imag + c.imag

	return
}
