package graphics

import (
	"strings"
	"sync"
	"syscall/js"

	"../matrix"
	net "../network"
)

// Texture struct
type Texture struct {
	TextureID js.Value
	Image     js.Value
	Width     int
	Height    int
}

// RenderSystem struct
type RenderSystem struct {
	View                  []float32
	ModelView             []float32
	ModelViewProject      []float32
	ModelViewProjectTyped js.TypedArray
	programID             js.Value
	programName           string
	mvpIds                map[string]js.Value
	uniforms              map[string]map[string]js.Value
	textureIds            map[string]js.Value
	shaders               map[string]js.Value
	Textures              map[string]*Texture
}

// RenderSystemInit func
func RenderSystemInit() *RenderSystem {
	g := &RenderSystem{}
	g.View = make([]float32, 16)
	g.ModelView = make([]float32, 16)
	g.ModelViewProject = make([]float32, 16)
	g.ModelViewProjectTyped = js.TypedArrayOf(g.ModelViewProject)
	g.mvpIds = make(map[string]js.Value)
	g.uniforms = make(map[string]map[string]js.Value)
	g.textureIds = make(map[string]js.Value)
	g.shaders = make(map[string]js.Value)
	g.Textures = make(map[string]*Texture)
	return g
}

// SetTexture func
func (me *RenderSystem) SetTexture(gl js.Value, name string) {
	gl.Call("activeTexture", GLxTexture0)
	gl.Call("bindTexture", GLxTexture2D, me.Textures[name].TextureID)
	gl.Call("uniform1i", me.textureIds[me.programName], 0)
}

// SetTextureDirect func
func (me *RenderSystem) SetTextureDirect(gl js.Value, textureID js.Value) {
	gl.Call("activeTexture", GLxTexture0)
	gl.Call("bindTexture", GLxTexture2D, textureID)
	gl.Call("uniform1i", me.textureIds[me.programName], 0)
}

// SetIndexTextureDirect func
func (me *RenderSystem) SetIndexTextureDirect(gl js.Value, textureUnit js.Value, unitIndex int, name string, textureID js.Value) {
	gl.Call("activeTexture", textureUnit)
	gl.Call("bindTexture", GLxTexture2D, textureID)
	loc, ok := me.uniforms[me.programName][name]
	if !ok {
		loc = gl.Call("getUniformLocation", me.programID, name)
		me.uniforms[me.programName][name] = loc
	}
	gl.Call("uniform1i", loc, unitIndex)
}

// SetProgram func
func (me *RenderSystem) SetProgram(gl js.Value, name string) {
	me.programID = me.shaders[name]
	me.programName = name
	gl.Call("useProgram", me.programID)
}

// SetUniformVec2 func
func (me *RenderSystem) SetUniformVec2(gl js.Value, name string, x, y float32) {
	loc, ok := me.uniforms[me.programName][name]
	if !ok {
		loc = gl.Call("getUniformLocation", me.programID, name)
		me.uniforms[me.programName][name] = loc
	}
	gl.Call("uniform2f", loc, x, y)
}

// SetUniformMatrix4 func
func (me *RenderSystem) SetUniformMatrix4(gl js.Value, name string, matrix []float32) {
	loc, ok := me.uniforms[me.programName][name]
	if !ok {
		loc = gl.Call("getUniformLocation", me.programID, name)
		me.uniforms[me.programName][name] = loc
	}
	gl.Call("uniformMatrix4fv", loc, false, matrix)
}

// MakeProgram func
func (me *RenderSystem) MakeProgram(gl js.Value, name string) {
	code := net.Get("shaders/" + name + ".glsl")
	parts := strings.Split(code, "===========================================================")
	vertex := parts[0]
	fragment := strings.TrimSpace(parts[1])
	program := RenderSystemCompileProgram(gl, vertex, fragment)
	me.shaders[name] = program
	me.uniforms[name] = make(map[string]js.Value)
	me.mvpIds[name] = gl.Call("getUniformLocation", program, "u_mvp")
	me.textureIds[name] = gl.Call("getUniformLocation", program, "u_texture0")
}

// SetOrthographic func
func (me *RenderSystem) SetOrthographic(orthographic []float32, x, y float32) {
	matrix.Identity(me.ModelView)
	matrix.Translate(me.ModelView, x, y, -1)
	matrix.Multiply(me.ModelViewProject, orthographic, me.ModelView)
}

// SetPerspective func
func (me *RenderSystem) SetPerspective(perspective []float32, x, y, z, rx, ry float32) {
	matrix.Identity(me.View)
	matrix.RotateX(me.View, rx)
	matrix.RotateY(me.View, ry)
	matrix.TranslateFromView(me.ModelView, me.View, x, y, z)
	matrix.Multiply(me.ModelViewProject, perspective, me.ModelView)
}

// MakeImage func
func (me *RenderSystem) MakeImage(gl js.Value, name string, wrap js.Value) {
	textureID := gl.Call("createTexture")

	var join sync.WaitGroup
	join.Add(1)
	image := js.Global().Get("Image").New()
	var call js.Func
	call = js.FuncOf(func(self js.Value, args []js.Value) interface{} {
		defer join.Done()
		defer call.Release()
		return nil
	})
	image.Set("onload", call)
	image.Set("src", "images/"+name+".png")
	join.Wait()

	texture := &Texture{}
	texture.TextureID = textureID
	texture.Image = image
	texture.Width = image.Get("width").Int()
	texture.Height = image.Get("height").Int()

	gl.Call("bindTexture", GLxTexture2D, textureID)
	gl.Call("texImage2D", GLxTexture2D, 0, GLxRGBA, GLxRGBA, GLxUnsignedByte, texture.Image)
	gl.Call("texParameteri", GLxTexture2D, GLxTextureMinFilter, GLxNearest)
	gl.Call("texParameteri", GLxTexture2D, GLxTextureMagFilter, GLxNearest)
	gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapS, wrap)
	gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapT, wrap)
	gl.Call("bindTexture", GLxTexture2D, nil)

	me.Textures[name] = texture
}

// UpdateMvp func
func (me *RenderSystem) UpdateMvp(gl js.Value) {
	gl.Call("uniformMatrix4fv", me.mvpIds[me.programName], false, me.ModelViewProjectTyped)
}

// RenderSystemSetFrameBuffer func
func RenderSystemSetFrameBuffer(gl js.Value, fbo js.Value) {
	gl.Call("bindFramebuffer", GLxFrameBuffer, fbo)
}

// RenderSystemSetView func
func RenderSystemSetView(gl js.Value, x, y, width, height int) {
	gl.Call("viewport", x, y, width, height)
	gl.Call("scissor", x, y, width, height)
}

// RenderSystemBindVao func
func RenderSystemBindVao(gl js.Value, buffer *RenderBuffer) {
	gl.Call("bindVertexArray", buffer.vao)
}

// RenderSystemUpdateVao func
func RenderSystemUpdateVao(gl js.Value, buffer *RenderBuffer) {
	gl.Call("bindVertexArray", buffer.vao)
	gl.Call("bindBuffer", GLxArrayBuffer, buffer.vbo)
	gl.Call("bufferData", GLxArrayBuffer, buffer.verticesJs, GLxStaticDraw)
	gl.Call("bindBuffer", GLxElementArrayBuffer, buffer.ebo)
	gl.Call("bufferData", GLxElementArrayBuffer, buffer.indicesJs, GLxStaticDraw)
}

// RenderSystemBindAndDraw func
func RenderSystemBindAndDraw(gl js.Value, buffer *RenderBuffer) {
	gl.Call("bindVertexArray", buffer.vao)
	gl.Call("drawElements", GLxTriangles, buffer.IndexPos, GLxUnsignedInt, 0)
}

// RenderSystemDrawRange func
func RenderSystemDrawRange(gl js.Value, start, count int) {
	gl.Call("drawElements", GLxTriangles, count, GLxUnsignedInt, start)
}

// RenderSystemUpdateAndDraw func
func RenderSystemUpdateAndDraw(gl js.Value, buffer *RenderBuffer) {
	if buffer.VertexPos == 0 {
		return
	}
	gl.Call("bindVertexArray", buffer.vao)
	gl.Call("bindBuffer", GLxArrayBuffer, buffer.vbo)
	gl.Call("bufferData", GLxArrayBuffer, buffer.verticesJs, GLxDynamicDraw)
	gl.Call("bindBuffer", GLxElementArrayBuffer, buffer.ebo)
	gl.Call("bufferData", GLxElementArrayBuffer, buffer.indicesJs, GLxDynamicDraw)
	gl.Call("drawElements", GLxTriangles, buffer.IndexPos, GLxUnsignedInt, 0)
}

// RenderSystemMakeVao func
func RenderSystemMakeVao(gl js.Value, buffer *RenderBuffer, position, color, texture int) {
	buffer.vao = gl.Call("createVertexArray")
	buffer.vbo = gl.Call("createBuffer")
	buffer.ebo = gl.Call("createBuffer")
	gl.Call("bindVertexArray", buffer.vao)
	gl.Call("bindBuffer", GLxArrayBuffer, buffer.vbo)
	gl.Call("bindBuffer", GLxElementArrayBuffer, buffer.ebo)

	stride := (position + color + texture) * 4
	index := 0
	offset := 0
	if position > 0 {
		gl.Call("vertexAttribPointer", index, position, GLxFloat, false, stride, offset)
		gl.Call("enableVertexAttribArray", index)
		index++
		offset += position * 4
	}
	if color > 0 {
		gl.Call("vertexAttribPointer", index, color, GLxFloat, false, stride, offset)
		gl.Call("enableVertexAttribArray", index)
		index++
		offset += color * 4
	}
	if texture > 0 {
		gl.Call("vertexAttribPointer", index, texture, GLxFloat, false, stride, offset)
		gl.Call("enableVertexAttribArray", index)
	}
}

// RenderSystemUpdateFrameBuffer func
func RenderSystemUpdateFrameBuffer(gl js.Value, frame *FrameBuffer) {
	for i := 0; i < len(frame.format); i++ {
		gl.Call("bindTexture", GLxTexture2D, frame.Textures[i])
		gl.Call("texImage2D", GLxTexture2D, 0, frame.internalFormat[i], frame.Width, frame.Height, 0, frame.format[i], frame.typeOf[i], nil)
	}
	if frame.depth {
		gl.Call("bindTexture", GLxTexture2D, frame.DepthTexture)
		gl.Call("texImage2D", GLxTexture2D, 0, GLxDepth24Stencil8, frame.Width, frame.Height, 0, GLxDepthStencil, GLxUnsignedInt24x8, nil)
	}
}

// RenderSystemTextureFrameBuffer func
func RenderSystemTextureFrameBuffer(gl js.Value, frame *FrameBuffer) {
	for i := 0; i < len(frame.format); i++ {
		texture := gl.Call("createTexture")
		gl.Call("bindTexture", GLxTexture2D, texture)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapS, GLxClampToEdge)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapT, GLxClampToEdge)
		if frame.linear {
			gl.Call("texParameteri", GLxTexture2D, GLxTextureMinFilter, GLxLinear)
			gl.Call("texParameteri", GLxTexture2D, GLxTextureMagFilter, GLxLinear)
		} else {
			gl.Call("texParameteri", GLxTexture2D, GLxTextureMinFilter, GLxNearest)
			gl.Call("texParameteri", GLxTexture2D, GLxTextureMagFilter, GLxNearest)
		}
		color := GLxColorAttachment[i]
		gl.Call("framebufferTexture2D", GLxFrameBuffer, color, GLxTexture2D, texture, 0)
		frame.Textures[i] = texture
		frame.drawBuffers[i] = color
	}
	gl.Call("drawBuffers", frame.drawBuffers)
	if frame.depth {
		frame.DepthTexture = gl.Call("createTexture")
		gl.Call("bindTexture", GLxTexture2D, frame.DepthTexture)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapS, GLxClampToEdge)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureWrapT, GLxClampToEdge)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureMinFilter, GLxNearest)
		gl.Call("texParameteri", GLxTexture2D, GLxTextureMagFilter, GLxNearest)
		gl.Call("framebufferTexture2D", GLxFrameBuffer, GLxDepthStencilAttachment, GLxTexture2D, frame.DepthTexture, 0)
	}
	RenderSystemUpdateFrameBuffer(gl, frame)

	status := gl.Call("checkFramebufferStatus", GLxFrameBuffer)
	if status != GLxFrameBufferComplete {
		js.Global().Get("console").Call("error", "framebuffer error:", status)
	}
}

// RenderSystemMakeFrameBuffer func
func RenderSystemMakeFrameBuffer(gl js.Value, frame *FrameBuffer) {
	frame.Fbo = gl.Call("createFramebuffer")
	gl.Call("bindFramebuffer", GLxFrameBuffer, frame.Fbo)
	RenderSystemTextureFrameBuffer(gl, frame)
}

// RenderSystemCompileProgram func
func RenderSystemCompileProgram(gl js.Value, v, f string) js.Value {
	vert := RenderSystemCompileShader(gl, v, GLxVertexShader)
	frag := RenderSystemCompileShader(gl, f, GLxFragmentShader)
	program := gl.Call("createProgram")
	gl.Call("attachShader", program, vert)
	gl.Call("attachShader", program, frag)
	gl.Call("linkProgram", program)
	if !gl.Call("getProgramParameter", program, GLxLinkStatus).Truthy() {
		js.Global().Get("console").Call("error", v+", "+f)
		js.Global().Get("console").Call("error", gl.Call("getProgramInfoLog", program))
	}
	return program
}

// RenderSystemCompileShader func
func RenderSystemCompileShader(gl js.Value, source string, typeOf js.Value) js.Value {
	shader := gl.Call("createShader", typeOf)
	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)
	if !gl.Call("getShaderParameter", shader, GLxCompileStatus).Truthy() {
		js.Global().Get("console").Call("error", source)
		js.Global().Get("console").Call("error", gl.Call("getShaderInfoLog", shader))
	}
	return shader
}
