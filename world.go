package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
)

type World struct {
	Map         [20][20]int
	Width       int
	Height      int
	BaseTexture image.Image
}

func (world *World) LoadTextures() {
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

func (world World) Draw(m *image.RGBA) {
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
}
