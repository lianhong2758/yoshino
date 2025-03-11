package yoshino

import (
	"bytes"
	"image/color"
	"log"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino/file"
)

type GameUI struct {
	ui           *ebitenui.UI
	Background   *ebiten.Image
	Character    []*ebiten.Image
	Words        string       //完整的台词
	DrawString   func(string) //修改正在显示的文字
	newString    func() string
	counter      int
	stringchange bool

	AudioContext *audio.Context
	BgmPlayer    *audio.Player // 背景音乐播放器
	VoicePlayer  *audio.Player // 语音播放器
}

var q = 0

func (gu *GameUI) Init(g *Game) {
	gu.stringchange = true
	gu.Words = "丛雨: 早安,我是丛雨."
	c, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(file.Congyu))
	if err != nil {
		log.Println(err)
		return
	}
	gu.Character = append(gu.Character, c)
	gu.Background, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(file.Background))
	if err != nil {
		log.Println(err)
		return
	}

	//ui
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	label1 := widget.NewText(
		widget.TextOpts.Text("", g.FontFace[0], color.Black),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	gu.DrawString = func(s string) { label1.Label = s }
	textbox := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
				Left:   20,
				Right:  20,
			}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
			widget.WidgetOpts.MinSize(900, 200),
		),
	)
	textbox.AddChild(label1)

	// 角色头像
	avatar := widget.NewGraphic(
		widget.GraphicOpts.Image(gu.Character[0]),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	// 主菜单按钮
	menubu := widget.NewButton(
		widget.ButtonOpts.Image(LoadButtonImage()),
		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text("主菜单", g.FontFace[0], LoadButtonTextColor()),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("主菜单按钮被点击")
			g.Next(StatusMenu)
		}),
		widget.ButtonOpts.DisableDefaultKeys(),
	)
	// 创建网格布局容器（单行）
	menu := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(5),     // 单列布局
			widget.GridLayoutOpts.Spacing(0, 10), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 50), // 设置最小宽度 100px
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	menu.AddChild(menubu)
	// 布局设置
	rootContainer.AddChild(textbox)
	rootContainer.AddChild(avatar)
	rootContainer.AddChild(menu)

	gu.ui = &ebitenui.UI{
		Container: rootContainer,
	}
	//time
	g.startTime = time.Now()
}
func (*GameUI) Clear(g *Game) {}
func (gu *GameUI) Update(g *Game) {
	if gu.stringchange {
		gu.newString = StreamString(gu.Words)
		gu.stringchange = false
	}
	gu.counter++
	if gu.counter == 1 {
		gu.counter = 0
		gu.DrawString(gu.newString())
	}

	by := time.Since(g.startTime)
	if by.Seconds() > 2.5 && q == 0 {
		gu.Words = "丛雨: Ciallo~"
		gu.stringchange = true
		q = 1
	}
	gu.ui.Update()
}

func (gu *GameUI) Draw(g *Game, screen *ebiten.Image) {
	//背景
	screen.DrawImage(gu.Background, DrawImageCentreOption(gu.Background))
	//人物
	gu.ui.Draw(screen)
}
