package main

import (
	"time"

	"../fast"
)

// World constants
const (
	WorldTickRate  = 50
	WorldPositiveX = 0
	WorldPositiveY = 1
	WorldPositiveZ = 2
	WorldNegativeX = 3
	WorldNegativeY = 4
	WorldNegativeZ = 5
)

// Broadcast constants
const (
	BroadcastNew    = uint8(0)
	BroadcastDelete = uint8(1)
	BroadcastChat   = uint8(2)
)

// World variables
var (
	WorldThreads = []string{"ai", "pathing"}
)

// World struct
type World struct {
	Width                           int
	Height                          int
	Length                          int
	Slice                           int
	All                             int
	Blocks                          []block
	thingCount                      int
	itemCount                       int
	missileCount                    int
	things                          []*thing
	items                           []*item
	missiles                        []*missile
	ThreadIndex                     int
	ThreadID                        string
	SpawnYouX, SpawnYouY, SpawnYouZ float32
	broadcastCount                  uint8
	broadcast                       *fast.ByteWriter
}

// NewWorld func
func NewWorld() *World {
	world := &World{}
	world.broadcast = fast.ByteWriterInit(64)
	world.SpawnYouX = 12
	world.SpawnYouY = 12
	world.SpawnYouZ = 12
	return world
}

// Load func
func (me *World) Load(data []byte) {
	content := ParserRead(data)

	width := ParseInt(content["w"].(string))
	height := ParseInt(content["h"].(string))
	length := ParseInt(content["l"].(string))

	blocks := content["b"].(*Array).data
	num := len(blocks)

	things := content["t"].(*Array).data
	items := content["i"].(*Array).data

	me.Width = width
	me.Height = height
	me.Length = length
	me.Slice = width * height
	me.All = me.Slice * length
	me.Blocks = make([]block, me.All)

	me.thingCount = 0
	me.itemCount = 0
	me.missileCount = 0

	me.things = make([]*thing, 5)
	me.items = make([]*item, 5)
	me.missiles = make([]*missile, 5)

	bx := 0
	by := 0
	bz := 0
	for b := 0; b < num; b++ {
		bdata := blocks[b].(map[string]interface{})
		tiles := bdata["t"].(*Array).data
		lights := bdata["c"].(*Array).data

		block := &me.Blocks[bx+by*me.Width+bz*me.Slice]
		block.blockInit(bx, by, bz)
		if len(tiles) > 0 {
			for t := 0; t < BlockAll; t++ {
				block.Tiles[t] = ParseInt(tiles[t].(string))
			}
		}

		for t := 0; t < len(lights); t++ {
			light := lights[t].(map[string]interface{})
			x := ParseInt(light["x"].(string))
			y := ParseInt(light["y"].(string))
			z := ParseInt(light["z"].(string))
			rgb := ParseInt(light["v"].(string))
			block.addLight(lightInit(x, y, z, rgb))
		}

		bx++
		if bx == width {
			bx = 0
			by++
			if by == height {
				by = 0
				bz++
			}
		}
	}

	for t := 0; t < len(things); t++ {
		thing := things[t].(map[string]interface{})
		uid := ParseInt(thing["u"].(string))
		x := ParseFloat(thing["x"].(string))
		y := ParseFloat(thing["y"].(string))
		z := ParseFloat(thing["z"].(string))
		LoadNewthing(me, uint16(uid), x, y, z)
	}

	for t := 0; t < len(items); t++ {
		item := items[t].(map[string]interface{})
		uid := ParseInt(item["u"].(string))
		x := ParseFloat(item["x"].(string))
		y := ParseFloat(item["y"].(string))
		z := ParseFloat(item["z"].(string))
		LoadNewItem(me, uint16(uid), x, y, z)
	}
}

// NewPlayer func
func (me *World) NewPlayer(person *Person) *You {
	return NewYou(me, person, me.SpawnYouX, me.SpawnYouY, me.SpawnYouZ)
}

// Save func
func (me *World) Save(person *Person) []byte {
	data := fast.ByteWriterInit(256)

	data.PutUint16(person.Character.NID)
	data.PutUint16(uint16(me.Width))
	data.PutUint16(uint16(me.Height))
	data.PutUint16(uint16(me.Length))

	for i := 0; i < me.All; i++ {
		me.Blocks[i].Save(data)
	}

	numthings := me.thingCount
	data.PutUint16(uint16(numthings))
	for i := 0; i < numthings; i++ {
		me.things[i].Save(data)
	}

	numItems := me.itemCount
	data.PutUint16(uint16(numItems))
	for i := 0; i < numItems; i++ {
		me.items[i].Save(data)
	}

	numMissiles := me.missileCount
	data.PutUint16(uint16(numMissiles))
	for i := 0; i < numMissiles; i++ {
		me.missiles[i].Snap(data)
	}

	return data.Bytes()
}

// BuildSnapshots func
func (me *World) BuildSnapshots(people []*Person) {
	peopleCount := len(people)
	time := uint32(time.Now().UnixNano()/1000000 - 1552330000000)
	things := me.thingCount

	var broadcastBinary []byte
	broadcastSize := me.broadcastCount
	if broadcastSize > 0 {
		binary := me.broadcast.Bytes()
		broadcastBinary = make([]byte, len(binary))
		copy(broadcastBinary, binary)
		me.broadcast.Reset()
		me.broadcastCount = 0
	}

	full := &fast.ByteWriter{}
	body := &fast.ByteWriter{}
	for i := 0; i < things; i++ {
		me.things[i].Snap(body)
	}

	body.Reset()
	spriteSet := make(map[*thing]bool)
	updatedThings := uint16(0)
	for i := 0; i < things; i++ {
		thing := me.things[i]
		if _, has := spriteSet[thing]; !has {
			spriteSet[thing] = true
			if thing.Binary != nil {
				body.PutBytes(thing.Binary)
				updatedThings++
			}
		}
	}

	full.PutUint32(time)
	full.PutUint8(broadcastSize)
	if broadcastSize > 0 {
		full.PutBytes(broadcastBinary)
	}
	full.PutUint16(updatedThings)
	full.PutBytes(body.Bytes())

	for i := 0; i < peopleCount; i++ {
		person := people[i]
		binary := full.Bytes()
		person.snap = make([]byte, len(binary))
		copy(person.snap, binary)
	}
}

// FindBlock func
func (me *World) FindBlock(x, y, z float32) int {
	gx := int(x)
	gy := int(y)
	gz := int(z)
	bx := int(x * InverseBlockSize)
	by := int(y * InverseBlockSize)
	bz := int(z * InverseBlockSize)
	tx := gx - bx*BlockSize
	ty := gy - by*BlockSize
	tz := gz - bz*BlockSize
	block := &me.Blocks[bx+by*me.Width+bz*me.Slice]
	return block.Tiles[tx+ty*BlockSize+tz*BlockSlice]
}

// GetTileType func
func (me *World) GetTileType(bx, by, bz, tx, ty, tz int) int {
	for tx < 0 {
		tx += BlockSize
		bx--
	}
	for tx >= BlockSize {
		tx -= BlockSize
		bx++
	}
	for ty < 0 {
		ty += BlockSize
		by--
	}
	for ty >= BlockSize {
		ty -= BlockSize
		by++
	}
	for tz < 0 {
		tz += BlockSize
		bz--
	}
	for bz >= BlockSize {
		tz -= BlockSize
		bz++
	}
	block := me.getBlock(bx, by, bz)
	if block == nil {
		return TileNone
	}
	return block.GetTileTypeUnsafe(tx, ty, tz)
}

// getBlock func
func (me *World) getBlock(x, y, z int) *block {
	if x < 0 || x >= me.Width {
		return nil
	}
	if y < 0 || y >= me.Height {
		return nil
	}
	if z < 0 || z >= me.Length {
		return nil
	}
	return &me.Blocks[x+y*me.Width+z*me.Slice]
}

// addThing func
func (me *World) addThing(t *thing) {
	if me.thingCount == len(me.things) {
		array := make([]*thing, me.thingCount+5)
		copy(array, me.things)
		me.things = array
	}
	me.things[me.thingCount] = t
	me.thingCount++
}

// removeThing func
func (me *World) removeThing(t *thing) {
	for i := 0; i < me.thingCount; i++ {
		if me.things[i] == t {
			me.things[i] = me.things[me.thingCount-1]
			me.thingCount--
			return
		}
	}
}

// addItem func
func (me *World) addItem(t *item) {
	if me.itemCount == len(me.items) {
		array := make([]*item, me.itemCount+5)
		copy(array, me.items)
		me.items = array
	}
	me.items[me.itemCount] = t
	me.itemCount++
}

// removeItem func
func (me *World) removeItem(t *item) {
	for i := 0; i < me.itemCount; i++ {
		if me.items[i] == t {
			me.items[i] = me.items[me.itemCount-1]
			me.itemCount--
			return
		}
	}
}

// addMissile func
func (me *World) addMissile(t *missile) {
	if me.missileCount == len(me.missiles) {
		array := make([]*missile, me.missileCount+5)
		copy(array, me.missiles)
		me.missiles = array
	}
	me.missiles[me.missileCount] = t
	me.missileCount++
}

// Update func
func (me *World) Update() {
	me.ThreadID = WorldThreads[me.ThreadIndex]
	me.ThreadIndex++
	if me.ThreadIndex == len(WorldThreads) {
		me.ThreadIndex = 0
	}
	num := me.thingCount
	for i := 0; i < num; i++ {
		thing := me.things[i]
		if thing.Update() {
			me.things[i] = me.things[num-1]
			me.things[num-1] = nil
			me.thingCount--
			num--
			i--
		}
	}
	num = me.missileCount
	for i := 0; i < num; i++ {
		missile := me.missiles[i]
		if missile.Update() {
			me.missiles[i] = me.missiles[num-1]
			me.missiles[num-1] = nil
			me.missileCount--
			num--
			i--
		}
	}
}
