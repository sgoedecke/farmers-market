package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"image/color"
	"math"
	"time"
)

var (
	fullscreen = false
	scale      = 50.0
	world      = World{}
	fps        = 60.0
	player     = Player{}
)

func setup() {
	player.Pos = pixel.Vec{15, 15}
	player.Dir = pixel.Vec{0, 0}
	player.Speed = 0.12
	world.Height = 20
	world.Width = 20

	player.LoadTextures()
	world.LoadTextures()
}

func moveEntities() {
	newX := player.Pos.X + (player.Dir.X * player.Speed)
	newY := player.Pos.Y + (player.Dir.Y * player.Speed)

	if (newX >= 0 && newX < float64(world.Width)) && (newY >= 0 && newY < float64(world.Height)) {
		player.Pos.X = float64(newX)
		player.Pos.Y = float64(newY)
	}
}

// from the global variables that hold the game state, draw a frame of the game into
// the image.RGBA buffer
func frame() *image.RGBA {

	m := image.NewRGBA(image.Rect(0, 0, world.Width*int(scale), world.Height*int(scale)))

	world.Draw(m)

	// draw player
	player.Draw(m)
	return m
}

func run() {

	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, 1000, 1000),
		VSync:       true,
		Undecorated: false,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	for !win.Closed() {
		// advance game

		// if we're running faster than our fps, wait
		dt := time.Since(last).Seconds()
		if dt < 1/fps {
			//fmt.Println("Sleeping!")
			//fmt.Println(1/fps - dt)
			time.Sleep(time.Duration((1/fps - dt) * 1000000000)) // have to convert seconds to ns for the cast
		} else {
			fmt.Println("Running slow")
			fmt.Println(time.Duration((1/fps - dt) * 1000000000))
		}
		last = time.Now()

		// advance world tick - used for animations
		world.Tick += 1
		if world.Tick > 29 {
			world.Tick = 1
		}

		player.SetActiveTextureCoord(world.Tick)

		if math.Abs(float64(world.HighlightedTile.Tick-world.Tick)) > 10 {
			world.HighlightedTile.Tick = 0
		}

		moveEntities()
		haltPlayer()

		// handle keys
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			moveUp()
		}

		if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyLeft) {
			moveLeft()
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			moveDown()
		}

		if win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
			moveRight()
		}

		if win.Pressed(pixelgl.KeyE) {
			actOnTile()
		}

		// draw
		win.Clear(color.Black)
		p := pixel.PictureDataFromImage(frame())

		// offset the center of the screen so camera-following works
		c := win.Bounds().Center().Add(pixel.Vec{float64(world.Width) * 0.5, float64(world.Height) * -0.5}.Scaled(scale))

		// since we store player coordinates on the 2d array (0,0 is top left) but draw with
		// Euclidean coordinates, we need to flip the X coord
		playerPos := pixel.Vec{player.Pos.X * -1.0, player.Pos.Y}.Scaled(scale)

		d := playerPos.Add(c)

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(d))

		win.Update()
	}
}

func actOnTile() {
	x := int(player.Pos.X) + int(player.Dir.X)
	y := int(player.Pos.Y) + int(player.Dir.Y)
	world.HighlightedTile.Pos = image.Point{x, y}
	world.HighlightedTile.Tick = world.Tick
}

func haltPlayer() {
	player.Dir.X = 0
	player.Dir.Y = 0
}
func moveUp() {
	player.Dir.Y = -1.0
}
func moveDown() {
	player.Dir.Y = 1.0
}
func moveLeft() {
	player.Dir.X = -1.0
}
func moveRight() {
	player.Dir.X = 1.0
}

func main() {
	setup()
	pixelgl.Run(run)
}
