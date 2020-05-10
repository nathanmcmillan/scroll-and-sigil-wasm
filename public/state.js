class WorldState {
    constructor(app) {
        this.app = app
        this.snapshotTime = new Date().getTime()
        this.previousUpdate = new Date().getTime()
        this.chatbox = ["poop"]
    }
    serverUpdates() {
        let world = this.app.world

        for (let i = 0; i < SocketQueue.length; i += 1) {
            let dat = new DataView(SocketQueue[i])
            let dex = 0

            let serverTime = dat.getUint32(dex, true)
            dex += 4

            this.snapshotTime = serverTime + 1552330000000
            this.previousUpdate = new Date().getTime()

            let broadcastCount = dat.getUint8(dex, true)
            dex += 1
            for (let b = 0; b < broadcastCount; b += 1) {
                let broadcastType = dat.getUint8(dex, true)
                dex += 1
                switch (broadcastType) {
                    case BroadcastNew:
                        {
                            let uid = dat.getUint16(dex, true)
                            dex += 2
                            let nid = dat.getUint16(dex, true)
                            dex += 2
                            if (world.netLookup.has(nid))
                                break
                            let x = dat.getFloat32(dex, true)
                            dex += 4
                            let y = dat.getFloat32(dex, true)
                            dex += 4
                            let z = dat.getFloat32(dex, true)
                            dex += 4
                            if (uid === PlasmaUID) {
                                let dx = dat.getFloat32(dex, true)
                                dex += 4
                                let dy = dat.getFloat32(dex, true)
                                dex += 4
                                let dz = dat.getFloat32(dex, true)
                                dex += 4
                                let damage = dat.getUint16(dex, true)
                                dex += 2
                                new Plasma(world, nid, damage, x, y, z, dx, dy, dz)
                            } else if (uid === HumanUID) {
                                let angle = dat.getFloat32(dex, true)
                                dex += 4
                                let health = dat.getUint16(dex, true)
                                dex += 2
                                let status = dat.getUint8(dex, true)
                                dex += 1
                                new Human(world, nid, x, y, z, angle, health, status)
                            }
                        }
                        break
                    case BroadcastDelete:
                        {
                            let nid = dat.getUint16(dex, true)
                            dex += 2
                            let entity = world.netLookup.get(nid)
                            if (entity) entity.Cleanup()
                        }
                        break
                    case BroadcastChat:
                        {
                            let num = dat.getUint8(dex, true)
                            dex += 1
                            let chat = ""
                            for (let ch = 0; ch < num; ch += 1) {
                                chat += String.fromCharCode(dat.getUint8(dex, true))
                                dex += 1
                            }
                            this.chatbox.push(chat)
                        }
                        break
                }
            }

            let thingCount = dat.getUint16(dex, true)
            dex += 2
            for (let t = 0; t < thingCount; t += 1) {
                let nid = dat.getUint16(dex, true)
                dex += 2
                let delta = dat.getUint8(dex, true)
                dex += 1
                let thing = world.netLookup.get(nid)
                if (thing) {
                    if (delta & 0x1) {
                        thing.NetX = dat.getFloat32(dex, true)
                        thing.DeltaNetX = (thing.NetX - thing.X) * InverseNetRate
                        dex += 4
                        thing.NetZ = dat.getFloat32(dex, true)
                        thing.DeltaNetZ = (thing.NetZ - thing.Z) * InverseNetRate
                        dex += 4
                    }
                    if (delta & 0x2) {
                        thing.NetY = dat.getFloat32(dex, true)
                        thing.DeltaNetY = (thing.NetY - thing.Y) * InverseNetRate
                        dex += 4
                    }
                    if (delta & 0x4) {
                        let health = dat.getUint16(dex, true)
                        dex += 2
                        thing.NetUpdateHealth(health)
                    }
                    if (delta & 0x8) {
                        let status = dat.getUint8(dex, true)
                        dex += 1
                        thing.NetUpdateState(status)
                    }
                    switch (thing.UID) {
                        case HumanUID:
                            if (delta & 0x10) {
                                thing.Angle = dat.getFloat32(dex, true)
                                dex += 4
                            }
                            break
                        case BaronUID:
                            if (delta & 0x10) {
                                let direction = dat.getUint8(dex, true)
                                dex += 1
                                if (direction !== DirectionNone)
                                    thing.Angle = DirectionToAngle[direction]
                            }
                            break
                    }
                } else {
                    throw new Error("missing thing nid " + nid)
                }
            }
        }
    }
    update() {
        let world = this.app.world

        if (SocketQueue.length > 0) {
            this.serverUpdates()
            SocketQueue = []
        }

        world.update()

        let socketSendOperations = 0
        for (var [op, value] of SocketSendSet) {
            SocketSend.setUint8(SocketSendIndex, op, true)
            SocketSendIndex += 1
            if (op === InputOpNewMove) {
                SocketSend.setFloat32(SocketSendIndex, value, true)
                SocketSendIndex += 4
            } else if (op === InputOpChat) {
                let chat = Math.min(value.length, 255)
                SocketSend.setUint8(SocketSendIndex, chat, true)
                SocketSendIndex += 1
                for (let ch = 0; ch < chat; ch += 1) {
                    SocketSend.setUint8(SocketSendIndex, value.charCodeAt(ch), true)
                    SocketSendIndex += 1
                }
            }
            socketSendOperations += 1
        }
        if (socketSendOperations > 0) {
            let buffer = SocketSend.buffer.slice(0, SocketSendIndex)
            let view = new DataView(buffer)
            view.setUint8(0, socketSendOperations, true)
            SocketConnection.send(buffer)
            SocketSendIndex = 1
            SocketSendSet.clear()
        }
    }
    render() {
        let g = this.app.g
        let gl = this.app.gl
        let frameGeo = this.app.frameGeo
        let frame = this.app.frame
        let canvas = this.app.canvas
        let canvasOrtho = this.app.canvasOrtho
        let drawPerspective = this.app.drawPerspective
        let drawOrtho = this.app.drawOrtho
        let drawImages = this.app.drawImages
        let screen = this.app.screen
        let world = this.app.world
        let cam = this.app.camera

        cam.update(world)

        RenderSystem.SetFrameBuffer(gl, frameGeo.fbo)
        RenderSystem.SetView(gl, 0, 0, frame.width, frame.height)

        gl.clear(gl.COLOR_BUFFER_BIT)
        gl.clear(gl.DEPTH_BUFFER_BIT)

        g.SetProgram(gl, "copy")
        g.SetOrthographic(drawOrtho, 0, 0)
        g.UpdateMvp(gl)
        g.SetTexture(gl, "sky")
        drawImages.Zero()

        let turnX = frame.width * 2
        let skyX = cam.ry / Tau * turnX;
        if (skyX >= turnX) skyX -= turnX
        let skyYOffset = cam.rx / Tau * frame.height;
        let skTrueHeight = g.textures.get("sky").image.height
        let skyHeight = skTrueHeight * 2
        let skyTop = frame.height - skyHeight
        let skyY = skyTop * 0.5 + skyYOffset
        if (skyY > skyTop) skyY = skyTop
        Render.Image(drawImages, -skyX, skyY, turnX * 2, skyHeight, 0, 0, 2, 1)
        RenderSystem.UpdateAndDraw(gl, drawImages)

        gl.enable(gl.DEPTH_TEST)
        gl.enable(gl.CULL_FACE)

        g.SetPerspective(drawPerspective, -cam.x, -cam.y, -cam.z, cam.rx, cam.ry)

        let camBlockX = Math.floor(cam.x * InverseBlockSize)
        let camBlockY = Math.floor(cam.y * InverseBlockSize)
        let camBlockZ = Math.floor(cam.z * InverseBlockSize)

        world.render(g, camBlockX, camBlockY, camBlockZ, cam.x, cam.z, cam.ry)

        gl.disable(gl.CULL_FACE)
        gl.disable(gl.DEPTH_TEST)

        const noShade = 0
        const motionBlur = 1
        const antiAlias = 2
        let shading = noShade

        if (shading === motionBlur) {
            let frame2 = this.app.frame2
            let frameScreen = this.app.frameScreen

            let drawInversePerspective = this.app.drawInversePerspective
            let drawInverseMv = this.app.drawInverseMv
            let drawPreviousMvp = this.app.drawPreviousMvp
            let drawCurrentToPreviousMvp = this.app.drawCurrentToPreviousMvp
            Matrix.Inverse(drawInverseMv, g.mv)
            Matrix.Multiply(drawCurrentToPreviousMvp, drawPreviousMvp, drawInverseMv)
            for (let i = 0; i < 16; i++) {
                drawPreviousMvp[i] = g.mvp[i]
            }

            RenderSystem.SetFrameBuffer(gl, frame2.fbo)
            RenderSystem.SetView(gl, 0, 0, frame2.width, frame2.height)
            g.SetProgram(gl, "motion")
            g.SetUniformMatrix4(gl, "inverse_projection", drawInversePerspective);
            g.SetUniformMatrix4(gl, "current_to_previous_matrix", drawCurrentToPreviousMvp);
            g.SetOrthographic(drawOrtho, 0, 0)
            g.UpdateMvp(gl)
            g.SetTextureDirect(gl, frameGeo.textures[0])
            g.SetIndexTextureDirect(gl, 1, "u_texture1", frameGeo.depthTexture)
            RenderSystem.BindAndDraw(gl, frameScreen)

            RenderSystem.SetFrameBuffer(gl, null)
            RenderSystem.SetView(gl, 0, 0, canvas.width, canvas.height)
            g.SetProgram(gl, "screen")
            g.SetOrthographic(canvasOrtho, 0, 0)
            g.UpdateMvp(gl)
            g.SetTextureDirect(gl, frame2.textures[0])
            RenderSystem.BindAndDraw(gl, screen)
        } else if (shading === antiAlias) {
            RenderSystem.SetFrameBuffer(gl, null)
            RenderSystem.SetView(gl, 0, 0, canvas.width, canvas.height)
            g.SetProgram(gl, "fxaa")
            g.SetUniformVec2(gl, "texel", 1.0 / canvas.width, 1.0 / canvas.height)
            g.SetOrthographic(canvasOrtho, 0, 0)
            g.UpdateMvp(gl)
            g.SetTextureDirect(gl, frameGeo.textures[0])
            RenderSystem.BindAndDraw(gl, screen)
        } else {
            g.SetProgram(gl, "texture2d")
            g.SetOrthographic(drawOrtho, 0, 0)
            g.UpdateMvp(gl)
            g.SetTexture(gl, "font")
            drawImages.Zero()
            let chatbox = this.chatbox
            let y = 10
            for (let ch = 0; ch < chatbox.length; ch++) {
                Render.Print(drawImages, 10, y, chatbox[ch], 2)
                y += FontHeight * 2
            }
            RenderSystem.UpdateAndDraw(gl, drawImages)

            RenderSystem.SetFrameBuffer(gl, null)
            RenderSystem.SetView(gl, 0, 0, canvas.width, canvas.height)
            g.SetProgram(gl, "screen")
            g.SetOrthographic(canvasOrtho, 0, 0)
            g.UpdateMvp(gl)
            g.SetTextureDirect(gl, frameGeo.textures[0])
            RenderSystem.BindAndDraw(gl, screen)
        }
    }
}
