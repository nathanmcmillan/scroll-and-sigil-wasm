package main

import (
	"math"
	"sync"
	"syscall/js"

	"./graphics"
	"./matrix"
	net "./network"
	"./render"
)

type canvas struct {
	element js.Value
	width   int
	height  int
}

type app struct {
	on                       bool
	canvas                   *canvas
	call                     js.Func
	gl                       js.Value
	g                        *graphics.RenderSystem
	screen                   *graphics.RenderBuffer
	frameScreen              *graphics.RenderBuffer
	drawImages               *graphics.RenderBuffer
	frame                    *graphics.FrameBuffer
	frame2                   *graphics.FrameBuffer
	frameGeo                 *graphics.FrameBuffer
	camera                   *camera
	state                    *worldState
	world                    *world
	player                   *you
	socket                   js.Value
	socketQueue              [][]byte
	socketSend               map[uint8]interface{}
	canvasOrtho              []float32
	drawOrtho                []float32
	drawPerspective          []float32
	drawInversePerspective   []float32
	drawInverseMv            []float32
	drawPreviousMvp          []float32
	drawCurrentToPreviousMvp []float32
}

func canvasInit() *canvas {
	width := dom.window.Get("innerWidth").Int()
	height := dom.window.Get("innerHeight").Int()

	element := dom.doc.Call("createElement", "canvas")
	style := element.Get("style")

	style.Set("display", "block")
	style.Set("position", "absolute")
	style.Set("left", "0")
	style.Set("right", "0")
	style.Set("top", "0")
	style.Set("bottom", "0")
	style.Set("margin", "auto")

	canvas := &canvas{}
	canvas.element = element
	canvas.update(width, height)
	return canvas
}

func (me *canvas) update(width, height int) {
	me.width = width
	me.height = height
	me.element.Set("width", width)
	me.element.Set("height", height)
}

func (me *app) resize() {
	gl := me.gl
	canvas := me.canvas
	window := dom.window

	width := window.Get("innerWidth").Int()
	height := window.Get("innerHeight").Int()

	canvas.update(width, height)

	const drawFraction = 1.0

	drawWidth := int(math.Floor(float64(width) / drawFraction))
	drawHeight := int(math.Floor(float64(canvas.height) / drawFraction))
	ratio := float64(drawWidth) / float64(drawHeight)
	fov := float32(2.0 * math.Atan(math.Tan(60.0*(math.Pi/180.0)/2.0)/ratio) * (180.0 / math.Pi))

	me.canvasOrtho = make([]float32, 16)
	me.drawOrtho = make([]float32, 16)
	me.drawPerspective = make([]float32, 16)
	me.drawInversePerspective = make([]float32, 16)
	me.drawInverseMv = make([]float32, 16)
	me.drawPreviousMvp = make([]float32, 16)
	me.drawCurrentToPreviousMvp = make([]float32, 16)

	matrix.Orthographic(me.canvasOrtho, 0.0, float32(width), 0.0, float32(height), 0.0, 1.0)
	matrix.Orthographic(me.drawOrtho, 0.0, float32(drawWidth), 0.0, float32(drawHeight), 0.0, 1.0)
	matrix.Perspective(me.drawPerspective, float32(fov), 0.01, 100.0, float32(ratio))
	matrix.Inverse(me.drawInversePerspective, me.drawPerspective)

	if me.frame == nil {
		rgb := []js.Value{graphics.GLxRGB}
		unsignedByte := []js.Value{graphics.GLxUnsignedByte}
		me.frame = graphics.FrameBufferMake(gl, drawWidth, drawHeight, rgb, rgb, unsignedByte, true, true)
		me.frameGeo = graphics.FrameBufferMake(gl, drawWidth, drawHeight, rgb, rgb, unsignedByte, true, true)
		me.frame2 = graphics.FrameBufferMake(gl, drawWidth, drawHeight, rgb, rgb, unsignedByte, true, true)
	} else {
		me.frame.Resize(gl, drawWidth, drawHeight)
		me.frameGeo.Resize(gl, drawWidth, drawHeight)
		me.frame2.Resize(gl, drawWidth, drawHeight)
	}

	me.screen.Zero()
	render.Screen(me.screen, 0.0, 0.0, float32(canvas.width), float32(canvas.height))
	graphics.RenderSystemUpdateVao(gl, me.screen)

	me.frameScreen.Zero()
	render.Screen(me.frameScreen, 0.0, 0.0, float32(drawWidth), float32(drawHeight))
	graphics.RenderSystemUpdateVao(gl, me.frameScreen)
}

func (me *app) switchState(state *worldState) {
	me.state = state
}

func (me *app) init() {
	g := me.g
	gl := me.gl
	world := me.world

	wadRead(g, gl, net.Get("wad"))

	socket := net.Socket("websocket")
	me.socket = socket
	me.socketQueue = make([][]byte, 0)
	me.socketSend = make(map[uint8]interface{})

	socket.Set("binaryType", "arraybuffer")
	socket.Set("onclose", js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		me.on = false
		panic("lost connection to server")
		return nil
	}))

	var raw []byte

	var join sync.WaitGroup
	join.Add(1)
	getworld := js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		defer join.Done()
		data := args[0].Get("data")
		uint8data := js.Global().Get("Uint8Array").New(data)
		raw = make([]byte, uint8data.Get("byteLength").Int())
		typeArray := js.TypedArrayOf(raw)
		typeArray.Call("set", uint8data)
		typeArray.Release()
		return nil
	})
	socket.Set("onmessage", getworld)
	join.Wait()
	getworld.Release()

	world.load(raw)

	socket.Set("onmessage", js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		data := args[0].Get("data")
		uint8data := js.Global().Get("Uint8Array").New(data)
		binary := make([]byte, uint8data.Get("byteLength").Int())
		typeArray := js.TypedArrayOf(binary)
		typeArray.Call("set", uint8data)
		typeArray.Release()
		me.socketQueue = append(me.socketQueue, binary)
		return nil
	}))

	player := world.netLookup[world.pid].(*you)
	me.camera = cameraInit(world, player.thing, 10.0)
	player.camera = me.camera
	player.socket = me.socket
	player.socketSend = me.socketSend
}

func (me *app) run() {
	me.init()
	dom.doc.Set("onkeyup", InputSetKeyUp())
	dom.doc.Set("onkeydown", InputSetKeyDown())
	dom.window.Set("onresize", js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		console("resize!")
		me.resize()
		return nil
	}))
	dom.body.Call("appendChild", me.canvas.element)
	me.resize()
	me.loop()
}

func (me *app) loop() {
	if me.on {
		me.state.update()
		me.state.render()
	}
	js.Global().Call("requestAnimationFrame", me.call)
}

func appInit() *app {
	canvas := canvasInit()

	gl := canvas.element.Call("getContext", "webgl2")
	if gl == js.Undefined() {
		js.Global().Call("alert", "webgl is not supported")
		panic("webgl is not supported")
	}

	graphics.SetupOpenGl(gl)

	gl.Call("clearColor", 0, 0, 0, 1)
	gl.Call("depthFunc", graphics.GLxLEqual)
	gl.Call("cullFace", graphics.GLxBack)
	gl.Call("blendFunc", graphics.GLxSrcAlpha, graphics.GLxOneMinusSrcAlpha)
	gl.Call("disable", graphics.GLxCullFace)
	gl.Call("disable", graphics.GLxBlend)
	gl.Call("disable", graphics.GLxDepthTest)

	g := graphics.RenderSystemInit()

	app := &app{}
	app.on = true
	app.canvas = canvas
	app.g = g
	app.gl = gl
	app.screen = graphics.RenderBufferInit(gl, 2, 0, 0, 4, 6)
	app.frameScreen = graphics.RenderBufferInit(gl, 2, 0, 0, 4, 6)
	app.drawImages = graphics.RenderBufferInit(gl, 2, 0, 2, 400, 600)
	app.state = worldStateInit(app)
	app.world = worldInit(g, gl)
	app.call = js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		app.loop()
		return nil
	})

	return app
}
