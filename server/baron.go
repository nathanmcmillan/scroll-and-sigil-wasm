package main

import (
	"math"

	"../fast"
)

// Animation constants
const (
	BaronWalkAnimation    int = 2 * AnimationRate
	BaronMeleeAnimation   int = 2 * AnimationRate
	BaronMissileAnimation int = 3 * AnimationRate
	BaronDeathAnimation   int = 2 * AnimationRate
)

// Baron constants
const (
	BaronSleep   = uint8(0)
	BaronDead    = uint8(1)
	BaronLook    = uint8(2)
	BaronChase   = uint8(3)
	BaronMelee   = uint8(4)
	BaronMissile = uint8(5)
)

// Baron struct
type Baron struct {
	*Npc
	Status       uint8
	Reaction     uint8
	MeleeRange   float32
	MissileRange float32
	DeltaHealth  bool
	DeltaStatus  bool
}

// NewBaron func
func NewBaron(world *World, x, y, z float32) *Baron {
	baron := &Baron{}
	baron.Npc = &Npc{}
	baron.thing = &thing{}
	baron.UID = BaronUID
	baron.NID = NextNID()
	baron.World = world
	baron.thing.Update = baron.Update
	baron.thing.Damage = baron.Damage
	baron.thing.Save = baron.Save
	baron.thing.Snap = baron.Snap
	baron.X = x
	baron.Y = y
	baron.Z = z
	baron.Radius = 0.4
	baron.Height = 1.0
	baron.Animation = BaronWalkAnimation
	baron.Health = 1
	baron.Group = DemonGroup
	baron.Speed = 0.1
	baron.MoveDirection = DirectionNone
	baron.Status = BaronLook
	baron.MeleeRange = 2.0
	baron.MissileRange = 10.0
	world.addThing(baron.thing)
	baron.blockBorders()
	baron.addToBlocks()
	return baron
}

// Save func
func (me *Baron) Save(data *fast.ByteWriter) {
	data.PutUint16(me.UID)
	data.PutUint16(me.NID)
	data.PutFloat32(me.X)
	data.PutFloat32(me.Y)
	data.PutFloat32(me.Z)
	data.PutUint8(me.MoveDirection)
	data.PutUint16(me.Health)
	data.PutUint8(me.Status)
}

// Snap func
func (me *Baron) Snap(data *fast.ByteWriter) {
	delta := uint8(0)
	if me.DeltaMoveXZ {
		delta |= 0x1
	}
	if me.DeltaMoveY {
		delta |= 0x2
	}
	if me.DeltaHealth {
		delta |= 0x4
	}
	if me.DeltaStatus {
		delta |= 0x8
	}
	if me.DeltaMoveDirection {
		delta |= 0x10
	}
	if delta == 0 {
		me.Binary = nil
		return
	}
	data.Reset()
	data.PutUint16(me.NID)
	data.PutUint8(delta)
	if me.DeltaMoveXZ {
		data.PutFloat32(me.X)
		data.PutFloat32(me.Z)
		me.DeltaMoveXZ = false
	}
	if me.DeltaMoveY {
		data.PutFloat32(me.Y)
		me.DeltaMoveY = false
	}
	if me.DeltaHealth {
		data.PutUint16(me.Health)
		me.DeltaHealth = false
	}
	if me.DeltaStatus {
		data.PutUint8(me.Status)
		me.DeltaStatus = false
	}
	if me.DeltaMoveDirection {
		data.PutUint8(me.MoveDirection)
		me.DeltaMoveDirection = false
	}
	binary := data.Bytes()
	me.Binary = make([]byte, len(binary))
	copy(me.Binary, binary)
}

// Damage func
func (me *Baron) Damage(amount uint16) {
	if me.Status == BaronDead {
		return
	}
	me.Health -= amount
	me.DeltaHealth = true
	if me.Health < 1 {
		me.Health = 0
		me.Status = BaronDead
		me.DeltaStatus = true
		me.AnimationFrame = 0
		me.Animation = BaronDeathAnimation
		me.removeFromBlocks()
	}
}

// Dead func
func (me *Baron) Dead() {
	if me.AnimationFrame == me.Animation-1 {
		me.thing.Update = me.NopUpdate
		me.thing.Snap = me.NopSnap
	} else {
		me.UpdateAnimation()
	}
}

// Look func
func (me *Baron) Look() {
	for i := 0; i < me.World.thingCount; i++ {
		thing := me.World.things[i]
		if thing.Group == HumanGroup && thing.Health > 0 {
			me.Target = thing
			me.Status = BaronChase
			me.DeltaStatus = true
			return
		}
	}
	if me.UpdateAnimation() == AnimationDone {
		me.AnimationFrame = 0
	}
}

// Melee func
func (me *Baron) Melee() {
	anim := me.UpdateAnimation()
	if anim == AnimationAlmostDone {
		me.Reaction = 40 + NextRandP()%220
		if me.ApproximateDistance(me.Target) <= me.MeleeRange {
			me.Target.Damage(uint16(1 + NextRandP()%3))
		}
	} else if anim == AnimationDone {
		me.Status = BaronChase
		me.DeltaStatus = true
		me.AnimationFrame = 0
		me.Animation = BaronWalkAnimation
	}
}

// missile func
func (me *Baron) missile() {
	anim := me.UpdateAnimation()
	if anim == AnimationAlmostDone {
		me.Reaction = 40 + NextRandP()%220
		const speed = 0.3
		angle := math.Atan2(float64(me.Target.Z-me.Z), float64(me.Target.X-me.X))
		dx := float32(math.Cos(angle))
		dz := float32(math.Sin(angle))
		dist := me.ApproximateDistance(me.Target)
		dy := (me.Target.Y + me.Target.Height*0.5 - me.Y - me.Height*0.5) / (dist / speed)
		x := me.X + dx*me.Radius*3.0
		y := me.Y + me.Height*0.75
		z := me.Z + dz*me.Radius*3.0
		NewPlasma(me.World, uint16(1+NextRandP()%3), x, y, z, dx*speed, dy, dz*speed)
	} else if anim == AnimationDone {
		me.Status = BaronChase
		me.DeltaStatus = true
		me.AnimationFrame = 0
		me.Animation = BaronWalkAnimation
	}
}

// Chase func
func (me *Baron) Chase() {
	if me.Reaction > 0 {
		me.Reaction--
	}
	if me.Target == nil || me.Target.Health <= 0 {
		me.Target = nil
		me.DeltaStatus = true
		me.Status = BaronLook
	} else {
		dist := me.ApproximateDistance(me.Target)
		if me.Reaction == 0 && dist < me.MeleeRange {
			me.Status = BaronMelee
			me.DeltaStatus = true
			me.AnimationFrame = 0
			me.Animation = BaronMeleeAnimation
		} else if me.Reaction == 0 && dist <= me.MissileRange {
			me.Status = BaronMissile
			me.DeltaStatus = true
			me.AnimationFrame = 0
			me.Animation = BaronMissileAnimation
		} else {
			me.MoveCount--
			move := me.Move()
			if me.MoveCount < 0 || !move {
				me.NewChaseDirection()
			}
			if me.UpdateAnimation() == AnimationDone {
				me.AnimationFrame = 0
			}
		}
	}
}

// Update func
func (me *Baron) Update() bool {
	switch me.Status {
	case BaronDead:
		me.Dead()
	case BaronLook:
		me.Look()
	case BaronMelee:
		me.Melee()
	case BaronMissile:
		me.missile()
	case BaronChase:
		me.Chase()
	}
	me.IntegrateY()
	return false
}
