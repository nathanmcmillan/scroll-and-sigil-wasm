package render

import "../graphics"

// RendSprite func
func RendSprite(b *graphics.RenderBuffer, x, y, z, sin, cos float32, sprite *Sprite) {
	sine := sprite.HalfWidth * sin
	cosine := sprite.HalfWidth * cos
	pos := b.VertexPos

	b.Vertices[pos] = x - cosine
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z + sine
	b.Vertices[pos+3] = sprite.Left
	b.Vertices[pos+4] = sprite.Bottom

	b.Vertices[pos+5] = x + cosine
	b.Vertices[pos+6] = y
	b.Vertices[pos+7] = z - sine
	b.Vertices[pos+8] = sprite.Right
	b.Vertices[pos+9] = sprite.Bottom

	b.Vertices[pos+10] = x + cosine
	b.Vertices[pos+11] = y + sprite.Height
	b.Vertices[pos+12] = z - sine
	b.Vertices[pos+13] = sprite.Right
	b.Vertices[pos+14] = sprite.Top

	b.Vertices[pos+15] = x - cosine
	b.Vertices[pos+16] = y + sprite.Height
	b.Vertices[pos+17] = z + sine
	b.Vertices[pos+18] = sprite.Left
	b.Vertices[pos+19] = sprite.Top

	b.VertexPos += 20
	Index4(b)
}

// RendMirrorSprite func
func RendMirrorSprite(b *graphics.RenderBuffer, x, y, z, sin, cos float32, sprite *Sprite) {
	sine := sprite.HalfWidth * sin
	cosine := sprite.HalfWidth * cos
	pos := b.VertexPos

	b.Vertices[pos] = x - cosine
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z + sine
	b.Vertices[pos+3] = sprite.Right
	b.Vertices[pos+4] = sprite.Bottom

	b.Vertices[pos+5] = x + cosine
	b.Vertices[pos+6] = y
	b.Vertices[pos+7] = z - sine
	b.Vertices[pos+8] = sprite.Left
	b.Vertices[pos+9] = sprite.Bottom

	b.Vertices[pos+10] = x + cosine
	b.Vertices[pos+11] = y + sprite.Height
	b.Vertices[pos+12] = z - sine
	b.Vertices[pos+13] = sprite.Left
	b.Vertices[pos+14] = sprite.Top

	b.Vertices[pos+15] = x - cosine
	b.Vertices[pos+16] = y + sprite.Height
	b.Vertices[pos+17] = z + sine
	b.Vertices[pos+18] = sprite.Right
	b.Vertices[pos+19] = sprite.Top

	b.VertexPos += 20
	Index4(b)
}
