package render

import (
	"strings"

	"../graphics"
)

// Constants
const (
	Font       = "0123456789abcdefghijklmnopqrstuvwxyz%"
	FontWidth  = 9
	FontHeight = 9
	FontGrid   = 7
	FontColumn = float32(FontWidth) / 64.0
	FontRow    = float32(FontHeight) / 64.0
)

// Lumin func
func Lumin(rgb *[3]float32) float32 {
	return 0.2126*rgb[0] + 0.7152*rgb[1] + 0.0722*rgb[2]
}

// PackRgb func
func PackRgb(red, green, blue int32) int32 {
	return (red << 16) | (green << 8) | blue
}

// UnpackRgb func
func UnpackRgb(rgb int32) (int32, int32, int32) {
	red := (rgb >> 16) & 255
	green := (rgb >> 8) & 255
	blue := rgb & 255
	return red, green, blue
}

// Index4 func
func Index4(b *graphics.RenderBuffer) {
	pos := b.IndexPos
	offset := b.IndexOffset
	b.Indices[pos] = offset
	b.Indices[pos+1] = offset + 1
	b.Indices[pos+2] = offset + 2
	b.Indices[pos+3] = offset + 2
	b.Indices[pos+4] = offset + 3
	b.Indices[pos+5] = offset
	b.IndexPos += 6
	b.IndexOffset += 4
}

// MirrorIndex4 func
func MirrorIndex4(b *graphics.RenderBuffer) {
	pos := b.IndexPos
	offset := b.IndexOffset
	b.Indices[pos] = offset + 1
	b.Indices[pos+1] = offset + 2
	b.Indices[pos+2] = offset + 3
	b.Indices[pos+3] = offset + 3
	b.Indices[pos+4] = offset
	b.Indices[pos+5] = offset + 1
	b.IndexPos += 6
	b.IndexOffset += 4
}

// Screen func
func Screen(b *graphics.RenderBuffer, x, y, width, height float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y

	b.Vertices[pos+2] = x + width
	b.Vertices[pos+3] = y

	b.Vertices[pos+4] = x + width
	b.Vertices[pos+5] = y + height

	b.Vertices[pos+6] = x
	b.Vertices[pos+7] = y + height

	b.VertexPos += 8
	Index4(b)
}

// Image func
func Image(b *graphics.RenderBuffer, x, y, width, height, left, top, right, bottom float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = left
	b.Vertices[pos+3] = bottom

	b.Vertices[pos+4] = x + width
	b.Vertices[pos+5] = y
	b.Vertices[pos+6] = right
	b.Vertices[pos+7] = bottom

	b.Vertices[pos+8] = x + width
	b.Vertices[pos+9] = y + height
	b.Vertices[pos+10] = right
	b.Vertices[pos+11] = top

	b.Vertices[pos+12] = x
	b.Vertices[pos+13] = y + height
	b.Vertices[pos+14] = left
	b.Vertices[pos+15] = top

	b.VertexPos += 16
	Index4(b)
}

// Rectangle func
func Rectangle(b *graphics.RenderBuffer, x, y, width, height, red, green, blue float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = red
	b.Vertices[pos+3] = green
	b.Vertices[pos+4] = blue

	b.Vertices[pos+5] = x + width
	b.Vertices[pos+6] = y
	b.Vertices[pos+7] = red
	b.Vertices[pos+8] = green
	b.Vertices[pos+9] = blue

	b.Vertices[pos+10] = x + width
	b.Vertices[pos+11] = y + height
	b.Vertices[pos+12] = red
	b.Vertices[pos+13] = green
	b.Vertices[pos+14] = blue

	b.Vertices[pos+15] = x
	b.Vertices[pos+16] = y + height
	b.Vertices[pos+17] = red
	b.Vertices[pos+18] = green
	b.Vertices[pos+19] = blue

	b.VertexPos += 20
	Index4(b)
}

// Print func
func Print(b *graphics.RenderBuffer, x, y float32, text string, scale int) {
	widthScale := float32(FontWidth * scale)
	heightScale := float32(FontHeight * scale)
	xx := x
	yy := y
	num := len(text)
	for i := 0; i < num; i++ {
		c := text[i]
		if c == ' ' {
			xx += widthScale
			continue
		} else if c == '\n' {
			xx = x
			yy += heightScale
			continue
		}
		loc := strings.IndexByte(Font, c)
		tx1 := float32(loc%FontGrid) * FontColumn
		ty1 := float32(loc/FontGrid) * FontRow
		tx2 := tx1 + FontColumn
		ty2 := ty1 + FontRow
		Image(b, xx, yy, widthScale, heightScale, tx1, ty1, tx2, ty2)
		xx += widthScale
	}
}
