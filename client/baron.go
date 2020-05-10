package main

import (
	"math"
	"math/rand"
	"strconv"

	"../fast"

	"./render"
)

const (
	baronSleep   = uint8(0)
	baronDead    = uint8(1)
	baronLook    = uint8(2)
	baronChase   = uint8(3)
	baronMelee   = uint8(4)
	baronMissile = uint8(5)
)

var (
	baronAnimationIdle    [][]*render.Sprite
	baronAnimationWalk    [][]*render.Sprite
	baronAnimationMelee   [][]*render.Sprite
	baronAnimationMissile [][]*render.Sprite
	baronAnimationDeath   [][]*render.Sprite
)

type baron struct {
	*thing
	status uint8
}

func baronInit(world *world, nid uint16, x, y, z float32, direction uint8, health uint16, status uint8) *baron {
	baron := &baron{}
	baron.thing = &thing{}
	baron.thing.update = baron.updateFn
	baron.world = world
	baron.uid = BaronUID
	baron.sid = "baron"
	baron.nid = nid
	baron.animateInit(status)
	baron.x = x
	baron.y = y
	baron.z = z
	if direction != DirectionNone {
		baron.angle = DirectionToAngle[direction]
	}
	baron.oldX = x
	baron.oldY = y
	baron.oldZ = z
	baron.radius = 0.4
	baron.height = 1.0
	baron.speed = 0.1
	baron.health = health
	baron.status = status
	world.addThing(baron.thing)
	world.netLookup[baron.nid] = baron
	baron.blockBorders()
	baron.addToBlocks()
	return baron
}

func (me *baron) animateInit(status uint8) {
	switch status {
	case baronDead:
		me.animation = baronAnimationDeath
	case baronMelee:
		me.animation = baronAnimationMelee
	case baronMissile:
		me.animation = baronAnimationMissile
	default:
		me.animation = baronAnimationWalk
	}
}

func (me *baron) netUpdate(data *fast.ByteReader, delta uint8) {
	me.thingNetUpdate(data, delta)
	if delta&0x4 != 0 {
		me.netUpdateHealth(data.GetUint16())
	}
	if delta&0x8 != 0 {
		me.netUpdateState(data.GetUint8())
	}
	if delta&0x10 != 0 {
		direction := data.GetUint8()
		if direction != DirectionNone {
			me.angle = DirectionToAngle[direction]
		}
	}
}

func (me *baron) netUpdateState(status uint8) {
	if me.status == status {
		return
	}
	me.animationMod = 0
	me.animationFrame = 0
	switch status {
	case baronDead:
		me.animation = baronAnimationDeath
	case baronMelee:
		me.animation = baronAnimationMelee
		playWadSound("baron-melee")
	case baronMissile:
		me.animation = baronAnimationMissile
		playWadSound("baron-missile")
	case baronChase:
		if rand.Float32() < 0.1 {
			playWadSound("baron-scream")
		}
	default:
		me.animation = baronAnimationWalk
	}
	me.status = status
}

func (me *baron) netUpdateHealth(health uint16) {
	if health < me.health {
		if health < 1 {
			playWadSound("baron-death")
		} else {
			playWadSound("baron-pain")
		}
		for i := 0; i < 20; i++ {
			spriteName := "blood-" + strconv.Itoa(int(math.Floor(rand.Float64()*3)))
			x := me.x + me.radius*(1-rand.Float32()*2)
			y := me.y + me.height*rand.Float32()
			z := me.z + me.radius*(1-rand.Float32()*2)
			const spread = 0.2
			dx := spread * (1 - rand.Float32()*2)
			dy := spread * rand.Float32()
			dz := spread * (1 - rand.Float32()*2)
			bloodInit(me.world, x, y, z, dx, dy, dz, spriteName)
		}

	}
	me.health = health
}

func (me *baron) dead() {
	if me.animationFrame == len(me.animation)-1 {
		me.update = me.emptyUpdate
	} else {
		me.updateAnimation()
	}
}

func (me *baron) look() {
	if me.updateAnimation() == AnimationDone {
		me.animationFrame = 0
	}
}

func (me *baron) melee() {
	if me.updateAnimation() == AnimationDone {
		me.animationFrame = 0
		me.animation = baronAnimationWalk
	}
}

func (me *baron) missile() {
	if me.updateAnimation() == AnimationDone {
		me.animationFrame = 0
		me.animation = baronAnimationWalk
	}
}

func (me *baron) chase() {
	if me.updateAnimation() == AnimationDone {
		me.animationFrame = 0
	}
}

func (me *baron) updateFn() {
	switch me.status {
	case baronDead:
		me.dead()
	case baronLook:
		me.look()
	case baronMelee:
		me.melee()
	case baronMissile:
		me.missile()
	case baronChase:
		me.chase()
	}
	me.updateNetworkDelta()
}

func (me *baron) emptyUpdate() {
}
