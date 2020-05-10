package main

import (
	"../fast"
)

// block constants
const (
	BlockSize        = 8
	InverseBlockSize = 1.0 / BlockSize
	BlockSlice       = BlockSize * BlockSize
	BlockAll         = BlockSlice * BlockSize
)

type block struct {
	X, Y, Z      int
	Tiles        [BlockAll]int
	thingCount   int
	itemCount    int
	missileCount int
	things       []*thing
	items        []*item
	missiles     []*missile
	lights       []*light
	lightCount   int
}

// NewBlock func
func (me *block) blockInit(x, y, z int) {
	me.X = x
	me.Y = y
	me.Z = z
}

// Save func
func (me *block) Save(data *fast.ByteWriter) {
	notEmpty := me.NotEmpty()
	data.PutUint8(notEmpty)
	if notEmpty == 1 {
		for i := 0; i < BlockAll; i++ {
			data.PutUint8(uint8(me.Tiles[i]))
		}
	}
	data.PutUint8(uint8(me.lightCount))
	for i := 0; i < me.lightCount; i++ {
		me.lights[i].Save(data)
	}
}

// NotEmpty func
func (me *block) NotEmpty() uint8 {
	for i := 0; i < BlockAll; i++ {
		if me.Tiles[i] != TileNone {
			return 1
		}
	}
	return 0
}

// GetTileTypeUnsafe func
func (me *block) GetTileTypeUnsafe(x, y, z int) int {
	return me.Tiles[x+y*BlockSize+z*BlockSlice]
}

func (me *block) addThing(t *thing) {
	if me.thingCount == len(me.things) {
		array := make([]*thing, me.thingCount+5)
		copy(array, me.things)
		me.things = array
	}
	me.things[me.thingCount] = t
	me.thingCount++
}

func (me *block) addItem(t *item) {
	if me.itemCount == len(me.items) {
		array := make([]*item, me.itemCount+5)
		copy(array, me.items)
		me.items = array
	}
	me.items[me.itemCount] = t
	me.itemCount++
}

func (me *block) addMissile(t *missile) {
	if me.missileCount == len(me.missiles) {
		array := make([]*missile, me.missileCount+5)
		copy(array, me.missiles)
		me.missiles = array
	}
	me.missiles[me.missileCount] = t
	me.missileCount++
}

func (me *block) addLight(t *light) {
	if me.lightCount == len(me.lights) {
		array := make([]*light, me.lightCount+5)
		copy(array, me.lights)
		me.lights = array
	}
	me.lights[me.lightCount] = t
	me.lightCount++
}

func (me *block) removeThing(t *thing) {
	size := me.thingCount
	for i := 0; i < size; i++ {
		if me.things[i] == t {
			me.things[i] = me.things[size-1]
			me.things[size-1] = nil
			me.thingCount--
			break
		}
	}
}

func (me *block) removeItem(t *item) {
	size := me.itemCount
	for i := 0; i < size; i++ {
		if me.items[i] == t {
			me.items[i] = me.items[size-1]
			me.items[size-1] = nil
			me.itemCount--
			break
		}
	}
}

func (me *block) removeMissile(t *missile) {
	size := me.missileCount
	for i := 0; i < size; i++ {
		if me.missiles[i] == t {
			me.missiles[i] = me.missiles[size-1]
			me.missiles[size-1] = nil
			me.missileCount--
			break
		}
	}
}

func (me *block) removeLight(t *light) {
	size := me.lightCount
	for i := 0; i < size; i++ {
		if me.lights[i] == t {
			me.lights[i] = me.lights[size-1]
			me.lights[size-1] = nil
			me.lightCount--
			break
		}
	}
}
