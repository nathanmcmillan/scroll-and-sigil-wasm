package main

import (
	"math"

	"../fast"
)

// thing constants
const (
	AnimationRate       = 5
	Gravity             = 0.01
	AnimationNotDone    = 0
	AnimationAlmostDone = 1
	AnimationDone       = 2
)

// UID constants
const (
	HumanUID  = uint16(0)
	BaronUID  = uint16(1)
	TreeUID   = uint16(2)
	PlasmaUID = uint16(3)
	MedkitUID = uint16(4)
)

// Group constants
const (
	NoGroup    = 0
	HumanGroup = 1
	DemonGroup = 2
)

// thing variables
var (
	thingNetworkNum = uint16(0)
)

// thing struct
type thing struct {
	World          *World
	UID            uint16
	NID            uint16
	Animation      int
	AnimationFrame int
	X              float32
	Y              float32
	Z              float32
	Angle          float32
	DeltaX         float32
	DeltaY         float32
	DeltaZ         float32
	OldX           float32
	OldZ           float32
	MinBX          int
	MinBY          int
	MinBZ          int
	MaxBX          int
	MaxBY          int
	MaxBZ          int
	Ground         bool
	Radius         float32
	Height         float32
	Speed          float32
	Health         uint16
	Group          int
	DeltaMoveXZ    bool
	DeltaMoveY     bool
	Update         func() bool
	Damage         func(uint16)
	Save           func(data *fast.ByteWriter)
	Snap           func(data *fast.ByteWriter)
	Binary         []byte
}

// NextNID func
func NextNID() uint16 {
	thingNetworkNum++
	return thingNetworkNum
}

// LoadNewthing func
func LoadNewthing(world *World, uid uint16, x, y, z float32) {
	switch uid {
	case BaronUID:
		NewBaron(world, x, y, z)
	case TreeUID:
		NewTree(world, x, y, z)
	}
}

// NopUpdate func
func (me *thing) NopUpdate() bool {
	return false
}

// NopSnap func
func (me *thing) NopSnap(data *fast.ByteWriter) {
	me.Binary = nil
}

// NopDamage func
func (me *thing) NopDamage(amount uint16) {
}

// blockBorders func
func (me *thing) blockBorders() {
	me.MinBX = int((me.X - me.Radius) * InverseBlockSize)
	me.MinBY = int(me.Y * InverseBlockSize)
	me.MinBZ = int((me.Z - me.Radius) * InverseBlockSize)
	me.MaxBX = int((me.X + me.Radius) * InverseBlockSize)
	me.MaxBY = int((me.Y + me.Height) * InverseBlockSize)
	me.MaxBZ = int((me.Z + me.Radius) * InverseBlockSize)
}

// addToBlocks func
func (me *thing) addToBlocks() {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				me.World.getBlock(gx, gy, gz).addThing(me)
			}
		}
	}
}

// removeFromBlocks func
func (me *thing) removeFromBlocks() {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				me.World.getBlock(gx, gy, gz).removeThing(me)
			}
		}
	}
}

// UpdateAnimation func
func (me *thing) UpdateAnimation() int {
	me.AnimationFrame++
	if me.AnimationFrame == me.Animation-AnimationRate {
		return AnimationAlmostDone
	} else if me.AnimationFrame == me.Animation {
		return AnimationDone
	}
	return AnimationNotDone
}

// TerrainCollisionXZ func
func (me *thing) TerrainCollisionXZ() {
	minGX := int((me.X - me.Radius))
	minGY := int(me.Y)
	minGZ := int((me.Z - me.Radius))
	maxGX := int((me.X + me.Radius))
	maxGY := int((me.Y + me.Height))
	maxGZ := int((me.Z + me.Radius))

	minBX := int(float32(minGX) * InverseBlockSize)
	minBY := int(float32(minGY) * InverseBlockSize)
	minBZ := int(float32(minGZ) * InverseBlockSize)

	minTX := minGX - minBX*BlockSize
	minTY := minGY - minBY*BlockSize
	minTZ := minGZ - minBZ*BlockSize

	world := me.World

	bx := minBX
	tx := minTX
	for gx := minGX; gx <= maxGX; gx++ {
		by := minBY
		ty := minTY
		for gy := minGY; gy <= maxGY; gy++ {
			bz := minBZ
			tz := minTZ
			for gz := minGZ; gz <= maxGZ; gz++ {
				block := world.getBlock(bx, by, bz)
				tile := block.GetTileTypeUnsafe(tx, ty, tz)
				if TileClosed[tile] {
					xx := float32(gx)
					closeX := me.X
					if closeX < xx {
						closeX = xx
					} else if closeX > xx+1 {
						closeX = xx + 1
					}

					zz := float32(gz)
					closeZ := me.Z
					if closeZ < zz {
						closeZ = zz
					} else if closeZ > zz+1 {
						closeZ = zz + 1
					}

					dxx := me.X - closeX
					dzz := me.Z - closeZ
					dist := dxx*dxx + dzz*dzz

					if dist > me.Radius*me.Radius {
						continue
					}
					dist = float32(math.Sqrt(float64(dist)))
					if dist == 0 {
						dist = 1
					}
					mult := me.Radius / dist
					me.X += dxx*mult - dxx
					me.Z += dzz*mult - dzz
				}
				tz++
				if tz == BlockSize {
					tz = 0
					bz++
				}
			}
			ty++
			if ty == BlockSize {
				ty = 0
				by++
			}
		}
		tx++
		if tx == BlockSize {
			tx = 0
			bx++
		}
	}
}

// TerrainCollisionY func
func (me *thing) TerrainCollisionY() {
	if me.DeltaY < 0 {
		gx := int(me.X)
		gy := int(me.Y)
		gz := int(me.Z)
		bx := int(me.X * InverseBlockSize)
		by := int(me.Y * InverseBlockSize)
		bz := int(me.Z * InverseBlockSize)
		tx := gx - bx*BlockSize
		ty := gy - by*BlockSize
		tz := gz - bz*BlockSize

		tile := me.World.GetTileType(bx, by, bz, tx, ty, tz)
		if TileClosed[tile] {
			me.Y = float32(gy + 1)
			me.Ground = true
			me.DeltaY = 0
		}
	}
}

// Resolve func
func (me *thing) Resolve(b *thing) {
	square := me.Radius + b.Radius
	absx := Abs(me.X - b.X)
	absz := Abs(me.Z - b.Z)
	if absx > square || absz > square {
		return
	}
	x := me.OldX - b.X
	z := me.OldZ - b.Z
	if Abs(x) > Abs(z) {
		if x < 0 {
			me.X = b.X - square
		} else {
			me.X = b.X + square
		}
		me.DeltaX = 0.0
	} else {
		if z < 0 {
			me.Z = b.Z - square
		} else {
			me.Z = b.Z + square
		}
		me.DeltaZ = 0.0
	}
}

// Overlap func
func (me *thing) Overlap(b *thing) bool {
	square := me.Radius + b.Radius
	return Abs(me.X-b.X) <= square && Abs(me.Z-b.Z) <= square
}

// TryOverlap func
func (me *thing) TryOverlap(x, z float32, b *thing) bool {
	square := me.Radius + b.Radius
	return Abs(x-b.X) <= square && Abs(z-b.Z) <= square
}

// ApproximateDistance func
func (me *thing) ApproximateDistance(other *thing) float32 {
	dx := Abs(me.X - other.X)
	dz := Abs(me.Z - other.Z)
	if dx > dz {
		return dx + dz - dz*0.5
	}
	return dx + dz - dx*0.5
}

// IntegrateXZ func
func (me *thing) IntegrateXZ() {
	me.OldX = me.X
	me.OldZ = me.Z

	me.X += me.DeltaX
	me.Z += me.DeltaZ
	me.DeltaMoveXZ = true

	collided := make([]*thing, 0)
	searched := make(map[*thing]bool)

	me.removeFromBlocks()
	me.blockBorders()

	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
				for t := 0; t < block.thingCount; t++ {
					thing := block.things[t]
					if _, has := searched[thing]; !has {
						searched[thing] = true
						if me.Overlap(thing) {
							collided = append(collided, thing)
						}
					}
				}
			}
		}
	}

	num := len(collided)
	for num > 0 {
		closest := 0
		manhattan := float32(math.MaxFloat32)
		for i := 0; i < num; i++ {
			thing := collided[i]
			dist := Abs(me.OldX-thing.X) + Abs(me.OldZ-thing.Z)
			if dist < manhattan {
				manhattan = dist
				closest = i
			}
		}
		me.Resolve(collided[closest])
		collided[closest] = collided[num-1]
		num--
	}

	me.TerrainCollisionXZ()

	me.blockBorders()
	me.addToBlocks()

	me.DeltaX = 0.0
	me.DeltaZ = 0.0
}

// IntegrateY func
func (me *thing) IntegrateY() {
	if !me.Ground {
		me.DeltaY -= Gravity
		me.Y += me.DeltaY
		me.DeltaMoveY = true
		me.TerrainCollisionY()

		me.removeFromBlocks()
		me.blockBorders()
		me.addToBlocks()
	}
}
