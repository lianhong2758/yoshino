package yoshino

import (
	"bytes"
	"encoding/gob"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quasilyte/gdata"
)

type SaveUI struct {
	Players       []Player
	ui            *ebitenui.UI
	buttons       []*widget.Button
	confirmwindow *widget.Window
}

func (s *SaveUI) Init(g *Game) {
	//g.SavePlayers([]Player{g.Player})
	//ui
	clear(s.Players)
	s.buttons = []*widget.Button{}
	s.Players, _ = g.LoadPlayers()
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	for k, v := range s.Players {
		if v.ID != "" {
			//非空按钮
			s.buttons = append(s.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadButtonImageByImage(g, v)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
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
					log.Println("覆盖保存", k)
					s.OpenWindows(g, func() { s.Players[k] = g.Player; g.SavePlayers(s.Players); g.Next(StatusSave) })
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		} else {
			s.buttons = append(s.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadNoDataButtonImage(g)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
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
					log.Println("保存", k)
					s.OpenWindows(g, func() { s.Players[k] = g.Player; g.SavePlayers(s.Players); g.Next(StatusSave) })
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		}

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
func (s *SaveUI) Clear(g *Game) {
	s.confirmwindow = nil
}
func (s *SaveUI) Update(g *Game) {
	s.ui.Update()
}
func (s *SaveUI) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 250, 250, 255})
	s.ui.Draw(screen)
}

type LoadUI struct {
	Players       []Player
	ui            *ebitenui.UI
	buttons       []*widget.Button
	confirmwindow *widget.Window
}

func (l *LoadUI) Init(g *Game) {
	//ui
	clear(l.Players)
	l.buttons = []*widget.Button{}
	l.Players, _ = g.LoadPlayers()
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	for k, v := range l.Players {
		if v.ID != "" {
			//非空按钮
			l.buttons = append(l.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadButtonImageByImage(g, v)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("点击存档", k)
					l.OpenWindows(g, func() {
						FistID = l.Players[k].ID
						g.Player = l.Players[k]
						g.Transition(func() {
							g.Next(StatusGame)
						},
							AnimationTransparent(g),
						)
					}, "确定加载此存档?")
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		} else {
			l.buttons = append(l.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadNoDataButtonImage(g)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("点击空存档", k)
					l.OpenWindows(g, func() {}, "此存档为空! ")
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		}

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
	for _, v := range l.buttons {
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
			g.Next(g.lastState)
		}),
		widget.ButtonOpts.DisableDefaultKeys(),
	))

	rootContainer.AddChild(btcont)
	rootContainer.AddChild(menu)
	l.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}
func (l *LoadUI) Clear(g *Game) {
	l.confirmwindow = nil
}
func (l *LoadUI) Update(g *Game) {
	l.ui.Update()
}
func (l *LoadUI) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 250, 250, 255})
	l.ui.Draw(screen)
}

func (g *Game) LoadPlayers() ([]Player, error) {
	if !g.FileSystem.ItemExists("players.gob") {
		return make([]Player, 12), nil
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
		// players[k].screenContent
		players[k].screenEbitenImage, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(players[k].ScreenData))
	}
	//如果不够12就补齐
	for range 12 - len(players) {
		players = append(players, Player{})
	}
	return players, nil
}

func (g *Game) SavePlayers(p []Player) error {
	for range 12 - len(p) {
		p = append(p, Player{})
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

func createWindow(g *Game, text string, actionf func(), closef *func()) *widget.Window {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.RGBA{255, 218, 185, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(600, 350), // 设置最小宽度 100px
		),
	)
	// Create the contents of the window
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),       //  并列
			widget.GridLayoutOpts.Spacing(100, 10), // 按钮间距 5px
			widget.GridLayoutOpts.Padding(widget.Insets{Top: 120}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 100), // 设置最小宽度 100px
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	for _, v := range []string{"确认", "取消"} {
		ok := v
		windowContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(LoadConfirmButtonImage()),
			widget.ButtonOpts.Text(v, g.FontFace[0].Face(40), LoadBlueButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   20,
				Right:  20,
				Top:    5,
				Bottom: 5,
			}),
			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("用户保存选择了:", v)
				if ok == "确认" {
					//执行操作
					actionf()
				}else{
					(*closef)()
				}
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		))
	}
	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Insets(widget.Insets{Bottom: 100}),
		widget.TextOpts.Text(text, g.FontFace[0].Face(40), color.RGBA{135, 206, 250, 255}),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
	))
	rootContainer.AddChild(windowContainer)
	return widget.NewWindow(
		// Set the main contents of the window
		widget.WindowOpts.Contents(rootContainer),
		// Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
	)
}
func (s *SaveUI) OpenWindows(g *Game, actionf func()) {
	var closef *func() = new(func())
	if s.confirmwindow == nil {
		s.confirmwindow = createWindow(g, "确定要保存在这里吗?", actionf, closef) //Colse 如果确认跳转会导致空指针,所以作为else选项
	}
	(*closef) = func() { s.confirmwindow.Close(); s.confirmwindow = nil }
	if !s.ui.IsWindowOpen(s.confirmwindow) {
		log.Println("打开确认窗口")
		// Get the preferred size of the content
		x, y := s.confirmwindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		s.confirmwindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		s.ui.AddWindow(s.confirmwindow)
	}
}
func (l *LoadUI) OpenWindows(g *Game, actionf func(), text string) {
	var closef *func() = new(func())
	if l.confirmwindow == nil {
		l.confirmwindow = createWindow(g, text, actionf, closef) // s.confirmwindow.Close
	}
	(*closef) = func() { l.confirmwindow.Close(); l.confirmwindow = nil }
	if !l.ui.IsWindowOpen(l.confirmwindow) {
		log.Println("打开确认窗口")
		// Get the preferred size of the content
		x, y := l.confirmwindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		l.confirmwindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		l.ui.AddWindow(l.confirmwindow)
	}
}
