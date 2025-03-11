package yoshino

import "github.com/hajimehoshi/ebiten/v2"

type TreeUI struct{}

func (*TreeUI) Init(g *Game)                       {}
func (*TreeUI) Clear(g *Game)                      {}
func (*TreeUI) Update(g *Game)                     {}
func (*TreeUI) Draw(g *Game, screen *ebiten.Image) {}
