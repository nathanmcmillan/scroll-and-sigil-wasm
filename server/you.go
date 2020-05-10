package main

import (
	"math"

	"../fast"
)

// Animation constants
const (
	HumanWalkAnimation    int = 2 * AnimationRate
	HumanMeleeAnimation   int = 2 * AnimationRate
	HumanMissileAnimation int = 3 * AnimationRate
	HumanDeathAnimation   int = 2 * AnimationRate
)

// Human constants
const (
	HumanDead    = uint8(0)
	HumanIdle    = uint8(1)
	HumanWalk    = uint8(2)
	HumanMelee   = uint8(3)
	HumanMissile = uint8(4)
)

// Input constants
const (
	InputOpNewMove      = uint8(0)
	InputOpContinueMove = uint8(1)
	InputOpMissile      = uint8(2)
	InputOpSearch       = uint8(3)
	InputOpChat         = uint8(4)
)

// You struct
type You struct {
	*thing
	Status      uint8
	DeltaAngle  bool
	DeltaHealth bool
	DeltaStatus bool
	Person      *Person
}

// NewYou func
func NewYou(world *World, person *Person, x, y, z float32) *You {
	you := &You{}
	you.thing = &thing{}
	you.UID = HumanUID
	you.NID = NextNID()
	you.World = world
	you.thing.Update = you.Update
	you.thing.Damage = you.Damage
	you.thing.Save = you.Save
	you.thing.Snap = you.Snap
	you.X = x
	you.Y = y
	you.Z = z
	you.Radius = 0.4
	you.Height = 1.0
	you.Speed = 0.1
	you.Health = 8
	you.Group = HumanGroup
	you.Status = HumanIdle
	you.Person = person
	world.addThing(you.thing)
	you.blockBorders()
	you.addToBlocks()
	return you
}

// Save func
func (me *You) Save(data *fast.ByteWriter) {
	data.PutUint16(me.UID)
	data.PutUint16(me.NID)
	data.PutFloat32(me.X)
	data.PutFloat32(me.Y)
	data.PutFloat32(me.Z)
	data.PutFloat32(me.Angle)
	data.PutUint16(me.Health)
	data.PutUint8(me.Status)
}

// Snap func
func (me *You) Snap(data *fast.ByteWriter) {
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
	if me.DeltaAngle {
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
	if me.DeltaAngle {
		data.PutFloat32(me.Angle)
		me.DeltaAngle = false
	}
	binary := data.Bytes()
	me.Binary = make([]byte, len(binary))
	copy(me.Binary, binary)
}

// Search func
func (me *You) Search() {
	gy := me.MinBY
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
			block := me.World.getBlock(gx, gy, gz)
			for i := 0; i < block.itemCount; i++ {
				item := block.items[i]
				if item.Overlap(me.thing) {
					item.Cleanup()
					return
				}
			}
		}
	}

}

// Damage func
func (me *You) Damage(amount uint16) {
	if me.Status == HumanDead {
		return
	}
	me.Health -= amount
	me.DeltaHealth = true
	if me.Health < 1 {
		me.Health = 0
		me.Status = HumanDead
		me.DeltaStatus = true
		me.AnimationFrame = 0
		me.Animation = HumanDeathAnimation
		me.removeFromBlocks()
	}
}

// Dead func
func (me *You) Dead() {
	if me.AnimationFrame == me.Animation-1 {
		me.thing.Update = me.NopUpdate
		me.thing.Snap = me.NopSnap
	} else {
		me.UpdateAnimation()
	}
}

// missile func
func (me *You) missile() {
	anim := me.UpdateAnimation()
	if anim == AnimationAlmostDone {
		const speed = 0.5
		angle := float64(me.Angle)
		dx := float32(math.Sin(angle))
		dz := -float32(math.Cos(angle))
		x := me.X + dx*me.Radius*3.0
		y := me.Y + me.Height*0.75
		z := me.Z + dz*me.Radius*3.0
		NewPlasma(me.World, uint16(1+NextRandP()%3), x, y, z, dx*speed, 0.0, dz*speed)
	} else if anim == AnimationDone {
		me.Status = HumanIdle
		me.DeltaStatus = true
	}
}

// Walk func
func (me *You) Walk() {
	person := me.Person
	if person.InputCount == 0 {
		return
	}
	move := false
	attack := false
gotoRead:
	for i := 0; i < person.InputCount; i++ {
		input := person.InputQueue[i]
		data := fast.ByteReaderInit(input)
		if data.NotSafe(1) {
			break gotoRead
		}
		opCount := data.GetUint8()
		for c := uint8(0); c < opCount; c++ {
			if data.NotSafe(1) {
				break gotoRead
			}
			opUint8 := data.GetUint8()
			switch opUint8 {
			case InputOpSearch:
				me.Search()
			case InputOpMissile:
				attack = true
			case InputOpContinueMove:
				move = true
			case InputOpNewMove:
				if data.NotSafe(4) {
					break gotoRead
				}
				me.Angle = data.GetFloat32()
				me.DeltaAngle = true
				move = true
			case InputOpChat:
				if data.NotSafe(1) {
					break gotoRead
				}
				num := data.GetUint8()
				chat := make([]uint8, num)
				if data.NotSafe(int(num)) {
					break gotoRead
				}
				for ch := uint8(0); ch < num; ch++ {
					chat[ch] = data.GetUint8()
				}
				me.World.broadcastCount++
				me.World.broadcast.PutUint8(BroadcastChat)
				me.World.broadcast.PutUint8(num)
				for ch := uint8(0); ch < num; ch++ {
					me.World.broadcast.PutUint8(chat[ch])
				}
			}
		}
	}
	person.InputCount = 0

	if attack {
		me.Status = HumanMissile
		me.DeltaStatus = true
		me.AnimationFrame = 0
		me.Animation = HumanMissileAnimation
	} else if move {
		me.DeltaX += float32(math.Sin(float64(me.Angle))) * me.Speed
		me.DeltaZ -= float32(math.Cos(float64(me.Angle))) * me.Speed
		me.IntegrateXZ()
		if me.Status == HumanIdle {
			me.Status = HumanWalk
			me.DeltaStatus = true
		}
	} else if me.Status == HumanWalk {
		me.Status = HumanIdle
		me.DeltaStatus = true
	}
}

// Update func
func (me *You) Update() bool {
	switch me.Status {
	case HumanDead:
		me.Dead()
	case HumanMissile:
		me.missile()
	default:
		me.Walk()
	}
	me.IntegrateY()
	return false
}
