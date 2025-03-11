package yoshino

import (
	"bytes"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LoadFont(ttf []byte, size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func (g *Game) Next(state GameStatus) {
	g.GameUI[g.Status].Clear(g)
	g.Status = state
	g.GameUI[g.Status].Init(g)
}

func LoadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func LoadButtonTextColor() *widget.ButtonTextColor {
	return &widget.ButtonTextColor{
		Idle:    color.NRGBA{0xdf, 0xf4, 0xff, 0xff}, //闲置
		Hover:   color.NRGBA{0, 255, 128, 255},       //徘徊
		Pressed: color.NRGBA{255, 0, 0, 255},         //按下
	}
}

func DrawImageCentreOption(img *ebiten.Image) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	scaleFactor := max(float64(Width)/float64(img.Bounds().Dx()), float64(Height)/float64(img.Bounds().Dy()))
	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(
		(float64(Width)-(float64(img.Bounds().Dx())*scaleFactor))/2,
		(float64(Height)-(float64(img.Bounds().Dy())*scaleFactor))/2,
	)
	return op
}

// 创建会逐步增加输出的函数
func StreamString(b string) func() string {
	runes := []rune(b) // 将字符串转为rune切片
	max := len(runes)  // 实际字符数量
	current := 0       // 记录当前字符位置

	return func() (bt string) {
		if current >= max {
			return (b) // 超过字符数时直接输出原始字符串
		}
		// 输出从开始到当前字符位置的切片
		bt = string(runes[:current])
		current++ // 移动到下一个字符位置
		return
	}
}
