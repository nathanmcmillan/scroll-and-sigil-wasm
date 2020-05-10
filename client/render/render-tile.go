package render

import "../graphics"

// RendTilePosX func
func RendTilePosX(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x + 1.0
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x + 1.0
	b.Vertices[pos+9] = y + 1.0
	b.Vertices[pos+10] = z
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]

	b.Vertices[pos+16] = x + 1.0
	b.Vertices[pos+17] = y + 1.0
	b.Vertices[pos+18] = z + 1.0
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x + 1.0
	b.Vertices[pos+25] = y
	b.Vertices[pos+26] = z + 1.0
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}

// RendTileNegX func
func RendTileNegX(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x
	b.Vertices[pos+9] = y
	b.Vertices[pos+10] = z + 1.0
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]

	b.Vertices[pos+16] = x
	b.Vertices[pos+17] = y + 1.0
	b.Vertices[pos+18] = z + 1.0
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x
	b.Vertices[pos+25] = y + 1.0
	b.Vertices[pos+26] = z
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}

// RendTilePosY func
func RendTilePosY(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y + 1.0
	b.Vertices[pos+2] = z
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x
	b.Vertices[pos+9] = y + 1.0
	b.Vertices[pos+10] = z + 1.0
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]

	b.Vertices[pos+16] = x + 1.0
	b.Vertices[pos+17] = y + 1.0
	b.Vertices[pos+18] = z + 1.0
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x + 1.0
	b.Vertices[pos+25] = y + 1.0
	b.Vertices[pos+26] = z
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}

// RendTileNegY func
func RendTileNegY(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x + 1.0
	b.Vertices[pos+9] = y
	b.Vertices[pos+10] = z
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]
	b.Vertices[pos+16] = x + 1.0
	b.Vertices[pos+17] = y
	b.Vertices[pos+18] = z + 1.0
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x
	b.Vertices[pos+25] = y
	b.Vertices[pos+26] = z + 1.0
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}

// RendTilePosZ func
func RendTilePosZ(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x + 1.0
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z + 1.0
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x + 1.0
	b.Vertices[pos+9] = y + 1.0
	b.Vertices[pos+10] = z + 1.0
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]

	b.Vertices[pos+16] = x
	b.Vertices[pos+17] = y + 1.0
	b.Vertices[pos+18] = z + 1.0
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x
	b.Vertices[pos+25] = y
	b.Vertices[pos+26] = z + 1.0
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}

// RendTileNegZ func
func RendTileNegZ(b *graphics.RenderBuffer, x, y, z float32, texture *[4]float32, rgb *[4][3]float32) {
	pos := b.VertexPos
	b.Vertices[pos] = x
	b.Vertices[pos+1] = y
	b.Vertices[pos+2] = z
	b.Vertices[pos+3] = rgb[0][0]
	b.Vertices[pos+4] = rgb[0][1]
	b.Vertices[pos+5] = rgb[0][2]
	b.Vertices[pos+6] = texture[0]
	b.Vertices[pos+7] = texture[1]

	b.Vertices[pos+8] = x
	b.Vertices[pos+9] = y + 1.0
	b.Vertices[pos+10] = z
	b.Vertices[pos+11] = rgb[1][0]
	b.Vertices[pos+12] = rgb[1][1]
	b.Vertices[pos+13] = rgb[1][2]
	b.Vertices[pos+14] = texture[2]
	b.Vertices[pos+15] = texture[1]

	b.Vertices[pos+16] = x + 1.0
	b.Vertices[pos+17] = y + 1.0
	b.Vertices[pos+18] = z
	b.Vertices[pos+19] = rgb[2][0]
	b.Vertices[pos+20] = rgb[2][1]
	b.Vertices[pos+21] = rgb[2][2]
	b.Vertices[pos+22] = texture[2]
	b.Vertices[pos+23] = texture[3]

	b.Vertices[pos+24] = x + 1.0
	b.Vertices[pos+25] = y
	b.Vertices[pos+26] = z
	b.Vertices[pos+27] = rgb[3][0]
	b.Vertices[pos+28] = rgb[3][1]
	b.Vertices[pos+29] = rgb[3][2]
	b.Vertices[pos+30] = texture[0]
	b.Vertices[pos+31] = texture[3]

	b.VertexPos += 32

	if Lumin(&rgb[0])+Lumin(&rgb[2]) < Lumin(&rgb[1])+Lumin(&rgb[3]) {
		MirrorIndex4(b)
	} else {
		Index4(b)
	}
}
