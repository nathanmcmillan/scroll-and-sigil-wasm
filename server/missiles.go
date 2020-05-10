package main

import (
	"../fast"
)

// missile struct
type missile struct {
	World                  *World
	UID                    uint16
	NID                    uint16
	X, Y, Z                float32
	DeltaX, DeltaY, DeltaZ float32
	MinBX, MinBY, MinBZ    int
	MaxBX, MaxBY, MaxBZ    int
	Radius                 float32
	Height                 float32
	DamageAmount           uint16
	Hit                    func(thing *thing)
}

// blockBorders func
func (me *missile) blockBorders() {
	me.MinBX = int((me.X - me.Radius) * InverseBlockSize)
	me.MinBY = int(me.Y * InverseBlockSize)
	me.MinBZ = int((me.Z - me.Radius) * InverseBlockSize)
	me.MaxBX = int((me.X + me.Radius) * InverseBlockSize)
	me.MaxBY = int((me.Y + me.Height) * InverseBlockSize)
	me.MaxBZ = int((me.Z + me.Radius) * InverseBlockSize)
}

// addToBlocks func
func (me *missile) addToBlocks() bool {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
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

// removeFromBlocks func
func (me *missile) removeFromBlocks() {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
				if block != nil {
					block.removeMissile(me)
				}
			}
		}
	}
}

// Overlap func
func (me *missile) Overlap(b *thing) bool {
	square := me.Radius + b.Radius
	return Abs(me.X-b.X) <= square && Abs(me.Z-b.Z) <= square
}

// Collision func
func (me *missile) Collision() bool {
	searched := make(map[*thing]bool)
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
				for t := 0; t < block.thingCount; t++ {
					thing := block.things[t]
					if thing.Health > 0 {
						if _, has := searched[thing]; !has {
							searched[thing] = true
							if me.Overlap(thing) {
								me.Hit(thing)
								return true
							}
						}
					}
				}
			}
		}
	}
	minGX := int((me.X - me.Radius))
	minGY := int(me.Y)
	minGZ := int((me.Z - me.Radius))
	maxGX := int((me.X + me.Radius))
	maxGY := int((me.Y + me.Height))
	maxGZ := int((me.Z + me.Radius))
	for gx := minGX; gx <= maxGX; gx++ {
		for gy := minGY; gy <= maxGY; gy++ {
			for gz := minGZ; gz <= maxGZ; gz++ {
				bx := int(float32(gx) * InverseBlockSize)
				by := int(float32(gy) * InverseBlockSize)
				bz := int(float32(gz) * InverseBlockSize)
				tx := gx - bx*BlockSize
				ty := gy - by*BlockSize
				tz := gz - bz*BlockSize
				tile := me.World.GetTileType(bx, by, bz, tx, ty, tz)
				if TileClosed[tile] {
					me.Hit(nil)
					return true
				}
			}
		}
	}
	return false
}

// Update func
func (me *missile) Update() bool {
	if me.Collision() {
		return true
	}
	me.removeFromBlocks()
	me.X += me.DeltaX
	me.Y += me.DeltaY
	me.Z += me.DeltaZ
	me.blockBorders()
	if me.addToBlocks() {
		return true
	}
	return me.Collision()
}

// Snap func
func (me *missile) Snap(data *fast.ByteWriter) {
	data.PutUint16(me.UID)
	data.PutUint16(me.NID)
	data.PutFloat32(me.X)
	data.PutFloat32(me.Y)
	data.PutFloat32(me.Z)
	data.PutFloat32(me.DeltaX)
	data.PutFloat32(me.DeltaY)
	data.PutFloat32(me.DeltaZ)
	data.PutUint16(me.DamageAmount)
}

// BroadcastNew func
func (me *missile) BroadcastNew() {
	me.World.broadcastCount++
	me.World.broadcast.PutUint8(BroadcastNew)
	me.Snap(me.World.broadcast)
}

// BroadcastDelete func
func (me *missile) BroadcastDelete() {
	me.World.broadcastCount++
	me.World.broadcast.PutUint8(BroadcastDelete)
	me.World.broadcast.PutUint16(me.NID)
}

// NewPlasma func
func NewPlasma(world *World, damage uint16, x, y, z, dx, dy, dz float32) {
	me := &missile{}
	me.World = world
	me.X = x
	me.Y = y
	me.Z = z
	me.Radius = 0.2
	me.Height = 0.2
	me.blockBorders()
	if me.addToBlocks() {
		return
	}
	me.UID = PlasmaUID
	me.NID = NextNID()
	me.DeltaX = dx
	me.DeltaY = dy
	me.DeltaZ = dz
	me.DamageAmount = damage
	me.Hit = me.PlasmaHit

	world.addMissile(me)
	me.BroadcastNew()
}

// PlasmaHit func
func (me *missile) PlasmaHit(thing *thing) {
	me.X -= me.DeltaX
	me.Y -= me.DeltaY
	me.Z -= me.DeltaZ
	if thing != nil {
		thing.Damage(me.DamageAmount)
	}
	me.removeFromBlocks()
	me.BroadcastDelete()
}
