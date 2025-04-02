package yoshino

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"unsafe"

	ebimg "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/nfnt/resize"
)

type FontSource text.GoTextFaceSource

func LoadFont(ttf []byte) (f *FontSource, err error) {
	gf := new(text.GoTextFaceSource)
	gf, err = text.NewGoTextFaceSource(bytes.NewReader(ttf))
	f = (*FontSource)(gf)
	return
}

func (f *FontSource) Face(size float64) *text.GoTextFace {
	return &text.GoTextFace{
		Source: (*text.GoTextFaceSource)(f),
		Size:   size,
	}
}

// 透明按钮
func LoadRransparentButtonImage() *widget.ButtonImage {
	transparentImage := ebimg.NewNineSliceColor(color.RGBA{0, 0, 0, 0})
	return &widget.ButtonImage{
		Idle:    transparentImage,
		Hover:   transparentImage,
		Pressed: transparentImage,
	}
}

// left 50
func LoadNoDataButtonImage(g *Game) *widget.GraphicImage {
	img := ebiten.NewImage(300, 200)
	img.Fill(color.RGBA{255, 235, 205, 255})
	op := &text.DrawOptions{}
	op.GeoM.Translate((300-200)/2, (200-40)/2)
	op.ColorScale.ScaleWithColor(color.RGBA{135, 206, 250, 255})
	text.Draw(img, "No Data", g.FontFace[0].Face(40), op)
	return &widget.GraphicImage{
		Idle:    img,
		Hover:   img,
		Pressed: img,
	}
}

func LoadButtonImageByImage(g *Game, p Player) *widget.GraphicImage {
	img := ebiten.NewImage(300, 200)
	img.Fill(color.RGBA{255, 235, 205, 255})
	// op := &text.DrawOptions{}
	// op.GeoM.Translate((300-200)/2, (200-40)/2)
	// op.ColorScale.ScaleWithColor(color.RGBA{135, 206, 250, 255})
	// text.Draw(img, "No Data", g.FontFace[0].Face(40), op)
	//260宽
	var op = &ebiten.DrawImageOptions{}
	scaleFactor := 260 / float64(p.screenEbitenImage.Bounds().Dx())
	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(
		(300-(float64(p.screenEbitenImage.Bounds().Dx())*scaleFactor))/2,
		(300-(float64(p.screenEbitenImage.Bounds().Dx())*scaleFactor))/2,
	)
	img.DrawImage(p.screenEbitenImage, op)
	return &widget.GraphicImage{
		Idle:    img,
		Hover:   img,
		Pressed: img,
	}
}

// 选择窗口的按钮背景,米黄色?
func LoadConfirmButtonImage() *widget.ButtonImage {
	idle := ebimg.NewNineSliceColor(color.NRGBA{255, 165, 0, 255})

	hover := ebimg.NewNineSliceColor(color.NRGBA{255, 140, 0, 255})

	pressed := ebimg.NewNineSliceColor(color.NRGBA{255, 127, 80, 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func LoadBlackButtonTextColor() *widget.ButtonTextColor {
	return &widget.ButtonTextColor{
		Idle:    color.Black,                //闲置
		Hover:   color.RGBA{0, 0, 255, 255}, //徘徊
		Pressed: color.RGBA{0, 255, 0, 255}, //按下
	}
}

func LoadBlueButtonTextColor() *widget.ButtonTextColor {
	return &widget.ButtonTextColor{
		Idle:    color.RGBA{135, 206, 250, 255}, //闲置
		Hover:   color.RGBA{100, 149, 237, 255}, //徘徊
		Pressed: color.RGBA{255, 215, 0, 255},   //按下
	}
}

// 绘制背景
func DrawBackgroundOption(img *ebiten.Image) *ebiten.DrawImageOptions {
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

// 创建一个函数,以star为开头,逐步增加输出b的函数
func StreamStringWithString(star, b string) func(bool) (string, bool) {
	runes := []rune(b) // 将字符串转为rune切片
	max := len(runes)  // 实际字符数量
	current := 0       // 记录当前字符位置

	return func(isAll bool) (bt string, ok bool) {
		if current >= max || isAll {
			return fmt.Sprint(star, b), true // 超过字符数时直接输出原始字符串
		}
		// 输出从开始到当前字符位置的切片
		bt = string(runes[:current])
		current++ // 移动到下一个字符位置
		return fmt.Sprint(star, bt), false
	}
}

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// 返回指定大小的图片
func NewImageFromReader(width uint, height uint, imgdata []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		return nil, err
	}
	img = resize.Resize(width, height, img, resize.Bilinear) //改比例
	img2 := ebiten.NewImageFromImage(img)
	return img2, err
}
