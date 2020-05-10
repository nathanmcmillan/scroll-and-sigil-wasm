package main

import (
	"math"

	"../fast"
	"./graphics"
	"./render"
)

type missile struct {
	world     *world
	uid       uint16
	sid       string
	nid       uint16
	sprite    *render.Sprite
	x         float32
	y         float32
	z         float32
	deltaX    float32
	deltaY    float32
	deltaZ    float32
	minBX     int
	minBY     int
	minBZ     int
	maxBX     int
	maxBY     int
	maxBZ     int
	radius    float32
	height    float32
	cleanupFn func()
}

func (me *missile) blockBorders() {
	me.minBX = int((me.x - me.radius) * InverseBlockSize)
	me.minBY = int(me.y * InverseBlockSize)
	me.minBZ = int((me.z - me.radius) * InverseBlockSize)
	me.maxBX = int((me.x + me.radius) * InverseBlockSize)
	me.maxBY = int((me.y + me.height) * InverseBlockSize)
	me.maxBZ = int((me.z + me.radius) * InverseBlockSize)
}

func (me *missile) addToBlocks() bool {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				if block == nil {
					me.removeFromBlocks()
					return true
				}
				block.addMissile(me)
			}
		}
	}
	return false
}

func (me *missile) removeFromBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				if block != nil {
					block.removeMissile(me)
				}
			}
		}
	}
}

func (me *missile) netUpdate(data *fast.ByteReader, delta uint8) {
}

func (me *missile) cleanup() {
	me.world.removeMissile(me)
	me.removeFromBlocks()
	me.cleanupFn()
}

func (me *missile) update() bool {
	me.removeFromBlocks()
	me.x += me.deltaX
	me.y += me.deltaY
	me.z += me.deltaZ
	me.blockBorders()
	return me.addToBlocks()
}

func (me *missile) render(spriteBuffer map[string]*graphics.RenderBuffer, camX, camZ, camAngle float32) {
	sin := float64(camX - me.x)
	cos := float64(camZ - me.z)
	length := math.Sqrt(sin*sin + cos*cos)
	sin /= length
	cos /= length
	render.RendSprite(spriteBuffer[me.sid], me.x, me.y, me.z, float32(sin), float32(cos), me.sprite)
}

func plasmaInit(world *world, nid uint16, damage uint16, x, y, z, dx, dy, dz float32) {
	me := &missile{}
	me.world = world
	me.x = x
	me.y = y
	me.z = z
	me.blockBorders()
	if me.addToBlocks() {
		return
	}
	me.uid = PlasmaUID
	me.sid = "missiles"
	me.nid = nid
	me.sprite = wadSpriteData[me.sid]["baron-missile-front-1"]
	me.deltaX = dx * InverseNetRate
	me.deltaY = dy * InverseNetRate
	me.deltaZ = dz * InverseNetRate
	me.radius = 0.2
	me.height = 0.2
	me.cleanupFn = me.plasmaCleanup
	world.addMissile(me)
}

func (me *missile) plasmaCleanup() {
	playWadSound("plasma-impact")
	plasmaExplosionInit(me.world, me.x, me.y, me.z)
}
