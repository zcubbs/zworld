package main

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zcubbs/zworld"
	"github.com/zcubbs/zworld/engine/asset"
	"github.com/zcubbs/zworld/engine/render"
	"github.com/zcubbs/zworld/engine/tilemap"
	"log"
	"nhooyr.io/websocket"
	"os"
	"time"
)

func main() {
	// Setup Network
	url := "ws://localhost:8000"

	ctx := context.Background()
	c, resp, err := websocket.Dial(ctx, url, nil)
	check(err)

	log.Println("Connection response:", resp)

	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)

	go func() {
		counter := byte(0)
		for {
			time.Sleep(1 * time.Second)
			n, err := conn.Write([]byte{counter})
			if err != nil {
				log.Println("Error sending:", err)
				return
			}

			log.Println("Sent n bytes:", n)
			counter++
		}
	}()

	// Start Pixel
	pixelgl.Run(runGame)
}

func runGame() {
	cfg := pixelgl.WindowConfig{
		Title:     "zMMO",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
		Maximized: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)

	win.SetSmooth(false)

	load := asset.NewLoad(os.DirFS("./cmd/client"))
	spriteSheet, err := load.SpriteSheet("packed.json")
	check(err)

	// Create Tilemap
	seed := time.Now().UTC().UnixNano()
	tileSize := 16
	mapSize := 1000

	tmap := zworld.CreateTilemap(seed, mapSize, tileSize)

	grassTile, err := spriteSheet.Get("grass.png")
	check(err)
	dirtTile, err := spriteSheet.Get("dirt.png")
	check(err)
	waterTile, err := spriteSheet.Get("water.png")
	check(err)

	tmapRender := render.NewTilemapRender(spriteSheet, map[tilemap.TileType]*pixel.Sprite{
		zworld.GrassTile: grassTile,
		zworld.DirtTile:  dirtTile,
		zworld.WaterTile: waterTile,
	})

	tmapRender.Batch(tmap)

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
		tmapRender.Draw(win)

		// rendering
		for i := range pawns {
			pawns[i].Draw(win)
		}
		win.SetMatrix(pixel.IM)

		win.Update()
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
