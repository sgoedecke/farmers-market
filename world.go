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

type SelectedTile struct {
	Pos     image.Point
	Tick    int
	Texture image.Image
	Active  bool
}

type World struct {
	Map             [20][20]int
	Width           int
	Height          int
	BaseTexture     image.Image
	Tick            int
	HighlightedTile SelectedTile
}

func (world *World) LoadTextures() {
	// load textures
	textureFile, err := os.Open("./assets/nature-tileset.png")
	defer textureFile.Close()

	textures, err := png.Decode(textureFile)
	if err != nil {
		panic(err)
	}

	resizedTextures := resize.Resize(uint(scale*15), 0, textures, resize.NearestNeighbor)

	// load highlightedtile texture
	tileTextureFile, err := os.Open("./assets/selectedtile.png") // 20px/20px
	if err != nil {
		panic(err)
	}
	defer tileTextureFile.Close()

	textureMagnification := uint(scale)
	tileTexture, err := png.Decode(tileTextureFile)
	highlightedTileTexture := resize.Resize(textureMagnification, 0, tileTexture, resize.NearestNeighbor)
	world.HighlightedTile.Texture = highlightedTileTexture

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

	// draw shrubs and flowers
	// 12 tiles right, 3 tiles down
	shrubTextureCoords := image.Pt(395, 39)
	shrubTextureSize := 40

	shrubCoordSet := []image.Point{
		image.Pt(400, 400),
		image.Pt(450, 410),
		image.Pt(350, 480),
		image.Pt(150, 710),
	}

	for _, pt := range shrubCoordSet {
		draw.Draw(baseTexture,
			image.Rect(pt.X, pt.Y, pt.X+shrubTextureSize, pt.Y+shrubTextureSize),
			resizedTextures,
			shrubTextureCoords,
			draw.Over)
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

	tileX := int(float64(world.HighlightedTile.Pos.X) * scale)
	tileY := int(float64(world.HighlightedTile.Pos.Y) * scale)
	//draw highlighted tile
	if world.HighlightedTile.Active {
		draw.Draw(m,
			image.Rect(tileX, tileY, tileX+int(scale), tileY+int(scale)),
			world.HighlightedTile.Texture,
			image.Pt(0, 0),
			draw.Over)
	}
}
