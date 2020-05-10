package main

// Constants
const (
	TileNone = 0

	AmbientLow  = 100.0
	AmbientHalf = 175.0
	AmbientFull = 255.0
)

// Variables
var (
	TileLookup  = map[string]int{}
	TileTexture [][4]float32
	TileClosed  []bool
)

// TileAmbient func
func TileAmbient(side1, side2, corner bool) float32 {
	if side1 && side2 {
		return AmbientLow
	} else if side1 || side2 || corner {
		return AmbientHalf
	}
	return AmbientFull
}

type tile struct {
	typeOf int
	red    uint8
	green  uint8
	blue   uint8
}

func tileInit() *tile {
	return &tile{}
}
