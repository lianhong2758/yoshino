package yoshino

import (
	"bytes"
	"encoding/gob"
	"image/color"
	"image/png"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quasilyte/gdata"
)

type SaveUI struct {
	Players []Player
	ui      *ebitenui.UI
	buttons []*widget.Button
}

func (s *SaveUI) Init(g *Game) {
	//g.SavePlayers([]Player{g.Player})
	//ui
	clear(s.Players)
	s.buttons = []*widget.Button{}
	g.LoadPlayers()
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	// nilbutton := widget.NewButton(
	// 	widget.ButtonOpts.Graphic(LoadNoDataButtonImage(g)),
	// 	widget.ButtonOpts.Image(LoadRransparentButtonImage()),
	// 	// specify the button's text, the font face, and the color
	// 	//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
	// 	widget.ButtonOpts.Text("", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
	// 	widget.ButtonOpts.TextProcessBBCode(true),
	// 	// specify that the button's text needs some padding for correct display
	// 	widget.ButtonOpts.TextPadding(widget.Insets{
	// 		Left:   20,
	// 		Right:  20,
	// 		Top:    5,
	// 		Bottom: 5,
	// 	}),
	// 	// add a handler that reacts to clicking the button
	// 	widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
	// 		log.Println(" 按钮被点击")

	// 	}),
	// 	widget.ButtonOpts.DisableDefaultKeys(),
	// )

	for range s.Players {
		//...
	}
	for range 12 - len(s.Players) {
		s.buttons = append(s.buttons, widget.NewButton(
			widget.ButtonOpts.Graphic(LoadNoDataButtonImage(g)),
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			//widget.ButtonOpts.Text("", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
			//widget.ButtonOpts.TextProcessBBCode(true),
			// specify that the button's text needs some padding for correct display
			// widget.ButtonOpts.TextPadding(widget.Insets{
			// 	Left:   20,
			// 	Right:  20,
			// 	Top:    5,
			// 	Bottom: 5,
			// }),
			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println(" 按钮被点击")

			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		))
	}
	/*
	   x x x x
	   x x x x
	   x x x x
	*/
	// 创建网格布局容器 4*3
	btcont := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(4),      // 单列布局
			widget.GridLayoutOpts.Spacing(20, 30), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(Width-200, Height-400), // 设置最小宽度
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	for _, v := range s.buttons {
		btcont.AddChild(v)
	}
	// 创建网格布局容器（单行）
	menu := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),    // 单列布局
			widget.GridLayoutOpts.Spacing(0, 0), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(40, 30), // 设置最小宽度 100px
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	menu.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(LoadRransparentButtonImage()),
		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text("返回", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   20,
			Right:  20,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("返回按钮被点击")
			//g.Next(StatusSave)
			g.Next(StatusGame)
		}),
		widget.ButtonOpts.DisableDefaultKeys(),
	))

	rootContainer.AddChild(btcont)
	rootContainer.AddChild(menu)
	s.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}
func (s *SaveUI) Clear(g *Game) {}
func (s *SaveUI) Update(g *Game) {
	s.ui.Update()
}
func (s *SaveUI) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 250, 250, 255})
	s.ui.Draw(screen)
}

type LoadUI struct {
	Players []Player
}

func (*LoadUI) Init(g *Game)                       {}
func (*LoadUI) Clear(g *Game)                      {}
func (*LoadUI) Update(g *Game)                     {}
func (*LoadUI) Draw(g *Game, screen *ebiten.Image) {}

func (g *Game) LoadPlayers() ([]Player, error) {
	if !g.FileSystem.ItemExists("players.gob") {
		return []Player{}, nil
	}
	data, err := g.FileSystem.LoadItem("players.gob")
	if err != nil {
		log.Println("load err: ", err)
		return nil, err
	}
	players := make([]Player, 0)
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&players)
	if err != nil {
		log.Println("load err: ", err)
		return nil, err
	}
	for k := range len(players) {
		players[k].screenContent, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(players[k].ScreenData))
	}
	return players, nil
}

func (g *Game) SavePlayers(p []Player) error {
	for k := range len(p) {
		var picbuff bytes.Buffer
		_ = png.Encode(&picbuff, g.Player.screenContent)
		p[k].ScreenData = picbuff.Bytes()
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(p)
	if err != nil {
		log.Println("save err: ", err)
		return err
	}
	err = g.FileSystem.SaveItem("players.gob", buff.Bytes())
	if err != nil {
		log.Println("loaderr: ", err)
		return err
	}
	return nil
}

func (g *Game) LoadFileSystem() {
	m, err := gdata.Open(gdata.Config{
		AppName: "yoshino",
	})
	if err != nil {
		panic(err)
	}
	g.FileSystem = m
}
