package main

const wgt1 = 0.129633
const wgt2 = 0.175068
const w1 = -wgt1
const w2 = wgt1 + 0.5
const w3 = -wgt2
const w4 = wgt2 + 0.5

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func df(a, b int) int {
	return abs(a - b)
}

func clamp(x, floor, ceil int) int {
	return max(min(x, ceil), floor)
}

func clamp32(x, floor, ceil float32) float32 {
	return max32(min32(x, ceil), floor)
}

func matrix4() [][]int {
	m := make([][]int, 4)
	m[0] = make([]int, 4)
	m[1] = make([]int, 4)
	m[2] = make([]int, 4)
	m[3] = make([]int, 4)
	return m
}

func diagonalEdge(mat [][]int, wp []int) int {
	dw1 := wp[0]*(df(mat[0][2], mat[1][1])+df(mat[1][1], mat[2][0])+
		df(mat[1][3], mat[2][2])+df(mat[2][2], mat[3][1])) +
		wp[1]*(df(mat[0][3], mat[1][2])+df(mat[2][1], mat[3][0])) +
		wp[2]*(df(mat[0][3], mat[2][1])+df(mat[1][2], mat[3][0])) +
		wp[3]*df(mat[1][2], mat[2][1]) +
		wp[4]*(df(mat[0][2], mat[2][0])+df(mat[1][3], mat[3][1])) +
		wp[5]*(df(mat[0][1], mat[1][0])+df(mat[2][3], mat[3][2]))

	dw2 := wp[0]*(df(mat[0][1], mat[1][2])+df(mat[1][2], mat[2][3])+
		df(mat[1][0], mat[2][1])+df(mat[2][1], mat[3][2])) +
		wp[1]*(df(mat[0][0], mat[1][1])+df(mat[2][2], mat[3][3])) +
		wp[2]*(df(mat[0][0], mat[2][2])+df(mat[1][1], mat[3][3])) +
		wp[3]*df(mat[1][1], mat[2][2]) +
		wp[4]*(df(mat[1][0], mat[3][2])+df(mat[0][1], mat[2][3])) +
		wp[5]*(df(mat[0][2], mat[1][3])+df(mat[2][0], mat[3][1]))

	return (dw1 - dw2)
}

func superXBR(data []int, w, h int) []int {
	const f = 2
	outw := w * f
	outh := h * f
	wp := []int{2.0, 1.0, -1.0, 4.0, -1.0, 1.0}
	out := make([]int, outw*outh)

	r := matrix4()
	g := matrix4()
	b := matrix4()
	a := matrix4()
	Y := matrix4()
	var rf, gf, bf, af float32
	var ri, gi, bi, ai int
	var dEdge int
	var minRSample, maxRSample int
	var minGSample, maxGSample int
	var minBSample, maxBSample int
	var minASample, maxASample int
	for y := 0; y < outh; y++ {
		for x := 0; x < outw; x++ {
			cx := x / f
			cy := y / f
			for sx := -1; sx < 2; sx++ {
				for sy := -1; sy < 2; sy++ {
					csy := clamp(sy+cy, 0, h-1)
					csx := clamp(sx+cx, 0, w-1)
					sample := data[csy*w+csx]
					r[sx+1][sy+1] = sample & 0xff
					g[sx+1][sy+1] = (sample >> 8) & 0xff
					b[sx+1][sy+1] = (sample >> 16) & 0xff
					a[sx+1][sy+1] = (sample >> 24) & 0xff
					Y[sx+1][sy+1] = int(0.2126*float32(r[sx+1][sy+1]) + 0.7152*float32(g[sx+1][sy+1]) + 0.0722*float32(b[sx+1][sy+1]))
				}
			}
			minRSample = min(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			minGSample = min(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			minBSample = min(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			minASample = min(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			maxRSample = max(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			maxGSample = max(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			maxBSample = max(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			maxASample = max(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			dEdge = diagonalEdge(Y, wp)
			if dEdge <= 0 {
				rf = w1*float32(r[0][3]+r[3][0]) + w2*float32(r[1][2]+r[2][1])
				gf = w1*float32(g[0][3]+g[3][0]) + w2*float32(g[1][2]+g[2][1])
				bf = w1*float32(b[0][3]+b[3][0]) + w2*float32(b[1][2]+b[2][1])
				af = w1*float32(a[0][3]+a[3][0]) + w2*float32(a[1][2]+a[2][1])
			} else {
				rf = w1*float32(r[0][0]+r[3][3]) + w2*float32(r[1][1]+r[2][2])
				gf = w1*float32(g[0][0]+g[3][3]) + w2*float32(g[1][1]+g[2][2])
				bf = w1*float32(b[0][0]+b[3][3]) + w2*float32(b[1][1]+b[2][2])
				af = w1*float32(a[0][0]+a[3][3]) + w2*float32(a[1][1]+a[2][2])
			}
			rf = clamp32(rf, float32(minRSample), float32(maxRSample))
			gf = clamp32(gf, float32(minGSample), float32(maxGSample))
			bf = clamp32(bf, float32(minBSample), float32(maxBSample))
			af = clamp32(af, float32(minASample), float32(maxASample))
			ri = clamp(int(rf), 0, 255)
			gi = clamp(int(gf), 0, 255)
			bi = clamp(int(bf), 0, 255)
			ai = clamp(int(af), 0, 255)
			dat := data[cy*w+cx]
			out[y*outw+x] = dat
			out[y*outw+x+1] = dat
			out[(y+1)*outw+x] = dat
			out[(y+1)*outw+x+1] = (ai << 24) | (bi << 16) | (gi << 8) | ri
			x++
		}
		y++
	}

	wp[0] = 2.0
	wp[1] = 0.0
	wp[2] = 0.0
	wp[3] = 0.0
	wp[4] = 0.0
	wp[5] = 0.0

	for y := 0; y < outh; y++ {
		for x := 0; x < outw; x++ {
			for sx := -1; sx < 2; sx++ {
				for sy := -1; sy < 2; sy++ {
					csy := clamp(sx-sy+y, 0, f*h-1)
					csx := clamp(sx+sy+x, 0, f*w-1)
					sample := out[csy*outw+csx]
					r[sx+1][sy+1] = sample & 0xff
					g[sx+1][sy+1] = (sample >> 8) & 0xff
					b[sx+1][sy+1] = (sample >> 16) & 0xff
					a[sx+1][sy+1] = (sample >> 24) & 0xff
					Y[sx+1][sy+1] = int(0.2126*float32(r[sx+1][sy+1]) + 0.7152*float32(g[sx+1][sy+1]) + 0.0722*float32(b[sx+1][sy+1]))
				}
			}
			minRSample = min(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			minGSample = min(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			minBSample = min(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			minASample = min(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			maxRSample = max(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			maxGSample = max(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			maxBSample = max(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			maxASample = max(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			dEdge = diagonalEdge(Y, wp)
			if dEdge <= 0 {
				rf = w3*float32(r[0][3]+r[3][0]) + w4*float32(r[1][2]+r[2][1])
				gf = w3*float32(g[0][3]+g[3][0]) + w4*float32(g[1][2]+g[2][1])
				bf = w3*float32(b[0][3]+b[3][0]) + w4*float32(b[1][2]+b[2][1])
				af = w3*float32(a[0][3]+a[3][0]) + w4*float32(a[1][2]+a[2][1])
			} else {
				rf = w3*float32(r[0][0]+r[3][3]) + w4*float32(r[1][1]+r[2][2])
				gf = w3*float32(g[0][0]+g[3][3]) + w4*float32(g[1][1]+g[2][2])
				bf = w3*float32(b[0][0]+b[3][3]) + w4*float32(b[1][1]+b[2][2])
				af = w3*float32(a[0][0]+a[3][3]) + w4*float32(a[1][1]+a[2][2])
			}
			rf = clamp32(rf, float32(minRSample), float32(maxRSample))
			gf = clamp32(gf, float32(minGSample), float32(maxGSample))
			bf = clamp32(bf, float32(minBSample), float32(maxBSample))
			af = clamp32(af, float32(minASample), float32(maxASample))
			ri = clamp(int(rf), 0, 255)
			gi = clamp(int(gf), 0, 255)
			bi = clamp(int(bf), 0, 255)
			ai = clamp(int(af), 0, 255)
			out[y*outw+x+1] = (ai << 24) | (bi << 16) | (gi << 8) | ri

			for sx := -1; sx < 2; sx++ {
				for sy := -1; sy < 2; sy++ {
					csy := clamp(sx-sy+1+y, 0, f*h-1)
					csx := clamp(sx+sy-1+x, 0, f*w-1)
					sample := out[csy*outw+csx]
					r[sx+1][sy+1] = ((sample) >> 0) & 0xff
					g[sx+1][sy+1] = ((sample) >> 8) & 0xff
					b[sx+1][sy+1] = ((sample) >> 16) & 0xff
					a[sx+1][sy+1] = ((sample) >> 24) & 0xff
					Y[sx+1][sy+1] = int(0.2126*float32(r[sx+1][sy+1]) + 0.7152*float32(g[sx+1][sy+1]) + 0.0722*float32(b[sx+1][sy+1]))
				}
			}
			dEdge = diagonalEdge(Y, wp)
			if dEdge <= 0 {
				rf = w3*float32(r[0][3]+r[3][0]) + w4*float32(r[1][2]+r[2][1])
				gf = w3*float32(g[0][3]+g[3][0]) + w4*float32(g[1][2]+g[2][1])
				bf = w3*float32(b[0][3]+b[3][0]) + w4*float32(b[1][2]+b[2][1])
				af = w3*float32(a[0][3]+a[3][0]) + w4*float32(a[1][2]+a[2][1])
			} else {
				rf = w3*float32(r[0][0]+r[3][3]) + w4*float32(r[1][1]+r[2][2])
				gf = w3*float32(g[0][0]+g[3][3]) + w4*float32(g[1][1]+g[2][2])
				bf = w3*float32(b[0][0]+b[3][3]) + w4*float32(b[1][1]+b[2][2])
				af = w3*float32(a[0][0]+a[3][3]) + w4*float32(a[1][1]+a[2][2])
			}
			rf = clamp32(rf, float32(minRSample), float32(maxRSample))
			gf = clamp32(gf, float32(minGSample), float32(maxGSample))
			bf = clamp32(bf, float32(minBSample), float32(maxBSample))
			af = clamp32(af, float32(minASample), float32(maxASample))
			ri = clamp(int(rf), 0, 255)
			gi = clamp(int(gf), 0, 255)
			bi = clamp(int(bf), 0, 255)
			ai = clamp(int(af), 0, 255)
			out[(y+1)*outw+x] = (ai << 24) | (bi << 16) | (gi << 8) | ri
			x++
		}
		y++
	}

	wp[0] = 2.0
	wp[1] = 1.0
	wp[2] = -1.0
	wp[3] = 4.0
	wp[4] = -1.0
	wp[5] = 1.0

	for y := outh - 1; y > 0; y-- {
		for x := outw - 1; x > 0; x-- {
			for sx := -2; sx < 1; sx++ {
				for sy := -2; sy < 1; sy++ {
					csy := clamp(sy+y, 0, f*h-1)
					csx := clamp(sx+x, 0, f*w-1)
					sample := out[csy*outw+csx]
					r[sx+2][sy+2] = ((sample) >> 0) & 0xff
					g[sx+2][sy+2] = ((sample) >> 8) & 0xff
					b[sx+2][sy+2] = ((sample) >> 16) & 0xff
					a[sx+2][sy+2] = ((sample) >> 24) & 0xff
					Y[sx+2][sy+2] = int(0.2126*float32(r[sx+2][sy+2]) + 0.7152*float32(g[sx+2][sy+2]) + 0.0722*float32(b[sx+2][sy+2]))
				}
			}
			minRSample = min(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			minGSample = min(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			minBSample = min(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			minASample = min(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			maxRSample = max(r[1][1], min(r[2][1], min(r[1][2], r[2][2])))
			maxGSample = max(g[1][1], min(g[2][1], min(g[1][2], g[2][2])))
			maxBSample = max(b[1][1], min(b[2][1], min(b[1][2], b[2][2])))
			maxASample = max(a[1][1], min(a[2][1], min(a[1][2], a[2][2])))
			dEdge := diagonalEdge(Y, wp)
			if dEdge <= 0 {
				rf = w1*float32(r[0][3]+r[3][0]) + w2*float32(r[1][2]+r[2][1])
				gf = w1*float32(g[0][3]+g[3][0]) + w2*float32(g[1][2]+g[2][1])
				bf = w1*float32(b[0][3]+b[3][0]) + w2*float32(b[1][2]+b[2][1])
				af = w1*float32(a[0][3]+a[3][0]) + w2*float32(a[1][2]+a[2][1])
			} else {
				rf = w1*float32(r[0][0]+r[3][3]) + w2*float32(r[1][1]+r[2][2])
				gf = w1*float32(g[0][0]+g[3][3]) + w2*float32(g[1][1]+g[2][2])
				bf = w1*float32(b[0][0]+b[3][3]) + w2*float32(b[1][1]+b[2][2])
				af = w1*float32(a[0][0]+a[3][3]) + w2*float32(a[1][1]+a[2][2])
			}
			rf = clamp32(rf, float32(minRSample), float32(maxRSample))
			gf = clamp32(gf, float32(minGSample), float32(maxGSample))
			bf = clamp32(bf, float32(minBSample), float32(maxBSample))
			af = clamp32(af, float32(minASample), float32(maxASample))
			ri = clamp(int(rf), 0, 255)
			gi = clamp(int(gf), 0, 255)
			bi = clamp(int(bf), 0, 255)
			ai = clamp(int(af), 0, 255)
			out[y*outw+x] = (ai << 24) | (bi << 16) | (gi << 8) | ri
		}
	}

	return out
}
