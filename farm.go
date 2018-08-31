package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"
)

type Player struct {
	Pos     pixel.Vec
	Dir     pixel.Vec
	Speed   float64
	Texture image.Image
}

type World struct {
	Map         [20][20]int
	Width       int
	Height      int
	BaseTexture image.Image
}

var (
	fullscreen = false
	scale      = 50.0
	world      = World{}
	fps        = 60.0
	player     = Player{}
)

func setup() {
  player.Pos = pixel.Vec{15,15}
  player.Dir = pixel.Vec{0,0}
  player.Speed = 0.2
	world.Height = 20
	world.Width = 20

	playerTextureFile, err := os.Open("./assets/player.png")
	defer playerTextureFile.Close()

	playerTextures, err := png.Decode(playerTextureFile)
	resizedPlayerTextures := resize.Resize(uint(scale * 5), 0, playerTextures, resize.NearestNeighbor)
	player.Texture = resizedPlayerTextures

	// load textures
	textureFile, err := os.Open("./assets/nature-tileset.png")
	defer textureFile.Close()

	textures, err := png.Decode(textureFile)
	if err != nil {
		panic(err)
	}

	resizedTextures := resize.Resize(uint(scale*20), 0, textures, resize.NearestNeighbor)

    // generate base texture by tiling the grass tile to the size of the world
	baseTexture := image.NewRGBA(image.Rect(0, 0, world.Width*int(scale), world.Height*int(scale)))
	baseTextureSize := 50 // size of the grass tile in nature-tileset.png
	for x := 0; x <= baseTexture.Rect.Dx(); x += baseTextureSize {
		for y := 0; y <= baseTexture.Rect.Dy(); y += baseTextureSize {
			draw.Draw(baseTexture,
				image.Rect(x, y, x+baseTextureSize, y+baseTextureSize),
				resizedTextures,
				image.Pt(0, 0),
				draw.Src)
		}
	}

	world.BaseTexture = baseTexture
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

	// OPTIMIZATION: start by drawing the default texture across everything. this lets us avoid drawing it
	// every single tile one at a time
	draw.Draw(m,
		m.Bounds(),
		world.BaseTexture,
		image.Pt(0, 0),
		draw.Src)

	// draw tiles
	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			tile := world.Map[x][y]

			// OPTIMIZATION: if the tile is not in view, don't bother drawing it to the buffer
			if math.Abs(player.Pos.X-float64(x)) > 15 || math.Abs(player.Pos.Y-float64(y)) > 15 {
				continue
			}

			if tile == 0 {
				// don't draw, since we've already got the standard texture

			} else {
				fmt.Println(tile)
                // TODO: add non-player entities
			}

		}
	}

	// draw player

    tx := 18
    ty := 25

    if player.Dir.X == -1 {
      ty = 360
    }
    if player.Dir.X == 1 {
      ty = 140
    }

    if player.Dir.Y == -1 {
      ty = 25
    }
    if player.Dir.Y == 1 {
      ty = 250
    }
	draw.Draw(m,
		image.Rect(int(player.Pos.X*scale), int(player.Pos.Y*scale), int((player.Pos.X+1)*scale), int((player.Pos.Y+2)*scale)),
		player.Texture,
		image.Pt(tx, ty),
		draw.Over) // need to use Over rather than Src here to respect transparent background

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
		} else {
			fmt.Println("Running slow")
			fmt.Println(time.Duration((1/fps - dt) * 1000000000))
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
