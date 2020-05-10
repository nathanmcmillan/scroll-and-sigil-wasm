const NetUpdateRate = 50
const InverseNetRate = 16.67 / NetUpdateRate

const WorldPositiveX = 0
const WorldPositiveY = 1
const WorldPositiveZ = 2
const WorldNegativeX = 3
const WorldNegativeY = 4
const WorldNegativeZ = 5

const BroadcastNew = 0
const BroadcastDelete = 1
const BroadcastChat = 2

class World {
    constructor(g, gl) {
        this.g = g
        this.gl = gl
        this.width
        this.height
        this.length
        this.tileWidth
        this.tileHeight
        this.tileLength
        this.slice
        this.all
        this.blocks
        this.viewable
        this.spriteSet
        this.spriteBuffer = new Map()
        this.spriteCount
        this.thingCount
        this.itemCount
        this.missileCount
        this.particleCount
        this.things
        this.items
        this.missiles
        this.particles
        this.netLookup
        this.occluder = new Occluder()
        this.PID
    }
    Reset() {
        this.blocks = []
        this.viewable = []
        this.spriteSet = new Set()
        this.spriteCount = new Map()
        this.thingCount = 0
        this.itemCount = 0
        this.missileCount = 0
        this.particleCount = 0
        this.things = []
        this.items = []
        this.missiles = []
        this.particles = []
        this.netLookup = new Map()
    }
    Load(binary) {
        this.Reset()

        let dat = new DataView(binary)
        let dex = 0

        this.PID = dat.getUint16(dex, true)
        dex += 2

        this.width = dat.getUint16(dex, true)
        dex += 2
        this.height = dat.getUint16(dex, true)
        dex += 2
        this.length = dat.getUint16(dex, true)
        dex += 2
        this.slice = this.width * this.height
        this.all = this.slice * this.length

        this.tileWidth = this.width * BlockSize
        this.tileHeight = this.height * BlockSize
        this.tileLength = this.length * BlockSize

        let bx = 0
        let by = 0
        let bz = 0
        for (let b = 0; b < this.all; b++) {
            this.blocks[b] = new Block(bx, by, bz)
            bx++
            if (bx === this.width) {
                bx = 0
                by++
                if (by === this.height) {
                    by = 0
                    bz++
                }
            }
        }

        for (let b = 0; b < this.all; b++) {
            let block = this.blocks[b]
            let notEmpty = dat.getUint8(dex, true)
            dex += 1

            if (notEmpty) {
                for (let t = 0; t < BlockAll; t++) {
                    let tileType = dat.getUint8(dex, true)
                    dex += 1
                    block.tiles[t].type = tileType
                }
            }

            let lightCount = dat.getUint8(dex, true)
            dex += 1
            for (let t = 0; t < lightCount; t++) {
                let x = dat.getUint8(dex, true)
                dex += 1
                let y = dat.getUint8(dex, true)
                dex += 1
                let z = dat.getUint8(dex, true)
                dex += 1
                let rgb = dat.getInt32(dex, true)
                dex += 4
                block.AddLight(new Light(x, y, z, rgb))
            }
        }

        let thingCount = dat.getUint16(dex, true)
        dex += 2
        for (let t = 0; t < thingCount; t++) {
            let uid = dat.getUint16(dex, true)
            dex += 2
            let nid = dat.getUint16(dex, true)
            dex += 2
            let x = dat.getFloat32(dex, true)
            dex += 4
            let y = dat.getFloat32(dex, true)
            dex += 4
            let z = dat.getFloat32(dex, true)
            dex += 4
            switch (uid) {
                case HumanUID:
                    {
                        let angle = dat.getFloat32(dex, true)
                        dex += 4
                        let health = dat.getUint16(dex, true)
                        dex += 2
                        let status = dat.getUint8(dex, true)
                        dex += 1
                        if (nid === this.PID) new You(this, nid, x, y, z, angle, health, status)
                        else new Human(this, nid, x, y, z, angle, health, status)
                    }
                    break
                case BaronUID:
                    {
                        let direction = dat.getUint8(dex, true)
                        dex += 1
                        let health = dat.getUint16(dex, true)
                        dex += 2
                        let status = dat.getUint8(dex, true)
                        dex += 1
                        new Baron(this, nid, x, y, z, direction, health, status)
                    }
                    break
                case TreeUID:
                    new Tree(this, nid, x, y, z)
                    break
            }
        }

        let itemCount = dat.getUint16(dex, true)
        dex += 2
        for (let t = 0; t < itemCount; t++) {
            let uid = dat.getUint16(dex, true)
            dex += 2
            let nid = dat.getUint16(dex, true)
            dex += 2
            let x = dat.getFloat32(dex, true)
            dex += 4
            let y = dat.getFloat32(dex, true)
            dex += 4
            let z = dat.getFloat32(dex, true)
            dex += 4
            switch (uid) {
                case MedkitUID:
                    new Medkit(this, nid, x, y, z)
                    break
            }
        }

        let missileCount = dat.getUint16(dex, true)
        dex += 2
        for (let t = 0; t < missileCount; t++) {
            let uid = dat.getUint16(dex, true)
            dex += 2
            let nid = dat.getUint16(dex, true)
            dex += 2
            let x = dat.getFloat32(dex, true)
            dex += 4
            let y = dat.getFloat32(dex, true)
            dex += 4
            let z = dat.getFloat32(dex, true)
            dex += 4
            let dx = dat.getFloat32(dex, true)
            dex += 4
            let dy = dat.getFloat32(dex, true)
            dex += 4
            let dz = dat.getFloat32(dex, true)
            dex += 4
            switch (uid) {
                case PlasmaUID:
                    {
                        let damage = dat.getUint16(dex, true)
                        dex += 2
                        new Plasma(this, nid, damage, x, y, z, dx, dy, dz)
                    }
                    break
            }
        }

        this.Build()
    }
    Build() {
        for (let i = 0; i < this.all; i++) {
            let block = this.blocks[i]
            for (let j = 0; j < block.lightCount; j++)
                Light.Add(this, block, block.lights[j])
            Occluder.SetBlockVisible(block)
        }
        for (let i = 0; i < this.all; i++)
            this.blocks[i].BuildMesh(this)
    }
    FindBlock(x, y, z) {
        let gx = Math.floor(x)
        let gy = Math.floor(y)
        let gz = Math.floor(z)
        let bx = Math.floor(gx * InverseBlockSize)
        let by = Math.floor(gy * InverseBlockSize)
        let bz = Math.floor(gz * InverseBlockSize)
        let tx = gx - bx * BlockSize
        let ty = gy - by * BlockSize
        let tz = gz - bz * BlockSize
        let block = this.blocks[bx + by * this.width + bz * this.slice]
        return block.tiles[tx + ty * BlockSize + tz * BlockSlice].type
    }
    GetTilePointer(cx, cy, cz, tx, ty, tz) {
        while (tx < 0) {
            tx += BlockSize
            cx--
        }
        while (tx >= BlockSize) {
            tx -= BlockSize
            cx++
        }
        while (ty < 0) {
            ty += BlockSize
            cy--
        }
        while (ty >= BlockSize) {
            ty -= BlockSize
            cy++
        }
        while (tz < 0) {
            tz += BlockSize
            cz--
        }
        while (tz >= BlockSize) {
            tz -= BlockSize
            cz++
        }
        let block = this.GetBlock(cx, cy, cz)
        if (block === null)
            return null
        return block.GetTilePointerUnsafe(tx, ty, tz)
    }
    GetTileType(cx, cy, cz, tx, ty, tz) {
        while (tx < 0) {
            tx += BlockSize
            cx--
        }
        while (tx >= BlockSize) {
            tx -= BlockSize
            cx++
        }
        while (ty < 0) {
            ty += BlockSize
            cy--
        }
        while (ty >= BlockSize) {
            ty -= BlockSize
            cy++
        }
        while (tz < 0) {
            tz += BlockSize
            cz--
        }
        while (tz >= BlockSize) {
            tz -= BlockSize
            cz++
        }
        let block = this.GetBlock(cx, cy, cz)
        if (block === null) {
            return TileNone
        }
        return block.GetTileTypeUnsafe(tx, ty, tz)
    }
    GetBlock(x, y, z) {
        if (x < 0 || x >= this.width)
            return null
        if (y < 0 || y >= this.height)
            return null
        if (z < 0 || z >= this.length)
            return null
        return this.blocks[x + y * this.width + z * this.slice]
    }
    AddThing(thing) {
        this.things[this.thingCount] = thing
        this.thingCount++
        this.netLookup.set(thing.NID, thing)

        let count = this.spriteCount.get(thing.SID)
        if (count) {
            this.spriteCount.set(thing.SID, count + 1)
            let buffer = this.spriteBuffer.get(thing.SID)
            if ((count + 2) * 16 > buffer.vertices.length) {
                this.spriteBuffer.set(thing.SID, RenderBuffer.Expand(this.gl, buffer))
            }
        } else {
            this.spriteCount.set(thing.SID, 1)
            this.spriteBuffer.set(thing.SID, RenderBuffer.Init(this.gl, 3, 0, 2, 40, 60))
        }
    }
    AddItem(item) {
        this.items[this.itemCount] = item
        this.itemCount++
        this.netLookup.set(item.NID, item)

        let count = this.spriteCount.get(item.SID)
        if (count) {
            this.spriteCount.set(item.SID, count + 1)
            let buffer = this.spriteBuffer.get(item.SID)
            if ((count + 2) * 16 > buffer.vertices.length) {
                this.spriteBuffer.set(item.SID, RenderBuffer.Expand(this.gl, buffer))
            }
        } else {
            this.spriteCount.set(item.SID, 1)
            this.spriteBuffer.set(item.SID, RenderBuffer.Init(this.gl, 3, 0, 2, 40, 60))
        }
    }
    AddMissile(missile) {
        this.missiles[this.missileCount] = missile
        this.missileCount++
        this.netLookup.set(missile.NID, missile)

        let count = this.spriteCount.get(missile.SID)
        if (count) {
            this.spriteCount.set(missile.SID, count + 1)
            let buffer = this.spriteBuffer.get(missile.SID)
            if ((count + 2) * 16 > buffer.vertices.length) {
                this.spriteBuffer.set(missile.SID, RenderBuffer.Expand(this.gl, buffer))
            }
        } else {
            this.spriteCount.set(missile.SID, 1)
            this.spriteBuffer.set(missile.SID, RenderBuffer.Init(this.gl, 3, 0, 2, 40, 60))
        }
    }
    AddParticle(particle) {
        this.particles[this.particleCount] = particle
        this.particleCount++

        let count = this.spriteCount.get(particle.SID)
        if (count) {
            this.spriteCount.set(particle.SID, count + 1)
            let buffer = this.spriteBuffer.get(particle.SID)
            if ((count + 2) * 16 > buffer.vertices.length) {
                this.spriteBuffer.set(particle.SID, RenderBuffer.Expand(this.gl, buffer))
            }
        } else {
            this.spriteCount.set(particle.SID, 1)
            this.spriteBuffer.set(particle.SID, RenderBuffer.Init(this.gl, 3, 0, 2, 120, 180))
        }
    }
    RemoveThing(thing) {
        let len = this.thingCount
        for (let i = 0; i < len; i++) {
            if (this.things[i] === thing) {
                this.things[i] = this.things[len - 1]
                this.things[len - 1] = null
                this.thingCount--
                this.spriteCount.set(thing.SID, this.spriteCount.get(thing.SID) - 1)
                this.netLookup.delete(thing.NID)
                break
            }
        }
    }
    RemoveItem(item) {
        let len = this.itemCount
        for (let i = 0; i < len; i++) {
            if (this.items[i] === item) {
                this.items[i] = this.items[len - 1]
                this.items[len - 1] = null
                this.itemCount--
                this.spriteCount.set(item.SID, this.spriteCount.get(item.SID) - 1)
                this.netLookup.delete(item.NID)
                break
            }
        }
    }
    RemoveMissile(missile) {
        let len = this.missileCount
        for (let i = 0; i < len; i++) {
            if (this.missiles[i] === missile) {
                this.missiles[i] = this.missiles[len - 1]
                this.missiles[len - 1] = null
                this.missileCount--
                this.spriteCount.set(missile.SID, this.spriteCount.get(missile.SID) - 1)
                this.netLookup.delete(missile.NID)
                break
            }
        }
    }
    RemoveParticle(particle) {
        let len = this.particleCount
        for (let i = 0; i < len; i++) {
            if (this.particles[i] === particle) {
                this.particles[i] = this.particles[len - 1]
                this.particles[len - 1] = null
                this.particleCount--
                this.spriteCount.set(particle.SID, this.spriteCount.get(particle.SID) - 1)
                break
            }
        }
    }
    update() {
        let len = this.thingCount
        for (let i = 0; i < len; i++)
            this.things[i].Update()
        len = this.missileCount
        for (let i = 0; i < len; i++) {
            if (this.missiles[i].Update()) {
                this.missiles[i] = this.missiles[len - 1]
                this.missiles[len - 1] = null
                this.missileCount--
                len--
                i--
            }
        }
        len = this.particleCount
        for (let i = 0; i < this.particleCount; i++) {
            if (this.particles[i].Update()) {
                this.particles[i] = this.particles[len - 1]
                this.particles[len - 1] = null
                this.particleCount--
                len--
                i--

            }
        }
    }
    render(g, x, y, z, camX, camZ, camAngle) {
        let gl = this.gl
        let spriteSet = this.spriteSet
        let spriteBuffer = this.spriteBuffer

        this.occluder.PrepareFrustum(g)
        this.occluder.Occlude(this, x, y, z)

        spriteSet.clear()

        for (let buffer of spriteBuffer.values())
            buffer.Zero()

        g.SetProgram(gl, "texture-color3d")
        g.UpdateMvp(gl)
        g.SetTexture(gl, "tiles")

        for (let i = 0; i < OcclusionViewNum; i++) {
            let block = this.viewable[i]

            block.RenderThings(spriteSet, spriteBuffer, camX, camZ, camAngle)

            let mesh = block.mesh
            if (mesh.vertexPos === 0)
                continue

            RenderSystem.BindVao(gl, mesh)

            if (x == block.x) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveX], block.countSide[WorldPositiveX])
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeX], block.countSide[WorldNegativeX])
            } else if (x > block.x) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveX], block.countSide[WorldPositiveX])
            } else {
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeX], block.countSide[WorldNegativeX])
            }

            if (y == block.y) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveY], block.countSide[WorldPositiveY])
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeY], block.countSide[WorldNegativeY])
            } else if (y > block.y) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveY], block.countSide[WorldPositiveY])
            } else {
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeY], block.countSide[WorldNegativeY])
            }

            if (z == block.z) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveZ], block.countSide[WorldPositiveZ])
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeZ], block.countSide[WorldNegativeZ])
            } else if (z > block.z) {
                RenderSystem.DrawRange(gl, block.beginSide[WorldPositiveZ], block.countSide[WorldPositiveZ])
            } else {
                RenderSystem.DrawRange(gl, block.beginSide[WorldNegativeZ], block.countSide[WorldNegativeZ])
            }
        }

        g.SetProgram(gl, "texture3d")
        g.UpdateMvp(gl)
        for (let [name, buffer] of spriteBuffer) {
            if (buffer.vertexPos > 0) {
                g.SetTexture(gl, name)
                RenderSystem.UpdateAndDraw(gl, buffer)
            }
        }
    }
}
