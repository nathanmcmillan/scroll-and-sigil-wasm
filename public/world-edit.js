class WorldEdit extends World {
    constructor(g, gl) {
        super(g, gl)
    }
    LoadRaw(raw) {
        this.Reset()

        let content = Parser.read(raw)

        let width = parseInt(content["w"])
        let height = parseInt(content["h"])
        let length = parseInt(content["l"])

        let blocks = content["b"]

        let things = content["t"]
        let items = content["i"]

        this.width = width
        this.height = height
        this.length = length
        this.slice = width * height
        this.all = this.slice * length

        this.tileWidth = this.width * BlockSize
        this.tileHeight = this.height * BlockSize
        this.tileLength = this.length * BlockSize

        let bx = 0
        let by = 0
        let bz = 0
        for (let b = 0; b < blocks.length; b++) {
            let bdata = blocks[b]
            let tiles = bdata["t"]
            let lights = bdata["c"]

            let block = new Block(bx, by, bz)
            if (tiles.length > 0) {
                for (let t = 0; t < BlockAll; t++)
                    block.tiles[t].type = parseInt(tiles[t])
            }

            for (let t = 0; t < lights.length; t++) {
                let light = lights[t]
                let x = parseInt(light["x"])
                let y = parseInt(light["y"])
                let z = parseInt(light["z"])
                let rgb = parseInt(light["v"])
                block.AddLight(new Light(x, y, z, rgb))
            }

            this.blocks[bx + by * this.width + bz * this.slice] = block

            bx++
            if (bx === this.width) {
                bx = 0
                by++
                if (by === this.height) {
                    by = 0
                    bz++
                }
            }
        }

        for (let t = 0; t < things.length; t++) {
            let thing = things[t]
            let uid = parseInt(thing["u"])
            let x = parseFloat(thing["x"])
            let y = parseFloat(thing["y"])
            let z = parseFloat(thing["z"])
            switch (uid) {
                case BaronUID:
                    new Baron(this, null, x, y, z, DirectionNorth, 1, BaronLook)
                    break
                case TreeUID:
                    new Tree(this, null, x, y, z)
                    break
            }
        }

        for (let t = 0; t < items.length; t++) {
            let item = items[t]
            let uid = parseInt(item["u"])
            console.log("item", uid)
        }

        this.Build()
    }
    Save() {
        // TODO this.Compress()
        let data = "w:" + this.width
        data += ",h:" + this.height
        data += ",l:" + this.length
        data += ",b["
        for (let i = 0; i < this.all; i++) {
            let block = this.blocks[i]
            data += block.Save()
            data += ","
        }
        data += "],t["
        for (let i = 0; i < this.thingCount; i++) {
            data += this.things[i].Save()
            data += ","
        }
        data += "],i["
        for (let i = 0; i < this.itemCount; i++) {
            data += this.items[i].Save()
            data += ","
        }
        data += "]"
        return data
    }
}
