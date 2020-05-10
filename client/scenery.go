package main

import (
	"math"

	"../fast"
	"./graphics"
	"./render"
)

type scenery struct {
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

func treeInit(world *world, nid uint16, x, y, z float32) *scenery {
	tree := &scenery{}
	tree.world = world
	tree.uid = TreeUID
	tree.sid = "scenery"
	tree.nid = nid
	tree.sprite = wadSpriteData[tree.sid]["dead-tree"]
	tree.x = x
	tree.y = y
	tree.z = z
	tree.radius = 0.4
	tree.height = 1.0
	world.addScenery(tree)
	world.netLookup[tree.nid] = tree
	tree.blockBorders()
	tree.addToBlocks()
	return tree
}

func (me *scenery) blockBorders() {
	me.minBX = int((me.x - me.radius) * InverseBlockSize)
	me.minBY = int(me.y * InverseBlockSize)
	me.minBZ = int((me.z - me.radius) * InverseBlockSize)
	me.maxBX = int((me.x + me.radius) * InverseBlockSize)
	me.maxBY = int((me.y + me.height) * InverseBlockSize)
	me.maxBZ = int((me.z + me.radius) * InverseBlockSize)
}

func (me *scenery) addToBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.addScenery(me)
			}
		}
	}
}

func (me *scenery) removeFromBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.removeScenery(me)
			}
		}
	}
}

func (me *scenery) netUpdate(data *fast.ByteReader, delta uint8) {
}

func (me *scenery) cleanup() {
	me.world.removeScenery(me)
	me.removeFromBlocks()
}

func (me *scenery) render(spriteBuffer map[string]*graphics.RenderBuffer, camX, camZ float32) {
	sin := float64(camX - me.x)
	cos := float64(camZ - me.z)
	length := math.Sqrt(sin*sin + cos*cos)
	sin /= length
	cos /= length
	render.RendSprite(spriteBuffer[me.sid], me.x, me.y, me.z, float32(sin), float32(cos), me.sprite)
}
