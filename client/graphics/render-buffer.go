package graphics

import (
	"syscall/js"
)

// RenderBuffer struct
type RenderBuffer struct {
	position    int
	color       int
	texture     int
	vao         js.Value
	vbo         js.Value
	ebo         js.Value
	VertexPos   int
	IndexPos    int
	IndexOffset uint32
	Vertices    []float32
	Indices     []uint32
	verticesJs  js.TypedArray
	indicesJs   js.TypedArray
}

// RenderBufferInit func
func RenderBufferInit(gl js.Value, position, color, texture, vertexLimit, indexLimit int) *RenderBuffer {
	b := &RenderBuffer{}
	b.position = position
	b.color = color
	b.texture = texture
	b.VertexPos = 0
	b.IndexPos = 0
	b.IndexOffset = 0
	b.Vertices = make([]float32, vertexLimit*(position+color+texture))
	b.Indices = make([]uint32, indexLimit)
	b.verticesJs = js.TypedArrayOf(b.Vertices)
	b.indicesJs = js.TypedArrayOf(b.Indices)
	RenderSystemMakeVao(gl, b, position, color, texture)
	return b
}

// InitCopy func
func (me *RenderBuffer) InitCopy(gl js.Value) *RenderBuffer {
	b := &RenderBuffer{}
	b.Vertices = make([]float32, me.VertexPos)
	b.Indices = make([]uint32, me.IndexPos)
	b.verticesJs = js.TypedArrayOf(b.Vertices)
	b.indicesJs = js.TypedArrayOf(b.Indices)
	me.CopyTo(b)
	RenderSystemMakeVao(gl, b, me.position, me.color, me.texture)
	RenderSystemUpdateVao(gl, b)
	return b
}

// CopyTo func
func (me *RenderBuffer) CopyTo(to *RenderBuffer) {
	for i := 0; i < me.VertexPos; i++ {
		to.Vertices[i] = me.Vertices[i]
	}
	for i := 0; i < me.IndexPos; i++ {
		to.Indices[i] = me.Indices[i]
	}
	to.VertexPos = me.VertexPos
	to.IndexPos = me.IndexPos
	to.IndexOffset = me.IndexOffset
}

// Zero func
func (me *RenderBuffer) Zero() {
	me.VertexPos = 0
	me.IndexPos = 0
	me.IndexOffset = 0
}

// RenderBufferExpand func
func (me *RenderBuffer) RenderBufferExpand(gl js.Value) {
	Vertices := me.Vertices
	Indices := me.Indices

	me.Vertices = make([]float32, len(me.Vertices)*2)
	me.Indices = make([]uint32, len(me.Indices)*2)

	me.verticesJs.Release()
	me.indicesJs.Release()

	me.verticesJs = js.TypedArrayOf(me.Vertices)
	me.indicesJs = js.TypedArrayOf(me.Indices)

	for i := 0; i < me.VertexPos; i++ {
		me.Vertices[i] = Vertices[i]
	}

	for i := 0; i < me.IndexPos; i++ {
		me.Indices[i] = Indices[i]
	}

	RenderSystemUpdateVao(gl, me)
}

// RenderCopyInit func
func RenderCopyInit(position, color, texture, vertexLimit, indexLimit int) *RenderBuffer {
	b := &RenderBuffer{}
	b.position = position
	b.color = color
	b.texture = texture
	b.VertexPos = 0
	b.IndexPos = 0
	b.IndexOffset = 0
	b.Vertices = make([]float32, vertexLimit*(position+color+texture))
	b.Indices = make([]uint32, indexLimit)
	return b
}
