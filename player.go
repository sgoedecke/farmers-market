package main

import (
	"github.com/faiface/pixel"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Player struct {
	Pos       pixel.Vec
	Dir       pixel.Vec
	Speed     float64
	Texture   image.Image
	TexturePt image.Point
}

func (p *Player) LoadTextures() {
	playerTextureFile, err := os.Open("./assets/hero.png") // 20px/20px for each sprite
	if err != nil {
		panic(err)
	}
	defer playerTextureFile.Close()

	textureMagnification := uint(scale * 10)

	playerTextures, err := png.Decode(playerTextureFile)
	resizedPlayerTextures := resize.Resize(textureMagnification, 0, playerTextures, resize.NearestNeighbor)
	player.Texture = resizedPlayerTextures
	player.Dir.X = 1 // set initial direction so a texture sprite loads
}

// based on the player's direction and the current tick, returns the top-left point
// for the current image to draw.
func (p *Player) SetActiveTextureCoord(tick int) {
	if player.Dir.X == 0 && player.Dir.Y == 0 {
		// leave the player's texture pt what it was before
	} else {
		ty := 10
		tx := 0

		texWidth := 37 // the width of each sprite on our rescaled texture sheet. nfi why.

		// our sprite sheet has three sprites per walking frame, and max tick is 29. so we divide by 10 to get a
		// smooth three-frame walking animation
		if player.Dir.X == -1 {
			tx = (9 + tick/10) * texWidth // start at the 9th sprite on the sheet (0-indexed)
		}
		if player.Dir.X == 1 {
			tx = (0 + tick/10) * texWidth
		}

		if player.Dir.Y == -1 {
			tx = (3 + tick/10) * texWidth
		}
		if player.Dir.Y == 1 {
			tx = (6 + tick/10) * texWidth
		}
		player.TexturePt = image.Pt(tx, ty)
	}
}

func (p Player) Draw(m *image.RGBA) {
	width := 0.8
	// magic numbers 1 and 2 here correspond to tile width
	draw.Draw(m,
		image.Rect(int(player.Pos.X*scale), int(player.Pos.Y*scale), int((player.Pos.X+width)*scale), int((player.Pos.Y+1)*scale)),
		player.Texture,
		player.TexturePt,
		draw.Over) // need to use Over rather than Src here to respect transparent background
}
