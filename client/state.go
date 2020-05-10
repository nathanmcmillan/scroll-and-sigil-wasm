package main

import (
	"strings"
	"syscall/js"
	"time"

	"../fast"
	"./graphics"
	"./matrix"
	"./render"
)

type worldState struct {
	app            *app
	snapshotTime   int64
	previousUpdate int64
	chatbox        []string
}

func worldStateInit(a *app) *worldState {
	s := &worldState{}
	s.app = a
	return s
}

func (me *worldState) serverUpdates() {
	world := me.app.world
	socketQueue := me.app.socketQueue

	for i := 0; i < len(socketQueue); i++ {
		data := fast.ByteReaderInit(socketQueue[i])

		me.snapshotTime = int64(data.GetUint32()) + 1552330000000
		me.previousUpdate = time.Now().UnixNano()

		broadcastCount := data.GetUint8()
		for b := uint8(0); b < broadcastCount; b++ {
			broadcastType := data.GetUint8()
			switch broadcastType {
			case BroadcastNew:
				uid := data.GetUint16()
				nid := data.GetUint16()
				if _, ok := world.netLookup[nid]; ok {
					break
				}
				x := data.GetFloat32()
				y := data.GetFloat32()
				z := data.GetFloat32()
				switch uid {
				case PlasmaUID:
					dx := data.GetFloat32()
					dy := data.GetFloat32()
					dz := data.GetFloat32()
					damage := data.GetUint16()
					plasmaInit(world, nid, damage, x, y, z, dx, dy, dz)
				case HumanUID:
					angle := data.GetFloat32()
					health := data.GetUint16()
					status := data.GetUint8()
					humanInit(world, nid, x, y, z, angle, health, status)
				default:
					panic("unknown UID")
				}
			case BroadcastDelete:
				nid := data.GetUint16()
				if thing, ok := world.netLookup[nid]; ok {
					thing.cleanup()
				}
			case BroadcastChat:
				size := data.GetUint8()
				chat := &strings.Builder{}
				for i := uint8(0); i < size; i++ {
					ch := data.GetUint8()
					chat.WriteByte(ch)
				}
				me.chatbox = append(me.chatbox, chat.String())
			}
		}

		thingCount := data.GetUint16()
		for t := uint16(0); t < thingCount; t++ {
			nid := data.GetUint16()
			thing, ok := world.netLookup[nid]
			if !ok {
				panic("missing thing nid")
			}
			delta := data.GetUint8()
			thing.netUpdate(data, delta)
		}
	}
}

func (me *worldState) update() {
	world := me.app.world
	socketQueue := me.app.socketQueue

	if len(socketQueue) > 0 {
		me.serverUpdates()
		me.app.socketQueue = make([][]byte, 0)
	}

	world.update()

	socketSend := me.app.socketSend
	socketSendOperations := uint8(0)
	size := len(socketSend)
	if size > 0 {
		data := fast.ByteWriterInit(64)
		data.Position(1)
		for op, value := range socketSend {
			data.PutUint8(op)
			switch op {
			case inputOpNewMove:
				data.PutFloat32(value.(float32))
			case inputOpChat:
				chat := value.(string)
				chatSize := uint8(len(chat))
				if chatSize > 255 {
					chatSize = 255
				}
				data.PutUint8(chatSize)
				for ch := uint8(0); ch < chatSize; ch++ {
					data.PutUint8(chat[ch])
				}
			}
			socketSendOperations++
		}
		data.SetUint8(0, socketSendOperations)
		me.app.socket.Call("send", js.TypedArrayOf(data.Bytes()))

		for key := range socketSend {
			delete(socketSend, key)
		}
	}
}

func (me *worldState) render() {
	app := me.app
	g := app.g
	gl := app.gl
	frameGeo := app.frameGeo
	frame := app.frame
	canvas := app.canvas
	canvasOrtho := app.canvasOrtho
	drawPerspective := app.drawPerspective
	drawOrtho := app.drawOrtho
	drawImages := app.drawImages
	screen := app.screen
	world := app.world
	cam := app.camera

	cam.update(world)

	graphics.RenderSystemSetFrameBuffer(gl, frameGeo.Fbo)
	graphics.RenderSystemSetView(gl, 0, 0, frame.Width, frame.Height)

	gl.Call("clear", graphics.GLxColorBufferBit)
	gl.Call("clear", graphics.GLxDepthBufferBit)

	g.SetProgram(gl, "copy")
	g.SetOrthographic(drawOrtho, 0, 0)
	g.UpdateMvp(gl)
	g.SetTexture(gl, "sky")
	drawImages.Zero()

	turnX := float32(frame.Width * 2)
	skyX := cam.ry / Tau * turnX
	if skyX >= turnX {
		skyX -= turnX
	}
	frameHeight := float32(frame.Height)
	skyYOffset := cam.rx / Tau * frameHeight
	skTrueHeight := float32(g.Textures["sky"].Height)
	skyHeight := skTrueHeight * 2
	skyTop := frameHeight - skyHeight
	skyY := float32(skyTop)*0.5 + skyYOffset
	if skyY > skyTop {
		skyY = skyTop
	}
	render.Image(drawImages, -skyX, skyY, turnX*2, skyHeight, 0, 0, 2, 1)
	graphics.RenderSystemUpdateAndDraw(gl, drawImages)

	gl.Call("enable", graphics.GLxDepthTest)
	gl.Call("enable", graphics.GLxCullFace)

	g.SetPerspective(drawPerspective, -cam.x, -cam.y, -cam.z, cam.rx, cam.ry)

	camBlockX := int(cam.x * InverseBlockSize)
	camBlockY := int(cam.y * InverseBlockSize)
	camBlockZ := int(cam.z * InverseBlockSize)

	world.render(g, camBlockX, camBlockY, camBlockZ, cam.x, cam.z, cam.ry)

	gl.Call("disable", graphics.GLxCullFace)
	gl.Call("disable", graphics.GLxDepthTest)

	const noShade = 0
	const motionBlur = 1
	const antiAlias = 2
	shading := noShade

	if shading == motionBlur {
		frame2 := app.frame2
		frameScreen := app.frameScreen

		drawInversePerspective := app.drawInversePerspective
		drawInverseMv := app.drawInverseMv
		drawPreviousMvp := app.drawPreviousMvp
		drawCurrentToPreviousMvp := app.drawCurrentToPreviousMvp
		matrix.Inverse(drawInverseMv, g.ModelView)
		matrix.Multiply(drawCurrentToPreviousMvp, drawPreviousMvp, drawInverseMv)
		for i := 0; i < 16; i++ {
			drawPreviousMvp[i] = g.ModelViewProject[i]
		}

		graphics.RenderSystemSetFrameBuffer(gl, frame2.Fbo)
		graphics.RenderSystemSetView(gl, 0, 0, frame2.Width, frame2.Height)
		g.SetProgram(gl, "motion")
		g.SetUniformMatrix4(gl, "inverse_projection", drawInversePerspective)
		g.SetUniformMatrix4(gl, "current_to_previous_matrix", drawCurrentToPreviousMvp)
		g.SetOrthographic(drawOrtho, 0, 0)
		g.UpdateMvp(gl)
		g.SetTextureDirect(gl, frameGeo.Textures[0])
		g.SetIndexTextureDirect(gl, graphics.GLxTexture1, 1, "u_texture1", frameGeo.DepthTexture)
		graphics.RenderSystemBindAndDraw(gl, frameScreen)

		graphics.RenderSystemSetFrameBuffer(gl, js.Null())
		graphics.RenderSystemSetView(gl, 0, 0, canvas.width, canvas.height)
		g.SetProgram(gl, "screen")
		g.SetOrthographic(canvasOrtho, 0, 0)
		g.UpdateMvp(gl)
		g.SetTextureDirect(gl, frame2.Textures[0])
		graphics.RenderSystemBindAndDraw(gl, screen)
	} else if shading == antiAlias {
		graphics.RenderSystemSetFrameBuffer(gl, js.Null())
		graphics.RenderSystemSetView(gl, 0, 0, canvas.width, canvas.height)
		g.SetProgram(gl, "fxaa")
		g.SetUniformVec2(gl, "texel", 1.0/float32(canvas.width), 1.0/float32(canvas.height))
		g.SetOrthographic(canvasOrtho, 0, 0)
		g.UpdateMvp(gl)
		g.SetTextureDirect(gl, frameGeo.Textures[0])
		graphics.RenderSystemBindAndDraw(gl, screen)
	} else {
		g.SetProgram(gl, "texture2d")
		g.SetOrthographic(drawOrtho, 0, 0)
		g.UpdateMvp(gl)
		g.SetTexture(gl, "font")
		drawImages.Zero()
		chatbox := me.chatbox
		y := float32(10)
		for ch := 0; ch < len(chatbox); ch++ {
			render.Print(drawImages, 10, y, chatbox[ch], 2)
			y += render.FontHeight * 2
		}
		graphics.RenderSystemUpdateAndDraw(gl, drawImages)

		graphics.RenderSystemSetFrameBuffer(gl, js.Null())
		graphics.RenderSystemSetView(gl, 0, 0, canvas.width, canvas.height)

		g.SetProgram(gl, "screen")
		g.SetOrthographic(canvasOrtho, 0, 0)
		g.UpdateMvp(gl)
		g.SetTextureDirect(gl, frameGeo.Textures[0])
		graphics.RenderSystemBindAndDraw(gl, screen)
	}
}
