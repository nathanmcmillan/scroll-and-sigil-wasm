const BlockSize = 8
const InverseBlockSize = 1.0 / BlockSize

const BlockSlice = BlockSize * BlockSize
const BlockAll = BlockSlice * BlockSize
const BlockMesh = new RenderCopy(3, 3, 2, BlockAll * 6 * 4, BlockAll * 6 * 6)
const BlockMeshAmbient = new Array(BlockAll)
for (let i = 0; i < BlockAll; i++) {
    BlockMeshAmbient[i] = new Array(6)
    for (let j = 0; j < 6; j++) {
        BlockMeshAmbient[i][j] = new Uint8Array(4)
    }
}
const BlockColorDim = BlockSize + 1
const BlockColorSlice = BlockColorDim * BlockColorDim
const BlockMeshColor = new Array(BlockColorDim * BlockColorSlice)
for (let i = 0; i < BlockMeshColor.length; i++) {
    BlockMeshColor[i] = new Uint8Array(3)
}
const SliceX = [2, 1, 0, 2, 1, 0]
const SliceY = [0, 2, 1, 0, 2, 1]
const SliceZ = [1, 0, 2, 1, 0, 2]
const SliceTowards = [1, 1, 1, -1, -1, -1]
const Slice = new Array(3)
const SliceTemp = new Array(3)

class Block {
    constructor(px, py, pz) {
        this.x = px
        this.y = py
        this.z = pz
        this.mesh
        this.visibility = new Uint8Array(36)
        this.beginSide = new Array(6)
        this.countSide = new Array(6)
        this.thingCount = 0
        this.itemCount = 0
        this.missileCount = 0
        this.particleCount = 0
        this.lightCount = 0
        this.things = []
        this.items = []
        this.missiles = []
        this.particles = []
        this.lights = []
        this.tiles = []
        for (let t = 0; t < BlockAll; t++)
            this.tiles[t] = new Tile()
    }
    Save() {
        let data = "{t["
        if (this.NotEmpty()) {
            for (let i = 0; i < BlockAll; i++) {
                data += this.tiles[i].type
                data += ","
            }
        }
        data += "],c["
        for (let i = 0; i < this.lightCount; i++) {
            data += this.lights[i].Save()
            data += ","
        }
        data += "]}"
        return data
    }
    NotEmpty() {
        for (let i = 0; i < BlockAll; i++) {
            if (this.tiles[i].type !== TileNone) {
                return true
            }
        }
        return false
    }
    GetTilePointerUnsafe(x, y, z) {
        return this.tiles[x + y * BlockSize + z * BlockSlice]
    }
    GetTileTypeUnsafe(x, y, z) {
        return this.tiles[x + y * BlockSize + z * BlockSlice].type
    }
    AddThing(thing) {
        this.things[this.thingCount] = thing
        this.thingCount++
    }
    AddItem(item) {
        this.items[this.itemCount] = item
        this.itemCount++
    }
    AddMissile(missile) {
        this.missiles[this.missileCount] = missile
        this.missileCount++
    }
    AddParticle(particle) {
        this.particles[this.particleCount] = particle
        this.particleCount++
    }
    RemoveThing(thing) {
        let len = this.thingCount
        for (let i = 0; i < len; i++) {
            if (this.things[i] === thing) {
                this.things[i] = this.things[len - 1]
                this.things[len - 1] = null
                this.thingCount--
                return
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
                break
            }
        }
    }
    AddLight(light) {
        this.lights[this.lightCount] = light
        this.lightCount++
    }
    RemoveLight(light) {
        for (let i = 0; i < this.lightCount; i++) {
            if (this.lights[i] === light) {
                for (let j = i; j < this.lightCount - 1; j++)
                    this.lights[j] = this.lights[j + 1]
                this.lightCount--
                return
            }
        }
    }
    AmbientMesh(world) {
        for (let bz = 0; bz < BlockSize; bz++) {
            for (let by = 0; by < BlockSize; by++) {
                for (let bx = 0; bx < BlockSize; bx++) {
                    let index = bx + by * BlockSize + bz * BlockSlice
                    if (this.tiles[index].type === TileNone)
                        continue

                    let ao_mmz = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by - 1, bz)]
                    let ao_mmm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by - 1, bz - 1)]
                    let ao_mmp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by - 1, bz + 1)]
                    let ao_mzp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by, bz + 1)]
                    let ao_mzm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by, bz - 1)]
                    let ao_mpz = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by + 1, bz)]
                    let ao_mpp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by + 1, bz + 1)]
                    let ao_mpm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx - 1, by + 1, bz - 1)]
                    let ao_zpp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx, by + 1, bz + 1)]
                    let ao_zmp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx, by - 1, bz + 1)]
                    let ao_zpm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx, by + 1, bz - 1)]
                    let ao_zmm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx, by - 1, bz - 1)]
                    let ao_ppz = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by + 1, bz)]
                    let ao_pmz = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by - 1, bz)]
                    let ao_pzp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by, bz + 1)]
                    let ao_pzm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by, bz - 1)]
                    let ao_pmm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by - 1, bz - 1)]
                    let ao_ppm = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by + 1, bz - 1)]
                    let ao_ppp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by + 1, bz + 1)]
                    let ao_pmp = TileClosed[world.GetTileType(this.x, this.y, this.z, bx + 1, by - 1, bz + 1)]

                    BlockMeshAmbient[index][WorldPositiveX][0] = Tile.Ambient(ao_pmz, ao_pzm, ao_pmm)
                    BlockMeshAmbient[index][WorldPositiveX][1] = Tile.Ambient(ao_ppz, ao_pzm, ao_ppm)
                    BlockMeshAmbient[index][WorldPositiveX][2] = Tile.Ambient(ao_ppz, ao_pzp, ao_ppp)
                    BlockMeshAmbient[index][WorldPositiveX][3] = Tile.Ambient(ao_pmz, ao_pzp, ao_pmp)

                    BlockMeshAmbient[index][WorldNegativeX][0] = Tile.Ambient(ao_mmz, ao_mzm, ao_mmm)
                    BlockMeshAmbient[index][WorldNegativeX][1] = Tile.Ambient(ao_mmz, ao_mzp, ao_mmp)
                    BlockMeshAmbient[index][WorldNegativeX][2] = Tile.Ambient(ao_mpz, ao_mzp, ao_mpp)
                    BlockMeshAmbient[index][WorldNegativeX][3] = Tile.Ambient(ao_mpz, ao_mzm, ao_mpm)

                    BlockMeshAmbient[index][WorldPositiveY][0] = Tile.Ambient(ao_mpz, ao_zpm, ao_mpm)
                    BlockMeshAmbient[index][WorldPositiveY][1] = Tile.Ambient(ao_mpz, ao_zpp, ao_mpp)
                    BlockMeshAmbient[index][WorldPositiveY][2] = Tile.Ambient(ao_ppz, ao_zpp, ao_ppp)
                    BlockMeshAmbient[index][WorldPositiveY][3] = Tile.Ambient(ao_ppz, ao_zpm, ao_ppm)

                    BlockMeshAmbient[index][WorldNegativeY][0] = Tile.Ambient(ao_mmz, ao_zmm, ao_mmm)
                    BlockMeshAmbient[index][WorldNegativeY][1] = Tile.Ambient(ao_pmz, ao_zmm, ao_pmm)
                    BlockMeshAmbient[index][WorldNegativeY][2] = Tile.Ambient(ao_pmz, ao_zmp, ao_pmp)
                    BlockMeshAmbient[index][WorldNegativeY][3] = Tile.Ambient(ao_mmz, ao_zmp, ao_mmp)

                    BlockMeshAmbient[index][WorldPositiveZ][0] = Tile.Ambient(ao_pzp, ao_zmp, ao_pmp)
                    BlockMeshAmbient[index][WorldPositiveZ][1] = Tile.Ambient(ao_pzp, ao_zpp, ao_ppp)
                    BlockMeshAmbient[index][WorldPositiveZ][2] = Tile.Ambient(ao_mzp, ao_zpp, ao_mpp)
                    BlockMeshAmbient[index][WorldPositiveZ][3] = Tile.Ambient(ao_mzp, ao_zmp, ao_mmp)

                    BlockMeshAmbient[index][WorldNegativeZ][0] = Tile.Ambient(ao_mzm, ao_zmm, ao_mmm)
                    BlockMeshAmbient[index][WorldNegativeZ][1] = Tile.Ambient(ao_mzm, ao_zpm, ao_mpm)
                    BlockMeshAmbient[index][WorldNegativeZ][2] = Tile.Ambient(ao_pzm, ao_zpm, ao_ppm)
                    BlockMeshAmbient[index][WorldNegativeZ][3] = Tile.Ambient(ao_pzm, ao_zmm, ao_pmm)
                }
            }
        }
    }
    ColorMesh(world) {
        for (let bz = 0; bz < BlockColorDim; bz++) {
            for (let by = 0; by < BlockColorDim; by++) {
                for (let bx = 0; bx < BlockColorDim; bx++) {
                    let color = [0, 0, 0, 0]

                    let block_zzz = world.GetTilePointer(this.x, this.y, this.z, bx, by, bz)
                    let block_mzz = world.GetTilePointer(this.x, this.y, this.z, bx - 1, by, bz)
                    let block_mzm = world.GetTilePointer(this.x, this.y, this.z, bx - 1, by, bz - 1)
                    let block_zzm = world.GetTilePointer(this.x, this.y, this.z, bx, by, bz - 1)
                    let block_zmz = world.GetTilePointer(this.x, this.y, this.z, bx, by - 1, bz)
                    let block_mmz = world.GetTilePointer(this.x, this.y, this.z, bx - 1, by - 1, bz)
                    let block_mmm = world.GetTilePointer(this.x, this.y, this.z, bx - 1, by - 1, bz - 1)
                    let block_zmm = world.GetTilePointer(this.x, this.y, this.z, bx, by - 1, bz - 1)

                    if (block_zzz === null || TileClosed[block_zzz.type]) {
                        this.DetermineLight(block_mzz, color)
                        this.DetermineLight(block_zmz, color)
                        this.DetermineLight(block_zzm, color)
                    }
                    if (block_mzz === null || TileClosed[block_mzz.type]) {
                        this.DetermineLight(block_zzz, color)
                        this.DetermineLight(block_zmz, color)
                        this.DetermineLight(block_zzm, color)
                    }
                    if (block_mzm === null || TileClosed[block_mzm.type]) {
                        this.DetermineLight(block_mzz, color)
                        this.DetermineLight(block_zzm, color)
                        this.DetermineLight(block_mmm, color)
                    }
                    if (block_zzm === null || TileClosed[block_zzm.type]) {
                        this.DetermineLight(block_zzz, color)
                        this.DetermineLight(block_mzm, color)
                        this.DetermineLight(block_zmm, color)
                    }
                    if (block_zmz === null || TileClosed[block_zmz.type]) {
                        this.DetermineLight(block_zzz, color)
                        this.DetermineLight(block_mmz, color)
                        this.DetermineLight(block_zmm, color)
                    }
                    if (block_mmz === null || TileClosed[block_mmz.type]) {
                        this.DetermineLight(block_mzz, color)
                        this.DetermineLight(block_mmm, color)
                        this.DetermineLight(block_zmz, color)
                    }
                    if (block_mmm === null || TileClosed[block_mmm.type]) {
                        this.DetermineLight(block_mzm, color)
                        this.DetermineLight(block_zmm, color)
                        this.DetermineLight(block_mmz, color)
                    }
                    if (block_zmm === null || TileClosed[block_zmm.type]) {
                        this.DetermineLight(block_zzm, color)
                        this.DetermineLight(block_zmz, color)
                        this.DetermineLight(block_mmm, color)
                    }

                    let index = bx + by * BlockColorDim + bz * BlockColorSlice
                    if (color[3] > 0) {
                        BlockMeshColor[index][0] = color[0] / color[3]
                        BlockMeshColor[index][1] = color[1] / color[3]
                        BlockMeshColor[index][2] = color[2] / color[3]
                    } else {
                        BlockMeshColor[index][0] = 255
                        BlockMeshColor[index][1] = 255
                        BlockMeshColor[index][2] = 255
                    }
                }
            }
        }
    }
    DetermineLight(tile, color) {
        if (tile === null)
            return
        if (!TileClosed[tile.type]) {
            color[0] += tile.red
            color[1] += tile.green
            color[2] += tile.blue
            color[3]++
        }
    }
    LightOfSide(xs, ys, zs, side) {
        switch (side) {
            case WorldPositiveX:
                return [
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + (zs + 1) * BlockColorSlice]
                ]
            case WorldNegativeX:
                return [
                    BlockMeshColor[xs + ys * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + ys * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + zs * BlockColorSlice]
                ]
            case WorldPositiveY:
                return [
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + zs * BlockColorSlice]
                ]
            case WorldNegativeY:
                return [
                    BlockMeshColor[xs + ys * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + ys * BlockColorDim + (zs + 1) * BlockColorSlice]
                ]
            case WorldPositiveZ:
                return [
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + (zs + 1) * BlockColorSlice],
                    BlockMeshColor[xs + ys * BlockColorDim + (zs + 1) * BlockColorSlice]
                ]
            default:
                return [
                    BlockMeshColor[xs + ys * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + (ys + 1) * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + (ys + 1) * BlockColorDim + zs * BlockColorSlice],
                    BlockMeshColor[xs + 1 + ys * BlockColorDim + zs * BlockColorSlice]
                ]
        }
    }
    BuildMesh(world) {
        this.AmbientMesh(world)
        this.ColorMesh(world)
        BlockMesh.Zero()
        for (let side = 0; side < 6; side++) {
            let mesh_begin_index = BlockMesh.index_pos
            let pointerX = SliceX[side]
            let pointerY = SliceY[side]
            let pointerZ = SliceZ[side]
            let toward = SliceTowards[side]
            for (Slice[2] = 0; Slice[2] < BlockSize; Slice[2]++) {
                for (Slice[1] = 0; Slice[1] < BlockSize; Slice[1]++) {
                    for (Slice[0] = 0; Slice[0] < BlockSize; Slice[0]++) {
                        let type = this.GetTileTypeUnsafe(Slice[pointerX], Slice[pointerY], Slice[pointerZ])
                        if (type === TileNone)
                            continue
                        SliceTemp[0] = Slice[0]
                        SliceTemp[1] = Slice[1]
                        SliceTemp[2] = Slice[2] + toward
                        if (TileClosed[world.GetTileType(this.x, this.y, this.z, SliceTemp[pointerX], SliceTemp[pointerY], SliceTemp[pointerZ])])
                            continue
                        let xs = Slice[pointerX]
                        let ys = Slice[pointerY]
                        let zs = Slice[pointerZ]
                        let index = xs + ys * BlockSize + zs * BlockSlice

                        let texture = TileTexture[type]
                        let bx = xs + BlockSize * this.x
                        let by = ys + BlockSize * this.y
                        let bz = zs + BlockSize * this.z

                        let light = this.LightOfSide(xs, ys, zs, side)
                        let ambient = BlockMeshAmbient[index][side]

                        let rgb_a = Light.Colorize(light[0], ambient[0])
                        let rgb_b = Light.Colorize(light[1], ambient[1])
                        let rgb_c = Light.Colorize(light[2], ambient[2])
                        let rgb_d = Light.Colorize(light[3], ambient[3])

                        RenderTile.Side(BlockMesh, side, bx, by, bz, texture, rgb_a, rgb_b, rgb_c, rgb_d)
                    }
                }
            }
            this.beginSide[side] = mesh_begin_index * 4
            this.countSide[side] = BlockMesh.index_pos - mesh_begin_index
        }
        this.mesh = RenderBuffer.InitCopy(world.gl, BlockMesh)
    }
    RenderThings(spriteSet, spriteBuffer, camX, camZ, camAngle) {
        for (let i = 0; i < this.thingCount; i++) {
            let thing = this.things[i]
            if (spriteSet.has(thing)) continue
            spriteSet.add(thing)
            thing.Render(spriteBuffer, camX, camZ, camAngle)
        }
        for (let i = 0; i < this.itemCount; i++) {
            let item = this.items[i]
            if (spriteSet.has(item)) continue
            spriteSet.add(item)
            item.Render(spriteBuffer, camX, camZ, camAngle)
        }
        for (let i = 0; i < this.missileCount; i++) {
            let missile = this.missiles[i]
            if (spriteSet.has(missile)) continue
            spriteSet.add(missile)
            missile.Render(spriteBuffer, camX, camZ, camAngle)
        }
        for (let i = 0; i < this.particleCount; i++) {
            let particle = this.particles[i]
            if (spriteSet.has(particle)) continue
            spriteSet.add(particle)
            particle.Render(spriteBuffer, camX, camZ)
        }
    }
}
