package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zcubbs/zmmo/engine/asset"
	"github.com/zcubbs/zmmo/engine/pgen"
	"github.com/zcubbs/zmmo/engine/render"
	"github.com/zcubbs/zmmo/engine/tilemap"
	"math"
	"os"
	"time"
)

func main() {
	pixelgl.Run(runGame)
}

func runGame() {
	cfg := pixelgl.WindowConfig{
		Title:     "zMMO",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)

	win.SetSmooth(false)

	load := asset.NewLoad(os.DirFS("./"))
	spriteSheet, err := load.SpriteSheet("packed.json")
	check(err)

	// Create Tilemap
	seed := time.Now().UTC().UnixNano()
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
	tileSize := 16
	mapSize := 1000
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
				tiles[x][y] = GetTile(spriteSheet, WaterTile)
			} else if height < landLevel {
				tiles[x][y] = GetTile(spriteSheet, DirtTile)
			} else {
				tiles[x][y] = GetTile(spriteSheet, GrassTile)
			}
		}
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet.Picture())
	tmap := tilemap.New(tiles, batch, tileSize)
	tmap.Rebatch()

	// Create Pawns
	spawnPoint := pixel.V(float64(tileSize*mapSize/2), float64(tileSize*mapSize/2))

	gopher1, err := spriteSheet.Get("zopher.png")
	check(err)
	gopher2, err := spriteSheet.Get("gopher2.png")
	check(err)

	pawns := make([]Person, 0)

	pawns = append(pawns, NewPerson(gopher1, spawnPoint, KeyBinds{
		Up:    pixelgl.KeyUp,
		Down:  pixelgl.KeyDown,
		Left:  pixelgl.KeyLeft,
		Right: pixelgl.KeyRight,
	}))
	pawns = append(pawns, NewPerson(gopher2, spawnPoint, KeyBinds{
		Left:  pixelgl.KeyA,
		Right: pixelgl.KeyD,
		Up:    pixelgl.KeyW,
		Down:  pixelgl.KeyS,
	}))

	camera := render.NewCamera(win, 0, 0)
	zoomSpeed := 0.1

	for !win.JustPressed(pixelgl.KeyEscape) {
		win.Clear(pixel.RGB(0, 0, 0))

		scroll := win.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		// inputs
		for i := range pawns {
			pawns[i].HandleInput(win)
		}

		camera.Position = pawns[0].Position
		camera.Update()
		win.SetMatrix(camera.Mat())
		// Draw tilemap
		tmap.Draw(win)

		// rendering
		for i := range pawns {
			pawns[i].Draw(win)
		}
		win.SetMatrix(pixel.IM)

		win.Update()
	}
}

const (
	GrassTile tilemap.TileType = iota
	DirtTile
	WaterTile
)

func GetTile(ss *asset.SpriteSheet, t tilemap.TileType) tilemap.Tile {
	spriteName := ""

	switch t {
	case GrassTile:
		spriteName = "grass.png"
	case DirtTile:
		spriteName = "dirt.png"
	case WaterTile:
		spriteName = "water.png"
	default:
		panic(fmt.Sprintf("unknown tile type: %v", t))
	}

	sprite, err := ss.Get(spriteName)
	check(err)
	return tilemap.Tile{
		Type:   t,
		Sprite: sprite,
	}
}

type KeyBinds struct {
	Up, Down, Left, Right pixelgl.Button
}

type Person struct {
	Sprite   *pixel.Sprite
	Position pixel.Vec
	KeyBinds
}

func NewPerson(sprite *pixel.Sprite, position pixel.Vec, kBinds KeyBinds) Person {
	return Person{
		Sprite:   sprite,
		Position: position,
		KeyBinds: kBinds,
	}
}

func (p *Person) Draw(win *pixelgl.Window) {
	p.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(p.Position))
}

func (p *Person) HandleInput(win *pixelgl.Window) {
	if win.Pressed(p.KeyBinds.Left) {
		p.Position.X -= 2.0
	}
	if win.Pressed(p.KeyBinds.Right) {
		p.Position.X += 2.0
	}
	if win.Pressed(p.KeyBinds.Up) {
		p.Position.Y += 2.0
	}
	if win.Pressed(p.KeyBinds.Down) {
		p.Position.Y -= 2.0
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
