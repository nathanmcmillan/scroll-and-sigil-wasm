class EditApp {
    constructor(manager) {
        this.manager = manager

        let canvas = document.createElement("canvas")
        canvas.style.display = "block"
        canvas.style.position = "absolute"
        canvas.style.left = "0"
        canvas.style.right = "0"
        canvas.style.top = "0"
        canvas.style.bottom = "0"
        canvas.style.margin = "auto"
        canvas.width = window.innerWidth
        canvas.height = window.innerHeight

        let gl = canvas.getContext("webgl2")
        let g = new RenderSystem()

        gl.clearColor(0, 0, 0, 1)
        gl.depthFunc(gl.LEQUAL)
        gl.cullFace(gl.BACK)
        gl.blendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
        gl.disable(gl.CULL_FACE)
        gl.disable(gl.BLEND)
        gl.disable(gl.DEPTH_TEST)

        this.on = true
        this.canvas = canvas
        this.gl = gl
        this.g = g
        this.screen = RenderBuffer.Init(gl, 2, 0, 2, 4, 6)
        this.drawImages = RenderBuffer.Init(gl, 2, 0, 2, 400, 600)
        this.singleBlock = RenderBuffer.Init(gl, 3, 3, 2, 24, 36)
        this.world = new WorldEdit(g, gl)
        this.frame = null
        this.camera = null
        this.state = new EditState(this)

        document.onkeyup = Input.SetKeyUp
        document.onkeydown = Input.SetKeyDown
        document.onmouseup = Input.SetMouseUp
        document.onmousedown = Input.SetMouseDown
        document.onmousemove = Input.SetMouseMove

        let self = this

        window.onblur = function () {
            self.on = false
        }

        window.onfocus = function () {
            self.on = true
        }
    }
    resize() {
        let gl = this.gl
        let canvas = this.canvas
        let screen = this.screen

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight

        let canvasOrtho = []
        let drawOrtho = []
        let drawPerspective = []

        let scale = 1.0
        let drawWidth = canvas.width * scale
        let drawHeight = canvas.height * scale
        let ratio = drawWidth / drawHeight
        let fov = 2 * Math.atan(Math.tan(60 * (Math.PI / 180) / 2) / ratio) * (180 / Math.PI)

        Matrix.Orthographic(canvasOrtho, 0.0, canvas.width, 0.0, canvas.height, 0.0, 1.0)
        Matrix.Orthographic(drawOrtho, 0.0, drawWidth, 0.0, drawHeight, 0.0, 1.0)
        Matrix.Perspective(drawPerspective, fov, 0.01, 100.0, ratio)

        if (this.frame === null) {
            this.frame = FrameBuffer.Make(gl, drawWidth, drawHeight, [gl.RGB], [gl.RGB], [gl.UNSIGNED_BYTE], "nearest", "depth")
        } else {
            this.frame.Resize(gl, drawWidth, drawHeight)
        }

        screen.Zero()
        Render.Image(screen, 0, 0, canvas.width, canvas.height, 0.0, 1.0, 1.0, 0.0)
        RenderSystem.UpdateVao(gl, screen)

        this.canvasOrtho = canvasOrtho
        this.drawOrtho = drawOrtho
        this.drawPerspective = drawPerspective
    }
    async init() {
        let g = this.g
        let gl = this.gl

        let data = await Net.Request("wad")
        await Wad.Load(g, gl, data)
        await this.manager.init(this)
    }
    async run() {
        await this.init()
        let self = this
        window.onresize = function () {
            self.resize()
        }
        document.body.appendChild(this.canvas)
        this.resize()
        this.loop()
    }
    switch (state) {
        this.state = state
    }
    loop() {
        if (this.on) {
            this.state.update()
            this.state.render()
        }
        requestAnimationFrame(loop)
    }
}

function PlaySound(name) {
    let sound = Sounds[name]
    sound.pause()
    sound.volume = 0.25
    sound.currentTime = 0
    let promise = sound.play()
    if (promise) promise.then(_ => {}).catch(_ => {})
}

class Editor {
    constructor() {}
    async init(app) {
        let raw = await Net.Request("test.map")
        app.world.LoadRaw(raw)
        app.camera = new SimpleCamera(10.0, 10.0, 10.0)
    }
}

let editor = new Editor()
let app = new EditApp(editor)

app.run()

function loop() {
    app.loop()
}
