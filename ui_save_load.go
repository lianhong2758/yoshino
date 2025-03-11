package yoshino

import "github.com/hajimehoshi/ebiten/v2"

type SaveUI struct{}

func (*SaveUI) Init(g *Game)                       {}
func (*SaveUI) Clear(g *Game)                      {}
func (*SaveUI) Update(g *Game)                     {}
func (*SaveUI) Draw(g *Game, screen *ebiten.Image) {}

type LoadUI struct{}

func (*LoadUI) Init(g *Game)                       {}
func (*LoadUI) Clear(g *Game)                      {}
func (*LoadUI) Update(g *Game)                     {}
func (*LoadUI) Draw(g *Game, screen *ebiten.Image) {}
