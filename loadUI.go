package yoshino

import (
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type SaveUI struct {
	BaseScene

	Player Player //加载或导出时使用的存档

	Players       []Player
	ui            *ebitenui.UI
	buttons       []*widget.Button
	confirmwindow *widget.Window
}

func (s *SaveUI) Load(st State, sm stagehand.SceneController[State]) {
	s.State = State{Page: PageLoad}
	s.sm = sm.(*stagehand.SceneManager[State])
	s.Players, _ = StdFileSystem.LoadPlayers()
	s.makeui()
}

func (s *SaveUI) Unload() State {
	s.confirmwindow = nil
	s.Players = nil
	s.ui = nil
	s.buttons = nil
	return s.State
}

func (s *SaveUI) Update() error {
	s.ui.Update()
	return nil
}

func (s *SaveUI) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 250, 250, 255})
	s.ui.Draw(screen)
}

func (s *SaveUI) makeui() {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	for k, v := range s.Players {
		if v.ID != "" {
			//非空按钮
			s.buttons = append(s.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadButtonImageByImage(v.screenEbitenImage)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("覆盖保存", k)
					s.OpenWindows(func() {
						s.Players[k] = s.Player
						StdFileSystem.SavePlayers(s.Players)
						s.sm.SwitchTo(&SaveUI{Player: s.Player})
					})
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		} else {
			s.buttons = append(s.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadNoDataButtonImage()),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("保存", k)
					s.OpenWindows(func() {
						s.Players[k] = s.Player
						StdFileSystem.SavePlayers(s.Players)
						s.sm.SwitchTo(&SaveUI{Player: s.Player})
					})
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
		widget.ButtonOpts.Text("返回", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
		widget.ButtonOpts.TextProcessBBCode(true),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   20,
			Right:  20,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("返回按钮被点击")
			s.sm.SwitchTo(&GameUI{Player: s.Player})
		}),
		widget.ButtonOpts.DisableDefaultKeys(),
	))

	rootContainer.AddChild(btcont)
	rootContainer.AddChild(menu)
	s.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}

type LoadUI struct {
	BaseScene

	Player Player //加载或导出时使用的存档

	LastPage      int
	Players       []Player
	ui            *ebitenui.UI
	buttons       []*widget.Button
	confirmwindow *widget.Window
}

func (l *LoadUI) Load(st State, sm stagehand.SceneController[State]) {
	l.LastPage = st.Page
	l.State = State{Page: PageSave}
	l.sm = sm.(*stagehand.SceneManager[State])
	l.Players, _ = StdFileSystem.LoadPlayers()
	l.makeui()
}

func (l *LoadUI) Unload() State {
	l.confirmwindow = nil
	l.Players = nil
	l.ui = nil
	l.buttons = nil
	return l.State
}

func (l *LoadUI) Update() error {
	l.ui.Update()
	return nil
}

func (l *LoadUI) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 250, 250, 255})
	l.ui.Draw(screen)
}

func (l *LoadUI) makeui() {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	for k, v := range l.Players {
		if v.ID != "" {
			//非空按钮
			l.buttons = append(l.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadButtonImageByImage(v.screenEbitenImage)),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("点击存档", k)
					l.OpenWindows(func() {
						l.sm.SwitchTo(&GameUI{Player: l.Players[k]})
					}, "确定加载此存档?")
				}),
				widget.ButtonOpts.DisableDefaultKeys(),
			))
		} else {
			l.buttons = append(l.buttons, widget.NewButton(
				widget.ButtonOpts.Graphic(LoadNoDataButtonImage()),
				widget.ButtonOpts.GraphicPadding(widget.Insets{Left: 10, Top: 10}),
				widget.ButtonOpts.Image(LoadRransparentButtonImage()),
				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					log.Println("点击空存档", k)
					l.OpenWindows(func() {}, "此存档为空! ")
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
		widget.ButtonOpts.Text("返回", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   20,
			Right:  20,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("返回按钮被点击")
			switch l.LastPage {
			case PageMenu:
				l.sm.SwitchTo(&MenuUI{})
			case PageGame:
				l.sm.SwitchTo(&GameUI{Player: l.Player})
			}
		}),
		widget.ButtonOpts.DisableDefaultKeys(),
	))

	rootContainer.AddChild(btcont)
	rootContainer.AddChild(menu)
	l.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}

func createWindow(text string, actionf func(), closef *func()) *widget.Window {
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
			widget.GridLayoutOpts.Padding(&widget.Insets{Top: 120}),
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
			widget.ButtonOpts.Text(v, StdFonts[0].FacePointer(40), LoadBlueButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(&widget.Insets{
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
				} else {
					(*closef)()
				}
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		))
	}
	rootContainer.AddChild(widget.NewText(
		widget.TextOpts.Padding(&widget.Insets{Bottom: 100}),
		widget.TextOpts.Text(text, StdFonts[0].FacePointer(40), color.RGBA{135, 206, 250, 255}),
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
func (s *SaveUI) OpenWindows(actionf func()) {
	var closef *func() = new(func())
	if s.confirmwindow == nil {
		s.confirmwindow = createWindow("确定要保存在这里吗?", actionf, closef) //Colse 如果确认跳转会导致空指针,所以作为else选项
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
func (l *LoadUI) OpenWindows(actionf func(), text string) {
	var closef *func() = new(func())
	if l.confirmwindow == nil {
		l.confirmwindow = createWindow(text, actionf, closef) // s.confirmwindow.Close
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
