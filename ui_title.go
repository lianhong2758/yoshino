package yoshino

import (
	"bytes"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino/file"
)

type TitleUI struct {
	TitleFile []*ebiten.Image //开屏图片列表
}

func (t *TitleUI) Init(g *Game) {
	//img, _, err := ebitenutil.NewImageFromFile("../data/title.png")
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial("title.png")))
	if err != nil {
		log.Println(err)
		return
	}
	t.TitleFile = append(t.TitleFile, img)
	//	img2, _, err := ebitenutil.NewImageFromFile("../data/logo.png")
	img2, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(file.ReadMaterial("logo.png")))
	if err != nil {
		log.Println(err)
		return
	}
	t.TitleFile = append(t.TitleFile, img2)
	//加载字体
	f, _ := LoadFont(file.MaoKenTTF, 20)
	g.FontFace = append(g.FontFace, f)
	//计算开屏时间需要
	g.startTime = time.Now()
}
func (t *TitleUI) Clear(g *Game) { clear(t.TitleFile) }
func (*TitleUI) Update(g *Game)  {}

func (ti *TitleUI) Draw(g *Game, screen *ebiten.Image) {
	// 计算当前时间与动画开始时间的差值
	elapsed := time.Since(g.startTime)
	if elapsed.Seconds() >= 4 {
		g.Next(StatusMenu)
		return
	}
	i := 0
	t := elapsed.Seconds()
	if t >= 2 {
		t -= 2
		i = 1
	}
	//alpha := 1 - math.Abs(t-1) //一共4s,取1,3s为最亮点绘制两张图
	if t > 1 {
		t = 1
	}
	alpha := t
	var op *ebiten.DrawImageOptions
	if i == 0 {
		op = DrawImageCentreOption(ti.TitleFile[i])
		op.ColorScale.ScaleAlpha(float32(alpha)) //  // 调整透明度
	} else {
		screen.Fill(color.White)
		op = &ebiten.DrawImageOptions{}
		scaleFactor := max(float64(Width/2)/float64(ti.TitleFile[i].Bounds().Dx()), float64(Height/2)/float64(ti.TitleFile[i].Bounds().Dy()))
		op.GeoM.Scale(scaleFactor, scaleFactor)
		op.GeoM.Translate(
			(float64(Width)-(float64(ti.TitleFile[i].Bounds().Dx())*scaleFactor))/2,
			(float64(Height)-(float64(ti.TitleFile[i].Bounds().Dy())*scaleFactor))/2,
		)
		op.ColorScale.ScaleAlpha(float32(alpha)) //  // 调整透明度
	}
	screen.DrawImage(ti.TitleFile[i], op)
}
