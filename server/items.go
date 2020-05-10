package main

import (
	"../fast"
)

// item struct
type item struct {
	World               *World
	UID                 uint16
	NID                 uint16
	X, Y, Z             float32
	MinBX, MinBY, MinBZ int
	MaxBX, MaxBY, MaxBZ int
	Radius              float32
	Height              float32
}

// blockBorders func
func (me *item) blockBorders() {
	me.MinBX = int((me.X - me.Radius) * InverseBlockSize)
	me.MinBY = int(me.Y * InverseBlockSize)
	me.MinBZ = int((me.Z - me.Radius) * InverseBlockSize)
	me.MaxBX = int((me.X + me.Radius) * InverseBlockSize)
	me.MaxBY = int((me.Y + me.Height) * InverseBlockSize)
	me.MaxBZ = int((me.Z + me.Radius) * InverseBlockSize)
}

// addToBlocks func
func (me *item) addToBlocks() {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
				block.addItem(me)
			}
		}
	}
}

// removeFromBlocks func
func (me *item) removeFromBlocks() {
	for gx := me.MinBX; gx <= me.MaxBX; gx++ {
		for gy := me.MinBY; gy <= me.MaxBY; gy++ {
			for gz := me.MinBZ; gz <= me.MaxBZ; gz++ {
				block := me.World.getBlock(gx, gy, gz)
				if block != nil {
					block.removeItem(me)
				}
			}
		}
	}
}

// Overlap func
func (me *item) Overlap(b *thing) bool {
	square := me.Radius + b.Radius
	return Abs(me.X-b.X) <= square && Abs(me.Z-b.Z) <= square
}

// Cleanup func
func (me *item) Cleanup() {
	me.removeFromBlocks()
	me.World.removeItem(me)
	me.BroadcastDelete()
}

// BroadcastDelete func
func (me *item) BroadcastDelete() {
	me.World.broadcastCount++
	me.World.broadcast.PutUint8(BroadcastDelete)
	me.World.broadcast.PutUint16(me.NID)
}

// Save func
func (me *item) Save(data *fast.ByteWriter) {
	data.PutUint16(me.UID)
	data.PutUint16(me.NID)
	data.PutFloat32(me.X)
	data.PutFloat32(me.Y)
	data.PutFloat32(me.Z)
}

// LoadNewItem func
func LoadNewItem(world *World, uid uint16, x, y, z float32) {
	switch uid {
	case MedkitUID:
		NewMedkit(world, x, y, z)
	}
}

// NewMedkit func
func NewMedkit(world *World, x, y, z float32) *item {
	me := &item{}
	me.World = world
	me.X = x
	me.Y = y
	me.Z = z
	me.Radius = 0.3
	me.Height = 0.3
	me.blockBorders()
	me.addToBlocks()
	me.UID = MedkitUID
	me.NID = NextNID()
	world.addItem(me)
	return me
}
