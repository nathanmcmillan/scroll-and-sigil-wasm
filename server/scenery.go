package main

import (
	"../fast"
)

// NewTree func
func NewTree(world *World, x, y, z float32) *thing {
	tree := &thing{}
	tree.UID = TreeUID
	tree.NID = NextNID()
	tree.World = world
	tree.Update = tree.NopUpdate
	tree.Damage = tree.NopDamage
	tree.Save = tree.ScenerySave
	tree.Snap = tree.NopSnap
	tree.X = x
	tree.Y = y
	tree.Z = z
	tree.Radius = 0.4
	tree.Height = 1.0
	tree.Health = 1
	world.addThing(tree)
	tree.blockBorders()
	tree.addToBlocks()
	return tree
	// TODO make scenery its own entity
}

// ScenerySave func
func (me *thing) ScenerySave(data *fast.ByteWriter) {
	data.PutUint16(me.UID)
	data.PutUint16(me.NID)
	data.PutFloat32(me.X)
	data.PutFloat32(me.Y)
	data.PutFloat32(me.Z)
}
