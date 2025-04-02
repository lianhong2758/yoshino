package yoshino

import (
	"bytes"
	"image/color"
	"image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 渐白过渡动画
func AnimationTransparent(g *Game) func(screen *ebiten.Image) bool {
	g.startTime = time.Now()
	return func(screen *ebiten.Image) bool {
		elapsed := time.Since(g.startTime)
		if elapsed.Seconds() >= 1 {
			return true
		}
		overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		overlay.Fill(color.White)
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(elapsed.Seconds() / 1.0)) //  // 调整透明度
		screen.DrawImage(overlay, op)
		return false
	}
}

// 截屏并保存至g.Player.ScreenData
func ScreeCapture(g *Game) func(screen *ebiten.Image) bool {
	return func(screen *ebiten.Image) bool {
		var picbuff bytes.Buffer
		_ = png.Encode(&picbuff, screen.SubImage(screen.Bounds()))
		g.Player.ScreenData = picbuff.Bytes()
		return true
	}
}