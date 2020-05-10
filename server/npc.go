package main

// Npc constants
const (
	DirectionNorth     = uint8(0)
	DirectionNorthEast = uint8(1)
	DirectionEast      = uint8(2)
	DirectionSouthEast = uint8(3)
	DirectionSouth     = uint8(4)
	DirectionSouthWest = uint8(5)
	DirectionWest      = uint8(6)
	DirectionNorthWest = uint8(7)
	DirectionNone      = uint8(8)
	DirectionCount     = uint8(8)
)

// Npc variables
var (
	OppositeDirection = []uint8{
		DirectionSouth,
		DirectionSouthWest,
		DirectionWest,
		DirectionNorthWest,
		DirectionNorth,
		DirectionNorthEast,
		DirectionEast,
		DirectionSouthEast,
		DirectionNone,
	}
	DiagonalDirection = []uint8{
		DirectionNorthEast,
		DirectionNorthWest,
		DirectionSouthEast,
		DirectionSouthWest,
	}
	NpcMoveX = []float32{
		0.0, 0.5, 1.0, 0.5, 0.0, -0.5, -1.0, -0.5,
	}
	NpcMoveZ = []float32{
		-1.0, -0.5, 0.0, 0.5, 1.0, 0.5, 0.0, -0.5,
	}
)

// Npc struct
type Npc struct {
	*thing
	Target             *thing
	MoveCount          uint8
	MoveDirection      uint8
	DeltaMoveDirection bool
}

// NewChaseDirection func
func (me *Npc) NewChaseDirection() {
	me.DeltaMoveDirection = true
	const epsilon = 0.32
	dx := me.X - me.Target.X
	dz := me.Z - me.Target.X
	old := me.MoveDirection
	opposite := OppositeDirection[old]

	var directionX uint8
	if dx > epsilon {
		directionX = DirectionWest
	} else if dx < -epsilon {
		directionX = DirectionEast
	} else {
		directionX = DirectionNone
	}

	var directionZ uint8
	if dz > epsilon {
		directionZ = DirectionNorth
	} else if dz < -epsilon {
		directionZ = DirectionSouth
	} else {
		directionZ = DirectionNone
	}

	if directionX != DirectionNone && directionZ != DirectionNone {
		d := 0
		if dz < 0 {
			d += 2
		}
		if dx > 0 {
			d++
		}
		me.MoveDirection = DiagonalDirection[d]
		if me.MoveDirection != opposite && me.TestMove() {
			return
		}
	}

	if NextRandP() > 200 || Abs(dz) > Abs(dx) {
		temp := directionX
		directionX = directionZ
		directionZ = temp
	}

	if directionX != opposite {
		me.MoveDirection = directionX
		if me.TestMove() {
			return
		}
	}

	if directionZ != opposite {
		me.MoveDirection = directionZ
		if me.TestMove() {
			return
		}
	}

	if old != DirectionNone {
		me.MoveDirection = old
		if me.TestMove() {
			return
		}
	}

	if NextRandP()&1 > 0 {
		for d := uint8(0); d < DirectionCount; d++ {
			if d == opposite {
				continue
			}
			me.MoveDirection = d
			if me.TestMove() {
				return
			}
		}
	} else {
		for d := DirectionCount - 1; d >= 0; d-- {
			if d == opposite {
				continue
			}
			me.MoveDirection = d
			if me.TestMove() {
				return
			}
		}
	}

	if opposite != DirectionNone {
		me.MoveDirection = opposite
		if me.TestMove() {
			return
		}
	}

	me.MoveDirection = DirectionNone
}

// TestMove func
func (me *Npc) TestMove() bool {
	if !me.Move() {
		return false
	}
	me.MoveCount = 16 + NextRandP()&32
	return true
}

// TryMove func
func (me *Npc) TryMove(x, z float32) bool {
	minGX := int((x - me.Radius))
	minGY := int(me.Y)
	minGZ := int((z - me.Radius))
	maxGX := int((x + me.Radius))
	maxGY := int((me.Y + me.Height))
	maxGZ := int((z + me.Radius))

	minBX := int(float32(minGX) * InverseBlockSize)
	minBY := int(float32(minGY) * InverseBlockSize)
	minBZ := int(float32(minGZ) * InverseBlockSize)
	maxBX := int(float32(maxGX) * InverseBlockSize)
	maxBY := int(float32(maxGY) * InverseBlockSize)
	maxBZ := int(float32(maxGZ) * InverseBlockSize)

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
					return false
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

	searched := make(map[*thing]bool)

	for bx := minBX; bx <= maxBX; bx++ {
		for by := minBY; by <= maxBY; by++ {
			for bz := minBZ; bz <= maxBZ; bz++ {
				block := world.getBlock(bx, by, bz)
				for t := 0; t < block.thingCount; t++ {
					thing := block.things[t]
					if me.thing == thing {
						continue
					}
					if _, has := searched[thing]; !has {
						searched[thing] = true
						if me.TryOverlap(x, z, thing) {
							return false
						}
					}
				}
			}
		}
	}

	return true
}

// Move func
func (me *Npc) Move() bool {
	if me.MoveDirection == DirectionNone {
		return true
	}
	tryX := me.X + NpcMoveX[me.MoveDirection]*me.Speed
	tryZ := me.Z + NpcMoveZ[me.MoveDirection]*me.Speed
	if me.TryMove(tryX, tryZ) {
		me.removeFromBlocks()
		me.OldX = me.X
		me.OldZ = me.Z
		me.X = tryX
		me.Z = tryZ
		me.DeltaMoveXZ = true
		me.blockBorders()
		me.addToBlocks()
		return true
	}
	return false
}
