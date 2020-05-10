class SimpleCamera {
    constructor(x, y, z) {
        this.X = x
        this.Y = y
        this.Z = z
        this.RX = 0
        this.RY = 0
    }
    update() {
        if (Input.KeyDown("ArrowLeft")) {
            this.RY -= 0.05
            if (this.RY < 0)
                this.RY += Tau
        }

        if (Input.KeyDown("ArrowRight")) {
            this.RY += 0.05
            if (this.RY >= Tau)
                this.RY -= Tau
        }

        if (this.RX > -Math.PI && Input.KeyDown("ArrowUp"))
            this.RX -= 0.05

        if (this.RX < Math.PI && Input.KeyDown("ArrowDown"))
            this.RX += 0.05

        const speed = 0.1

        if (Input.KeyDown("w")) {
            this.X += Math.sin(this.RY) * speed
            this.Z -= Math.cos(this.RY) * speed
        }

        if (Input.KeyDown("s")) {
            this.X -= Math.sin(this.RY) * speed
            this.Z += Math.cos(this.RY) * speed
        }

        if (Input.KeyDown("a")) {
            this.X -= Math.cos(this.RY) * speed
            this.Z -= Math.sin(this.RY) * speed
        }

        if (Input.KeyDown("d")) {
            this.X += Math.cos(this.RY) * speed
            this.Z += Math.sin(this.RY) * speed
        }

        if (Input.KeyDown("q")) {
            this.Y += speed
        }

        if (Input.KeyDown("e")) {
            this.Y -= speed
        }
    }
}
