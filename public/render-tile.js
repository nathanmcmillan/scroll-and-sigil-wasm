class RenderTile {
	static Side(buffer, side, x, y, z, texture, rgb_a, rgb_b, rgb_c, rgb_d) {
		switch (side) {
			case WorldPositiveX:
				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]
				break
			case WorldNegativeX:
				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]
				break
			case WorldPositiveY:
				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]
				break
			case WorldNegativeY:
				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]
				break
			case WorldPositiveZ:
				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z + 1.0
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]
				break
			case WorldNegativeZ:
				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_a[0]
				buffer.vertices[buffer.vertexPos++] = rgb_a[1]
				buffer.vertices[buffer.vertexPos++] = rgb_a[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[1]

				buffer.vertices[buffer.vertexPos++] = x
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_b[0]
				buffer.vertices[buffer.vertexPos++] = rgb_b[1]
				buffer.vertices[buffer.vertexPos++] = rgb_b[2]
				buffer.vertices[buffer.vertexPos++] = texture[0]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y + 1.0
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_c[0]
				buffer.vertices[buffer.vertexPos++] = rgb_c[1]
				buffer.vertices[buffer.vertexPos++] = rgb_c[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[3]

				buffer.vertices[buffer.vertexPos++] = x + 1.0
				buffer.vertices[buffer.vertexPos++] = y
				buffer.vertices[buffer.vertexPos++] = z
				buffer.vertices[buffer.vertexPos++] = rgb_d[0]
				buffer.vertices[buffer.vertexPos++] = rgb_d[1]
				buffer.vertices[buffer.vertexPos++] = rgb_d[2]
				buffer.vertices[buffer.vertexPos++] = texture[2]
				buffer.vertices[buffer.vertexPos++] = texture[1]
				break
		}

		if (Render.Lumin(rgb_a) + Render.Lumin(rgb_c) < Render.Lumin(rgb_b) + Render.Lumin(rgb_d))
			Render.MirrorIndex4(buffer)
		else
			Render.Index4(buffer)
	}
}
