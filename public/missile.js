class Missile {
    constructor() {
        this.World = null
        this.UID = 0
        this.SID = "missiles"
        this.NID = 0
        this.Sprite = null
        this.X = 0
        this.Y = 0
        this.Z = 0
        this.DeltaX = 0
        this.DeltaY = 0
        this.DeltaZ = 0
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
                    if (block === null) {
                        this.RemoveFromBlocks()
                        return true
                    }
                    block.AddMissile(this)
                }
            }
        }
        return false
    }
    RemoveFromBlocks() {
        for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
            for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
                for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
                    let block = this.World.GetBlock(gx, gy, gz)
                    if (block !== null)
                        block.RemoveMissile(this)
                }
            }
        }
    }
    Cleanup() {
        this.World.RemoveMissile(this)
        this.RemoveFromBlocks()
    }
    Update() {
        this.RemoveFromBlocks()
        this.X += this.DeltaX
        this.Y += this.DeltaY
        this.Z += this.DeltaZ
        this.BlockBorders()
        if (this.AddToBlocks())
            return true
        return false
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

class Plasma extends Missile {
    constructor(world, nid, damage, x, y, z, dx, dy, dz) {
        super()
        this.World = world
        this.X = x
        this.Y = y
        this.Z = z
        this.BlockBorders()
        if (this.AddToBlocks())
            return this
        this.UID = PlasmaUID
        this.NID = nid
        this.Sprite = SpriteData[this.SID]["baron-missile-front-1"]
        this.DeltaX = dx * InverseNetRate
        this.DeltaY = dy * InverseNetRate
        this.DeltaZ = dz * InverseNetRate
        this.Radius = 0.2
        this.Height = 0.2
        this.DamageAmount = damage
        this.Hit = this.PlasmaHit
        world.AddMissile(this)
        return this
    }
    Cleanup() {
        super.Cleanup()
        PlaySound("plasma-impact")
        new PlasmaExplosion(this.World, this.X, this.Y, this.Z)
    }
}
