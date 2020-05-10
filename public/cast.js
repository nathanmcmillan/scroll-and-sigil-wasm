let CastX = 0
let CastY = 0
let CastZ = 0
let CastSide = 0
let CastTileType = null

class Cast {
    static IntGrid(fromX, fromY, toX, toY, visit) {
        let dx = Math.abs(toX - fromX)
        let dy = Math.abs(toY - fromY)
        let x = fromX
        let y = fromY
        let n = 1 + dx + dy
        let stepX = (toX > fromX) ? 1 : -1
        let stepY = (toY > fromY) ? 1 : -1
        let error = dx - dy
        dx *= 2
        dy *= 2
        for (; n > 0; --n) {
            visit(x, y)
            if (error > 0) {
                x += stepX
                error -= dy
            } else {
                y += stepY
                error += dx
            }
        }
    }
    static Block(block, fromX, fromY, fromZ, toX, toY, toZ) {
        let x = Math.floor(fromX)
        let y = Math.floor(fromY)
        let z = Math.floor(fromZ)
        let deltaX, deltaY, deltaZ
        let stepX, stepY, stepZ
        let nextX, nextY, nextZ
        let dx = toX - fromX
        if (dx === 0) {
            stepX = 0
            nextX = Number.MAX_VALUE
        } else if (dx > 0) {
            stepX = 1
            deltaX = 1.0 / dx
            nextX = (1.0 + x - fromX) * deltaX
        } else {
            stepX = -1
            deltaX = 1.0 / -dx
            nextX = (fromX - x) * deltaX
        }
        let dy = toY - fromY
        if (dy === 0) {
            stepY = 0
            nextY = Number.MAX_VALUE
        } else if (dy > 0) {
            stepY = 1
            deltaY = 1.0 / dy
            nextY = (1.0 + y - fromY) * deltaY
        } else {
            stepY = -1
            deltaY = 1.0 / -dy
            nextY = (fromY - y) * deltaY
        }
        let dz = toZ - fromZ
        if (dz === 0) {
            stepZ = 0
            nextZ = Number.MAX_VALUE
        } else if (dz > 0) {
            stepZ = 1
            deltaZ = 1.0 / dz
            nextZ = (1.0 + z - fromZ) * deltaZ
        } else {
            stepZ = -1
            deltaZ = 1.0 / -dz
            nextZ = (fromZ - z) * deltaZ
        }
        while (true) {
            if (TileClosed[block.GetTileTypeUnsafe(x, y, z)]) {
                return false
            } else if (x === Math.floor(toX) && y === Math.floor(toY) && z === Math.floor(toZ)) {
                return true
            }
            if (nextX < nextY) {
                if (nextX < nextZ) {
                    x += stepX
                    if (x < 0 || x >= BlockSize) {
                        return false
                    }
                    nextX += deltaX
                } else {
                    z += stepZ
                    if (z < 0 || z >= BlockSize) {
                        return false
                    }
                    nextZ += deltaZ
                }
            } else {
                if (nextY < nextZ) {
                    y += stepY
                    if (y < 0 || y >= BlockSize) {
                        return false
                    }
                    nextY += deltaY
                } else {
                    z += stepZ
                    if (z < 0 || z >= BlockSize) {
                        return false
                    }
                    nextZ += deltaZ
                }
            }
        }
    }
    static World(world, fromX, fromY, fromZ, toX, toY, toZ) {
        let x = Math.floor(fromX)
        let y = Math.floor(fromY)
        let z = Math.floor(fromZ)
        let deltaX, deltaY, deltaZ
        let stepX, stepY, stepZ
        let nextX, nextY, nextZ
        let dx = toX - fromX
        if (dx === 0) {
            stepX = 0
            deltaX = 0
            nextX = Number.MAX_VALUE
        } else if (dx > 0) {
            stepX = 1
            deltaX = 1.0 / dx
            nextX = (1.0 + x - fromX) * deltaX
        } else {
            stepX = -1
            deltaX = 1.0 / -dx
            nextX = (fromX - x) * deltaX
        }
        let dy = toY - fromY
        if (dy === 0) {
            stepY = 0
            deltaY = 0
            nextY = Number.MAX_VALUE
        } else if (dy > 0) {
            stepY = 1
            deltaY = 1.0 / dy
            nextY = (1.0 + y - fromY) * deltaY
        } else {
            stepY = -1
            deltaY = 1.0 / -dy
            nextY = (fromY - y) * deltaY
        }
        let dz = toZ - fromZ
        if (dz === 0) {
            stepZ = 0
            deltaZ = 0
            nextZ = Number.MAX_VALUE
        } else if (dz > 0) {
            stepZ = 1
            deltaZ = 1.0 / dz
            nextZ = (1.0 + z - fromZ) * deltaZ
        } else {
            stepZ = -1
            deltaZ = 1.0 / -dz
            nextZ = (fromZ - z) * deltaZ
        }
        let goalX = Math.floor(toX)
        let goalY = Math.floor(toY)
        let goalZ = Math.floor(toZ)
        while (true) {
            if (x === goalX && y === goalY && z === goalZ) {
                CastTileType = null
                return
            }
            if (nextX < nextY) {
                if (nextX < nextZ) {
                    x += stepX
                    if (x < 0 || x >= world.tileWidth) {
                        CastTileType = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        CastSide = stepX < 0 ? WorldPositiveX : WorldNegativeX
                        CastTileType = tileType
                        return
                    }
                    nextX += deltaX
                } else {
                    z += stepZ
                    if (z < 0 || z >= world.tileLength) {
                        CastTileType = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        CastSide = stepZ < 0 ? WorldPositiveZ : WorldNegativeZ
                        CastTileType = tileType
                        return
                    }
                    nextZ += deltaZ
                }
            } else {
                if (nextY < nextZ) {
                    y += stepY
                    if (y < 0 || y >= world.tileHeight) {
                        CastTileType = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        CastSide = stepY < 0 ? WorldPositiveY : WorldNegativeY
                        CastTileType = tileType
                        return
                    }
                    nextY += deltaY
                } else {
                    z += stepZ
                    if (z < 0 || z >= world.tileLength) {
                        CastTileType = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        CastSide = stepZ < 0 ? WorldPositiveZ : WorldNegativeZ
                        CastTileType = tileType
                        return
                    }
                    nextZ += deltaZ
                }
            }
        }
    }
    static Exact(world, fromX, fromY, fromZ, toX, toY, toZ) {
        let x = Math.floor(fromX)
        let y = Math.floor(fromY)
        let z = Math.floor(fromZ)
        let deltaX, deltaY, deltaZ
        let stepX, stepY, stepZ
        let nextX, nextY, nextZ
        let dx = toX - fromX
        if (dx === 0) {
            stepX = 0
            deltaX = 0
            nextX = Number.MAX_VALUE
        } else if (dx > 0) {
            stepX = 1
            deltaX = 1.0 / dx
            nextX = (1.0 + x - fromX) * deltaX
        } else {
            stepX = -1
            deltaX = 1.0 / -dx
            nextX = (fromX - x) * deltaX
        }
        let dy = toY - fromY
        if (dy === 0) {
            stepY = 0
            deltaY = 0
            nextY = Number.MAX_VALUE
        } else if (dy > 0) {
            stepY = 1
            deltaY = 1.0 / dy
            nextY = (1.0 + y - fromY) * deltaY
        } else {
            stepY = -1
            deltaY = 1.0 / -dy
            nextY = (fromY - y) * deltaY
        }
        let dz = toZ - fromZ
        if (dz === 0) {
            stepZ = 0
            deltaZ = 0
            nextZ = Number.MAX_VALUE
        } else if (dz > 0) {
            stepZ = 1
            deltaZ = 1.0 / dz
            nextZ = (1.0 + z - fromZ) * deltaZ
        } else {
            stepZ = -1
            deltaZ = 1.0 / -dz
            nextZ = (fromZ - z) * deltaZ
        }
        let goalX = Math.floor(toX)
        let goalY = Math.floor(toY)
        let goalZ = Math.floor(toZ)
        // console.log("---------------------------")
        // console.log("x data", fromX, x, deltaX, nextX, toX)
        // console.log("y data", fromY, y, deltaY, nextY, toY)
        // console.log("z data", fromZ, z, deltaZ, nextZ, toZ)
        while (true) {
            if (x === goalX && y === goalY && z === goalZ) {
                CastX = null
                return
            }
            if (nextX < nextY) {
                if (nextX < nextZ) {
                    x += stepX
                    if (x < 0 || x >= world.tileWidth) {
                        CastX = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        if (stepX < 0) x += 1

                        let vecX = toX - fromX
                        let vecY = toY - fromY
                        let vecZ = toZ - fromZ
                        let magnitude = Math.sqrt(vecX * vecX + vecY * vecY + vecZ * vecZ)
                        vecX /= magnitude
                        vecY /= magnitude
                        vecZ /= magnitude

                        // console.log("xyz", x, y, z)
                        // console.log("vec", vecX, vecY, vecZ, "magnitude", magnitude)
                        let distX = (x - fromX + (1.0 - stepX) * 0.5) / vecX
                        let distY = (y - fromY + (1.0 - stepY) * 0.5) / vecY
                        let distZ = (z - fromZ + (1.0 - stepZ) * 0.5) / vecZ
                        // console.log("distance", distX, distY, distZ)

                        // given known X position and vector
                        // get Y and Z position

                        if (nextY === Number.MAX_VALUE) {
                            y = toY
                        } else {
                            // y = fromY + distY * vecY
                            y += Math.floor(distY) - distY
                        }

                        if (nextZ === Number.MAX_VALUE) {
                            z = toZ
                        } else {
                            // z = fromZ + distZ * vecZ
                            z += Math.floor(distZ) - distZ
                            //z += Math.floor(nextZ) - nextZ
                        }

                        // let cast = new Line(fromX, fromZ, toX, toZ)
                        // let wall = new Line(x, z, x, z + 1)
                        // let point = cast.Intersect(wall)
                        // if (point !== null) {
                        //     console.log("intersect", point[0], point[1])
                        //     z = point[1]

                        //     let vecX = x - fromX
                        //     let vecY = y - fromY
                        //     let vecZ = z - fromZ
                        //     let distance = Math.sqrt(vecX * vecX + vecY * vecY + vecZ * vecZ)
                        //     console.log("true distance", distance)
                        // }

                        if (stepX < 0) x += 0.01
                        else x -= 0.01

                        CastX = x
                        CastY = y
                        CastZ = z
                        return
                    }
                    nextX += deltaX
                } else {
                    z += stepZ
                    if (z < 0 || z >= world.tileLength) {
                        CastX = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        return
                    }
                    nextZ += deltaZ
                }
            } else {
                if (nextY < nextZ) {
                    y += stepY
                    if (y < 0 || y >= world.tileHeight) {
                        CastX = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        return
                    }
                    nextY += deltaY
                } else {
                    z += stepZ
                    if (z < 0 || z >= world.tileLength) {
                        CastX = null
                        return
                    }
                    let bx = Math.floor(x * InverseBlockSize)
                    let by = Math.floor(y * InverseBlockSize)
                    let bz = Math.floor(z * InverseBlockSize)
                    let tx = x - bx * BlockSize
                    let ty = y - by * BlockSize
                    let tz = z - bz * BlockSize
                    let tileType = world.GetTileType(bx, by, bz, tx, ty, tz)
                    if (TileClosed[tileType]) {
                        CastX = x
                        CastY = y
                        CastZ = z
                        return
                    }
                    nextZ += deltaZ
                }
            }
        }
    }
}
