package yoshino

import (
	"fmt"
	_ "image/jpeg"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	Width, Height = 1600, 900
)

type Game struct {
	startTime time.Time
	Status    GameStatus      //状态机
	Player    Player          //存档
	GameUI    [StatusTree + 1]UI //期望是与GameStatus对应的
	FontFace  []text.Face
}

func (g *Game) Update() error {
	g.GameUI[g.Status].Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.GameUI[g.Status].Draw(g, screen)
	//fps文字图层
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func NewGame() *Game {
	g := Game{Status: 0, GameUI: [StatusTree + 1]UI{
		&TitleUI{}, &MenuUI{}, &GameUI{}, &CGUI{}, &SettingUI{}, &SaveUI{}, &LoadUI{}, &TreeUI{}},
	}
	g.Next(StatusTitle)
	return &g
}
