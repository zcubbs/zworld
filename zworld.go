package zworld

import (
	"github.com/zcubbs/zworld/engine/pgen"
	"github.com/zcubbs/zworld/engine/tilemap"
	"math"
)

const (
	GrassTile tilemap.TileType = iota
	DirtTile
	WaterTile
)

func CreateTilemap(seed int64, mapSize, tileSize int) *tilemap.Tilemap {
	octaves := []pgen.Octave{
		{0.01, 0.6},
		{0.05, 0.3},
		{0.1, 0.07},
		{0.2, 0.02},
		{0.4, 0.01},
	}
	exponent := 0.8
	terrain := pgen.NewNoiseMap(seed, octaves, exponent)

	waterLevel := 0.5
	landLevel := waterLevel + 0.1

	islandExponent := 2
	tiles := make([][]tilemap.Tile, mapSize, mapSize)
	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize, mapSize)
		for y := range tiles[x] {

			height := terrain.Get(x, y)

			// Force an island shape
			{
				dx := float64(x)/float64(mapSize) - 0.5
				dy := float64(y)/float64(mapSize) - 0.5
				d := math.Sqrt(dx*dx+dy*dy) * 2
				d = math.Pow(d, float64(islandExponent))
				height = (1 - d + height) / 2
			}

			if height < waterLevel {
				tiles[x][y] = tilemap.Tile{Type: WaterTile}
			} else if height < landLevel {
				tiles[x][y] = tilemap.Tile{Type: DirtTile}
			} else {
				tiles[x][y] = tilemap.Tile{Type: GrassTile}
			}
		}
	}

	tmap := tilemap.New(tiles, tileSize)

	return tmap
}
