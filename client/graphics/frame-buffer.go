package graphics

import (
	"syscall/js"
)

// FrameBuffer struct
type FrameBuffer struct {
	Fbo            js.Value
	internalFormat []js.Value
	format         []js.Value
	typeOf         []js.Value
	Width          int
	Height         int
	linear         bool
	depth          bool
	DepthTexture   js.Value
	Textures       []js.Value
	drawBuffers    []interface{}
}

// FrameBufferInit func
func FrameBufferInit(width, height int, internalFormat, format, typeOf []js.Value, linear, depth bool) *FrameBuffer {
	if len(internalFormat) != len(format) || len(internalFormat) != len(typeOf) {
		panic("frame format lengths differ")
	}
	f := &FrameBuffer{}
	f.internalFormat = internalFormat
	f.format = format
	f.typeOf = typeOf
	f.Width = width
	f.Height = height
	f.linear = linear
	f.depth = depth
	f.Textures = make([]js.Value, len(format))
	f.drawBuffers = make([]interface{}, len(format))
	return f
}

// FrameBufferMake func
func FrameBufferMake(gl js.Value, width, height int, internalFormat, format, typeOf []js.Value, linear, depth bool) *FrameBuffer {
	f := FrameBufferInit(width, height, internalFormat, format, typeOf, linear, depth)
	RenderSystemMakeFrameBuffer(gl, f)
	return f
}

// Resize func
func (me *FrameBuffer) Resize(gl js.Value, width, height int) {
	me.Width = width
	me.Height = height
	RenderSystemSetFrameBuffer(gl, me.Fbo)
	RenderSystemUpdateFrameBuffer(gl, me)
}
