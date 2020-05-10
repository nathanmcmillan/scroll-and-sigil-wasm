package matrix

import "math"

var (
	matrixTemp   = make([]float32, 16)
	matrixCopied = make([]float32, 16)
)

// Identity func
func Identity(matrix []float32) {
	matrix[0] = 1.0
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = 1.0
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = 0.0
	matrix[10] = 1.0
	matrix[11] = 0.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

// Orthographic func
func Orthographic(matrix []float32, left, right, bottom, top, near, far float32) {
	matrix[0] = 2.0 / (right - left)
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = 2.0 / (top - bottom)
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = 0.0
	matrix[10] = -2.0 / (far - near)
	matrix[11] = 0.0

	matrix[12] = -(right + left) / (right - left)
	matrix[13] = -(top + bottom) / (top - bottom)
	matrix[14] = -(far + near) / (far - near)
	matrix[15] = 1.0
}

// Perspective func
func Perspective(matrix []float32, fov, near, far, aspect float32) {
	top := near * float32(math.Tan(float64(fov)*math.Pi/360.0))
	bottom := -top
	left := bottom * aspect
	right := top * aspect

	Frustum(matrix, left, right, bottom, top, near, far)
}

// Frustum func
func Frustum(matrix []float32, left, right, bottom, top, near, far float32) {
	matrix[0] = (2.0 * near) / (right - left)
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = (2.0 * near) / (top - bottom)
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = (right + left) / (right - left)
	matrix[9] = (top + bottom) / (top - bottom)
	matrix[10] = -(far + near) / (far - near)
	matrix[11] = -1.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = -(2.0 * far * near) / (far - near)
	matrix[15] = 0.0
}

// Translate func
func Translate(matrix []float32, x, y, z float32) {
	matrix[12] = x*matrix[0] + y*matrix[4] + z*matrix[8] + matrix[12]
	matrix[13] = x*matrix[1] + y*matrix[5] + z*matrix[9] + matrix[13]
	matrix[14] = x*matrix[2] + y*matrix[6] + z*matrix[10] + matrix[14]
	matrix[15] = x*matrix[3] + y*matrix[7] + z*matrix[11] + matrix[15]
}

// TranslateFromView func
func TranslateFromView(matrix, view []float32, x, y, z float32) {
	matrix[0] = view[0]
	matrix[1] = view[1]
	matrix[2] = view[2]
	matrix[3] = view[3]
	matrix[4] = view[4]
	matrix[5] = view[5]
	matrix[6] = view[6]
	matrix[7] = view[7]
	matrix[8] = view[8]
	matrix[9] = view[9]
	matrix[10] = view[10]
	matrix[11] = view[11]
	matrix[12] = x*view[0] + y*view[4] + z*view[8] + view[12]
	matrix[13] = x*view[1] + y*view[5] + z*view[9] + view[13]
	matrix[14] = x*view[2] + y*view[6] + z*view[10] + view[14]
	matrix[15] = x*view[3] + y*view[7] + z*view[11] + view[15]
}

// RotateX func
func RotateX(matrix []float32, radian float32) {
	sin := float32(math.Sin(float64(radian)))
	cos := float32(math.Cos(float64(radian)))

	matrixTemp[0] = 1.0
	matrixTemp[1] = 0.0
	matrixTemp[2] = 0.0
	matrixTemp[3] = 0.0

	matrixTemp[4] = 0.0
	matrixTemp[5] = cos
	matrixTemp[6] = sin
	matrixTemp[7] = 0.0

	matrixTemp[8] = 0.0
	matrixTemp[9] = -sin
	matrixTemp[10] = cos
	matrixTemp[11] = 0.0

	matrixTemp[12] = 0.0
	matrixTemp[13] = 0.0
	matrixTemp[14] = 0.0
	matrixTemp[15] = 1.0

	for i := 0; i < 16; i++ {
		matrixCopied[i] = matrix[i]
	}

	Multiply(matrix, matrixCopied, matrixTemp)
}

// RotateY func
func RotateY(matrix []float32, radian float32) {
	sin := float32(math.Sin(float64(radian)))
	cos := float32(math.Cos(float64(radian)))

	matrixTemp[0] = cos
	matrixTemp[1] = 0.0
	matrixTemp[2] = -sin
	matrixTemp[3] = 0.0

	matrixTemp[4] = 0.0
	matrixTemp[5] = 1.0
	matrixTemp[6] = 0.0
	matrixTemp[7] = 0.0

	matrixTemp[8] = sin
	matrixTemp[9] = 0.0
	matrixTemp[10] = cos
	matrixTemp[11] = 0.0

	matrixTemp[12] = 0.0
	matrixTemp[13] = 0.0
	matrixTemp[14] = 0.0
	matrixTemp[15] = 1.0

	for i := 0; i < 16; i++ {
		matrixCopied[i] = matrix[i]
	}

	Multiply(matrix, matrixCopied, matrixTemp)
}

// RotateZ func
func RotateZ(matrix []float32, radian float32) {
	sin := float32(math.Sin(float64(radian)))
	cos := float32(math.Cos(float64(radian)))

	matrixTemp[0] = cos
	matrixTemp[1] = sin
	matrixTemp[2] = 0.0
	matrixTemp[3] = 0.0

	matrixTemp[4] = -sin
	matrixTemp[5] = cos
	matrixTemp[6] = 0.0
	matrixTemp[7] = 0.0

	matrixTemp[8] = 0.0
	matrixTemp[9] = 0.0
	matrixTemp[10] = 1.0
	matrixTemp[11] = 0.0

	matrixTemp[12] = 0.0
	matrixTemp[13] = 0.0
	matrixTemp[14] = 0.0
	matrixTemp[15] = 1.0

	for i := 0; i < 16; i++ {
		matrixCopied[i] = matrix[i]
	}

	Multiply(matrix, matrixCopied, matrixTemp)
}

// Multiply func
func Multiply(matrix, b, c []float32) {
	matrix[0] = b[0]*c[0] + b[4]*c[1] + b[8]*c[2] + b[12]*c[3]
	matrix[1] = b[1]*c[0] + b[5]*c[1] + b[9]*c[2] + b[13]*c[3]
	matrix[2] = b[2]*c[0] + b[6]*c[1] + b[10]*c[2] + b[14]*c[3]
	matrix[3] = b[3]*c[0] + b[7]*c[1] + b[11]*c[2] + b[15]*c[3]

	matrix[4] = b[0]*c[4] + b[4]*c[5] + b[8]*c[6] + b[12]*c[7]
	matrix[5] = b[1]*c[4] + b[5]*c[5] + b[9]*c[6] + b[13]*c[7]
	matrix[6] = b[2]*c[4] + b[6]*c[5] + b[10]*c[6] + b[14]*c[7]
	matrix[7] = b[3]*c[4] + b[7]*c[5] + b[11]*c[6] + b[15]*c[7]

	matrix[8] = b[0]*c[8] + b[4]*c[9] + b[8]*c[10] + b[12]*c[11]
	matrix[9] = b[1]*c[8] + b[5]*c[9] + b[9]*c[10] + b[13]*c[11]
	matrix[10] = b[2]*c[8] + b[6]*c[9] + b[10]*c[10] + b[14]*c[11]
	matrix[11] = b[3]*c[8] + b[7]*c[9] + b[11]*c[10] + b[15]*c[11]

	matrix[12] = b[0]*c[12] + b[4]*c[13] + b[8]*c[14] + b[12]*c[15]
	matrix[13] = b[1]*c[12] + b[5]*c[13] + b[9]*c[14] + b[13]*c[15]
	matrix[14] = b[2]*c[12] + b[6]*c[13] + b[10]*c[14] + b[14]*c[15]
	matrix[15] = b[3]*c[12] + b[7]*c[13] + b[11]*c[14] + b[15]*c[15]
}

// Inverse func
func Inverse(matrix, b []float32) {
	for i := 0; i < 4; i++ {
		matrixCopied[i+0] = b[i*4+0]
		matrixCopied[i+4] = b[i*4+1]
		matrixCopied[i+8] = b[i*4+2]
		matrixCopied[i+12] = b[i*4+3]
	}

	matrixTemp[0] = matrixCopied[10] * matrixCopied[15]
	matrixTemp[1] = matrixCopied[11] * matrixCopied[14]
	matrixTemp[2] = matrixCopied[9] * matrixCopied[15]
	matrixTemp[3] = matrixCopied[11] * matrixCopied[13]
	matrixTemp[4] = matrixCopied[9] * matrixCopied[14]
	matrixTemp[5] = matrixCopied[10] * matrixCopied[13]
	matrixTemp[6] = matrixCopied[8] * matrixCopied[15]
	matrixTemp[7] = matrixCopied[11] * matrixCopied[12]
	matrixTemp[8] = matrixCopied[8] * matrixCopied[14]
	matrixTemp[9] = matrixCopied[10] * matrixCopied[12]
	matrixTemp[10] = matrixCopied[8] * matrixCopied[13]
	matrixTemp[11] = matrixCopied[9] * matrixCopied[12]

	matrix[0] = matrixTemp[0]*matrixCopied[5] + matrixTemp[3]*matrixCopied[6] + matrixTemp[4]*matrixCopied[7]
	matrix[0] -= matrixTemp[1]*matrixCopied[5] + matrixTemp[2]*matrixCopied[6] + matrixTemp[5]*matrixCopied[7]
	matrix[1] = matrixTemp[1]*matrixCopied[4] + matrixTemp[6]*matrixCopied[6] + matrixTemp[9]*matrixCopied[7]
	matrix[1] -= matrixTemp[0]*matrixCopied[4] + matrixTemp[7]*matrixCopied[6] + matrixTemp[8]*matrixCopied[7]
	matrix[2] = matrixTemp[2]*matrixCopied[4] + matrixTemp[7]*matrixCopied[5] + matrixTemp[10]*matrixCopied[7]
	matrix[2] -= matrixTemp[3]*matrixCopied[4] + matrixTemp[6]*matrixCopied[5] + matrixTemp[11]*matrixCopied[7]
	matrix[3] = matrixTemp[5]*matrixCopied[4] + matrixTemp[8]*matrixCopied[5] + matrixTemp[11]*matrixCopied[6]
	matrix[3] -= matrixTemp[4]*matrixCopied[4] + matrixTemp[9]*matrixCopied[5] + matrixTemp[10]*matrixCopied[6]
	matrix[4] = matrixTemp[1]*matrixCopied[1] + matrixTemp[2]*matrixCopied[2] + matrixTemp[5]*matrixCopied[3]
	matrix[4] -= matrixTemp[0]*matrixCopied[1] + matrixTemp[3]*matrixCopied[2] + matrixTemp[4]*matrixCopied[3]
	matrix[5] = matrixTemp[0]*matrixCopied[0] + matrixTemp[7]*matrixCopied[2] + matrixTemp[8]*matrixCopied[3]
	matrix[5] -= matrixTemp[1]*matrixCopied[0] + matrixTemp[6]*matrixCopied[2] + matrixTemp[9]*matrixCopied[3]
	matrix[6] = matrixTemp[3]*matrixCopied[0] + matrixTemp[6]*matrixCopied[1] + matrixTemp[11]*matrixCopied[3]
	matrix[6] -= matrixTemp[2]*matrixCopied[0] + matrixTemp[7]*matrixCopied[1] + matrixTemp[10]*matrixCopied[3]
	matrix[7] = matrixTemp[4]*matrixCopied[0] + matrixTemp[9]*matrixCopied[1] + matrixTemp[10]*matrixCopied[2]
	matrix[7] -= matrixTemp[5]*matrixCopied[0] + matrixTemp[8]*matrixCopied[1] + matrixTemp[11]*matrixCopied[2]

	matrixTemp[0] = matrixCopied[2] * matrixCopied[7]
	matrixTemp[1] = matrixCopied[3] * matrixCopied[6]
	matrixTemp[2] = matrixCopied[1] * matrixCopied[7]
	matrixTemp[3] = matrixCopied[3] * matrixCopied[5]
	matrixTemp[4] = matrixCopied[1] * matrixCopied[6]
	matrixTemp[5] = matrixCopied[2] * matrixCopied[5]
	matrixTemp[6] = matrixCopied[0] * matrixCopied[7]
	matrixTemp[7] = matrixCopied[3] * matrixCopied[4]
	matrixTemp[8] = matrixCopied[0] * matrixCopied[6]
	matrixTemp[9] = matrixCopied[2] * matrixCopied[4]
	matrixTemp[10] = matrixCopied[0] * matrixCopied[5]
	matrixTemp[11] = matrixCopied[1] * matrixCopied[4]

	matrix[8] = matrixTemp[0]*matrixCopied[13] + matrixTemp[3]*matrixCopied[14] + matrixTemp[4]*matrixCopied[15]
	matrix[8] -= matrixTemp[1]*matrixCopied[13] + matrixTemp[2]*matrixCopied[14] + matrixTemp[5]*matrixCopied[15]
	matrix[9] = matrixTemp[1]*matrixCopied[12] + matrixTemp[6]*matrixCopied[14] + matrixTemp[9]*matrixCopied[15]
	matrix[9] -= matrixTemp[0]*matrixCopied[12] + matrixTemp[7]*matrixCopied[14] + matrixTemp[8]*matrixCopied[15]
	matrix[10] = matrixTemp[2]*matrixCopied[12] + matrixTemp[7]*matrixCopied[13] + matrixTemp[10]*matrixCopied[15]
	matrix[10] -= matrixTemp[3]*matrixCopied[12] + matrixTemp[6]*matrixCopied[13] + matrixTemp[11]*matrixCopied[15]
	matrix[11] = matrixTemp[5]*matrixCopied[12] + matrixTemp[8]*matrixCopied[13] + matrixTemp[11]*matrixCopied[14]
	matrix[11] -= matrixTemp[4]*matrixCopied[12] + matrixTemp[9]*matrixCopied[13] + matrixTemp[10]*matrixCopied[14]
	matrix[12] = matrixTemp[2]*matrixCopied[10] + matrixTemp[5]*matrixCopied[11] + matrixTemp[1]*matrixCopied[9]
	matrix[12] -= matrixTemp[4]*matrixCopied[11] + matrixTemp[0]*matrixCopied[9] + matrixTemp[3]*matrixCopied[10]
	matrix[13] = matrixTemp[8]*matrixCopied[11] + matrixTemp[0]*matrixCopied[8] + matrixTemp[7]*matrixCopied[10]
	matrix[13] -= matrixTemp[6]*matrixCopied[10] + matrixTemp[9]*matrixCopied[11] + matrixTemp[1]*matrixCopied[8]
	matrix[14] = matrixTemp[6]*matrixCopied[9] + matrixTemp[11]*matrixCopied[11] + matrixTemp[3]*matrixCopied[8]
	matrix[14] -= matrixTemp[10]*matrixCopied[11] + matrixTemp[2]*matrixCopied[8] + matrixTemp[7]*matrixCopied[9]
	matrix[15] = matrixTemp[10]*matrixCopied[10] + matrixTemp[4]*matrixCopied[8] + matrixTemp[9]*matrixCopied[9]
	matrix[15] -= matrixTemp[8]*matrixCopied[9] + matrixTemp[11]*matrixCopied[10] + matrixTemp[5]*matrixCopied[8]

	det := 1.0 / (matrixCopied[0]*matrix[0] + matrixCopied[1]*matrix[1] + matrixCopied[2]*matrix[2] + matrixCopied[3]*matrix[3])

	for i := 0; i < 16; i++ {
		matrix[i] = matrix[i] * det
	}
}

// MultiplyVector func
func MultiplyVector(matrix []float32, x, y, z float32) (float32, float32, float32) {
	a := matrix[0]*x + matrix[4]*y + matrix[8]*z + matrix[12]
	b := matrix[1]*x + matrix[5]*y + matrix[9]*z + matrix[13]
	c := matrix[2]*x + matrix[6]*y + matrix[10]*z + matrix[14]
	return a, b, c
}
