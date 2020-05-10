const TileNone = 0
const TileGrass = 1
const TilePlankFloor = 2

const AMBIENT_LOW = 100
const AMBIENT_HALF = 175
const AMBIENT_FULL = 255

const TileLookup = new Map()
const TileTexture = []
const TileClosed = []

class Tile {
    constructor() {
        this.type = TileNone
        this.red = 0
        this.green = 0
        this.blue = 0
    }
    static Ambient(side1, side2, corner) {
        if (side1 && side2)
            return AMBIENT_LOW
        if (side1 || side2 || corner)
            return AMBIENT_HALF
        return AMBIENT_FULL
    }
}
