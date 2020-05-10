class Scenery extends Thing {
    constructor(world, nid, x, y, z) {
        super()
        this.World = world
        this.UID = TreeUID
        this.SID = "scenery"
        this.NID = nid
        this.X = x
        this.Y = y
        this.Z = z
        this.OldX = x
        this.OldY = y
        this.OldZ = z
        this.Radius = 0.4
        this.Height = 1.0
        this.Health = 1
        world.AddThing(this)
        this.BlockBorders()
        this.AddToBlocks()
    }
    Save() {
        let data = "{u:" + this.UID
        data += ",x:" + this.X
        data += ",y:" + this.Y
        data += ",z:" + this.Z
        data += "}"
        return data
    }
    Update() {}
    Render(spriteBuffer, camX, camZ, camAngle) {
        let sin = camX - this.X
        let cos = camZ - this.Z
        let length = Math.sqrt(sin * sin + cos * cos)
        sin /= length
        cos /= length
        Render3.Sprite(spriteBuffer.get(this.SID), this.X, this.Y, this.Z, sin, cos, this.Sprite)
    }
}

class Tree extends Scenery {
    constructor(world, nid, x, y, z) {
        super(world, nid, x, y, z)
        this.Sprite = SpriteData[this.SID]["dead-tree"]
    }
}
