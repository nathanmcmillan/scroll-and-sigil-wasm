class Item {
    constructor() {
        this.World = null
        this.UID = 0
        this.SID = "item"
        this.NID = 0
        this.Sprite = null
        this.X = 0
        this.Y = 0
        this.Z = 0
        this.MinBX = 0
        this.MinBY = 0
        this.MinBZ = 0
        this.MaxBX = 0
        this.MaxBY = 0
        this.MaxBZ = 0
        this.Radius = 0
        this.Height = 0
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
                    let block = this.World.GetBlock(gx, gy, gz)
                    block.AddItem(this)
                }
            }
        }
    }
    RemoveFromBlocks() {
        for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
            for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
                for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
                    let block = this.World.GetBlock(gx, gy, gz)
                    block.RemoveItem(this)
                }
            }
        }
    }
    Cleanup() {
        this.World.RemoveItem(this)
        this.RemoveFromBlocks()
    }
    Render(spriteBuffer, camX, camZ, camAngle) {
        let sin = camX - this.X
        let cos = camZ - this.Z
        let length = Math.sqrt(sin * sin + cos * cos)
        sin /= length
        cos /= length
        Render3.Sprite(spriteBuffer.get(this.SID), this.X, this.Y, this.Z, sin, cos, this.Sprite)
    }
}

class Medkit extends Item {
    constructor(world, nid, x, y, z) {
        super()
        this.World = world
        this.X = x
        this.Y = y
        this.Z = z
        this.BlockBorders()
        this.AddToBlocks()
        this.UID = PlasmaUID
        this.NID = nid
        this.Sprite = SpriteData[this.SID]["medkit"]
        this.Radius = 0.2
        this.Height = 0.2
        world.AddItem(this)
        return this
    }
}
