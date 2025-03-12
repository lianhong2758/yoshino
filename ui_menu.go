package yoshino

import (
	"bytes"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino/file"
)

type MenuUI struct {
	ui       *ebitenui.UI
	btns     []*widget.Button //预定5个
	MenuFile []*ebiten.Image
}

func (m *MenuUI) Init(g *Game) {
	//	img, _, err := ebitenutil.NewImageFromFile("../data/menu.jpg")
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial("menu.jpg")))
	if err != nil {
		log.Println(err)
		return
	}
	m.MenuFile = append(m.MenuFile, img)
	//根容器使用锚点布局
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	m.btns = append(m.btns,
		//开始游戏
		widget.NewButton(
			widget.ButtonOpts.Image(LoadButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("新的游戏", g.FontFace[0], LoadButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			//Move the text down and right on press
			widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
				m.btns[0].Text().Inset.Top = 4
				m.btns[0].Text().Inset.Left = 4
			}),
			//Move the text back to start on press released
			widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
				m.btns[0].Text().Inset.Top = 0
				m.btns[0].Text().Inset.Left = 0
			}),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("新的游戏")
				g.Next(StatusGame)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		//继续游戏
		widget.NewButton(
			widget.ButtonOpts.Image(LoadButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("加载游戏", g.FontFace[0], LoadButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			//Move the text down and right on press
			widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
				m.btns[1].Text().Inset.Top = 4
				m.btns[1].Text().Inset.Left = 4
			}),
			//Move the text back to start on press released
			widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
				m.btns[1].Text().Inset.Top = 0
				m.btns[1].Text().Inset.Left = 0
			}),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("加载游戏")
				g.Next(StatusLoad)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		//设置
		widget.NewButton(
			widget.ButtonOpts.Image(LoadButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("设置", g.FontFace[0], LoadButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			//Move the text down and right on press
			widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
				m.btns[2].Text().Inset.Top = 4
				m.btns[2].Text().Inset.Left = 4
			}),
			//Move the text back to start on press released
			widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
				m.btns[2].Text().Inset.Top = 0
				m.btns[2].Text().Inset.Left = 0
			}),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("设置")
				g.Next(StatusSetting)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		//退出
		widget.NewButton(
			widget.ButtonOpts.Image(LoadButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("退出游戏", g.FontFace[0], LoadButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			//Move the text down and right on press
			widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
				m.btns[3].Text().Inset.Top = 4
				m.btns[3].Text().Inset.Left = 4
			}),
			//Move the text back to start on press released
			widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
				m.btns[3].Text().Inset.Top = 0
				m.btns[3].Text().Inset.Left = 0
			}),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("退出游戏")
				os.Exit(0)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
	)

	// 创建网格布局容器（单列）
	grid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),      // 单列布局
			widget.GridLayoutOpts.Spacing(10, 10), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 100), // 设置最小宽度 100px
		),
	)
	// 添加按钮到网格
	for _, btn := range m.btns {
		grid.AddChild(btn)
	}
	// 创建锚点布局容器（用于定位到左侧）
	anchorContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{Left: 150}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	anchorContainer.AddChild(grid)
	rootContainer.AddChild(anchorContainer)
	m.ui = &ebitenui.UI{
		Container: rootContainer,
	}

}
func (m *MenuUI) Clear(g *Game) {
	m.MenuFile = []*ebiten.Image{}
	m.btns = []*widget.Button{}
	m.ui = nil
}
func (m *MenuUI) Update(g *Game) { m.ui.Update() }
func (m *MenuUI) Draw(g *Game, screen *ebiten.Image) {
	//背景图层
	screen.DrawImage(m.MenuFile[0], DrawImageCentreOption(m.MenuFile[0]))
	//操作图层
	m.ui.Draw(screen)
}
