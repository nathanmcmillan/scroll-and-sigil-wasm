class Camera {
    constructor(world, thing, radius) {
        this.thing = thing
        this.radius = radius
        this.x = 0
        this.y = 0
        this.z = 0
        this.rx = 0
        this.ry = 0
        this.update(world)
    }
    update(world) {
        if (Input.KeyDown("ArrowLeft")) {
            this.ry -= 0.05
            if (this.ry < 0)
                this.ry += Tau
        }

        if (Input.KeyDown("ArrowRight")) {
            this.ry += 0.05
            if (this.ry >= Tau)
                this.ry -= Tau
        }

        if (this.rx > -0.25 && Input.KeyDown("ArrowUp"))
            this.rx -= 0.05

        if (this.rx < 0.25 && Input.KeyDown("ArrowDown"))
            this.rx += 0.05

        let sinX = Math.sin(this.rx)
        let cosX = Math.cos(this.rx)
        let sinY = Math.sin(this.ry)
        let cosY = Math.cos(this.ry)

        let thing = this.thing

        let dx = -cosX * sinY
        let dy = sinX
        let dz = cosX * cosY

        let x = thing.X + this.radius * dx
        let y = thing.Y + this.radius * dy + thing.Height
        let z = thing.Z + this.radius * dz

        Cast.Exact(world, thing.X, thing.Y + thing.Height, thing.Z, x, y, z)

        if (CastX === null) {
            this.x = x
            this.y = y
            this.z = z
        } else {
            this.x = CastX
            this.y = CastY
            this.z = CastZ
        }
    }
}
