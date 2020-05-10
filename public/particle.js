const PlasmaExplosionAnimation = []

class Particle {
    constructor() {
        this.World = null
        this.SID = "particles"
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
                    this.World.GetBlock(gx, gy, gz).AddParticle(this)
                }
            }
        }
    }
    RemoveFromBlocks() {
        for (let gx = this.MinBX; gx <= this.MaxBX; gx++) {
            for (let gy = this.MinBY; gy <= this.MaxBY; gy++) {
                for (let gz = this.MinBZ; gz <= this.MaxBZ; gz++) {
                    this.World.GetBlock(gx, gy, gz).RemoveParticle(this)
                }
            }
        }
    }
    Collision() {
        let minGX = Math.floor(this.X - this.Radius)
        let minGY = Math.floor(this.Y)
        let minGZ = Math.floor(this.Z - this.Radius)
        let maxGX = Math.floor(this.X + this.Radius)
        let maxGY = Math.floor(this.Y + this.Height)
        let maxGZ = Math.floor(this.Z + this.Radius)
        for (let gx = minGX; gx <= maxGX; gx++) {
            for (let gy = minGY; gy <= maxGY; gy++) {
                for (let gz = minGZ; gz <= maxGZ; gz++) {
                    let bx = Math.floor(gx * InverseBlockSize)
                    let by = Math.floor(gy * InverseBlockSize)
                    let bz = Math.floor(gz * InverseBlockSize)
                    let tx = gx - bx * BlockSize
                    let ty = gy - by * BlockSize
                    let tz = gz - bz * BlockSize
                    let tile = this.World.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tile]) {
                        return true
                    }
                }
            }
        }
        return false
    }
    Update() {}
    Render(spriteBuffer, camX, camZ) {
        let sin = camX - this.X
        let cos = camZ - this.Z
        let length = Math.sqrt(sin * sin + cos * cos)
        sin /= length
        cos /= length
        Render3.Sprite(spriteBuffer.get(this.SID), this.X, this.Y, this.Z, sin, cos, this.Sprite)
    }
}

class PlasmaExplosion extends Particle {
    constructor(world, x, y, z) {
        super()
        this.AnimationMod = 0
        this.AnimationFrame = 0
        this.Animation = PlasmaExplosionAnimation
        this.Sprite = this.Animation[0]
        this.World = world
        this.X = x
        this.Y = y
        this.Z = z
        this.DeltaX = 0
        this.DeltaY = 0
        this.DeltaZ = 0
        this.Radius = 0.4
        this.Height = 1.0
        world.AddParticle(this)
        this.BlockBorders()
        this.AddToBlocks()
        return this
    }
    UpdateAnimation() {
        const PlasmaAnimationRate = 8
        this.AnimationMod++
        if (this.AnimationMod === PlasmaAnimationRate) {
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
    Update() {
        if (this.UpdateAnimation() === AnimationDone) {
            this.RemoveFromBlocks()
            return true
        }
        this.Sprite = this.Animation[this.AnimationFrame]
        return false
    }
}

class Blood extends Particle {
    constructor(world, x, y, z, dx, dy, dz, spriteName) {
        super()
        this.Sprite = SpriteData[this.SID][spriteName]
        this.World = world
        this.X = x
        this.Y = y
        this.Z = z
        this.DeltaX = dx
        this.DeltaY = dy
        this.DeltaZ = dz
        this.Radius = 0.2
        this.Height = 0.2
        world.AddParticle(this)
        this.BlockBorders()
        this.AddToBlocks()
        return this
    }
    Update() {
        this.DeltaX *= 0.95
        this.DeltaY -= 0.01
        this.DeltaZ *= 0.95

        this.X += this.DeltaX
        this.Y += this.DeltaY
        this.Z += this.DeltaZ

        if (this.Collision()) {
            this.RemoveFromBlocks()
            return true
        }

        this.RemoveFromBlocks()
        this.BlockBorders()
        this.AddToBlocks()

        return false
    }
}
