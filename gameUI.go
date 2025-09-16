package yoshino

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"strconv"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/lianhong2758/yoshino/file"
	"github.com/tinne26/mpegg"
)

type GameUI struct {
	BaseScene

	Player Player //加载或导出时使用的存档

	ui              *ebitenui.UI
	selectionwindow *widget.Window
	historywindow   *widget.Window
	newString       func(bool) (string, bool) //输出的字符串,是否输出全句
	isAllString     bool                      //记录是否输出全句
	counter         int                       //计数器,逐字打印需要
	nowWord         string                    //当前显示的对话内容
	doingchange     bool                      //是否切换下一剧目
	nextid          string                    //下一个id
	rep             *Repertoire               //当前剧本
	history         []string                  //存放id?还可以扩展语音播放等
	//doingTransition bool                      // 用于主UI的过渡动画执行后改变needchange
	hide      bool   //隐藏ui
	lastMusic string //上一个背景音乐,用于判断是否切换背景音乐

	VideoPlayer  *mpegg.Player //视频播放器
	AudioContext *audio.Context
	MusicPlayer  *audio.Player // 背景音乐播放器
	VoicePlayer  *audio.Player // 语音播放器
}

/*
ui

	/背景
	/立绘
		角色1	角色2	角色3


	/文本区	 【角色名】
	头像op   ⌈对话内容⌋
						保存 读取 设置 主菜单
*/
func (g *GameUI) Load(st State, sm stagehand.SceneController[State]) {
	g.State = State{Page: PageGame}
	g.sm = sm.(*stagehand.SceneManager[State])
	ScriptInit()
	g.nextid = g.Player.ID
	g.doingchange = true
	fb := ebiten.NewImage(1320, 150)
	fb.Fill(color.RGBA{255, 182, 193, 255})
	StdImagePool.PostImage("fontback", fb)
	if t := audio.CurrentContext(); t == nil {
		g.AudioContext = audio.NewContext(48000)
	} else {
		g.AudioContext = t
	}
	g.makeui()
}

func (g *GameUI) Unload() State {
	StdImagePool.Clear()
	g.selectionwindow = nil
	g.history = []string{}
	if g.MusicPlayer != nil {
		g.MusicPlayer.Close()
	}
	if g.VoicePlayer != nil {
		g.VoicePlayer.Close()
	}
	return g.State
}

func (g *GameUI) Update() error {
	if g.doingchange {
		//change
		g.rep = LoadRepertoire(g.nextid)
		g.history = append(g.history, g.rep.ID)
		log.Println(g.rep.Types)
		switch g.rep.Types {
		case "A": //常规类型
			g.newString = StreamStringWithString(fmt.Sprintf("【%s】\n", g.rep.Role), g.rep.Text)
			g.isAllString = false
			g.LoadCreation(g.rep.Creation) //加载立绘
			g.LoadAvatar(g.rep.Avatar)
			g.LoadBackground(g.rep.Background)
			g.PlayMusic(g.rep.Music)
			g.PlayVoice(g.rep.Voice)
			g.doingchange = false
			g.nextid = g.rep.Next
		case "B": //CG
			g.LoadBackground(g.rep.Background)
			//g.LoadCreation(g.rep.Creation) //让立绘为空白
			g.PlayMusic(g.rep.Music) //清空播放器或者特殊音效?
			g.PlayVoice(g.rep.Voice)
			g.hide = true
			g.doingchange = false
			g.nextid = g.rep.Next
		case "C":
			//选择界面
			g.newString = StreamStringWithString(fmt.Sprintf("【%s】\n", g.rep.Role), g.rep.Text)
			g.isAllString = false
			g.LoadCreation(g.rep.Creation) //加载立绘
			g.LoadAvatar(g.rep.Avatar)
			g.LoadBackground(g.rep.Background)
			g.PlayMusic(g.rep.Music)
			g.PlayVoice(g.rep.Voice)
			g.doingchange = false
			g.OpenSelectWindows()
		case "D": //个人线判断
			// id := g.rep.Map[strconv.Itoa(g.Player.Token)]
			// if id == "" {
			// 	log.Print("存档损坏")
			// 	os.Exit(1)
			// }
			// g.nextid = id
			// return
		}

	}
	g.counter++
	if g.counter == 5 {
		g.counter = 0
		var s string
		s, g.isAllString = g.newString(g.isAllString)
		g.LoadString(s)
	}
	//处理输入
	g.Input()
	//用于判断视频结束后跳转
	if g.VideoPlayer != nil && !g.VideoPlayer.IsPlaying() {
		g.doingchange = true
		g.CloseVideo()
		g.hide = false
	}
	g.ui.Update()
	return nil
}

func (g *GameUI) Draw(screen *ebiten.Image) {
	if g.rep == nil {
		return
	}
	//背景
	switch g.rep.BackgroundType {
	case "image":
		screen.DrawImage(StdImagePool.GetImage("bg"), DrawBackgroundOption(StdImagePool.GetImage("bg")))
	case "mpg":
		if g.VideoPlayer != nil {
			mpegg.Draw(screen, g.VideoPlayer.CurrentFrame())
		}
		return
	}
	//立绘
	g.drawCreation(screen)
	if !g.hide {
		//word
		g.drawString(screen)
		//左下角图标
		g.drawAvatar(screen)
		//ui
		g.ui.Draw(screen)
	}
}
func (g *GameUI) Input() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	//用于跳过视频
	if g.VideoPlayer != nil {
		if g.VideoPlayer.IsPlaying() {
			g.VideoPlayer.Pause()
		}
		return
	}
	//避免调用按钮时误触
	_, my := ebiten.CursorPosition()
	if my > Height-30 ||
		g.historywindow != nil || g.selectionwindow != nil {
		return
	}
	//判断是否用于隐藏ui操作的解除
	if g.hide {
		g.hide = false
		return
	}
	//判断是否用于跳过逐字输出
	if !g.isAllString {
		g.isAllString = true
		return
	}
	//之后再用于切换剧目
	g.doingchange = true
}

func (g *GameUI) LoadVideo(name string) {
	f := (file.OpenMaterial(name))
	mpeggPlayer, err := mpegg.NewPlayer(struct {
		io.Reader
		io.Seeker
	}{f, f.(io.Seeker)})
	if err != nil {
		log.Println("Error:", err)
		return
	}
	g.VideoPlayer = mpeggPlayer
	g.VideoPlayer.Play()
}

func (g *GameUI) CloseVideo() {
	if g.VideoPlayer != nil {
		if g.VideoPlayer.IsPlaying() {
			g.VideoPlayer.Pause()
		}
		g.VideoPlayer = nil
	}
}

// 初始化角色立绘
func (g *GameUI) LoadCreation(c [3]Creation) {
	for i, v := range c {
		if v.Role != "" {
			img, err := NewImageFromReader(400, 0, file.OpenMaterial(v.Role))
			if err != nil {
				log.Println("Error:", err)
			}
			StdImagePool.PostImage("creation"+strconv.Itoa(i), img)
		}
	}
}

// 更改显示的文字
func (g *GameUI) LoadString(s string) {
	g.nowWord = s
}

// 更改显示的头像
func (g *GameUI) LoadAvatar(s string) {
	if s != "" {
		img, _ := NewImageFromReader(150, 0, file.OpenMaterial(s))
		StdImagePool.PostImage("avatar", img)
	} else {
		StdImagePool.PostImage("avatar", NilImage)
	}

}

// 更改显示的背景
func (g *GameUI) LoadBackground(bgName string) {
	switch g.rep.BackgroundType {
	case "image":
		if bgName != "" {
			img, _ := NewImageFromReader(1600, 0, file.OpenMaterial(bgName))
			StdImagePool.PostImage("bg", img)
		} else {
			StdImagePool.PostImage("bg", NilImage)
		}
	case "mpg":
		g.LoadVideo(bgName)
	default:
		StdImagePool.PostImage("bg", NilImage)
	}

}

// 绘制立绘和action
func (g *GameUI) drawCreation(screen *ebiten.Image) {
	for i := range 3 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(1600/6*float64(i+1), 900/3*1)
		screen.DrawImage(StdImagePool.GetImage("creation"+strconv.Itoa(i)), op)
	}
}

// 绘制文字和文字背景
func (g *GameUI) drawString(screen *ebiten.Image) {
	//文字背景
	bo := &ebiten.DrawImageOptions{}
	bo.GeoM.Translate(280, 900-190)
	bo.ColorScale.ScaleAlpha(0.4) //  // 调整透明度
	screen.DrawImage(StdImagePool.GetImage("fontback"), bo)
	//文字
	op := &text.DrawOptions{}
	op.GeoM.Translate(300, 900/5*4)
	op.ColorScale.ScaleWithColor(color.Black)
	op.LayoutOptions.LineSpacing = 30
	text.Draw(screen, g.nowWord, StdFonts[0].Face(25), op)
}

// 绘制左下角图标
func (g *GameUI) drawAvatar(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 900-160)
	screen.DrawImage(StdImagePool.GetImage("avatar"), op)
}

// 播放/切换/停止背景音乐
func (g *GameUI) PlayMusic(name string) {
	log.Println("playMusic:", name)
	if name == "" {
		g.lastMusic = ""
		if g.MusicPlayer != nil {
			g.MusicPlayer.Close()
		}
		return
	}
	if name != g.lastMusic {
		g.lastMusic = name
		stream, err := mp3.DecodeF32(file.OpenMaterial(name))
		if err != nil {
			log.Println("Error: mp3.DecodeF32 ", err)
		}
		g.MusicPlayer, _ = g.AudioContext.NewPlayerF32(stream)
		g.MusicPlayer.Play()
	}
	//结束后进入循环
	if !g.MusicPlayer.IsPlaying() {
		g.MusicPlayer.Rewind()
		g.MusicPlayer.Play()
	}
}

// 播放语音
func (g *GameUI) PlayVoice(name string) {
	if g.VoicePlayer != nil {
		g.VoicePlayer.Close()
	}
	if name == "" {
		return
	}
	stream, _ := mp3.DecodeF32(file.OpenMaterial(name))
	g.VoicePlayer, _ = g.AudioContext.NewPlayerF32(stream)
	g.VoicePlayer.Play()
}

// 创建窗口
func (g *GameUI) createSelectWindow() *widget.Window {
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

	for _, v := range g.rep.Select {
		windowContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			widget.ButtonOpts.Text(v.Text, StdFonts[0].FacePointer(35), LoadBlueButtonTextColor()),
			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(&widget.Insets{
				Left:   20,
				Right:  20,
				Top:    5,
				Bottom: 5,
			}),
			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("Select选择了:", v.Text)
				//	g.Player.Token += v.Token
				g.doingchange = true
				g.nextid = v.Next
				g.selectionwindow.Close()
				g.selectionwindow = nil
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
func (g *GameUI) OpenSelectWindows() {
	if g.selectionwindow == nil {
		g.selectionwindow = g.createSelectWindow()
	}
	if !g.ui.IsWindowOpen(g.selectionwindow) {
		log.Println("打开选择窗口")
		// Get the preferred size of the content
		x, y := g.selectionwindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		g.selectionwindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		g.ui.AddWindow(g.selectionwindow)
	}
}

func (g *GameUI) OpenHistory() {
	if g.historywindow == nil {
		g.historywindow = g.createHistoryWindow()
	}
	if !g.ui.IsWindowOpen(g.historywindow) {
		log.Println("打开选择窗口")
		// Get the preferred size of the content
		x, y := g.historywindow.Contents.PreferredSize()

		// Create a rect with the preferred size of the content
		r := image.Rect(0, 0, x, y)
		// Use the Add method to move the window to the specified point
		//左上角点
		r = r.Add(image.Pt((Width-x)/2, (Height-y)/2))
		// Set the windows location to the rect.
		g.historywindow.SetLocation(r)
		// Add the window to the UI.
		// Note: If the window is already added, this will just move the window and not add a duplicate.
		g.ui.AddWindow(g.historywindow)
	} else {
		g.historywindow.Close()
		g.historywindow = nil
	}
}
func (g *GameUI) createHistoryWindow() *widget.Window {
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
		widget.TextAreaOpts.FontFace(StdFonts[0].FacePointer(40)),
		//Set the initial text for the textarea
		//It will automatically line wrap and process newlines characters
		//If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text(""),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(*widget.NewInsetsSimple(20)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerImage(&widget.ScrollContainerImage{
			Idle: eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 200}),
			Mask: eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 200}),
		}),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderParams(&widget.SliderParams{
			TrackImage: &widget.SliderTrackImage{
				Idle:  eimage.NewNineSliceColor(color.NRGBA{255, 255, 255, 100}),
				Hover: eimage.NewNineSliceColor(color.NRGBA{225, 255, 255, 100}),
			},
			HandleImage: &widget.ButtonImage{
				Idle:    eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		}),
	)
	//Add text to the end of the textarea
	for _, v := range g.history {
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

func (g *GameUI) makeui() {
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
			widget.ButtonOpts.Text("读取", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
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
				log.Println("读取按钮被点击")
				//g.Next(StatusLoad)
				g.sm.SwitchTo(&LoadUI{Player: g.Player})
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("保存", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
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
				log.Println("保存按钮被点击")
				//g.Next(StatusSave)
				// FistID, g.Player.ID = g.rep.ID, g.rep.ID
				// g.Transition(func() { g.Next(StatusSave) }, ScreeCapture(g))
				g.sm.SwitchTo(&SaveUI{Player: g.Player})
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("设置", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
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
				log.Println("设置按钮被点击")
				//	g.Next(StatusSetting)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("历史", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
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
				log.Println("历史按钮被点击")
				//计划做个弹窗?
				g.OpenHistory()
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			// specify the button's text, the font face, and the color
			//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
			widget.ButtonOpts.Text("流程", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
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
				log.Println("流程按钮被点击")
				//g.Next(StatusTree)
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			widget.ButtonOpts.Text("主菜单", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
			widget.ButtonOpts.TextProcessBBCode(true),
			widget.ButtonOpts.TextPadding(&widget.Insets{
				Left:   20,
				Right:  20,
				Top:    5,
				Bottom: 5,
			}),

			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("主菜单按钮被点击")
				//g.Next(StatusMenu)
				g.sm.SwitchTo(&MenuUI{})
			}),
			widget.ButtonOpts.DisableDefaultKeys(),
		),
		widget.NewButton(
			widget.ButtonOpts.Image(LoadRransparentButtonImage()),
			widget.ButtonOpts.Text("X", StdFonts[0].FacePointer(20), LoadBlueButtonTextColor()),
			widget.ButtonOpts.TextProcessBBCode(true),
			widget.ButtonOpts.TextPadding(&widget.Insets{
				Left:   20,
				Right:  20,
				Top:    5,
				Bottom: 5,
			}),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				log.Println("x按钮被点击")
				g.hide = !g.hide
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
	rootContainer.AddChild(menu)

	g.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}
