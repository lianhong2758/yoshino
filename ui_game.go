package yoshino

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lianhong2758/yoshino/file"
)

type GameUI struct {
	ui              *ebitenui.UI
	Background      *ebiten.Image
	selectionwindow *widget.Window
	//Character  *ebiten.Image
	//Words      string       //完整的台词
	newString  func() string
	counter    int //计数器,逐字打印需要
	needchange bool
	nextid     string      //下一个id
	rep        *Repertoire //当前剧本

	DrawString     func(string) //修改正在显示的文字
	DrawAvatar     func(string) //修改头像
	DrawCreation   func(string) //修改立绘
	DrawBackground func(string) //修改背景
	PlayMusic      func(string)
	PlayVideo      func(string)

	AudioContext *audio.Context
	BgmPlayer    *audio.Player // 背景音乐播放器
	VoicePlayer  *audio.Player // 语音播放器
}

func (gu *GameUI) Init(g *Game) {
	ScriptInit()
	gu.nextid = "1"
	gu.needchange = true
	gu.Background = ebiten.NewImage(1, 1)

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
		widget.GraphicOpts.Image(ebiten.NewImage(1, 1)),
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

	gu.DrawString = func(s string) { label1.Label = s }
	gu.DrawAvatar = func(s string) {
		avatar.Image, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial(s)))
	}
	gu.DrawCreation = func(s string) {}
	gu.DrawBackground = func(s string) {
		gu.Background, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial(s)))
	}
	gu.PlayMusic = func(s string) {}
	gu.PlayVideo = func(s string) {}
}
func (*GameUI) Clear(g *Game) {}
func (gu *GameUI) Update(g *Game) {
	if gu.needchange {
		//change
		gu.rep = LoadNextRepertoire(gu.nextid)
		switch gu.rep.Types {
		case "A": //常规类型
			gu.newString = StreamString(gu.rep.Role + ": " + gu.rep.Text)
			gu.DrawAvatar(gu.rep.Avatar)
			gu.DrawCreation(gu.rep.Creation)
			gu.DrawBackground(gu.rep.Background)
			gu.PlayMusic(gu.rep.Music)
			gu.PlayVideo(gu.rep.Video)
			gu.needchange = false
			gu.nextid = gu.rep.Next
		case "B": //CG
		case "C":
			//选择界面
			gu.OpenWindows(g)
		case "D": //个人线判断
			id := gu.rep.Map[strconv.Itoa(g.Player.Token)]
			if id == "" {
				log.Print("存档损坏")
				os.Exit(1)
			}
			gu.nextid = id
			return
		}
	}
	gu.counter++
	if gu.counter == 5 {
		gu.counter = 0
		gu.DrawString(gu.newString())
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		gu.needchange = true
	}
	gu.ui.Update()
}

func (gu *GameUI) Draw(g *Game, screen *ebiten.Image) {
	if gu.rep == nil {
		return
	}
	//背景
	switch gu.rep.BackgroundType {
	case "image":
		screen.DrawImage(gu.Background, DrawImageCentreOption(gu.Background))
	}
	//人物
	gu.ui.Draw(screen)
}

func (gu *GameUI) createWindow(g *Game) *widget.Window {
	// Create the contents of the window
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),      // 单列布局
			widget.GridLayoutOpts.Spacing(10, 10), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 100), // 设置最小宽度 100px
		),
	)
	//gu.selectionwindow.Close()
	windowContainer.AddChild(widget.NewText(

		widget.TextOpts.Text(gu.rep.Select[0].Text, g.FontFace[0], color.NRGBA{254, 255, 255, 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	return widget.NewWindow(
		// Set the main contents of the window
		widget.WindowOpts.Contents(windowContainer),
		// Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
	)
}
func (gu *GameUI) OpenWindows(g *Game) {
	if gu.selectionwindow == nil {
		gu.selectionwindow = gu.createWindow(g)
	}
	if !gu.ui.IsWindowOpen(gu.selectionwindow) {
		log.Println("打开选择窗口")
		// Get the preferred size of the content
		x, y := gu.selectionwindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		gu.selectionwindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		gu.ui.AddWindow(gu.selectionwindow)
	}
}
