package yoshino

import (
	"encoding/json"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type SettingUI struct {
	ui *ebitenui.UI
}
type Config struct {
	Resolution struct{ Width, Height int }
}

func (s *SettingUI) Init(g *Game) {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		),
		))

	s.ui = &ebitenui.UI{
		Container: root,
	}
}
func (s *SettingUI) Clear(g *Game) {}
func (s *SettingUI) Update(g *Game) {
	s.ui.Update()
}
func (s *SettingUI) Draw(g *Game, screen *ebiten.Image) {
	s.ui.Draw(screen)
}

func NewConfig() *Config {
	return &Config{
		Resolution: struct {
			Width  int
			Height int
		}{1600, 900},
	}
}

func (g *Game) LoadConfig() {
	if !g.FileSystem.ItemExists("config.json") {
		g.Config = NewConfig()
		return
	}
	data, err := g.FileSystem.LoadItem("config.json")
	if err != nil {
		log.Println("load err: ", err)
		g.Config = NewConfig()
		return
	}
	g.Config = new(Config)
	err = json.Unmarshal(data, g.Config)
	if err != nil {
		log.Println("load err: ", err)
		g.Config = NewConfig()
		return
	}
}
func (g *Game) SaveConfig() {
	data, _ := json.Marshal(g.Config)
	err := g.FileSystem.SaveItem("config.json", data)
	if err != nil {
		log.Println("save err: ", err)
	}
}
