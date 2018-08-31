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
	Pos     pixel.Vec
	Dir     pixel.Vec
	Speed   float64
	Texture image.Image
}

func (p *Player) LoadTextures() {
	playerTextureFile, err := os.Open("./assets/player.png")
	if err != nil {
		panic(err)
	}
	defer playerTextureFile.Close()

	playerTextures, err := png.Decode(playerTextureFile)
	resizedPlayerTextures := resize.Resize(uint(scale*5), 0, playerTextures, resize.NearestNeighbor)
	player.Texture = resizedPlayerTextures
}

func (p Player) Draw(m *image.RGBA) {
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
}
