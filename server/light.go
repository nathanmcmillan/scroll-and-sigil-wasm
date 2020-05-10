package main

import (
	"../fast"
)

type light struct {
	X   int
	Y   int
	Z   int
	RGB int
}

func lightInit(x, y, z, rgb int) *light {
	light := &light{X: x, Y: y, Z: z, RGB: rgb}
	return light
}

// Save func
func (me *light) Save(data *fast.ByteWriter) {
	data.PutUint8(uint8(me.X))
	data.PutUint8(uint8(me.Y))
	data.PutUint8(uint8(me.Z))
	data.PutInt32(int32(me.RGB))
}
