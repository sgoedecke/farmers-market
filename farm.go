package main

import (
	"image"
	"image/color"
    "time"

	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Coord struct {
	X int
	Y int
}

type Player struct {
  Pos Coord
  Dir Coord
}

type World struct {
	Map    [140][140]int
	Width  int
	Height int
}


var (
	fullscreen = true
	scale      = 6.0
	world      = World{}

	player = Player{Coord{5,5},Coord{0,0}}
)

func setup() {
	world.Height = 140
	world.Width = 140
}

func moveEntities() {
  fmt.Println("a")
  newX := player.Pos.X + player.Dir.X
  newY := player.Pos.Y + player.Dir.Y
  if (newX > 0 && newX < world.Width) && (newY > 0 && newY < world.Height) {
    player.Pos.X = newX
    player.Pos.Y = newY
  }
}

// from the global variables that hold the game state, draw a frame of the game into
// the image.RGBA buffer
func frame() *image.RGBA {

	m := image.NewRGBA(image.Rect(0, 0, world.Width, world.Height))

	var c color.RGBA
	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			tile := world.Map[x][y]
			if tile == 0 {
				c = color.RGBA{200, 200, 200, 1} // black

			} else {
				c = color.RGBA{0, 0, 0, 1}
			}

			m.Set(x, y, c)
		}
	}

	c = color.RGBA{200, 0, 0, 1}
	m.Set(player.Pos.X, player.Pos.Y, c)

	return m
}

func run() {

	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(world.Width)*scale, float64(world.Height)*scale),
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

	c := win.Bounds().Center()

	last := time.Now()

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Clear(color.Black)

		dt := time.Since(last).Seconds()
        last = time.Now()
        if dt < 1 {
        moveEntities()
      }

        haltPlayer()
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


		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		win.Update()
	}
}

func haltPlayer() {
  player.Dir.X = 0
  player.Dir.Y = 0
}
func moveUp() {
	player.Dir.Y = -1
}
func moveDown() {
	player.Dir.Y = 1
}
func moveLeft() {
	player.Dir.X = -1
}
func moveRight() {
	player.Dir.X = 1
}

func main() {
	setup()
	pixelgl.Run(run)
}
