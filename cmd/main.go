package main

import (

	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino"
)

func main() {
	ebiten.SetWindowSize(yoshino.Width, yoshino.Height)
	ebiten.SetWindowTitle("yoshino(GAL)")
	icon, _, err := ebitenutil.NewImageFromFile("../data/icon.jpg")
	if err != nil {
		log.Println("load icon err:", err)
	} else {
		ebiten.SetWindowIcon([]image.Image{icon})
	}
	if err := ebiten.RunGame(yoshino.NewGame()); err != nil {
		log.Fatal(err)
	}
}
