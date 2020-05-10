const OcclusionSliceA = new Int32Array(3)
const OcclusionSliceB = new Int32Array(3)

const FullOcclusion = 0
const PartialOcclusion = 1
const NoOcclusion = 2

const OcclusionQueue = []
const OcclusionGoto = []
const OcclusionQueueFrom = []

let OcclusionViewNum = 0
let OcclusionQueuePos = 0
let OcclusionQueueNum = 0

class Occluder {
    constructor() {
        this.frustum = []
        for (let i = 0; i < 6; i++)
            this.frustum.push(new Float32Array(4))
    }
    static SetBlockVisible(block) {
        for (let sideA = 0; sideA < 6; sideA++) {
            let ax = SliceX[sideA]
            let ay = SliceY[sideA]
            let az = SliceZ[sideA]
            for (let sideB = sideA + 1; sideB < 6; sideB++) {
                let bx = SliceX[sideB]
                let by = SliceY[sideB]
                let bz = SliceZ[sideB]

                if (SliceTowards[sideA] > 0) OcclusionSliceA[2] = BlockSize - 1
                else OcclusionSliceA[2] = 0

                if (SliceTowards[sideB] > 0) OcclusionSliceB[2] = BlockSize - 1
                else OcclusionSliceB[2] = 0

                loop:
                    for (OcclusionSliceA[1] = 0; OcclusionSliceA[1] < BlockSize; OcclusionSliceA[1]++) {
                        for (OcclusionSliceA[0] = 0; OcclusionSliceA[0] < BlockSize; OcclusionSliceA[0]++) {
                            for (OcclusionSliceB[1] = 0; OcclusionSliceB[1] < BlockSize; OcclusionSliceB[1]++) {
                                for (OcclusionSliceB[0] = 0; OcclusionSliceB[0] < BlockSize; OcclusionSliceB[0]++) {
                                    let fromX = OcclusionSliceA[ax] + 0.5
                                    let fromY = OcclusionSliceA[ay] + 0.5
                                    let fromZ = OcclusionSliceA[az] + 0.5
                                    let toX = OcclusionSliceB[bx] + 0.5
                                    let toY = OcclusionSliceB[by] + 0.5
                                    let toZ = OcclusionSliceB[bz] + 0.5
                                    if (Cast.Block(block, fromX, fromY, fromZ, toX, toY, toZ)) {
                                        block.visibility[sideA * 6 + sideB] = 1
                                        break loop
                                    }
                                }
                            }
                        }
                    }
            }
        }
    }
    Visit(world, from, B, to) {
        let x = B.x
        let y = B.y
        let z = B.z
        switch (to) {
            case WorldPositiveX:
                x++
                if (x === world.width) return
                break
            case WorldNegativeX:
                x--
                if (x === -1) return
                break
            case WorldPositiveY:
                y++
                if (y === world.height) return
                break
            case WorldNegativeY:
                y--
                if (y === -1) return
                break
            case WorldPositiveZ:
                z++
                if (z === world.length) return
                break
            case WorldNegativeZ:
                z--
                if (z === -1) return
                break
        }
        let index = x + y * world.width + z * world.slice
        if (OcclusionGoto[index] === false)
            return
        if (from >= 0) {
            switch (from) {
                case WorldPositiveX:
                    from = WorldNegativeX
                    break
                case WorldNegativeX:
                    from = WorldPositiveX
                    break
                case WorldPositiveY:
                    from = WorldNegativeY
                    break
                case WorldNegativeY:
                    from = WorldPositiveY
                    break
                case WorldPositiveZ:
                    from = WorldNegativeZ
                    break
                case WorldNegativeZ:
                    from = WorldPositiveZ
                    break
            }
            let sideA, sideB
            if (from < to) {
                sideA = from
                sideB = to
            } else {
                sideA = to
                sideB = from
            }
            if (B.visibility[sideA * 6 + sideB] === 0)
                return
        }
        OcclusionGoto[index] = false
        let C = world.blocks[index]
        let pos_cx = C.x * BlockSize
        let pos_cy = C.y * BlockSize
        let pos_cz = C.z * BlockSize
        let box = this.InBox(
            pos_cx + BlockSize, pos_cy + BlockSize, pos_cz + BlockSize,
            pos_cx, pos_cy, pos_cz)
        if (box === NoOcclusion)
            return

        let queue = OcclusionQueuePos + OcclusionQueueNum
        if (queue >= world.all)
            queue -= world.all

        OcclusionQueue[queue] = C
        OcclusionQueueFrom[queue] = to
        OcclusionQueueNum++
    }
    InBox(pos_x, pos_y, pos_z, neg_x, neg_y, neg_z) {
        let pvx, pvy, pvz
        let nvx, nvy, nvz
        let result = FullOcclusion
        for (let i = 0; i < 6; i++) {
            let plane = this.frustum[i]
            if (plane[0] > 0) {
                pvx = pos_x
                nvx = neg_x
            } else {
                pvx = neg_x
                nvx = pos_x
            }
            if (plane[1] > 0) {
                pvy = pos_y
                nvy = neg_y
            } else {
                pvy = neg_y
                nvy = pos_y
            }
            if (plane[2] > 0) {
                pvz = pos_z
                nvz = neg_z
            } else {
                pvz = neg_z
                nvz = pos_z
            }
            if (pvx * plane[0] + pvy * plane[1] + pvz * plane[2] + plane[3] < 0)
                return NoOcclusion
            if (nvx * plane[0] + nvy * plane[1] + nvz * plane[2] + plane[3] < 0)
                result = PartialOcclusion
        }
        return result
    }
    PrepareFrustum(g) {
        // left
        this.frustum[0][0] = g.mvp[3] + g.mvp[0]
        this.frustum[0][1] = g.mvp[7] + g.mvp[4]
        this.frustum[0][2] = g.mvp[11] + g.mvp[8]
        this.frustum[0][3] = g.mvp[15] + g.mvp[12]
        this.NormalizePlane(0)

        // right
        this.frustum[1][0] = g.mvp[3] - g.mvp[0]
        this.frustum[1][1] = g.mvp[7] - g.mvp[4]
        this.frustum[1][2] = g.mvp[11] - g.mvp[8]
        this.frustum[1][3] = g.mvp[15] - g.mvp[12]
        this.NormalizePlane(1)

        // top
        this.frustum[2][0] = g.mvp[3] - g.mvp[1]
        this.frustum[2][1] = g.mvp[7] - g.mvp[5]
        this.frustum[2][2] = g.mvp[11] - g.mvp[9]
        this.frustum[2][3] = g.mvp[15] - g.mvp[13]
        this.NormalizePlane(2)

        // bottom
        this.frustum[3][0] = g.mvp[3] + g.mvp[1]
        this.frustum[3][1] = g.mvp[7] + g.mvp[5]
        this.frustum[3][2] = g.mvp[11] + g.mvp[9]
        this.frustum[3][3] = g.mvp[15] + g.mvp[13]
        this.NormalizePlane(3)

        // near
        this.frustum[4][0] = g.mvp[3] + g.mvp[2]
        this.frustum[4][1] = g.mvp[7] + g.mvp[6]
        this.frustum[4][2] = g.mvp[11] + g.mvp[10]
        this.frustum[4][3] = g.mvp[15] + g.mvp[14]
        this.NormalizePlane(4)

        // far
        this.frustum[5][0] = g.mvp[3] - g.mvp[2]
        this.frustum[5][1] = g.mvp[7] - g.mvp[6]
        this.frustum[5][2] = g.mvp[11] - g.mvp[10]
        this.frustum[5][3] = g.mvp[15] - g.mvp[14]
        this.NormalizePlane(5)
    }
    NormalizePlane(i) {
        let n = Math.sqrt(this.frustum[i][0] * this.frustum[i][0] + this.frustum[i][1] * this.frustum[i][1] + this.frustum[i][2] * this.frustum[i][2])
        this.frustum[i][0] /= n
        this.frustum[i][1] /= n
        this.frustum[i][2] /= n
        this.frustum[i][3] /= n
    }
    Occlude(world, lx, ly, lz) {
        OcclusionViewNum = 0

        if (lx < 0 || lx >= world.width || ly < 0 || ly >= world.height || lz < 0 || lz >= world.length) {
            while (OcclusionViewNum < world.all) {
                world.viewable[OcclusionViewNum] = world.blocks[OcclusionViewNum]
                OcclusionViewNum++
            }
            return
        }

        OcclusionQueuePos = 0
        OcclusionQueueNum = 1
        OcclusionQueue[0] = world.blocks[lx + ly * world.width + lz * world.slice]
        OcclusionQueueFrom[0] = -1

        for (let i = 0; i < world.all; i++)
            OcclusionGoto[i] = true

        while (OcclusionQueueNum > 0) {
            let B = OcclusionQueue[OcclusionQueuePos]
            let from = OcclusionQueueFrom[OcclusionQueuePos]

            world.viewable[OcclusionViewNum] = B
            OcclusionViewNum++

            OcclusionQueuePos++
            if (OcclusionQueuePos === world.all)
                OcclusionQueuePos = 0

            OcclusionQueueNum--

            if (from !== WorldNegativeX)
                this.Visit(world, from, B, WorldPositiveX)

            if (from !== WorldPositiveX)
                this.Visit(world, from, B, WorldNegativeX)

            if (from !== WorldNegativeY)
                this.Visit(world, from, B, WorldPositiveY)

            if (from !== WorldPositiveY)
                this.Visit(world, from, B, WorldNegativeY)

            if (from !== WorldNegativeZ)
                this.Visit(world, from, B, WorldPositiveZ)

            if (from !== WorldPositiveZ)
                this.Visit(world, from, B, WorldNegativeZ)
        }
    }
}
