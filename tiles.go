package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
)

type SelectedTile struct {
	Pos     image.Point
	Tick    int
	Texture image.Image
	Active  bool
}

func (tile *SelectedTile) LoadTextures() {
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
}
