package main

import (
	"image"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino"
	"github.com/lianhong2758/yoshino/file"
	_ "golang.org/x/image/webp"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	ebiten.SetWindowSize(yoshino.Width, yoshino.Height)
	ebiten.SetWindowTitle("yoshino(GAL)")
	iconf := file.OpenMaterial("icon.jpg")
	icon, _, err := ebitenutil.NewImageFromReader(iconf)
	if err != nil {
		log.Println("load icon err:", err)
	} else {
		ebiten.SetWindowIcon([]image.Image{icon})
	}
	iconf.Close()
	if err := ebiten.RunGame(yoshino.NewGame()); err != nil {
		log.Fatal(err)
	}
}
