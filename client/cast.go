package main

import (
	"math"
)

func castBlock(block *block, fromX, fromY, fromZ, toX, toY, toZ float32) bool {
	x := int(fromX)
	y := int(fromY)
	z := int(fromZ)
	var deltaX, deltaY, deltaZ float32
	var nextX, nextY, nextZ float32
	var stepX, stepY, stepZ int
	dx := toX - fromX
	if dx == 0 {
		stepX = 0
		nextX = math.MaxFloat32
	} else if dx > 0 {
		stepX = 1
		deltaX = 1.0 / dx
		nextX = (1.0 + float32(x) - fromX) * deltaX
	} else {
		stepX = -1
		deltaX = 1.0 / -dx
		nextX = (fromX - float32(x)) * deltaX
	}
	dy := toY - fromY
	if dy == 0 {
		stepY = 0
		nextY = math.MaxFloat32
	} else if dy > 0 {
		stepY = 1
		deltaY = 1.0 / dy
		nextY = (1.0 + float32(y) - fromY) * deltaY
	} else {
		stepY = -1
		deltaY = 1.0 / -dy
		nextY = (fromY - float32(y)) * deltaY
	}
	dz := toZ - fromZ
	if dz == 0 {
		stepZ = 0
		nextZ = math.MaxFloat32
	} else if dz > 0 {
		stepZ = 1
		deltaZ = 1.0 / dz
		nextZ = (1.0 + float32(z) - fromZ) * deltaZ
	} else {
		stepZ = -1
		deltaZ = 1.0 / -dz
		nextZ = (fromZ - float32(z)) * deltaZ
	}
	for {
		if TileClosed[block.tiles[x+(y<<BlockShift)+(z<<BlockShiftSlice)].typeOf] {
			return false
		} else if x == int(toX) && y == int(toY) && z == int(toZ) {
			return true
		}
		if nextX < nextY {
			if nextX < nextZ {
				x += stepX
				if x < 0 || x >= BlockSize {
					return false
				}
				nextX += deltaX
			} else {
				z += stepZ
				if z < 0 || z >= BlockSize {
					return false
				}
				nextZ += deltaZ
			}
		} else {
			if nextY < nextZ {
				y += stepY
				if y < 0 || y >= BlockSize {
					return false
				}
				nextY += deltaY
			} else {
				z += stepZ
				if z < 0 || z >= BlockSize {
					return false
				}
				nextZ += deltaZ
			}
		}
	}
}
