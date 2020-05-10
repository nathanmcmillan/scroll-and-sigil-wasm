const Font = "0123456789abcdefghijklmnopqrstuvwxyz%"
const FontWidth = 9
const FontHeight = 9
const FontGrid = Math.floor(64.0 / FontWidth)
const FontColumn = FontWidth / 64.0
const FontRow = FontHeight / 64.0

class Render {
    static Lumin(rgb) {
        return 0.2126 * rgb[0] + 0.7152 * rgb[1] + 0.0722 * rgb[2]
    }
    static PackRgb(red, green, blue) {
        return (red << 16) | (green << 8) | blue
    }
    static UnpackRgb(rgb) {
        let red = (rgb >> 16) & 255
        let green = (rgb >> 8) & 255
        let blue = rgb & 255
        return [red, green, blue]
    }
    static Index4(buffer) {
        buffer.indices[buffer.index_pos++] = buffer.index_offset
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 1
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 2
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 2
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 3
        buffer.indices[buffer.index_pos++] = buffer.index_offset
        buffer.index_offset += 4
    }
    static MirrorIndex4(buffer) {
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 1
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 2
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 3
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 3
        buffer.indices[buffer.index_pos++] = buffer.index_offset
        buffer.indices[buffer.index_pos++] = buffer.index_offset + 1
        buffer.index_offset += 4
    }
    static Screen(buffer, x, y, width, height) {
        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y + height

        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y + height

        Render.Index4(buffer)
    }
    static Image(buffer, x, y, width, height, left, top, right, bottom) {
        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y
        buffer.vertices[buffer.vertexPos++] = left
        buffer.vertices[buffer.vertexPos++] = bottom

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y
        buffer.vertices[buffer.vertexPos++] = right
        buffer.vertices[buffer.vertexPos++] = bottom

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y + height
        buffer.vertices[buffer.vertexPos++] = right
        buffer.vertices[buffer.vertexPos++] = top

        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y + height
        buffer.vertices[buffer.vertexPos++] = left
        buffer.vertices[buffer.vertexPos++] = top

        Render.Index4(buffer)
    }
    static Rectangle(buffer, x, y, width, height, red, green, blue) {
        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y
        buffer.vertices[buffer.vertexPos++] = red
        buffer.vertices[buffer.vertexPos++] = green
        buffer.vertices[buffer.vertexPos++] = blue

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y
        buffer.vertices[buffer.vertexPos++] = red
        buffer.vertices[buffer.vertexPos++] = green
        buffer.vertices[buffer.vertexPos++] = blue

        buffer.vertices[buffer.vertexPos++] = x + width
        buffer.vertices[buffer.vertexPos++] = y + height
        buffer.vertices[buffer.vertexPos++] = red
        buffer.vertices[buffer.vertexPos++] = green
        buffer.vertices[buffer.vertexPos++] = blue

        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y + height
        buffer.vertices[buffer.vertexPos++] = red
        buffer.vertices[buffer.vertexPos++] = green
        buffer.vertices[buffer.vertexPos++] = blue

        Render.Index4(buffer)
    }
    static Circle(buffer, x, y, radius, red, green, blue) {
        const points = 32
        const tau = Math.PI * 2.0
        const slice = tau / points

        let firstIndex = buffer.index_offset
        buffer.vertices[buffer.vertexPos++] = x
        buffer.vertices[buffer.vertexPos++] = y
        buffer.vertices[buffer.vertexPos++] = red
        buffer.vertices[buffer.vertexPos++] = green
        buffer.vertices[buffer.vertexPos++] = blue
        buffer.index_offset++

        let radian = 0
        while (radian < tau) {
            buffer.vertices[buffer.vertexPos++] = x + Math.cos(radian) * radius
            buffer.vertices[buffer.vertexPos++] = y + Math.sin(radian) * radius
            buffer.vertices[buffer.vertexPos++] = red
            buffer.vertices[buffer.vertexPos++] = green
            buffer.vertices[buffer.vertexPos++] = blue

            buffer.indices[buffer.index_pos++] = firstIndex
            buffer.indices[buffer.index_pos++] = buffer.index_offset
            buffer.indices[buffer.index_pos++] = buffer.index_offset + 1
            buffer.index_offset++

            radian += slice
        }
        buffer.indices[buffer.index_pos - 1] = firstIndex + 1
    }
    static Print(buffer, x, y, text, scale) {
        let xx = x
        let yy = y
        for (let i = 0; i < text.length; i++) {
            let c = text.charAt(i)
            if (c === " ") {
                xx += FontWidth * scale
                continue
            } else if (c === "\n") {
                xx = x
                yy += FontHeight * scale
                continue
            }
            let loc = Font.indexOf(c)
            let tx1 = Math.floor(loc % FontGrid) * FontColumn
            let ty1 = Math.floor(loc / FontGrid) * FontRow
            let tx2 = tx1 + FontColumn
            let ty2 = ty1 + FontRow
            Render.Image(buffer, xx, yy, FontWidth * scale, FontHeight * scale, tx1, ty1, tx2, ty2)
            xx += FontWidth * scale
        }
    }
}
