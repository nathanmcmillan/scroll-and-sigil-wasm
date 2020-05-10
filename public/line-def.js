class Line {
    constructor(fromX, fromY, toX, toY) {
        this.fromX = fromX
        this.fromY = fromY
        this.toX = toX
        this.toY = toY
    }
    Intersect(line) {
        let a1 = this.toY - this.fromY
        let b1 = this.fromX - this.toX
        let c1 = this.toX * this.fromY - this.fromX * this.toY

        let r3 = a1 * line.fromX + b1 * line.fromY + c1
        let r4 = a1 * line.toX + b1 * line.toY + c1

        if (r3 !== 0 && r4 !== 0 && r3 * r4 >= 0) {
            return null
        }

        let a2 = line.toY - line.fromY
        let b2 = line.fromX - line.toX
        let c2 = line.toX * line.fromY - line.fromX * line.toY

        let r1 = a2 * this.fromX + b2 * this.fromY + c2
        let r2 = a2 * this.toX + b2 * this.toY + c2

        if (r1 !== 0 && r2 !== 0 && r1 * r2 >= 0) {
            return null
        }

        let denom = a1 * b2 - a2 * b1
        if (denom === 0) {
            return null
        }

        let offset
        let x
        let y

        if (denom < 0) {
            offset = -denom / 2
        } else {
            offset = denom / 2
        }

        let num = b1 * c2 - b2 * c1
        if (num < 0) {
            x = (num - offset) / denom
        } else {
            x = (num + offset) / denom
        }

        num = a2 * c1 - a1 * c2
        if (num < 0) {
            y = (num - offset) / denom
        } else {
            y = (num + offset) / denom
        }

        return [x, y]
    }
}
