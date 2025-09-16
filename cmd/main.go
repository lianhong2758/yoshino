package main

import (
	"image"
	"log"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joelschutz/stagehand"
	"github.com/lianhong2758/yoshino"
	"github.com/lianhong2758/yoshino/file"
	_ "golang.org/x/image/webp"
)

func main() {
	ebiten.SetWindowSize(yoshino.Width, yoshino.Height)
	ebiten.SetWindowTitle("yoshino(GAL)")
	ebiten.SetTPS(60)
	iconf := file.OpenMaterial("icon.jpg")
	icon, _, err := ebitenutil.NewImageFromReader(iconf)
	if err != nil {
		log.Println("load icon err:", err)
	} else {
		ebiten.SetWindowIcon([]image.Image{icon})
	}
	iconf.Close()
	if err := ebiten.RunGame(stagehand.NewSceneManager(&yoshino.TitleUI{}, yoshino.State{})); err != nil {
		log.Fatal(err)
	}
}
