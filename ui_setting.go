package yoshino

import "github.com/hajimehoshi/ebiten/v2"

type SettingUI struct{}

func (*SettingUI) Init(g *Game)                       {}
func (*SettingUI) Clear(g *Game)                      {}
func (*SettingUI) Update(g *Game)                     {}
func (*SettingUI) Draw(g *Game, screen *ebiten.Image) {}
