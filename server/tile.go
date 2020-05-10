package main

// Tile constants
const (
	TileNone       = 0
	TileGrass      = 1
	TilePlankFloor = 2
	TilePlanks     = 3
	TileStone      = 4
	TileStoneFloor = 5
)

// Tile variables
var (
	TileClosed = []bool{
		false,
		true,
		true,
		true,
		true,
		true,
	}
)
