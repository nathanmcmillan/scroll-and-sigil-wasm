package fast

const (
	pi        = float32(3.14159265358979323846264338327950288419716939937510582097494459)
	tau       = float32(pi * 2.0)
	halfPi    = float32(pi * 0.5)
	b         = 4.0 / pi
	c         = -4.0 / (pi * pi)
	p         = 0.225
	precision = true
)

// Sin func
func Sin(x float32) float32 {
	if x > pi {
		x -= tau
	} else if x < -pi {
		x += tau
	}
	if precision {
		var y float32
		if x < 0 {
			y = b*x + (c * x * -x)
		} else {
			y = b*x + (c * x * x)
		}
		if y < 0 {
			return p*((y*-y)-y) + y
		}
		return p*((y*y)-y) + y
	}
	if x < 0 {
		return b*x + (c * x * -x)
	}
	return b*x + (c * x * x)
}

// Cos func
func Cos(x float32) float32 {
	return Sin(x + halfPi)
}
