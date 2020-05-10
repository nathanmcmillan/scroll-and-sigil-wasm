package main

import (
	"math"

	"./graphics"
)

const (
	fullOcclusion    = uint8(0)
	partialOcclusion = uint8(1)
	noOcclusion      = uint8(2)
)

var (
	occlusionSliceA [3]int32
	occlusionSliceB [3]int32
)

type occluder struct {
	frustum   [6][4]float32
	queue     []*block
	gotoBlock []bool
	queueFrom []int
	viewNum   int
	viewable  []*block
	queuePos  int
	queueNum  int
}

func occluderInit(size int) *occluder {
	occluder := &occluder{}
	occluder.viewable = make([]*block, size)
	occluder.queue = make([]*block, size)
	occluder.gotoBlock = make([]bool, size)
	occluder.queueFrom = make([]int, size)
	return occluder
}

func (me *block) occlusion() {
	for sideA := uint(0); sideA < 6; sideA++ {
		ax := blockSliceX[sideA]
		ay := blockSliceY[sideA]
		az := blockSliceZ[sideA]
		for sideB := sideA + 1; sideB < 6; sideB++ {
			bx := blockSliceX[sideB]
			by := blockSliceY[sideB]
			bz := blockSliceZ[sideB]
			if blockSliceTowards[sideA] > 0 {
				occlusionSliceA[2] = BlockSize - 1
			} else {
				occlusionSliceA[2] = 0
			}
			if blockSliceTowards[sideB] > 0 {
				occlusionSliceB[2] = BlockSize - 1
			} else {
				occlusionSliceB[2] = 0
			}
		loop:
			for occlusionSliceA[1] = 0; occlusionSliceA[1] < BlockSize; occlusionSliceA[1]++ {
				for occlusionSliceA[0] = 0; occlusionSliceA[0] < BlockSize; occlusionSliceA[0]++ {
					for occlusionSliceB[1] = 0; occlusionSliceB[1] < BlockSize; occlusionSliceB[1]++ {
						for occlusionSliceB[0] = 0; occlusionSliceB[0] < BlockSize; occlusionSliceB[0]++ {
							fromX := float32(occlusionSliceA[ax]) + 0.5
							fromY := float32(occlusionSliceA[ay]) + 0.5
							fromZ := float32(occlusionSliceA[az]) + 0.5
							toX := float32(occlusionSliceB[bx]) + 0.5
							toY := float32(occlusionSliceB[by]) + 0.5
							toZ := float32(occlusionSliceB[bz]) + 0.5
							if castBlock(me, fromX, fromY, fromZ, toX, toY, toZ) {
								me.visibility[sideA*6+sideB] = true
								break loop
							}
						}
					}
				}
			}
		}
	}
}

func (me *occluder) prepareFrustum(g *graphics.RenderSystem) {
	// left
	me.frustum[0][0] = g.ModelViewProject[3] + g.ModelViewProject[0]
	me.frustum[0][1] = g.ModelViewProject[7] + g.ModelViewProject[4]
	me.frustum[0][2] = g.ModelViewProject[11] + g.ModelViewProject[8]
	me.frustum[0][3] = g.ModelViewProject[15] + g.ModelViewProject[12]

	// right
	me.frustum[1][0] = g.ModelViewProject[3] - g.ModelViewProject[0]
	me.frustum[1][1] = g.ModelViewProject[7] - g.ModelViewProject[4]
	me.frustum[1][2] = g.ModelViewProject[11] - g.ModelViewProject[8]
	me.frustum[1][3] = g.ModelViewProject[15] - g.ModelViewProject[12]

	// top
	me.frustum[2][0] = g.ModelViewProject[3] - g.ModelViewProject[1]
	me.frustum[2][1] = g.ModelViewProject[7] - g.ModelViewProject[5]
	me.frustum[2][2] = g.ModelViewProject[11] - g.ModelViewProject[9]
	me.frustum[2][3] = g.ModelViewProject[15] - g.ModelViewProject[13]

	// bottom
	me.frustum[3][0] = g.ModelViewProject[3] + g.ModelViewProject[1]
	me.frustum[3][1] = g.ModelViewProject[7] + g.ModelViewProject[5]
	me.frustum[3][2] = g.ModelViewProject[11] + g.ModelViewProject[9]
	me.frustum[3][3] = g.ModelViewProject[15] + g.ModelViewProject[13]

	// near
	me.frustum[4][0] = g.ModelViewProject[3] + g.ModelViewProject[2]
	me.frustum[4][1] = g.ModelViewProject[7] + g.ModelViewProject[6]
	me.frustum[4][2] = g.ModelViewProject[11] + g.ModelViewProject[10]
	me.frustum[4][3] = g.ModelViewProject[15] + g.ModelViewProject[14]

	// far
	me.frustum[5][0] = g.ModelViewProject[3] - g.ModelViewProject[2]
	me.frustum[5][1] = g.ModelViewProject[7] - g.ModelViewProject[6]
	me.frustum[5][2] = g.ModelViewProject[11] - g.ModelViewProject[10]
	me.frustum[5][3] = g.ModelViewProject[15] - g.ModelViewProject[14]

	me.normalizePlanes()
}

func (me *occluder) normalizePlanes() {
	for i := 0; i < 6; i++ {
		n := float32(math.Sqrt(float64(me.frustum[i][0]*me.frustum[i][0] + me.frustum[i][1]*me.frustum[i][1] + me.frustum[i][2]*me.frustum[i][2])))
		me.frustum[i][0] /= n
		me.frustum[i][1] /= n
		me.frustum[i][2] /= n
		me.frustum[i][3] /= n
	}
}

func (me *occluder) visit(world *world, from int, B *block, to int) {
	x := B.x
	y := B.y
	z := B.z
	switch to {
	case WorldPositiveX:
		x++
		if x == world.width {
			return
		}
	case WorldNegativeX:
		x--
		if x == -1 {
			return
		}
	case WorldPositiveY:
		y++
		if y == world.height {
			return
		}
	case WorldNegativeY:
		y--
		if y == -1 {
			return
		}
	case WorldPositiveZ:
		z++
		if z == world.length {
			return
		}
	case WorldNegativeZ:
		z--
		if z == -1 {
			return
		}
	}
	index := x + y*world.width + z*world.slice
	if me.gotoBlock[index] == false {
		return
	}
	if from >= 0 {
		switch from {
		case WorldPositiveX:
			from = WorldNegativeX
		case WorldNegativeX:
			from = WorldPositiveX
		case WorldPositiveY:
			from = WorldNegativeY
		case WorldNegativeY:
			from = WorldPositiveY
		case WorldPositiveZ:
			from = WorldNegativeZ
		case WorldNegativeZ:
			from = WorldPositiveZ
		}
		var sideA, sideB int
		if from < to {
			sideA = from
			sideB = to
		} else {
			sideA = to
			sideB = from
		}
		if B.visibility[sideA*6+sideB] == false {
			return
		}
	}
	me.gotoBlock[index] = false
	C := &world.blocks[index]
	posX := float32(C.x * BlockSize)
	posY := float32(C.y * BlockSize)
	posZ := float32(C.z * BlockSize)
	box := me.inBox(
		posX+BlockSize, posY+BlockSize, posZ+BlockSize,
		posZ, posY, posZ)
	if box == noOcclusion {
		return
	}

	queue := me.queuePos + me.queueNum
	if queue >= world.all {
		queue -= world.all
	}

	me.queue[queue] = C
	me.queueFrom[queue] = to
	me.queueNum++
}

func (me *occluder) inBox(posX, posY, posZ, negX, negY, negZ float32) uint8 {
	var pvx, pvy, pvz float32
	var nvx, nvy, nvz float32
	result := fullOcclusion
	for i := 0; i < 6; i++ {
		plane := &me.frustum[i]
		if plane[0] > 0 {
			pvx = posX
			nvx = negX
		} else {
			pvx = negX
			nvx = posX
		}
		if plane[1] > 0 {
			pvy = posY
			nvy = negY
		} else {
			pvy = negY
			nvy = posY
		}
		if plane[2] > 0 {
			pvz = posZ
			nvz = negZ
		} else {
			pvz = negZ
			nvz = posZ
		}
		if pvx*plane[0]+pvy*plane[1]+pvz*plane[2]+plane[3] < 0 {
			return noOcclusion
		}
		if nvx*plane[0]+nvy*plane[1]+nvz*plane[2]+plane[3] < 0 {
			result = partialOcclusion
		}
	}
	return result
}

func (me *occluder) search(world *world, lx, ly, lz int) {
	me.viewNum = 0
	all := world.all

	if lx < 0 || lx >= world.width || ly < 0 || ly >= world.height || lz < 0 || lz >= world.length {
		for me.viewNum < all {
			me.viewable[me.viewNum] = &world.blocks[me.viewNum]
			me.viewNum++
		}
		return
	}

	me.queuePos = 0
	me.queueNum = 1
	me.queue[0] = &world.blocks[lx+ly*world.width+lz*world.slice]
	me.queueFrom[0] = -1

	for i := 0; i < all; i++ {
		me.gotoBlock[i] = true
	}

	for me.queueNum > 0 {
		B := me.queue[me.queuePos]
		from := me.queueFrom[me.queuePos]

		me.viewable[me.viewNum] = B
		me.viewNum++

		me.queuePos++
		if me.queuePos == all {
			me.queuePos = 0
		}

		me.queueNum--

		if from != WorldNegativeX {
			me.visit(world, from, B, WorldPositiveX)
		}

		if from != WorldPositiveX {
			me.visit(world, from, B, WorldNegativeX)
		}

		if from != WorldNegativeY {
			me.visit(world, from, B, WorldPositiveY)
		}

		if from != WorldPositiveY {
			me.visit(world, from, B, WorldNegativeY)
		}

		if from != WorldNegativeZ {
			me.visit(world, from, B, WorldPositiveZ)
		}

		if from != WorldPositiveZ {
			me.visit(world, from, B, WorldNegativeZ)
		}
	}
}
