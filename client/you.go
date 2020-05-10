package main

import (
	"math"
	"math/rand"
	"strconv"
	"syscall/js"

	"../fast"
)

const (
	inputOpNewMove      = uint8(0)
	inputOpContinueMove = uint8(1)
	inputOpMissile      = uint8(2)
	inputOpSearch       = uint8(3)
	inputOpChat         = uint8(4)
)

type you struct {
	*thing
	status     uint8
	camera     *camera
	socket     js.Value
	socketSend map[uint8]interface{}
}

func youInit(world *world, nid uint16, x, y, z, angle float32, health uint16, status uint8) *you {
	you := &you{}
	you.thing = &thing{}
	you.thing.update = you.updateFn
	you.world = world
	you.uid = HumanUID
	you.sid = "baron"
	you.nid = nid
	you.animation = humanAnimationWalk
	you.x = x
	you.y = y
	you.z = z
	you.angle = angle
	you.oldX = x
	you.oldY = y
	you.oldZ = z
	you.radius = 0.4
	you.height = 1.0
	you.speed = 0.1
	you.health = health
	you.status = status
	world.addThing(you.thing)
	world.netLookup[you.nid] = you
	you.blockBorders()
	you.addToBlocks()
	return you
}

func (me *you) netUpdate(data *fast.ByteReader, delta uint8) {
	me.thingNetUpdate(data, delta)
	if delta&0x4 != 0 {
		me.netUpdateHealth(data.GetUint16())
	}
	if delta&0x8 != 0 {
		me.netUpdateState(data.GetUint8())
	}
	if delta&0x10 != 0 {
		me.angle = data.GetFloat32()
	}
}

func (me *you) netUpdateState(status uint8) {
	if me.status == status {
		return
	}
	me.animationMod = 0
	me.animationFrame = 0
	switch status {
	case humanDead:
		me.animation = humanAnimationDeath
	case humanMissile:
		me.animation = humanAnimationMissile
	case humanIdle:
		me.animation = humanAnimationIdle
	default:
		me.animation = humanAnimationWalk
	}
	me.status = status
}

func (me *you) netUpdateHealth(health uint16) {
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

func (me *you) dead() {
	if me.animationFrame == len(me.animation)-1 {
		me.update = me.emptyUpdate
	} else {
		me.updateAnimation()
	}
}

func (me *you) missile() {
	if me.updateAnimation() == AnimationDone {
		me.animationFrame = 0
		me.animation = humanAnimationWalk
	}
}

func (me *you) walk() {
	if InputIsKeyPress("p") {
		me.socketSend[inputOpSearch] = true
		me.socketSend[inputOpChat] = "todo test chat"
	}

	if InputIsKeyDown(" ") {
		me.socketSend[inputOpMissile] = true
		me.status = humanMissile
		me.animationMod = 0
		me.animationFrame = 0
		me.animation = humanAnimationMissile
		playWadSound("baron-missile")
		return
	}

	const (
		none          = 0
		forward       = 1
		backward      = 2
		left          = 3
		right         = 4
		forwardLeft   = 5
		forwardRight  = 6
		backwardLeft  = 7
		backwardRight = 8
	)
	direction := none
	var goal float32

	if InputIsKeyDown("w") {
		direction = forward
		goal = me.camera.ry
	}

	if InputIsKeyDown("s") {
		if direction == none {
			direction = backward
			goal = me.camera.ry + Pi
		} else {
			direction = none
		}
	}

	if InputIsKeyDown("a") {
		if direction == none {
			direction = left
			goal = me.camera.ry - HalfPi
		} else if direction == forward {
			direction = forwardLeft
			goal -= QuarterPi
		} else if direction == backward {
			direction = backwardLeft
			goal += QuarterPi
		}
	}

	if InputIsKeyDown("d") {
		if direction == none {
			direction = right
			goal = me.camera.ry + HalfPi
		} else if direction == left {
			direction = none
		} else if direction == forwardLeft {
			goal = me.camera.ry
		} else if direction == backwardLeft {
			goal = me.camera.ry + Pi
		} else if direction == forward {
			goal += QuarterPi
		} else if direction == backward {
			goal -= QuarterPi
		}
	}

	if direction == none {
		me.animationMod = 0
		me.animationFrame = 0
		me.animation = humanAnimationIdle
	} else {
		if goal < 0 {
			goal += Tau
		} else if goal >= Tau {
			goal -= Tau
		}

		if me.angle != goal {
			me.angle = goal
			me.socketSend[inputOpNewMove] = goal
		} else {
			me.socketSend[inputOpContinueMove] = true
		}

		if &me.animation[0] != &humanAnimationWalk[0] {
			me.animation = humanAnimationWalk
		} else if me.updateAnimation() == AnimationDone {
			me.animationFrame = 0
		}
	}
}

func (me *you) updateFn() {
	switch me.status {
	case humanDead:
		me.dead()
	case humanMissile:
		me.missile()
	default:
		me.walk()
	}
	me.updateNetworkDelta()
}

func (me *you) emptyUpdate() {
}
