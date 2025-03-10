package graphics

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hellory4n/stellarthing/core"
	"github.com/hellory4n/stellarthing/entities"
)

// how down can you go
const MinLayer int64 = -16

// how up can you go
const MaxLayer int64 = 48

// mate
const TotalLayers int64 = 64

// chunk size (they're square)
const ChunkSize int64 = 16

// optimized tile. for actual data use TileData
type Tile struct {
	TileId TileId
	EntityRef entities.EntityRef
	Variation VariationId
}

// as the name implies, it gets the data
func (t *Tile) GetData() *TileData {
	return Tiles[t.TileId][t.Variation]
}

func (t *Tile) String() string {
	return fmt.Sprintf("%v (variation %v, owned by %v)", TileNames[t.TileId], t.Variation, t.EntityRef)
}

// world of tiles :D
type TileWorld struct {
	Seed int64
	randGen *rand.Rand
	// set with SetCameraPosition :)
	CameraPosition core.Vec3
	// mate
	CameraOffset core.Vec2
	// the top left corner
	StartPos core.Vec2i
	// the bottom right corner
	EndPos core.Vec2i
	LoadedGroundTiles map[core.Vec3i]*Tile
	LoadedObjectTiles map[core.Vec3i]*Tile
	LoadedChunks []core.Vec3i
	// this sucks
	TheyMightBeMoving []*Tile
}

// current world.
var CurrentWorld *TileWorld

// makes a new world. the startPos is the top left corner and the endPos is the bottom right corner.
func NewTileWorld(startPos core.Vec2i, endPos core.Vec2i, seed int64) *TileWorld {
	tilhjjh := &TileWorld{}
	tilhjjh.StartPos = startPos
	tilhjjh.EndPos = endPos
	tilhjjh.Seed = seed
	tilhjjh.randGen = rand.New(rand.NewSource(tilhjjh.Seed))
	tilhjjh.LoadedGroundTiles = make(map[core.Vec3i]*Tile)
	tilhjjh.LoadedObjectTiles = make(map[core.Vec3i]*Tile)
	tilhjjh.CameraOffset = core.RenderSize.Sdiv(2).ToVec2()

	fmt.Println("[TILEMAP] Created new world")

	// load some chunks :)
	tilhjjh.SetCameraPosition(core.NewVec3(-float64(ChunkSize), 0, 0))
	tilhjjh.SetCameraPosition(core.NewVec3(float64(ChunkSize), 0, 0))
	tilhjjh.SetCameraPosition(core.NewVec3(0, float64(ChunkSize), 0))
	tilhjjh.SetCameraPosition(core.NewVec3(0, -float64(ChunkSize), 0))
	tilhjjh.SetCameraPosition(core.NewVec3(float64(ChunkSize), float64(ChunkSize), 0))
	tilhjjh.SetCameraPosition(core.NewVec3(-float64(ChunkSize), -float64(ChunkSize), 0))
	tilhjjh.SetCameraPosition(core.NewVec3(0, 0, 0))

	return tilhjjh
}

// sets the camera position and loads chunks
func (t *TileWorld) SetCameraPosition(pos core.Vec3) {
	// so i can use early returns :)
	defer func() { t.CameraPosition = pos }()

	chunkX := int64(pos.X / float64(ChunkSize))
	chunkY := int64(pos.Y / float64(ChunkSize))
	chunkPos := core.NewVec3i(chunkX, chunkY, int64(pos.Z))

	// do we even have to load it at all?
	_, hasChunk := t.LoadedGroundTiles[core.NewVec3i(int64(pos.X), int64(pos.Y), int64(pos.Z))]
	if hasChunk {
		return
	}

	// TODO check if it's on the save

	// if it's not on the save we generate it
	// generate multiple chunks fuck it
	man := func(offset core.Vec3i) {
		fmt.Printf("[TILEMAP] Generating chunk at %v\n", chunkPos.Add(offset))
		newGround, newObjects := GenerateChunk(t.randGen, chunkPos.Add(offset))

		// copy crap
		for k, v := range newGround {
			t.LoadedGroundTiles[k] = v
		}
		for k, v := range newObjects {
			t.LoadedObjectTiles[k] = v
		}
		t.LoadedChunks = append(t.LoadedChunks, chunkPos.Add(offset))
	}
	man(core.NewVec3i(0, 0, 0))
	man(core.NewVec3i(-1, -1, 0))
	man(core.NewVec3i(0, -1, 0))
	man(core.NewVec3i(1, 0, 0))
	man(core.NewVec3i(-1, 0, 0))
	man(core.NewVec3i(0, 1, 0))
	man(core.NewVec3i(-1, 1, 0))
	man(core.NewVec3i(0, 1, 0))
	man(core.NewVec3i(1, 1, 0))
}

// as the name implies, it gets a tile
func (t *TileWorld) GetTile(pos core.Vec3i, ground bool) *Tile {
	if ground {
		return t.LoadedGroundTiles[pos]
	} else {
		return t.LoadedObjectTiles[pos]
	}
}

// as the name implies, it makes a new tile. if the variation doesn't exist, it's gonna copy the
// default variation too
func (t *TileWorld) NewTile(pos core.Vec3i, ground bool, tileId TileId, entity entities.EntityRef,
variation VariationId) *Tile {
	var letile *Tile = &Tile{
		TileId: tileId,
		EntityRef: entity,
		Variation: variation,
	}

	_, variationExists := Tiles[tileId][variation]
	if !variationExists {
		Tiles[tileId][variation] = Tiles[tileId][0]
	}

	if ground {
		t.LoadedGroundTiles[pos] = letile
	} else {
		t.LoadedObjectTiles[pos] = letile
	}

	// man
	if entity != 0 {
		t.TheyMightBeMoving = append(t.TheyMightBeMoving, letile)
	}

	return letile
}

func (t *TileWorld) drawTile(pos core.Vec2i, ground bool) {
	// grod ng tiles
	tile := t.GetTile(core.NewVec3i(pos.X, pos.Y, int64(t.CameraPosition.Z)), ground)
	if tile == nil {
		return
	}
	data := tile.GetData()
	// YOU SEE TEXTURES ARE CACHED SO ITS NOT TOO OUTRAGEOUS TO PUT SOMETHING IN A FUNCTION
	// RAN EVERY FRAME
	texture := LoadTexture(data.Texture)

	var pospos core.Vec2
	if data.UsingCustomPos {
		// im losing my mind
		// im going insane
		// im watching my life go down the drain
		pospos = core.NewVec2(
			((data.Position.X - t.CameraPosition.X) * texture.Size().ToVec2().X) + t.CameraOffset.X,
			((data.Position.Y - t.CameraPosition.Y) * texture.Size().ToVec2().Y) + t.CameraOffset.Y,
		)
	} else {
		pospos = core.NewVec2(
			((float64(pos.X) - t.CameraPosition.X) * texture.Size().ToVec2().X) + t.CameraOffset.X,
			((float64(pos.Y) - t.CameraPosition.Y) * texture.Size().ToVec2().Y) + t.CameraOffset.Y,
		)
	}

	DrawTexture(texture, pospos, 0, data.Tint)
}

// it draws the world. no shit.
func (t *TileWorld) Draw() {
	// we draw the neighbors of the current chunk so it doesn't look funny
	// when crossing chunk borders
	renderAreaStartX := int64(math.Floor(t.CameraPosition.X - float64(ChunkSize)))
	renderAreaStartY := int64(math.Floor(t.CameraPosition.Y - float64(ChunkSize)))
	renderAreaEndX := int64(math.Floor(t.CameraPosition.X + float64(ChunkSize)))
	renderAreaEndY := int64(math.Floor(t.CameraPosition.Y + float64(ChunkSize)))

	for x := renderAreaStartX; x < renderAreaEndX; x++ {
		for y := renderAreaStartY; y < renderAreaEndY; y++ {
			t.drawTile(core.NewVec2i(x, y), true)
		}
	}

	for x := renderAreaStartX; x < renderAreaEndX; x++ {
		for y := renderAreaStartY; y < renderAreaEndY; y++ {
			t.drawTile(core.NewVec2i(x, y), false)
		}
	}

	for _, tile := range t.TheyMightBeMoving {
		// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
		if tile == nil {
			return
		}
		data := tile.GetData()
		// YOU SEE TEXTURES ARE CACHED SO ITS NOT TOO OUTRAGEOUS TO PUT SOMETHING IN A FUNCTION
		// RAN EVERY FRAME
		texture := LoadTexture(data.Texture)

		var pospos core.Vec2
		if data.UsingCustomPos {
			// im losing my mind
			// im going insane
			// im watching my life go down the drain
			pospos = core.NewVec2(
				((data.Position.X - t.CameraPosition.X) * texture.Size().ToVec2().X) + t.CameraOffset.X,
				((data.Position.Y - t.CameraPosition.Y) * texture.Size().ToVec2().Y) + t.CameraOffset.Y,
			)
		} else {
			continue
		}

		DrawTexture(texture, pospos, 0, data.Tint)
	}
}

// gets a tile position from screen positions
func (t *TileWorld) ScreenToTile(pos core.Vec2, textureSize core.Vec2i) core.Vec3i {
	return core.NewVec3i(
		int64(math.Floor(((pos.X - t.CameraOffset.X) / textureSize.ToVec2().X) + t.CameraPosition.X)),
		int64(math.Floor(((pos.Y - t.CameraOffset.Y) / textureSize.ToVec2().Y) + t.CameraPosition.Y)),
		int64(t.CameraPosition.Z),
	)
}

// gets a screen position from tile positions
func (t *TileWorld) TileToScreen(pos core.Vec3i, textureSize core.Vec2i) core.Vec2 {
	return core.NewVec2(
		((float64(pos.X) - t.CameraPosition.X) * textureSize.ToVec2().X) + t.CameraOffset.X,
		((float64(pos.Y) - t.CameraPosition.Y) * textureSize.ToVec2().Y) + t.CameraOffset.Y,
	)
}