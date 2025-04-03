package yoshino

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/lianhong2758/yoshino/file"
)

var FistID string = "1"

type GameUI struct {
	ui *ebitenui.UI
	//临时参数,在剧目之间会刷新
	backgroundImage *ebiten.Image
	avatarImage     *ebiten.Image
	selectionwindow *widget.Window
	historywindow   *widget.Window
	creation        [3]creationOpt
	newString       func(bool) (string, bool) //输出的字符串,是否输出全句
	isAllString     bool                      //记录是否输出全句
	counter         int                       //计数器,逐字打印需要
	nowWord         string                    //当前显示的对话内容
	doingchange     bool                      //是否切换下一剧目
	nextid          string                    //下一个id
	rep             *Repertoire               //当前剧本
	history         []string                  //存放id?还可以扩展语音播放等
	doingTransition bool                      // 用于主UI的过渡动画执行后改变needchange
	hide            bool                      //隐藏ui
	lastMusic       string                    //上一个背景音乐,用于判断是否切换背景音乐

	LoadString     func(string) //修改正在显示的文字
	LoadAvatar     func(string) //修改头像
	LoadBackground func(string) //修改背景
	// PlayMusic      func(string)
	// PlayVideo      func(string)

	DoAction     [3]func() //在updata里面执行可能存在的action动画
	DoTransition func(screen *ebiten.Image) bool

	AudioContext *audio.Context
	MusicPlayer  *audio.Player // 背景音乐播放器
	VoicePlayer  *audio.Player // 语音播放器
}

type creationOpt struct {
	Image *ebiten.Image
	O     *ebiten.DrawImageOptions
}

func (gu *GameUI) Init(g *Game) {
	ScriptInit()
	gu.nextid = FistID
	FistID = "1" //重置,避免加载存档后利用firstid导致无法开启新游戏
	gu.doingchange = true
	gu.DoAction = [3]func(){nilActionFunc(), nilActionFunc(), nilActionFunc()}
	gu.backgroundImage = ebiten.NewImage(1, 1)
	if t := audio.CurrentContext(); t == nil {
		gu.AudioContext = audio.NewContext(48000)
	} else {
		gu.AudioContext = t
	}

	/*ui
	  /背景
	  /立绘
	  	角色1	角色2	角色3


	  /文本区	 【角色名】
	  头像op   ⌈对话内容⌋
	  					保存 读取 设置 主菜单
	*/

	//ui
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	//菜单区
	buttons := []*widget.Button{
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("读取", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("读取按钮被点击")
				g.Next(StatusLoad)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("保存", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("保存按钮被点击")
				//g.Next(StatusSave)
				FistID, g.Player.ID = gu.rep.ID, gu.rep.ID
				g.Transition(func() { g.Next(StatusSave) }, ScreeCapture(g))
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("设置", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("设置按钮被点击")
				g.Next(StatusSetting)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("历史", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("历史按钮被点击")
				//g.Next(StatusMenu)
				//计划做个弹窗?
				gu.OpenHistory(g)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		// widget.NewButton(
		// 	widget.ButtonOpts.Image(LoadRransparentButtonImage()),
		// 	// specify the button's text, the font face, and the color
		// 	//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		// 	widget.ButtonOpts.Text("流程", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
		// 		log.Println("流程按钮被点击")
		// 		g.Next(StatusTree)
		// 	}),
		// 	widget.ButtonOpts.DisableDefaultKeys(),
		// ),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("主菜单", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("主菜单按钮被点击")
				g.Next(StatusMenu)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("X", g.FontFace[0].Face(20), LoadBlueButtonTextColor()),
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
				log.Println("x按钮被点击")
				gu.hide = true
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
	}
	// 创建网格布局容器（单行）
	menu := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(7),    // 单列布局
			widget.GridLayoutOpts.Spacing(0, 0), // 按钮间距 5px
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 30), // 设置最小宽度 100px
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)
	for _, v := range buttons {
		menu.AddChild(v)
	}
	// 布局设置

	rootContainer.AddChild(menu)

	gu.ui = &ebitenui.UI{
		Container: rootContainer,
	}
	//time
	g.startTime = time.Now()

	gu.LoadString = func(s string) { gu.nowWord = s }
	gu.LoadAvatar = func(s string) {
		if s != "" {
			gu.avatarImage, _ = NewImageFromReader(150, 0, file.ReadMaterial(s))
		} else {
			gu.avatarImage = nil
		}
	}
	gu.LoadBackground = func(s string) {
		if s != "" {
			//gu.Background, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial(s)))
			gu.backgroundImage, _ = NewImageFromReader(1600, 0, file.ReadMaterial(s))
		} else {
			gu.backgroundImage = ebiten.NewImage(1, 1)
		}
	}
}
func (gu *GameUI) Clear(g *Game) {
	gu.selectionwindow = nil
	gu.history = []string{}
	if gu.MusicPlayer != nil {
		gu.MusicPlayer.Close()
	}
	if gu.VoicePlayer != nil {
		gu.VoicePlayer.Close()
	}
	gu.AudioContext.IsReady()
}
func (gu *GameUI) Update(g *Game) {
	if gu.doingchange {
		//change
		gu.rep = LoadRepertoire(gu.nextid)
		gu.history = append(gu.history, gu.rep.ID)
		switch gu.rep.Types {
		case "A": //常规类型
			gu.newString = StreamStringWithString(fmt.Sprintf("【%s】\n", gu.rep.Role), gu.rep.Text)
			gu.isAllString = false
			gu.LoadCreation(gu.rep.Creation) //加载立绘
			gu.LoadAvatar(gu.rep.Avatar)
			gu.LoadBackground(gu.rep.Background)
			gu.PlayMusic(gu.rep.Music)
			gu.PlayVideo(gu.rep.Video)
			gu.doingchange = false
			gu.nextid = gu.rep.Next
		case "B": //CG
		case "C":
			//选择界面
			gu.newString = StreamStringWithString(fmt.Sprintf("【%s】\n", gu.rep.Role), gu.rep.Text)
			gu.isAllString = false
			gu.LoadCreation(gu.rep.Creation) //加载立绘
			gu.LoadAvatar(gu.rep.Avatar)
			gu.LoadBackground(gu.rep.Background)
			gu.PlayMusic(gu.rep.Music)
			gu.PlayVideo(gu.rep.Video)
			gu.doingchange = false
			gu.OpenSelectWindows(g)
		case "D": //个人线判断
			id := gu.rep.Map[strconv.Itoa(g.Player.Token)]
			if id == "" {
				log.Print("存档损坏")
				os.Exit(1)
			}
			gu.nextid = id
			return
		}
		//过渡动画
		switch gu.rep.Transition {
		case "白色渐变":
			gu.DoTransition = AnimationTransparent(g)
		default:
			gu.DoTransition = nil
		}

	}
	for _, v := range gu.DoAction {
		v()
	}
	gu.counter++
	if gu.counter == 5 {
		gu.counter = 0
		var s string
		s, gu.isAllString = gu.newString(gu.isAllString)
		gu.LoadString(s)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		//避免调用按钮时误触
		mx, my := ebiten.CursorPosition()
		if mx > 0 && my < Height-30 && gu.historywindow == nil && gu.selectionwindow == nil {
			//判断是否用于隐藏ui操作的解除
			if gu.hide {
				gu.hide = false
			} else {
				//判断是否用于跳过逐字输出
				if !gu.isAllString {
					gu.isAllString = true
				} else {
					//之后再用于切换剧目
					if gu.rep.Transition == "" {
						gu.doingchange = true
					} else {
						g.startTime = time.Now() //修正过渡动画的时间
						gu.doingTransition = true
					}
				}
			}
		}
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
		screen.DrawImage(gu.backgroundImage, DrawBackgroundOption(gu.backgroundImage))
	}
	//立绘
	gu.drawCreation(screen)
	if gu.hide {
		return
	}
	//word
	gu.drawString(g, screen)
	//左下角图标
	gu.drawAvatar(screen)
	//ui
	gu.ui.Draw(screen)
	//过渡动画
	gu.DrawTransition(screen)
}

// 初始化角色立绘
func (gu *GameUI) LoadCreation(c [3]Creation) {
	for i, v := range c {
		if v.Role != "" {
			var err error
			gu.creation[i].Image, err = NewImageFromReader(400, 0, file.ReadMaterial(v.Role))
			if err != nil {
				log.Println("Error:", err)
			}
			gu.creation[i].O = &ebiten.DrawImageOptions{}
			gu.creation[i].O.GeoM.Translate(1600/6*float64(i+1), 900/3*1)
			//action
			switch v.Action {
			case "":
				gu.DoAction[i] = nilActionFunc()
			case "jump": //内置的jump
				current := 0 //top相对位置
				stage := 0   // 0: 0→15, 1:15→-15, 2:-15→0, 3:保持0
				gu.DoAction[i] = func() {
					switch stage {
					case 0:
						current++
						if current > 15 {
							stage = 1
						}
						gu.creation[i].O.GeoM.Translate(0, 1)
					case 1:
						current--
						if current < -15 {
							stage = 2
						}
						gu.creation[i].O.GeoM.Translate(0, -1)
					case 2:
						current++
						if current == 0 {
							stage = 3
							current = 0
						}
						gu.creation[i].O.GeoM.Translate(0, 1)
					default:
						gu.DoAction[i] = nilActionFunc() //结束action
					}
				}
			}
		} else {
			gu.creation[i].Image = nil
		}
	}
}

// 绘制立绘和action
func (gu *GameUI) drawCreation(screen *ebiten.Image) {
	for _, v := range gu.creation {
		if v.Image != nil {
			screen.DrawImage(v.Image, v.O)
		}
	}
}

// 绘制文字和文字背景
func (gu *GameUI) drawString(g *Game, screen *ebiten.Image) {
	//文字背景
	bo := &ebiten.DrawImageOptions{}
	bo.GeoM.Translate(280, 900-190)
	bo.ColorScale.ScaleAlpha(0.4) //  // 调整透明度
	back := ebiten.NewImage(1320, 150)
	back.Fill(color.RGBA{255, 182, 193, 255})
	screen.DrawImage(back, bo)
	//文字
	op := &text.DrawOptions{}
	op.GeoM.Translate(300, 900/5*4)
	op.ColorScale.ScaleWithColor(color.Black)
	op.LayoutOptions.LineSpacing = 30
	text.Draw(screen, gu.nowWord, g.FontFace[0].Face(25), op)
}

// 绘制左下角图标
func (gu *GameUI) drawAvatar(screen *ebiten.Image) {
	if gu.avatarImage == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 900-160)
	screen.DrawImage(gu.avatarImage, op)
}

// 绘制过渡动画在两个剧目之间
func (gu *GameUI) DrawTransition(screen *ebiten.Image) {
	if gu.doingTransition && gu.DoTransition != nil {
		if ok := gu.DoTransition(screen); ok {
			//动画播放完毕
			gu.doingchange = true
			gu.doingTransition = false
		}
	}
}

// 播放/切换/停止背景音乐
func (gu *GameUI) PlayMusic(name string) {
	log.Println("playMusic:", name)
	if name == "" {
		gu.lastMusic = ""
		if gu.MusicPlayer != nil {
			gu.MusicPlayer.Close()
		}
		return
	}
	if name != gu.lastMusic {
		gu.lastMusic = name
		stream, err := mp3.DecodeF32(bytes.NewReader(file.ReadMaterial(name)))
		if err != nil {
			log.Println("Error: mp3.DecodeF32 ", err)
		}
		gu.MusicPlayer, _ = gu.AudioContext.NewPlayerF32(stream)
		//gu.MusicPlayer = gu.AudioContext.NewPlayerFromBytes(file.ReadMaterial(name))
		gu.MusicPlayer.Play()
	}
	//结束后进入循环
	if !gu.MusicPlayer.IsPlaying() {
		gu.MusicPlayer.Rewind()
		gu.MusicPlayer.Play()
	}
}

// 播放语音
func (gu *GameUI) PlayVideo(name string) {
	if gu.VoicePlayer != nil {
		gu.VoicePlayer.Close()
	}
	if name == "" {
		return
	}
	stream, _ := mp3.DecodeF32(bytes.NewReader(file.ReadMaterial(name)))
	gu.MusicPlayer, _ = gu.AudioContext.NewPlayerF32(stream)
	gu.MusicPlayer.Play()
}

// 创建窗口
func (gu *GameUI) createSelectWindow(g *Game) *widget.Window {
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

	for _, v := range gu.rep.Select {
		windowContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			widget.ButtonOpts.Text(v.Text, g.FontFace[0].Face(35), LoadBlueButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   20,
				Right:  20,
				Top:    5,
				Bottom: 5,
			}),
			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("Select选择了:", v.Text)
				g.Player.Token += v.Token
				gu.doingchange = true
				gu.nextid = v.Next
				gu.selectionwindow.Close()
				gu.selectionwindow = nil
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		))
	}

	return widget.NewWindow(
		// Set the main contents of the window
		widget.WindowOpts.Contents(windowContainer),
		// Set the window above everything else and block input elsewhere
		//widget.WindowOpts.Modal(),
	)
}
func (gu *GameUI) OpenSelectWindows(g *Game) {
	if gu.selectionwindow == nil {
		gu.selectionwindow = gu.createSelectWindow(g)
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

func (gu *GameUI) OpenHistory(g *Game) {
	if gu.historywindow == nil {
		gu.historywindow = gu.createHistoryWindow(g)
	}
	if !gu.ui.IsWindowOpen(gu.historywindow) {
		log.Println("打开选择窗口")
		// Get the preferred size of the content
		x, y := gu.historywindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		gu.historywindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		gu.ui.AddWindow(gu.historywindow)
	} else {
		gu.historywindow.Close()
		gu.historywindow = nil
	}
}
func (gu *GameUI) createHistoryWindow(g *Game) *widget.Window {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0, 0, 0, 0})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  Width - 200,
					MaxHeight: Height - 200,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(Width-200, Height-200),
			),
		),
		//Set gap between scrollbar and text
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		//Tell the textarea to display bbcodes
		widget.TextAreaOpts.ProcessBBCode(true),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.Black),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(g.FontFace[0].Face(40)),
		//Set the initial text for the textarea
		//It will automatically line wrap and process newlines characters
		//If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text(""),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(20)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 200}),
				Mask: eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 200}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 100}),
					Hover: eimage.NewNineSliceColor(color.NRGBA{225, 255, 255, 100}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
	)
	//Add text to the end of the textarea
	for _, v := range gu.history {
		t := LoadRepertoire(v)
		if t.Types == "A" || t.Types == "C" {
			textarea.AppendText(fmt.Sprint("\n", t.Role, ": ", t.Text, "\n"))
		}
	}
	rootContainer.AddChild(textarea)

	return widget.NewWindow(

		// Set the main contents of the window
		widget.WindowOpts.Contents(rootContainer),
		// Set the window above everything else and block input elsewhere
		//widget.WindowOpts.Modal(),
	)
}

var nilFunc = func() {}

func nilActionFunc() func() {
	return nilFunc
}
