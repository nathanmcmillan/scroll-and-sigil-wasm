package main

import (
	"math"

	"../fast"
	"./graphics"
	"./render"
)

// Constants
const (
	Gravity = 0.01

	AnimationRate = 16

	AnimationNotDone    = 0
	AnimationAlmostDone = 1
	AnimationDone       = 2

	AnimationFront     = 0
	AnimationFrontSide = 1
	AnimationSide      = 2
	AnimationBackSide  = 3
	AnimationBack      = 4

	DirectionNorth     = 0
	DirectionNorthEast = 1
	DirectionEast      = 2
	DirectionSouthEast = 3
	DirectionSouth     = 4
	DirectionSouthWest = 5
	DirectionWest      = 6
	DirectionNorthWest = 7
	DirectionCount     = 8
	DirectionNone      = 8

	HumanUID  = uint16(0)
	BaronUID  = uint16(1)
	TreeUID   = uint16(2)
	PlasmaUID = uint16(3)
	MedkitUID = uint16(4)

	ThingAngleA = 337.5 * DegToRad
	ThingAngleB = 292.5 * DegToRad
	ThingAngleC = 247.5 * DegToRad
	ThingAngleD = 202.5 * DegToRad
	ThingAngleE = 157.5 * DegToRad
	ThingAngleF = 112.5 * DegToRad
	ThingAngleG = 67.5 * DegToRad
	ThingAngleH = 22.5 * DegToRad
)

// Variables
var (
	DirectionToAngle = []float32{
		0.0 * DegToRad,
		45.0 * DegToRad,
		90.0 * DegToRad,
		135.0 * DegToRad,
		180.0 * DegToRad,
		225.0 * DegToRad,
		270.0 * DegToRad,
		315.0 * DegToRad,
	}
)

type netThing interface {
	cleanup()
	netUpdate(*fast.ByteReader, uint8)
}

type thing struct {
	world          *world
	uid            uint16
	sid            string
	nid            uint16
	animation      [][]*render.Sprite
	animationMod   int
	animationFrame int
	x              float32
	y              float32
	z              float32
	angle          float32
	deltaX         float32
	deltaY         float32
	deltaZ         float32
	oldX           float32
	oldY           float32
	oldZ           float32
	netX           float32
	netY           float32
	netZ           float32
	deltaNetX      float32
	deltaNetY      float32
	deltaNetZ      float32
	minBX          int
	minBY          int
	minBZ          int
	maxBX          int
	maxBY          int
	maxBZ          int
	ground         bool
	radius         float32
	height         float32
	speed          float32
	health         uint16
	update         func()
	damage         func(uint16)
}

func (me *thing) blockBorders() {
	me.minBX = int((me.x - me.radius) * InverseBlockSize)
	me.minBY = int(me.y * InverseBlockSize)
	me.minBZ = int((me.z - me.radius) * InverseBlockSize)
	me.maxBX = int((me.x + me.radius) * InverseBlockSize)
	me.maxBY = int((me.y + me.height) * InverseBlockSize)
	me.maxBZ = int((me.z + me.radius) * InverseBlockSize)
}

func (me *thing) addToBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.addThing(me)
			}
		}
	}
}

func (me *thing) removeFromBlocks() {
	for gx := me.minBX; gx <= me.maxBX; gx++ {
		for gy := me.minBY; gy <= me.maxBY; gy++ {
			for gz := me.minBZ; gz <= me.maxBZ; gz++ {
				block := me.world.getBlock(gx, gy, gz)
				block.removeThing(me)
			}
		}
	}
}

func (me *thing) thingNetUpdate(data *fast.ByteReader, delta uint8) {
	if delta&0x1 != 0 {
		x := data.GetFloat32()
		z := data.GetFloat32()
		me.netX = x
		me.deltaNetX = (x - me.x) * InverseNetRate
		me.netZ = z
		me.deltaNetZ = (z - me.z) * InverseNetRate
	}
	if delta&0x2 != 0 {
		y := data.GetFloat32()
		me.netY = y
		me.deltaNetY = (y - me.y) * InverseNetRate
	}
}

func (me *thing) cleanup() {
	me.world.removeThing(me)
	me.removeFromBlocks()
}

func (me *thing) updateAnimation() int {
	me.animationMod++
	if me.animationMod == AnimationRate {
		me.animationMod = 0
		me.animationFrame++
		size := len(me.animation)
		if me.animationFrame == size-1 {
			return AnimationAlmostDone
		} else if me.animationFrame == size {
			return AnimationDone
		}
	}
	return AnimationNotDone
}

func (me *thing) updateNetworkDelta() {
	updateBlocks := false

	if me.deltaNetX > 0 {
		me.x += me.deltaNetX
		updateBlocks = true
		if me.x >= me.netX {
			me.x = me.netX
			me.deltaNetX = 0
		}
	} else if me.deltaNetX < 0 {
		me.x += me.deltaNetX
		updateBlocks = true
		if me.x <= me.netX {
			me.x = me.netX
			me.deltaNetX = 0
		}
	}

	if me.deltaNetY > 0 {
		me.y += me.deltaNetY
		updateBlocks = true
		if me.y >= me.netY {
			me.y = me.netY
			me.deltaNetY = 0
		}
	} else if me.deltaNetY < 0 {
		me.y += me.deltaNetY
		updateBlocks = true
		if me.y <= me.netY {
			me.y = me.netY
			me.deltaNetY = 0
		}
	}

	if me.deltaNetZ > 0 {
		me.z += me.deltaNetZ
		updateBlocks = true
		if me.z >= me.netZ {
			me.z = me.netZ
			me.deltaNetZ = 0
		}
	} else if me.deltaNetZ < 0 {
		me.z += me.deltaNetZ
		updateBlocks = true
		if me.z <= me.netZ {
			me.z = me.netZ
			me.deltaNetZ = 0
		}
	}

	if updateBlocks {
		me.removeFromBlocks()
		me.blockBorders()
		me.addToBlocks()
	}
}

func (me *thing) render(spriteBuffer map[string]*graphics.RenderBuffer, camX, camZ, camAngle float32) {
	sin := float64(camX - me.x)
	cos := float64(camZ - me.z)
	length := math.Sqrt(sin*sin + cos*cos)
	sin /= length
	cos /= length

	angle := camAngle - me.angle
	if angle < 0 {
		angle += Tau
	}

	var direction int
	var mirror bool

	if angle > ThingAngleA {
		direction = AnimationBack
		mirror = false
	} else if angle > ThingAngleB {
		direction = AnimationBackSide
		mirror = true
	} else if angle > ThingAngleC {
		direction = AnimationSide
		mirror = true
	} else if angle > ThingAngleD {
		direction = AnimationFrontSide
		mirror = true
	} else if angle > ThingAngleE {
		direction = AnimationFront
		mirror = false
	} else if angle > ThingAngleF {
		direction = AnimationFrontSide
		mirror = false
	} else if angle > ThingAngleG {
		direction = AnimationSide
		mirror = false
	} else if angle > ThingAngleH {
		direction = AnimationBackSide
		mirror = false
	} else {
		direction = AnimationBack
		mirror = false
	}

	sprite := me.animation[me.animationFrame][direction]

	if mirror {
		render.RendMirrorSprite(spriteBuffer[me.sid], me.x, me.y, me.z, float32(sin), float32(cos), sprite)
	} else {
		render.RendSprite(spriteBuffer[me.sid], me.x, me.y, me.z, float32(sin), float32(cos), sprite)
	}
}
