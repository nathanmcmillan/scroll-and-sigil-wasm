package main

import "../fast"

type camera struct {
	thing  *thing
	radius float32
	x      float32
	y      float32
	z      float32
	rx     float32
	ry     float32
}

func cameraInit(world *world, thing *thing, radius float32) *camera {
	c := &camera{}
	c.thing = thing
	c.radius = radius
	c.x = 0
	c.y = 0
	c.z = 0
	c.rx = 0
	c.ry = 0
	c.update(world)
	return c
}

func (me *camera) update(world *world) {
	if InputIsKeyDown("ArrowLeft") {
		me.ry -= 0.05
		if me.ry < 0 {
			me.ry += Tau
		}
	}

	if InputIsKeyDown("ArrowRight") {
		me.ry += 0.05
		if me.ry >= Tau {
			me.ry -= Tau
		}
	}

	if me.rx > -0.25 && InputIsKeyDown("ArrowUp") {
		me.rx -= 0.05
	}

	if me.rx < 0.25 && InputIsKeyDown("ArrowDown") {
		me.rx += 0.05
	}

	// sinX := math.Sin(float64(me.rx))
	// cosX := math.Cos(float64(me.rx))
	// sinY := math.Sin(float64(me.ry))
	// cosY := math.Cos(float64(me.ry))
	sinX := fast.Sin(me.rx)
	cosX := fast.Cos(me.rx)
	sinY := fast.Sin(me.ry)
	cosY := fast.Cos(me.ry)

	thing := me.thing

	dx := float32(-cosX * sinY)
	dy := float32(sinX)
	dz := float32(cosX * cosY)

	x := thing.x + me.radius*dx
	y := thing.y + me.radius*dy + thing.height
	z := thing.z + me.radius*dz

	// Cast.Exact(world, thing.X, thing.Y + thing.Height, thing.Z, x, y, z)

	CastX := float32(0.0)
	CastY := float32(0.0)
	CastZ := float32(0.0)

	if CastX >= 0 {
		me.x = x
		me.y = y
		me.z = z
	} else {
		me.x = CastX
		me.y = CastY
		me.z = CastZ
	}
}
