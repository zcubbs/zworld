package tilemap

type TileType uint8

type Tile struct {
	Type TileType
}

type Tilemap struct {
	TileSize int
	Tiles    [][]Tile
}

func New(tiles [][]Tile, tileSize int) *Tilemap {
	return &Tilemap{
		TileSize: tileSize,
		Tiles:    tiles,
	}
}

func (t *Tilemap) Width() int {
	return len(t.Tiles)
}

func (t *Tilemap) Height() int {
	return len(t.Tiles[0])
}

func (t *Tilemap) Get(x, y int) (Tile, bool) {
	if x < 0 || x >= len(t.Tiles) || y < 0 || y >= len(t.Tiles[x]) {
		return Tile{}, false
	}

	return t.Tiles[x][y], true
}
