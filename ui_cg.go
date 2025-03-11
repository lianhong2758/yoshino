package yoshino

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type CGUI struct{}

func (*CGUI) Init(g *Game)   {}
func (*CGUI) Clear(g *Game)  {}
func (*CGUI) Update(g *Game) {}
func (*CGUI) Draw(g *Game, screen *ebiten.Image) {}
