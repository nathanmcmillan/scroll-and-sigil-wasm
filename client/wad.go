package main

import (
	"strings"
	"syscall/js"

	"./graphics"
	"./render"
)

var (
	wadSounds           = map[string]js.Value{}
	wadImageData        = map[string]map[string]*render.Sprite{}
	wadSpriteData       = map[string]map[string]*render.Sprite{}
	wadSpriteAlias      = map[string]map[string][]string{}
	wadSpriteAnimations = map[string]map[string][]string{}
	wadDirectionPrefix  = [5]string{
		"front-", "front-side-", "side-", "back-side-", "back-",
	}
)

func wadRead(g *graphics.RenderSystem, gl js.Value, data string) {
	wad := ParserRead(data)

	sprites := wad["sprites"].(map[string]interface{})
	animations := wad["animations"].(map[string]interface{})
	tiles := wad["tiles"].(map[string]interface{})
	resources := wad["resources"].(map[string]interface{})

	shaders := resources["shaders"].(*Array).data
	textures := resources["images"].(*Array).data
	sounds := resources["sounds"].(*Array).data

	for i := 0; i < len(shaders); i++ {
		name := shaders[i].(string)
		g.MakeProgram(gl, name)
	}

	for i := 0; i < len(textures); i++ {
		name := textures[i].(string)
		if name == "sky" {
			console("todo repeating textures (sky)")
			g.MakeImage(gl, name, graphics.GLxRepeat)
		} else {
			g.MakeImage(gl, name, graphics.GLxClampToEdge)
		}
	}

	for i := 0; i < len(sounds); i++ {
		name := sounds[i].(string)
		key := name[:strings.LastIndex(name, ".")]
		wadSounds[key] = js.Global().Get("Audio").New("sounds/" + name)
	}

	for name, value := range sprites {
		sprite := value.(map[string]interface{})

		texture := g.Textures[name]
		width := 1.0 / float32(texture.Width)
		height := 1.0 / float32(texture.Height)

		wadImageData[name] = make(map[string]*render.Sprite)
		wadSpriteData[name] = make(map[string]*render.Sprite)

		for frame, array := range sprite {
			data := array.(*Array).data
			size := len(data)
			atlas := make([]int, size)
			for i := 0; i < size; i++ {
				atlas[i] = ParseInt(data[i].(string))
			}
			wadImageData[name][frame] = render.SpriteBuild(atlas, width, height)
			wadSpriteData[name][frame] = render.SpriteBuild3d(atlas, width, height)
		}
	}

	for name, value := range animations {
		animation := value.(map[string]interface{})
		animationList := animation["animations"].(map[string]interface{})

		wadSpriteAlias[name] = make(map[string][]string)
		wadSpriteAnimations[name] = make(map[string][]string)

		for key, ls := range animationList {
			data := ls.(*Array).data
			size := len(data)
			array := make([]string, size)
			for i := 0; i < size; i++ {
				array[i] = data[i].(string)
			}
			wadSpriteAnimations[name][key] = array
		}

		if a, ok := animation["alias"]; ok {
			alias := a.(map[string]interface{})
			for key, ls := range alias {
				data := ls.(*Array).data
				size := len(data)
				array := make([]string, size)
				for i := 0; i < size; i++ {
					array[i] = data[i].(string)
				}
				wadSpriteAlias[name][key] = array
			}
		}
	}

	tileSprites := sprites["tiles"].(map[string]interface{})
	tileTexture := g.Textures["tiles"]
	tileAtlasWidth := 1.0 / float32(tileTexture.Width)
	tileAtlasHeight := 1.0 / float32(tileTexture.Height)
	tileSize := len(tileSprites)
	TileTexture = make([][4]float32, tileSize+1)
	TileClosed = make([]bool, tileSize+1)
	TileLookup["none"] = 0
	TileClosed[0] = false
	for name, value := range tileSprites {
		data := value.(*Array).data
		x := float32(ParseInt(data[0].(string)))
		y := float32(ParseInt(data[1].(string)))
		w := float32(ParseInt(data[2].(string)))
		h := float32(ParseInt(data[3].(string)))
		tileData := tiles[name].(map[string]interface{})
		tileUID := ParseInt(tileData["uid"].(string))
		TileLookup[name] = tileUID
		TileTexture[tileUID] = render.SimpleSprite(x, y, w, h, tileAtlasWidth, tileAtlasHeight)
		TileClosed[tileUID] = tileData["closed"].(string) == "true"
	}

	baronAnimationIdle = spriteBuilderDirectional("baron", "idle")
	baronAnimationWalk = spriteBuilderDirectional("baron", "walk")
	baronAnimationMelee = spriteBuilderDirectional("baron", "melee")
	baronAnimationMissile = spriteBuilderDirectional("baron", "missile")
	baronAnimationDeath = spriteBuilderDirectional("baron", "death")

	humanAnimationIdle = spriteBuilderDirectional("baron", "idle")
	humanAnimationWalk = spriteBuilderDirectional("baron", "walk")
	humanAnimationMelee = spriteBuilderDirectional("baron", "melee")
	humanAnimationMissile = spriteBuilderDirectional("baron", "missile")
	humanAnimationDeath = spriteBuilderDirectional("baron", "death")

	plasmaExplosionAnimation = spriteBuilder("particles", "plasma-explosion")
}

func spriteBuilder(sid string, name string) []*render.Sprite {
	animationData := wadSpriteAnimations[sid][name]
	size := len(animationData)
	data := make([]*render.Sprite, size)
	for i := 0; i < size; i++ {
		frame := animationData[i]
		data[i] = wadSpriteData[sid][frame]
	}
	return data
}

func spriteBuilderDirectional(sid string, name string) [][]*render.Sprite {
	animationData := wadSpriteAnimations[sid][name]
	size := len(animationData)
	data := make([][]*render.Sprite, size)
	for i := 0; i < size; i++ {
		frame := animationData[i]
		directional := make([]*render.Sprite, len(wadDirectionPrefix))
		for d := 0; d < len(wadDirectionPrefix); d++ {
			direction := wadDirectionPrefix[d]
			sprite, ok := wadSpriteData[sid][direction+frame]
			if ok {
				directional[d] = sprite
			} else {
				directional[d] = wadSpriteData[sid]["front-"+frame]
			}
		}
		data[i] = directional
	}
	return data
}

func playWadSound(name string) {
	sound := wadSounds[name]
	sound.Call("pause")
	sound.Set("currentTime", 0)
	sound.Call("play")
}
