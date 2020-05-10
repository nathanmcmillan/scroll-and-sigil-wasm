class EditState {
    constructor(app) {
        this.app = app
        this.name = "temp"
        this.editMode = "add"
        this.editTileType = TileGrass
    }
    update() {
        let world = this.app.world
        let cam = this.app.camera

        cam.update()


        if (Input.KeyPress("1")) {
            this.editTileType = TileLookup.get("grass")
        } else if (Input.KeyPress("2")) {
            this.editTileType = TileLookup.get("planks")
        } else if (Input.KeyPress("3")) {
            this.editTileType = TileLookup.get("stone")
        }

        if (Input.KeyPress("i")) {
            this.editMode = "add"
        }

        if (Input.KeyPress("r")) {
            this.editMode = "replace"
        }

        if (Input.KeyPress("c")) {
            this.editMode = "delete"
        }

        let distance = 10.0
        let toX = cam.X + distance * Math.sin(cam.RY)
        let toY = cam.Y - distance * Math.sin(cam.RX)
        let toZ = cam.Z - distance * Math.cos(cam.RY)

        Cast.World(world, cam.X, cam.Y, cam.Z, toX, toY, toZ)

        if (CastTileType !== null) {
            if (this.editMode === "add") {
                switch (CastSide) {
                    case WorldNegativeX:
                        CastX--
                        if (CastX < 0) CastTileType = null
                        break
                    case WorldPositiveX:
                        CastX++
                        if (CastX >= world.tileWidth) CastTileType = null
                        break
                    case WorldNegativeY:
                        CastY--
                        if (CastY < 0) CastTileType = null
                        break
                    case WorldPositiveY:
                        CastY++
                        if (CastY >= world.tileHeigh) CastTileType = null
                        break
                    case WorldNegativeZ:
                        CastZ--
                        if (CastZ < 0) CastTileType = null
                        break
                    case WorldPositiveZ:
                        CastZ++
                        if (CastZ >= world.tileLength) CastTileType = null
                        break
                }
            }
            if (CastTileType !== null && Input.KeyPress(" ")) {
                let bx = Math.floor(CastX * InverseBlockSize)
                let by = Math.floor(CastY * InverseBlockSize)
                let bz = Math.floor(CastZ * InverseBlockSize)
                let tx = CastX - bx * BlockSize
                let ty = CastY - by * BlockSize
                let tz = CastZ - bz * BlockSize
                let block = world.blocks[bx + by * world.width + bz * world.slice]
                let tile = block.tiles[tx + ty * BlockSize + tz * BlockSlice]
                if (this.editMode === "delete") {
                    if (tile.type !== TileNone) {
                        tile.type = TileNone
                        block.BuildMesh(world)
                    }
                } else if (this.editMode === "replace") {
                    if (tile.type !== this.editTileType) {
                        tile.type = this.editTileType
                        block.BuildMesh(world)
                    }
                } else if (this.editMode === "add") {
                    if (tile.type !== this.editTileType) {
                        tile.type = this.editTileType
                        block.BuildMesh(world)
                    }
                }
            }
        }

        if (Input.KeyPress("m")) {
            this.pressSave = false
            if (this.name !== "") {
                let data = world.Save()
                console.log(data)
                Net.Send("map", this.name + ":" + data)
            }
        }
    }
    render() {
        let g = this.app.g
        let gl = this.app.gl
        let frame = this.app.frame
        let canvas = this.app.canvas
        let canvasOrtho = this.app.canvasOrtho
        let drawPerspective = this.app.drawPerspective
        let drawOrtho = this.app.drawOrtho
        let drawImages = this.app.drawImages
        let screen = this.app.screen
        let world = this.app.world
        let cam = this.app.camera

        RenderSystem.SetFrameBuffer(gl, frame.fbo)
        RenderSystem.SetView(gl, 0, 0, frame.width, frame.height)

        gl.clear(gl.COLOR_BUFFER_BIT)
        gl.clear(gl.DEPTH_BUFFER_BIT)
        gl.enable(gl.DEPTH_TEST)
        gl.enable(gl.CULL_FACE)

        g.SetPerspective(drawPerspective, -cam.X, -cam.Y, -cam.Z, cam.RX, cam.RY)

        if (CastTileType !== null) {
            let singleBlock = this.app.singleBlock
            g.SetProgram(gl, "texture-color3d")
            g.UpdateMvp(gl)
            g.SetTexture(gl, "tiles")
            singleBlock.Zero()
            let rgb = [1, 1, 1]
            for (let side = 0; side < 6; side++)
                RenderTile.Side(singleBlock, side, CastX, CastY, CastZ, TileTexture[this.editTileType], rgb, rgb, rgb, rgb)
            RenderSystem.UpdateAndDraw(gl, singleBlock)
        }

        let camBlockX = Math.floor(cam.X * InverseBlockSize)
        let camBlockY = Math.floor(cam.Y * InverseBlockSize)
        let camBlockZ = Math.floor(cam.Z * InverseBlockSize)

        world.render(g, camBlockX, camBlockY, camBlockZ, cam.X, cam.Z, cam.RY)

        gl.disable(gl.DEPTH_TEST)
        gl.disable(gl.CULL_FACE)

        g.SetProgram(gl, "copy")
        g.SetOrthographic(drawOrtho, 0, 0)
        g.UpdateMvp(gl)
        g.SetTexture(gl, "tiles")
        drawImages.Zero()
        let tileTexture = TileTexture[this.editTileType]
        Render.Image(drawImages, 10, 10, 32, 32, tileTexture[0], tileTexture[1], tileTexture[2], tileTexture[3])
        Render.Image(drawImages, Math.floor(frame.width * 0.5 - 5), Math.floor(frame.height * 0.5 - 5), 10, 10, tileTexture[0], tileTexture[1], tileTexture[2], tileTexture[3])
        RenderSystem.UpdateAndDraw(gl, drawImages)

        RenderSystem.SetFrameBuffer(gl, null)
        RenderSystem.SetView(gl, 0, 0, canvas.width, canvas.height)
        g.SetProgram(gl, "copy")
        g.SetOrthographic(canvasOrtho, 0, 0)
        g.UpdateMvp(gl)
        g.SetTextureDirect(gl, frame.textures[0])
        RenderSystem.BindAndDraw(gl, screen)
    }
}
