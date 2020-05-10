package main

import (
	"./graphics"
	"./render"
)

// Constants
const (
	BlockShift       = 3
	BlockShiftSlice  = BlockShift + BlockShift
	BlockSize        = 8
	InverseBlockSize = 1.0 / BlockSize
	BlockSlice       = BlockSize * BlockSize
	BlockAll         = BlockSlice * BlockSize

	BlockColorDim   = BlockSize + 1
	BlockColorSlice = BlockColorDim * BlockColorDim
)

var (
	blockMesh         = graphics.RenderCopyInit(3, 3, 2, BlockAll*6*4, BlockAll*6*6)
	blockMeshAmbient  [BlockAll][6][4]float32
	blockMeshColor    [BlockColorDim * BlockColorSlice][3]float32
	blockMeshRgbPlane [4][3]float32
	blockColorSum     [4]int

	blockSliceX       = [6]int{2, 1, 0, 2, 1, 0}
	blockSliceY       = [6]int{0, 2, 1, 0, 2, 1}
	blockSliceZ       = [6]int{1, 0, 2, 1, 0, 2}
	blockSliceTowards = [6]int{1, 1, 1, -1, -1, -1}
	blockSlice        [3]int
	blockSliceTemp    [3]int
)

type block struct {
	x             int
	y             int
	z             int
	mesh          *graphics.RenderBuffer
	visibility    [36]bool
	beginSide     [6]int
	countSide     [6]int
	things        []*thing
	thingCount    int
	scenery       []*scenery
	sceneryCount  int
	items         []*item
	itemCount     int
	missiles      []*missile
	missileCount  int
	particles     []*particle
	particleCount int
	lights        []*light
	lightCount    int
	tiles         [BlockAll]tile
}

func (me *block) blockInit(x, y, z int) {
	me.x = x
	me.y = y
	me.z = z
}

func (me *block) save() string {
	data := "{t["
	if me.notEmpty() == 1 {
		for i := 0; i < BlockAll; i++ {
			data += string(me.tiles[i].typeOf)
			data += ","
		}
	}
	data += "],c["
	for i := 0; i < me.lightCount; i++ {
		data += me.lights[i].save()
		data += ","
	}
	data += "]}"
	return data
}

func (me *block) notEmpty() uint8 {
	for i := 0; i < BlockAll; i++ {
		if me.tiles[i].typeOf != TileNone {
			return 1
		}
	}
	return 0
}

func (me *block) addThing(t *thing) {
	if me.thingCount == len(me.things) {
		array := make([]*thing, me.thingCount+5)
		copy(array, me.things)
		me.things = array
	}
	me.things[me.thingCount] = t
	me.thingCount++
}

func (me *block) addScenery(t *scenery) {
	if me.sceneryCount == len(me.scenery) {
		array := make([]*scenery, me.sceneryCount+5)
		copy(array, me.scenery)
		me.scenery = array
	}
	me.scenery[me.sceneryCount] = t
	me.sceneryCount++
}

func (me *block) addItem(t *item) {
	if me.itemCount == len(me.items) {
		array := make([]*item, me.itemCount+5)
		copy(array, me.items)
		me.items = array
	}
	me.items[me.itemCount] = t
	me.itemCount++
}

func (me *block) addMissile(t *missile) {
	if me.missileCount == len(me.missiles) {
		array := make([]*missile, me.missileCount+5)
		copy(array, me.missiles)
		me.missiles = array
	}
	me.missiles[me.missileCount] = t
	me.missileCount++
}

func (me *block) addParticle(t *particle) {
	if me.particleCount == len(me.particles) {
		array := make([]*particle, me.particleCount+5)
		copy(array, me.particles)
		me.particles = array
	}
	me.particles[me.particleCount] = t
	me.particleCount++
}

func (me *block) addLight(t *light) {
	if me.lightCount == len(me.lights) {
		array := make([]*light, me.lightCount+5)
		copy(array, me.lights)
		me.lights = array
	}
	me.lights[me.lightCount] = t
	me.lightCount++
}

func (me *block) removeThing(t *thing) {
	size := me.thingCount
	for i := 0; i < size; i++ {
		if me.things[i] == t {
			me.things[i] = me.things[size-1]
			me.things[size-1] = nil
			me.thingCount--
			break
		}
	}
}

func (me *block) removeScenery(t *scenery) {
	size := me.sceneryCount
	for i := 0; i < size; i++ {
		if me.scenery[i] == t {
			me.scenery[i] = me.scenery[size-1]
			me.scenery[size-1] = nil
			me.sceneryCount--
			break
		}
	}
}

func (me *block) removeItem(t *item) {
	size := me.itemCount
	for i := 0; i < size; i++ {
		if me.items[i] == t {
			me.items[i] = me.items[size-1]
			me.items[size-1] = nil
			me.itemCount--
			break
		}
	}
}

func (me *block) removeMissile(t *missile) {
	size := me.missileCount
	for i := 0; i < size; i++ {
		if me.missiles[i] == t {
			me.missiles[i] = me.missiles[size-1]
			me.missiles[size-1] = nil
			me.missileCount--
			break
		}
	}
}

func (me *block) removeParticle(t *particle) {
	size := me.particleCount
	for i := 0; i < size; i++ {
		if me.particles[i] == t {
			me.particles[i] = me.particles[size-1]
			me.particles[size-1] = nil
			me.particleCount--
			break
		}
	}
}

func (me *block) removeLight(t *light) {
	size := me.lightCount
	for i := 0; i < size; i++ {
		if me.lights[i] == t {
			me.lights[i] = me.lights[size-1]
			me.lights[size-1] = nil
			me.lightCount--
			break
		}
	}
}

func (me *block) ambientMesh(world *world) {
	for tz := 0; tz < BlockSize; tz++ {
		for ty := 0; ty < BlockSize; ty++ {
			for tx := 0; tx < BlockSize; tx++ {
				index := tx + ty*BlockSize + tz*BlockSlice
				if me.tiles[index].typeOf == TileNone {
					continue
				}
				mmz := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty-1, tz)]
				mmm := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty-1, tz-1)]
				mmp := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty-1, tz+1)]
				mzp := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty, tz+1)]
				mzm := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty, tz-1)]
				mpz := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty+1, tz)]
				mpp := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty+1, tz+1)]
				mpm := TileClosed[world.getTileType(me.x, me.y, me.z, tx-1, ty+1, tz-1)]
				zpp := TileClosed[world.getTileType(me.x, me.y, me.z, tx, ty+1, tz+1)]
				zmp := TileClosed[world.getTileType(me.x, me.y, me.z, tx, ty-1, tz+1)]
				zpm := TileClosed[world.getTileType(me.x, me.y, me.z, tx, ty+1, tz-1)]
				zmm := TileClosed[world.getTileType(me.x, me.y, me.z, tx, ty-1, tz-1)]
				ppz := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty+1, tz)]
				pmz := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty-1, tz)]
				pzp := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty, tz+1)]
				pzm := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty, tz-1)]
				pmm := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty-1, tz-1)]
				ppm := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty+1, tz-1)]
				ppp := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty+1, tz+1)]
				pmp := TileClosed[world.getTileType(me.x, me.y, me.z, tx+1, ty-1, tz+1)]

				blockMeshAmbient[index][WorldPositiveX][0] = TileAmbient(pmz, pzm, pmm)
				blockMeshAmbient[index][WorldPositiveX][1] = TileAmbient(ppz, pzm, ppm)
				blockMeshAmbient[index][WorldPositiveX][2] = TileAmbient(ppz, pzp, ppp)
				blockMeshAmbient[index][WorldPositiveX][3] = TileAmbient(pmz, pzp, pmp)

				blockMeshAmbient[index][WorldNegativeX][0] = TileAmbient(mmz, mzm, mmm)
				blockMeshAmbient[index][WorldNegativeX][1] = TileAmbient(mmz, mzp, mmp)
				blockMeshAmbient[index][WorldNegativeX][2] = TileAmbient(mpz, mzp, mpp)
				blockMeshAmbient[index][WorldNegativeX][3] = TileAmbient(mpz, mzm, mpm)

				blockMeshAmbient[index][WorldPositiveY][0] = TileAmbient(mpz, zpm, mpm)
				blockMeshAmbient[index][WorldPositiveY][1] = TileAmbient(mpz, zpp, mpp)
				blockMeshAmbient[index][WorldPositiveY][2] = TileAmbient(ppz, zpp, ppp)
				blockMeshAmbient[index][WorldPositiveY][3] = TileAmbient(ppz, zpm, ppm)

				blockMeshAmbient[index][WorldNegativeY][0] = TileAmbient(mmz, zmm, mmm)
				blockMeshAmbient[index][WorldNegativeY][1] = TileAmbient(pmz, zmm, pmm)
				blockMeshAmbient[index][WorldNegativeY][2] = TileAmbient(pmz, zmp, pmp)
				blockMeshAmbient[index][WorldNegativeY][3] = TileAmbient(mmz, zmp, mmp)

				blockMeshAmbient[index][WorldPositiveZ][0] = TileAmbient(pzp, zmp, pmp)
				blockMeshAmbient[index][WorldPositiveZ][1] = TileAmbient(pzp, zpp, ppp)
				blockMeshAmbient[index][WorldPositiveZ][2] = TileAmbient(mzp, zpp, mpp)
				blockMeshAmbient[index][WorldPositiveZ][3] = TileAmbient(mzp, zmp, mmp)

				blockMeshAmbient[index][WorldNegativeZ][0] = TileAmbient(mzm, zmm, mmm)
				blockMeshAmbient[index][WorldNegativeZ][1] = TileAmbient(mzm, zpm, mpm)
				blockMeshAmbient[index][WorldNegativeZ][2] = TileAmbient(pzm, zpm, ppm)
				blockMeshAmbient[index][WorldNegativeZ][3] = TileAmbient(pzm, zmm, pmm)
			}
		}
	}
}

func (me *block) colorMesh(world *world) {
	for tz := 0; tz < BlockColorDim; tz++ {
		for ty := 0; ty < BlockColorDim; ty++ {
			for tx := 0; tx < BlockColorDim; tx++ {
				blockColorSum[0] = 0
				blockColorSum[1] = 0
				blockColorSum[2] = 0
				blockColorSum[3] = 0
				zzz := world.getTilePointer(me.x, me.y, me.z, tx, ty, tz)
				mzz := world.getTilePointer(me.x, me.y, me.z, tx-1, ty, tz)
				mzm := world.getTilePointer(me.x, me.y, me.z, tx-1, ty, tz-1)
				zzm := world.getTilePointer(me.x, me.y, me.z, tx, ty, tz-1)
				zmz := world.getTilePointer(me.x, me.y, me.z, tx, ty-1, tz)
				mmz := world.getTilePointer(me.x, me.y, me.z, tx-1, ty-1, tz)
				mmm := world.getTilePointer(me.x, me.y, me.z, tx-1, ty-1, tz-1)
				zmm := world.getTilePointer(me.x, me.y, me.z, tx, ty-1, tz-1)

				if zzz == nil || TileClosed[zzz.typeOf] {
					me.determineLight(mzz)
					me.determineLight(zmz)
					me.determineLight(zzm)
				}
				if mzz == nil || TileClosed[mzz.typeOf] {
					me.determineLight(zzz)
					me.determineLight(zmz)
					me.determineLight(zzm)
				}
				if mzm == nil || TileClosed[mzm.typeOf] {
					me.determineLight(mzz)
					me.determineLight(zzm)
					me.determineLight(mmm)
				}
				if zzm == nil || TileClosed[zzm.typeOf] {
					me.determineLight(zzz)
					me.determineLight(mzm)
					me.determineLight(zmm)
				}
				if zmz == nil || TileClosed[zmz.typeOf] {
					me.determineLight(zzz)
					me.determineLight(mmz)
					me.determineLight(zmm)
				}
				if mmz == nil || TileClosed[mmz.typeOf] {
					me.determineLight(mzz)
					me.determineLight(mmm)
					me.determineLight(zmz)
				}
				if mmm == nil || TileClosed[mmm.typeOf] {
					me.determineLight(mzm)
					me.determineLight(zmm)
					me.determineLight(mmz)
				}
				if zmm == nil || TileClosed[zmm.typeOf] {
					me.determineLight(zzm)
					me.determineLight(zmz)
					me.determineLight(mmm)
				}

				index := tx + ty*BlockColorDim + tz*BlockColorSlice
				size := float32(blockColorSum[3])
				if size > 0 {
					blockMeshColor[index][0] = float32(blockColorSum[0]) / size
					blockMeshColor[index][1] = float32(blockColorSum[1]) / size
					blockMeshColor[index][2] = float32(blockColorSum[2]) / size
				} else {
					blockMeshColor[index][0] = 255.0
					blockMeshColor[index][1] = 255.0
					blockMeshColor[index][2] = 255.0
				}
			}
		}
	}
}

func (me *block) determineLight(t *tile) {
	if t == nil {
		return
	}
	if !TileClosed[t.typeOf] {
		blockColorSum[0] += int(t.red)
		blockColorSum[1] += int(t.green)
		blockColorSum[2] += int(t.blue)
		blockColorSum[3]++
	}
}

func (me *block) lightOfSide(xs, ys, zs, side int) {
	switch side {
	case WorldPositiveX:
		blockMeshRgbPlane[0] = blockMeshColor[xs+1+ys*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+1+ys*BlockColorDim+(zs+1)*BlockColorSlice]
	case WorldNegativeX:
		blockMeshRgbPlane[0] = blockMeshColor[xs+ys*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+ys*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+(ys+1)*BlockColorDim+zs*BlockColorSlice]
	case WorldPositiveY:
		blockMeshRgbPlane[0] = blockMeshColor[xs+(ys+1)*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+zs*BlockColorSlice]
	case WorldNegativeY:
		blockMeshRgbPlane[0] = blockMeshColor[xs+ys*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+1+ys*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+1+ys*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+ys*BlockColorDim+(zs+1)*BlockColorSlice]
	case WorldPositiveZ:
		blockMeshRgbPlane[0] = blockMeshColor[xs+1+ys*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+(ys+1)*BlockColorDim+(zs+1)*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+ys*BlockColorDim+(zs+1)*BlockColorSlice]
	default:
		blockMeshRgbPlane[0] = blockMeshColor[xs+ys*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[1] = blockMeshColor[xs+(ys+1)*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[2] = blockMeshColor[xs+1+(ys+1)*BlockColorDim+zs*BlockColorSlice]
		blockMeshRgbPlane[3] = blockMeshColor[xs+1+ys*BlockColorDim+zs*BlockColorSlice]
	}
}

func (me *block) buildMesh(world *world) {
	me.ambientMesh(world)
	me.colorMesh(world)
	blockMesh.Zero()
	for side := 0; side < 6; side++ {
		meshBeginIndex := blockMesh.IndexPos
		pointerX := blockSliceX[side]
		pointerY := blockSliceY[side]
		pointerZ := blockSliceZ[side]
		toward := blockSliceTowards[side]
		for blockSlice[2] = 0; blockSlice[2] < BlockSize; blockSlice[2]++ {
			for blockSlice[1] = 0; blockSlice[1] < BlockSize; blockSlice[1]++ {
				for blockSlice[0] = 0; blockSlice[0] < BlockSize; blockSlice[0]++ {
					typeOf := me.tiles[blockSlice[pointerX]+(blockSlice[pointerY]<<BlockShift)+(blockSlice[pointerZ]<<BlockShiftSlice)].typeOf
					if typeOf == TileNone {
						continue
					}
					blockSliceTemp[0] = blockSlice[0]
					blockSliceTemp[1] = blockSlice[1]
					blockSliceTemp[2] = blockSlice[2] + toward
					if TileClosed[world.getTileType(me.x, me.y, me.z, blockSliceTemp[pointerX], blockSliceTemp[pointerY], blockSliceTemp[pointerZ])] {
						continue
					}

					xs := blockSlice[pointerX]
					ys := blockSlice[pointerY]
					zs := blockSlice[pointerZ]
					index := xs + (ys << BlockShift) + (zs << BlockShiftSlice)

					texture := &TileTexture[typeOf]
					gx := float32(xs + BlockSize*me.x)
					gy := float32(ys + BlockSize*me.y)
					gz := float32(zs + BlockSize*me.z)

					me.lightOfSide(xs, ys, zs, side)
					lightPlane(&blockMeshRgbPlane, &blockMeshAmbient[index][side])

					switch side {
					case WorldPositiveX:
						render.RendTilePosX(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					case WorldNegativeX:
						render.RendTileNegX(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					case WorldPositiveY:
						render.RendTilePosY(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					case WorldNegativeY:
						render.RendTileNegY(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					case WorldPositiveZ:
						render.RendTilePosZ(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					case WorldNegativeZ:
						render.RendTileNegZ(blockMesh, gx, gy, gz, texture, &blockMeshRgbPlane)
					}
				}
			}
		}
		me.beginSide[side] = meshBeginIndex * 4
		me.countSide[side] = blockMesh.IndexPos - meshBeginIndex
	}
	me.mesh = blockMesh.InitCopy(world.gl)
}

func (me *block) renderThings(spriteSet map[interface{}]bool, spriteBuffer map[string]*graphics.RenderBuffer, camX, camZ, camAngle float32) {
	for i := 0; i < me.thingCount; i++ {
		thing := me.things[i]
		if _, ok := spriteSet[thing]; ok {
			continue
		}
		spriteSet[thing] = true
		thing.render(spriteBuffer, camX, camZ, camAngle)
	}
	for i := 0; i < me.sceneryCount; i++ {
		scene := me.scenery[i]
		if _, ok := spriteSet[scene]; ok {
			continue
		}
		spriteSet[scene] = true
		scene.render(spriteBuffer, camX, camZ)
	}
	for i := 0; i < me.itemCount; i++ {
		item := me.items[i]
		if _, ok := spriteSet[item]; ok {
			continue
		}
		spriteSet[item] = true
		item.render(spriteBuffer, camX, camZ)
	}
	for i := 0; i < me.missileCount; i++ {
		missile := me.missiles[i]
		if _, ok := spriteSet[missile]; ok {
			continue
		}
		spriteSet[missile] = true
		missile.render(spriteBuffer, camX, camZ, camAngle)
	}
	for i := 0; i < me.particleCount; i++ {
		particle := me.particles[i]
		if _, ok := spriteSet[particle]; ok {
			continue
		}
		spriteSet[particle] = true
		particle.render(spriteBuffer, camX, camZ)
	}
}
