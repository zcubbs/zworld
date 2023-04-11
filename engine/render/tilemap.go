package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zcubbs/zworld/engine/asset"
	"github.com/zcubbs/zworld/engine/tilemap"
)

type TilemapRender struct {
	spritesheet  *asset.SpriteSheet
	batch        *pixel.Batch
	tileToSprite map[tilemap.TileType]*pixel.Sprite
}

func NewTilemapRender(spritesheet *asset.SpriteSheet, tileToSprite map[tilemap.TileType]*pixel.Sprite) *TilemapRender {
	return &TilemapRender{
		spritesheet:  spritesheet,
		batch:        pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture()),
		tileToSprite: tileToSprite,
	}
}

func (r *TilemapRender) Clear() {

}

func (r *TilemapRender) Batch(t *tilemap.Tilemap) {
	for x := 0; x < t.Width(); x++ {
		for y := 0; y < t.Height(); y++ {
			tile, ok := t.Get(x, y)
			if !ok {
				continue
			}
			pos := pixel.V(float64(x*t.TileSize), float64(y*t.TileSize))

			sprite, ok := r.tileToSprite[tile.Type]
			if !ok {
				panic("Unable to find TileType")
			}
			mat := pixel.IM.Moved(pos)
			sprite.Draw(r.batch, mat)
		}
	}
}

func (r *TilemapRender) Draw(win *pixelgl.Window) {
	r.batch.Draw(win)
}
