const LightQueueLimit = 30 * 30 * 30
const LightFade = 0.95
const LightQueue = new Array(LightQueueLimit)
for (let i = 0; i < LightQueue.length; i++) {
    LightQueue[i] = new Int32Array(3)
}
let LightBlockX = 0
let LightBlockY = 0
let LightBlockZ = 0
let LightPos = 0
let LightNum = 0

class Light {
    constructor(x, y, z, rgb) {
        this.x = x
        this.y = y
        this.z = z
        this.rgb = rgb
    }
    Save() {
        let data = "{x:" + this.x
        data += ",y:" + this.y
        data += ",z:" + this.z
        data += ",v:" + this.rgb
        data += "}"
        return data
    }
    static Colorize(rgb, ambient) {
        return [
            rgb[0] * ambient / 65025.0,
            rgb[1] * ambient / 65025.0,
            rgb[2] * ambient / 65025.0
        ]
    }
    static Visit(world, bx, by, bz, red, green, blue) {
        let tile = world.GetTilePointer(LightBlockX, LightBlockY, LightBlockZ, bx, by, bz)
        if (tile === null || TileClosed[tile.type])
            return
        if (tile.red >= red || tile.green >= green || tile.blue >= blue)
            return
        tile.red = red
        tile.green = green
        tile.blue = blue

        let queue = LightPos + LightNum
        if (queue >= LightQueueLimit)
            queue -= LightQueueLimit
        LightQueue[queue][0] = bx
        LightQueue[queue][1] = by
        LightQueue[queue][2] = bz
        LightNum++
    }
    static Add(world, block, light) {
        LightBlockX = block.x
        LightBlockY = block.y
        LightBlockZ = block.z

        let color = Render.UnpackRgb(light.rgb)

        let index = light.x + light.y * BlockSize + light.z * BlockSlice
        let tile = block.tiles[index]
        tile.red = color[0]
        tile.green = color[1]
        tile.blue = color[2]

        LightQueue[0][0] = light.x
        LightQueue[0][1] = light.y
        LightQueue[0][2] = light.z
        LightPos = 0
        LightNum = 1

        while (LightNum > 0) {
            let x = LightQueue[LightPos][0]
            let y = LightQueue[LightPos][1]
            let z = LightQueue[LightPos][2]

            LightPos++
            if (LightPos === LightQueueLimit)
                LightPos = 0
            LightNum--

            let node = world.GetTilePointer(LightBlockX, LightBlockY, LightBlockZ, x, y, z)
            if (node === null)
                continue

            let r = Math.floor(node.red * LightFade)
            let g = Math.floor(node.green * LightFade)
            let b = Math.floor(node.blue * LightFade)

            Light.Visit(world, x - 1, y, z, r, g, b)
            Light.Visit(world, x + 1, y, z, r, g, b)
            Light.Visit(world, x, y - 1, z, r, g, b)
            Light.Visit(world, x, y + 1, z, r, g, b)
            Light.Visit(world, x, y, z - 1, r, g, b)
            Light.Visit(world, x, y, z + 1, r, g, b)
        }
    }
    static Remove(world, x, y, z) {
        let cx = Math.floor(x * InverseBlockSize)
        let cy = Math.floor(y * InverseBlockSize)
        let cz = Math.floor(z * InverseBlockSize)
        let bx = x % BlockSize
        let by = y % BlockSize
        let bz = z % BlockSize
        let block = world.blocks[cx + cy * world.width + cz * world.slice]
        for (let i = 0; i < block.lights.length; i++) {
            let light = block.lights[i]
            if (light.x === bx && light.y === by && light.z === bz)
                block.lights.splice(i, 1)
        }
    }
}
