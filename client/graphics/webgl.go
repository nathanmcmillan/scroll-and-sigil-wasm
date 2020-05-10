package graphics

import (
	"syscall/js"
)

// Webgl
var (
	GLxStaticDraw             js.Value
	GLxArrayBuffer            js.Value
	GLxElementArrayBuffer     js.Value
	GLxVertexShader           js.Value
	GLxFragmentShader         js.Value
	GLxFloat                  js.Value
	GLxDepthTest              js.Value
	GLxColorBufferBit         js.Value
	GLxDepthBufferBit         js.Value
	GLxTriangles              js.Value
	GLxUnsignedShort          js.Value
	GLxLEqual                 js.Value
	GLxBack                   js.Value
	GLxSrcAlpha               js.Value
	GLxOneMinusSrcAlpha       js.Value
	GLxCullFace               js.Value
	GLxBlend                  js.Value
	GLxTexture0               js.Value
	GLxTexture1               js.Value
	GLxTexture2D              js.Value
	GLxTextureWrapS           js.Value
	GLxTextureWrapT           js.Value
	GLxClampToEdge            js.Value
	GLxRepeat                 js.Value
	GLxLinear                 js.Value
	GLxNearest                js.Value
	GLxColorAttachment        [3]js.Value
	GLxTextureMinFilter       js.Value
	GLxTextureMagFilter       js.Value
	GLxDepthStencilAttachment js.Value
	GLxRGB                    js.Value
	GLxRGBA                   js.Value
	GLxUnsignedByte           js.Value
	GLxLinkStatus             js.Value
	GLxCompileStatus          js.Value
	GLxFrameBufferComplete    js.Value
	GLxFrameBuffer            js.Value
	GLxUnsignedInt            js.Value
	GLxDynamicDraw            js.Value
	GLxDepthStencil           js.Value
	GLxDepth24Stencil8        js.Value
	GLxUnsignedInt24x8        js.Value
)

// SetupOpenGl func
func SetupOpenGl(gl js.Value) {
	GLxStaticDraw = gl.Get("STATIC_DRAW")
	GLxArrayBuffer = gl.Get("ARRAY_BUFFER")
	GLxElementArrayBuffer = gl.Get("ELEMENT_ARRAY_BUFFER")
	GLxVertexShader = gl.Get("VERTEX_SHADER")
	GLxFragmentShader = gl.Get("FRAGMENT_SHADER")
	GLxFloat = gl.Get("FLOAT")
	GLxDepthTest = gl.Get("DEPTH_TEST")
	GLxColorBufferBit = gl.Get("COLOR_BUFFER_BIT")
	GLxDepthBufferBit = gl.Get("DEPTH_BUFFER_BIT")
	GLxTriangles = gl.Get("TRIANGLES")
	GLxUnsignedShort = gl.Get("UNSIGNED_SHORT")
	GLxLEqual = gl.Get("LEQUAL")
	GLxBack = gl.Get("BACK")
	GLxSrcAlpha = gl.Get("SRC_ALPHA")
	GLxOneMinusSrcAlpha = gl.Get("ONE_MINUS_SRC_ALPHA")
	GLxCullFace = gl.Get("CULL_FACE")
	GLxBlend = gl.Get("BLEND")
	GLxTexture0 = gl.Get("TEXTURE0")
	GLxTexture1 = gl.Get("TEXTURE1")
	GLxTexture2D = gl.Get("TEXTURE_2D")
	GLxTextureWrapS = gl.Get("TEXTURE_WRAP_S")
	GLxTextureWrapT = gl.Get("TEXTURE_WRAP_T")
	GLxClampToEdge = gl.Get("CLAMP_TO_EDGE")
	GLxRepeat = gl.Get("REPEAT")
	GLxLinear = gl.Get("LINEAR")
	GLxNearest = gl.Get("NEAREST")
	GLxTextureMinFilter = gl.Get("TEXTURE_MIN_FILTER")
	GLxTextureMagFilter = gl.Get("TEXTURE_MAG_FILTER")
	GLxDepthStencilAttachment = gl.Get("DEPTH_STENCIL_ATTACHMENT")
	GLxRGB = gl.Get("RGB")
	GLxRGBA = gl.Get("RGBA")
	GLxUnsignedByte = gl.Get("UNSIGNED_BYTE")
	GLxLinkStatus = gl.Get("LINK_STATUS")
	GLxCompileStatus = gl.Get("COMPILE_STATUS")
	GLxFrameBufferComplete = gl.Get("FRAMEBUFFER_COMPLETE")
	GLxFrameBuffer = gl.Get("FRAMEBUFFER")
	GLxUnsignedInt = gl.Get("UNSIGNED_INT")
	GLxDynamicDraw = gl.Get("DYNAMIC_DRAW")
	GLxDepthStencil = gl.Get("DEPTH_STENCIL")
	GLxDepth24Stencil8 = gl.Get("DEPTH24_STENCIL8")
	GLxUnsignedInt24x8 = gl.Get("UNSIGNED_INT_24_8")

	GLxColorAttachment[0] = gl.Get("COLOR_ATTACHMENT0")
	GLxColorAttachment[1] = gl.Get("COLOR_ATTACHMENT1")
	GLxColorAttachment[2] = gl.Get("COLOR_ATTACHMENT2")
}
