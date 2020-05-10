package main

import (
	"math"

	"../fast"
	"./graphics"
	"./render"
)

type item struct {
	world  *world
	uid    uint16
	sid    string
	nid    uint16
	sprite *render.Sprite
	x      float32
	y      float32
	z      float32
	minBX  int
	minBY  int
	minBZ  int
	maxBX  int
	maxBY  int
	maxBZ  int
	radius float32
	height float32
}

func (me *item) blockBorders() {
	me.minBX = int((me.x - me.radius) * InverseBlockSize)
	me.minBY = int(me.y * InverseBlockSize)
	me.minBZ = int((me.z - me.radius) * InverseBlockSize)
	me.maxBX = int((me.x + me.radius) * InverseBlockSize)
	me.maxBY = int((me.y + me.height) * InverseBlockSize)
	me.maxBZ = int((me.z + me.radius) * InverseBlockSize)
}

func (me *item) addToBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.addItem(me)
			}
		}
	}
}

func (me *item) removeFromBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.removeItem(me)
			}
		}
	}
}

func (me *item) netUpdate(data *fast.ByteReader, delta uint8) {
}

func (me *item) cleanup() {
	me.world.removeItem(me)
	me.removeFromBlocks()
}

func (me *item) render(spriteBuffer map[string]*graphics.RenderBuffer, camX, camZ float32) {
	sin := float64(camX - me.x)
	cos := float64(camZ - me.z)
	length := math.Sqrt(sin*sin + cos*cos)
	sin /= length
	cos /= length
	render.RendSprite(spriteBuffer[me.sid], me.x, me.y, me.z, float32(sin), float32(cos), me.sprite)
}

func medkitInit(world *world, nid uint16, x, y, z float32) *item {
	meds := &item{}
	meds.world = world
	meds.x = x
	meds.y = y
	meds.z = z
	meds.blockBorders()
	meds.addToBlocks()
	meds.uid = MedkitUID
	meds.sid = "item"
	meds.nid = nid
	meds.sprite = wadSpriteData[meds.sid]["medkit"]
	meds.radius = 0.2
	meds.height = 0.2
	world.addItem(meds)
	return meds
}
