const Gravity = 0.01

const AnimationRate = 16

const AnimationNotDone = 0
const AnimationAlmostDone = 1
const AnimationDone = 2

const AnimationFront = 0
const AnimationFrontSide = 1
const AnimationSide = 2
const AnimationBackSide = 3
const AnimationBack = 4

const DirectionNorth = 0
const DirectionNorthEast = 1
const DirectionEast = 2
const DirectionSouthEast = 3
const DirectionSouth = 4
const DirectionSouthWest = 5
const DirectionWest = 6
const DirectionNorthWest = 7
const DirectionCount = 8
const DirectionNone = 8

const DirectionToAngle = [
    0.0 * DegToRad,
    45.0 * DegToRad,
    90.0 * DegToRad,
    135.0 * DegToRad,
    180.0 * DegToRad,
    225.0 * DegToRad,
    270.0 * DegToRad,
    315.0 * DegToRad
]

const ThingAngleA = 337.5 * DegToRad
const ThingAngleB = 292.5 * DegToRad
const ThingAngleC = 247.5 * DegToRad
const ThingAngleD = 202.5 * DegToRad
const ThingAngleE = 157.5 * DegToRad
const ThingAngleF = 112.5 * DegToRad
const ThingAngleG = 67.5 * DegToRad
const ThingAngleH = 22.5 * DegToRad

const HumanUID = 0
const BaronUID = 1
const TreeUID = 2
const PlasmaUID = 3
const MedkitUID = 4

class Thing {
    constructor() {
        this.World = null
        this.UID = 0
        this.SID = ""
        this.NID = 0
        this.Animation = null
        this.AnimationMod = 0
        this.AnimationFrame = 0
        this.X = 0
        this.Y = 0
        this.Z = 0
        this.Angle = 0
        this.DeltaX = 0
        this.DeltaY = 0
        this.DeltaZ = 0
        this.OldX = 0
        this.OldY = 0
        this.OldZ = 0
        this.NetX = 0
        this.NetY = 0
        this.NetZ = 0
        this.DeltaNetX = 0
        this.DeltaNetY = 0
        this.DeltaNetZ = 0
        this.MinBX = 0
        this.MinBY = 0
        this.MinBZ = 0
        this.MaxBX = 0
        this.MaxBY = 0
        this.MaxBZ = 0
        this.Ground = false
        this.Radius = 0
        this.Height = 0
        this.Speed = 0
        this.Health = 0
    }
    BlockBorders() {
        this.MinBX = Math.floor((this.X - this.Radius) * InverseBlockSize)
        this.MinBY = Math.floor(this.Y * InverseBlockSize)
        this.MinBZ = Math.floor((this.Z - this.Radius) * InverseBlockSize)
        this.MaxBX = Math.floor((this.X + this.Radius) * InverseBlockSize)
        this.MaxBY = Math.floor((this.Y + this.Height) * InverseBlockSize)
        this.MaxBZ = Math.floor((this.Z + this.Radius) * InverseBlockSize)
    }
    AddToBlocks() {
        for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
            for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
                for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
                    this.World.GetBlock(gx, gy, gz).AddThing(this)
                }
            }
        }
    }
    RemoveFromBlocks() {
        for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
            for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
                for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
                    this.World.GetBlock(gx, gy, gz).RemoveThing(this)
                }
            }
        }
    }
    Cleanup() {
        this.World.RemoveThing(this)
        this.RemoveFromBlocks()
    }
    UpdateAnimation() {
        this.AnimationMod++
        if (this.AnimationMod === AnimationRate) {
            this.AnimationMod = 0
            this.AnimationFrame++
            let len = this.Animation.length
            if (this.AnimationFrame === len - 1)
                return AnimationAlmostDone
            else if (this.AnimationFrame === len)
                return AnimationDone
        }
        return AnimationNotDone
    }
    NetUpdateState(_) {}
    NetUpdateHealth(_) {}
    TerrainCollisionY(world) {
        if (this.DeltaY < 0) {
            let gx = Math.floor(this.X)
            let gy = Math.floor(this.Y)
            let gz = Math.floor(this.Z)
            let bx = Math.floor(gx * InverseBlockSize)
            let by = Math.floor(gy * InverseBlockSize)
            let bz = Math.floor(gz * InverseBlockSize)
            let tx = gx - bx * BlockSize
            let ty = gy - by * BlockSize
            let tz = gz - bz * BlockSize

            let tile = world.GetTileType(bx, by, bz, tx, ty, tz)
            if (TileClosed[tile]) {
                this.Y = gy + 1
                this.Ground = true
                this.DeltaY = 0
            }
        }
    }
    Resolve(b) {
        let square = this.Radius + b.Radius
        if (Math.abs(this.X - b.X) > square || Math.abs(this.Z - b.Z) > square)
            return
        if (Math.abs(this.OldX - b.X) > Math.abs(this.OldZ - b.Z)) {
            if (this.OldX - b.X < 0) this.X = b.X - square
            else this.X = b.X + square
            this.DeltaX = 0.0
        } else {
            if (this.OldZ - b.Z < 0) this.Z = b.Z - square
            else this.Z = b.Z + square
            this.DeltaZ = 0.0
        }
    }
    Overlap(b) {
        let square = this.Radius + b.Radius
        return Math.abs(this.X - b.X) <= square && Math.abs(this.Z - b.Z) <= square
    }
    LerpNetCode() {
        let updateBlocks = false

        if (this.DeltaNetX > 0) {
            this.X += this.DeltaNetX
            updateBlocks = true
            if (this.X >= this.NetX) {
                this.X = this.NetX
                this.DeltaNetX = 0
            }
        } else if (this.DeltaNetX < 0) {
            this.X += this.DeltaNetX
            updateBlocks = true
            if (this.X <= this.NetX) {
                this.X = this.NetX
                this.DeltaNetX = 0
            }
        }

        if (this.DeltaNetY > 0) {
            this.Y += this.DeltaNetY
            updateBlocks = true
            if (this.Y >= this.NetY) {
                this.Y = this.NetY
                this.DeltaNetY = 0
            }
        } else if (this.DeltaNetY < 0) {
            this.Y += this.DeltaNetY
            updateBlocks = true
            if (this.Y <= this.NetY) {
                this.Y = this.NetY
                this.DeltaNetY = 0
            }
        }

        if (this.DeltaNetZ > 0) {
            this.Z += this.DeltaNetZ
            updateBlocks = true
            if (this.Z >= this.NetZ) {
                this.Z = this.NetZ
                this.DeltaNetZ = 0
            }
        } else if (this.DeltaNetZ < 0) {
            this.Z += this.DeltaNetZ
            updateBlocks = true
            if (this.Z <= this.NetZ) {
                this.Z = this.NetZ
                this.DeltaNetZ = 0
            }
        }

        if (updateBlocks) {
            this.RemoveFromBlocks()
            this.BlockBorders()
            this.AddToBlocks()
        }
    }
    Integrate() {
        // OldX and snapshot need to be different things

        // this.OldX = this.X
        // this.OldY = this.Y
        // this.OldZ = this.Z

        // if (this.DeltaX != 0.0 || this.DeltaZ != 0.0) {
        //     this.X += this.DeltaX
        //     this.Z += this.DeltaZ

        //     let collided = []
        //     let searched = new Set()

        //     this.RemoveFromBlocks(world)

        //     for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
        //         for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
        //             for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
        //                 let block = world.GetBlock(gx, gy, gz)
        //                 for (let t = 0; t < block.thingCount; t++) {
        //                     let thing = block.things[t]
        //                     if (searched.has(thing)) continue
        //                     searched.add(thing)
        //                     if (this.Overlap(thing)) collided.push(thing)
        //                 }
        //             }
        //         }
        //     }

        //     while (collided.length > 0) {
        //         let closest = null
        //         let manhattan = Number.MAX_VALUE
        //         for (let i = 0; i < collided.length; i++) {
        //             let thing = collided[i]
        //             let dist = Math.abs(this.OldX - thing.X) + Math.abs(this.OldZ - thing.Z)
        //             if (dist < manhattan) {
        //                 manhattan = dist
        //                 closest = thing
        //             }
        //         }
        //         this.Resolve(closest)
        //         collided.splice(closest)
        //     }

        //     this.BlockBorders()
        //     this.AddToBlocks(world)

        //     this.DeltaX = 0.0
        //     this.DeltaZ = 0.0
        // }

        // if (!this.Ground || this.DeltaY != 0.0) {
        //     this.DeltaY -= GRAVITY
        //     this.Y += this.DeltaY
        //     this.TerrainCollisionY(world)

        //     this.RemoveFromBlocks(world)
        //     this.BlockBorders()
        //     this.AddToBlocks(world)
        // }
    }
    Render(spriteBuffer, camX, camZ, camAngle) {
        let sin = camX - this.X
        let cos = camZ - this.Z
        let length = Math.sqrt(sin * sin + cos * cos)
        sin /= length
        cos /= length

        let angle = camAngle - this.Angle
        if (angle < 0) angle += Tau

        let direction
        let mirror

        if (angle > ThingAngleA) {
            direction = AnimationBack
            mirror = false
        } else if (angle > ThingAngleB) {
            direction = AnimationBackSide
            mirror = true
        } else if (angle > ThingAngleC) {
            direction = AnimationSide
            mirror = true
        } else if (angle > ThingAngleD) {
            direction = AnimationFrontSide
            mirror = true
        } else if (angle > ThingAngleE) {
            direction = AnimationFront
            mirror = false
        } else if (angle > ThingAngleF) {
            direction = AnimationFrontSide
            mirror = false
        } else if (angle > ThingAngleG) {
            direction = AnimationSide
            mirror = false
        } else if (angle > ThingAngleH) {
            direction = AnimationBackSide
            mirror = false
        } else {
            direction = AnimationBack
            mirror = false
        }

        let sprite = this.Animation[this.AnimationFrame][direction]

        if (mirror) Render3.MirrorSprite(spriteBuffer.get(this.SID), this.X, this.Y, this.Z, sin, cos, sprite)
        else Render3.Sprite(spriteBuffer.get(this.SID), this.X, this.Y, this.Z, sin, cos, sprite)
    }
}
