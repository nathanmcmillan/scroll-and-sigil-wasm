const Sounds = {}
const ImageData = {}
const SpriteData = {}
const SpriteAlias = {}
const SpriteAnimations = {}
const DirectionPrefix = ["front-", "front-side-", "side-", "back-side-", "back-"]

class Wad {
    static async Load(g, gl, string) {
        let wad = Parser.read(string)

        console.log(wad)

        let resources = wad["resources"]
        let sprites = wad["sprites"]
        let animations = wad["animations"]
        let tiles = wad["tiles"]
        let shaders = resources["shaders"]
        let textures = resources["images"]
        let sounds = resources["sounds"]

        let promises = []

        for (let index = 0; index < shaders.length; index++)
            promises.push(g.makeProgram(gl, shaders[index]))

        for (let index = 0; index < textures.length; index++) {
            if (textures[index] === "sky") {
                console.log("todo", textures[index])
                promises.push(g.makeImage(gl, textures[index], gl.REPEAT))
            } else
                promises.push(g.makeImage(gl, textures[index], gl.CLAMP_TO_EDGE))
        }

        await Promise.all(promises)

        for (let s in sounds) {
            let name = sounds[s]
            let key = name.substring(0, name.lastIndexOf("."))
            Sounds[key] = new Audio("sounds/" + name)
        }

        for (let name in sprites) {
            let sprite = sprites[name]

            let texture = g.textures.get(name)
            let width = 1.0 / texture.image.width
            let height = 1.0 / texture.image.height

            ImageData[name] = {}
            SpriteData[name] = {}

            for (let frame in sprite) {
                let data = sprite[frame]
                let atlas = []
                for (let i in data)
                    atlas.push(parseInt(data[i]))
                ImageData[name][frame] = Sprite.Build(atlas, width, height)
                SpriteData[name][frame] = Sprite.Build3(atlas, width, height)
            }
        }

        for (let name in animations) {
            let animation = animations[name]
            let animation_list = animation["animations"]
            let alias = ("alias" in animation) ? animation["alias"] : null

            SpriteAlias[name] = {}
            SpriteAnimations[name] = {}

            for (let key in animation_list)
                SpriteAnimations[name][key] = animation_list[key]

            if (alias != null) {
                for (let key in alias)
                    SpriteAlias[name][key] = alias[key]
            }
        }

        let tileSprites = sprites["tiles"]
        let texture = g.textures.get("tiles")
        let width = 1.0 / texture.image.width
        let height = 1.0 / texture.image.height
        TileTexture[0] = null
        TileClosed[0] = false
        for (let tileName in tileSprites) {
            let data = tileSprites[tileName]
            let x = parseInt(data[0])
            let y = parseInt(data[1])
            let w = parseInt(data[2])
            let h = parseInt(data[3])
            let tileData = tiles[tileName]
            let tileUID = parseInt(tileData["uid"])
            TileLookup.set(tileName, tileUID)
            TileTexture[tileUID] = Sprite.Simple(x, y, w, h, width, height)
            TileClosed[tileUID] = tileData["closed"] === "true"
        }

        Wad.SpriteBuilderDirectional("baron", BaronAnimationIdle, "idle")
        Wad.SpriteBuilderDirectional("baron", BaronAnimationWalk, "walk")
        Wad.SpriteBuilderDirectional("baron", BaronAnimationMelee, "melee")
        Wad.SpriteBuilderDirectional("baron", BaronAnimationMissile, "missile")
        Wad.SpriteBuilderDirectional("baron", BaronAnimationDeath, "death")

        HumanAnimationIdle.push.apply(HumanAnimationIdle, BaronAnimationIdle)
        HumanAnimationWalk.push.apply(HumanAnimationWalk, BaronAnimationWalk)
        HumanAnimationMelee.push.apply(HumanAnimationMelee, BaronAnimationMelee)
        HumanAnimationMissile.push.apply(HumanAnimationMissile, BaronAnimationMissile)
        HumanAnimationDeath.push.apply(HumanAnimationDeath, BaronAnimationDeath)

        Wad.SpriteBuilder("particles", PlasmaExplosionAnimation, "plasma-explosion")
    }
    static SpriteBuilder(sid, array, name) {
        let animation = []
        let animationData = SpriteAnimations[sid][name]
        for (let a in animationData) {
            let name = animationData[a]
            animation.push(SpriteData[sid][name])
        }
        array.push.apply(array, animation)
    }
    static SpriteBuilderDirectional(sid, array, name) {
        let animation = []
        let animationData = SpriteAnimations[sid][name]
        for (let a in animationData) {
            let name = animationData[a]
            let slice = new Array(5)
            for (let d in DirectionPrefix) {
                let direction = DirectionPrefix[d]
                let fullname = direction + name
                let sprite = SpriteData[sid]["front-" + name]
                if (fullname in SpriteData[sid]) {
                    sprite = SpriteData[sid][fullname]
                }
                slice[d] = sprite
            }
            animation.push(slice)
        }
        array.push.apply(array, animation)
    }
}
