package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"image/color"
	"time"
)

type Player struct {
	Pos pixel.Vec
	Dir pixel.Vec
    Speed float64
}

type World struct {
	Map    [50][50]int
	Width  int
	Height int
}

var (
	fullscreen = false
	scale      = 10.0
	world      = World{}
	fps        = 100.0

	player = Player{pixel.Vec{25, 25}, pixel.Vec{0, 0}, 0.4}
)

func setup() {
	world.Height = 50
	world.Width = 50
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

	m := image.NewRGBA(image.Rect(0, 0, world.Width * int(scale), world.Height * int(scale)))

    // draw tiles
	var c color.RGBA
	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			tile := world.Map[x][y]
			if tile == 0 {
				c = color.RGBA{200, 200, 200, 1} // black

			} else {
				c = color.RGBA{0, 0, 0, 1}
			}

            for ix := 0; ix < int(scale); ix++ { // TODO: replace with Tile.Draw & a texture
              for iy := 0; iy < int(scale); iy++ {
                m.Set((x * int(scale)) + ix, (y * int(scale)) + iy, c)
              }
            }

		}
	}

    // draw player
    c = color.RGBA{200, 0, 0, 1}
    for ix := 0; ix < int(scale); ix++ { // TODO: replace with Tile.Draw & a texture/animation if walking
      for iy := 0; iy < int(scale); iy++ {
        m.Set(int((player.Pos.X * scale) + float64(ix)), int((player.Pos.Y * scale) + float64(iy)), c)
      }
    }

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
			fmt.Println("Sleeping!")
			fmt.Println(1/fps - dt)
			time.Sleep(time.Duration((1/fps - dt) * 1000000000)) // have to convert seconds to ns for the cast
		}
		last = time.Now()

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

		// draw
		win.Clear(color.Black)
		p := pixel.PictureDataFromImage(frame())

		// offset the center of the screen so camera-following works
		c := win.Bounds().Center().Add(pixel.Vec{float64(world.Width) * 0.5, float64(world.Height) * -0.5}.Scaled(scale))

		// since we store player coordinates on the 2d array (0,0 is top left) but draw with
		// Euclidean coordinates, we need to flip the X coord
		playerPos := pixel.Vec{player.Pos.X * -1.0, player.Pos.Y}.Scaled(scale)
        fmt.Println(playerPos)

		d := playerPos.Add(c)

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(d))

		win.Update()
	}
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
