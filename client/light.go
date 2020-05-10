package main

import (
	"./render"
)

const (
	lightQueueLimit = 30 * 30 * 30
	lightFade       = 0.95
)

var (
	lightQueue    [lightQueueLimit][3]int
	lightQueuePos int
	lightQueueNum int
)

type light struct {
	x     int
	y     int
	z     int
	red   uint8
	green uint8
	blue  uint8
}

func lightInit(x, y, z uint8, rgb int32) *light {
	light := &light{}
	light.x = int(x)
	light.y = int(y)
	light.z = int(z)
	red, green, blue := render.UnpackRgb(rgb)
	light.red = uint8(red)
	light.green = uint8(green)
	light.blue = uint8(blue)
	return light
}

func (me *light) save() string {
	return ""
}

func lightPlane(rgb *[4][3]float32, ambient *[4]float32) {
	rgb[0][0] *= ambient[0] / 65025.0
	rgb[0][1] *= ambient[0] / 65025.0
	rgb[0][2] *= ambient[0] / 65025.0

	rgb[1][0] *= ambient[1] / 65025.0
	rgb[1][1] *= ambient[1] / 65025.0
	rgb[1][2] *= ambient[1] / 65025.0

	rgb[2][0] *= ambient[2] / 65025.0
	rgb[2][1] *= ambient[2] / 65025.0
	rgb[2][2] *= ambient[2] / 65025.0

	rgb[3][0] *= ambient[3] / 65025.0
	rgb[3][1] *= ambient[3] / 65025.0
	rgb[3][2] *= ambient[3] / 65025.0
}

func lightVisit(world *world, bx, by, bz, tx, ty, tz int, red, green, blue uint8) {
	tile := world.getTilePointer(bx, by, bz, tx, ty, tz)
	if tile == nil || TileClosed[tile.typeOf] {
		return
	}
	if tile.red >= red || tile.green >= green || tile.blue >= blue {
		return
	}
	tile.red = red
	tile.green = green
	tile.blue = blue

	queue := lightQueuePos + lightQueueNum
	if queue >= lightQueueLimit {
		queue -= lightQueueLimit
	}
	lightQueue[queue][0] = tx
	lightQueue[queue][1] = ty
	lightQueue[queue][2] = tz
	lightQueueNum++
}

func (me *light) addToWorld(world *world, block *block) {
	origin := &block.tiles[me.x+(me.y<<BlockShift)+(me.z<<BlockShiftSlice)]
	origin.red = me.red
	origin.green = me.green
	origin.blue = me.blue

	bx := block.x
	by := block.y
	bz := block.z

	lightQueue[0][0] = me.x
	lightQueue[0][1] = me.y
	lightQueue[0][2] = me.z
	lightQueuePos = 0
	lightQueueNum = 1

	for lightQueueNum > 0 {
		x := lightQueue[lightQueuePos][0]
		y := lightQueue[lightQueuePos][1]
		z := lightQueue[lightQueuePos][2]

		lightQueuePos++
		if lightQueuePos == lightQueueLimit {
			lightQueuePos = 0
		}
		lightQueueNum--

		tile := world.getTilePointer(bx, by, bz, x, y, z)
		if tile == nil {
			continue
		}

		r := uint8(float32(tile.red) * lightFade)
		g := uint8(float32(tile.green) * lightFade)
		b := uint8(float32(tile.blue) * lightFade)

		lightVisit(world, bx, by, bz, x-1, y, z, r, g, b)
		lightVisit(world, bx, by, bz, x+1, y, z, r, g, b)
		lightVisit(world, bx, by, bz, x, y-1, z, r, g, b)
		lightVisit(world, bx, by, bz, x, y+1, z, r, g, b)
		lightVisit(world, bx, by, bz, x, y, z-1, r, g, b)
		lightVisit(world, bx, by, bz, x, y, z+1, r, g, b)
	}
}
