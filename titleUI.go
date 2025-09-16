package yoshino

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type TitleUI struct {
	BaseScene
	imageKey string
}

func (s *TitleUI) Load(st State, sm stagehand.SceneController[State]) {
	s.State = State{Page: PageTitle, Count: 0, Count2: 0}
	s.sm = sm.(*stagehand.SceneManager[State])
	LoadFontsFromFs("MaokenZhuyuanTi.ttf", "STLITI.TTF")
	StdImagePool.LoadImageArray("title", "title.png", "logo", "logo.png")
	s.imageKey = "title"
}

func (s *TitleUI) Unload() State {
	StdImagePool.Clear()
	return s.State
}

func (s *TitleUI) Update() error {
	s.Count++
	//2s
	if s.Count > 120 {
		s.Count = 0
		s.Count2++
		s.imageKey = "logo"
	}
	if s.Count2 == 2 {
		s.sm.SwitchTo(&MenuUI{})
	}
	return nil
}

func (s *TitleUI) Draw(screen *ebiten.Image) {
	alpha := min(float32(s.Count)/60.0, 1)
	var op *ebiten.DrawImageOptions
	if s.imageKey == "title" {
		op = DrawBackgroundOption(StdImagePool.GetImage(s.imageKey))
		op.ColorScale.ScaleAlpha(float32(alpha)) //  // 调整透明度
	} else {
		screen.Fill(color.White)
		op = &ebiten.DrawImageOptions{}
		scaleFactor := max(float64(Width/2)/float64(StdImagePool.GetImage(s.imageKey).Bounds().Dx()), float64(Height/2)/float64(StdImagePool.GetImage(s.imageKey).Bounds().Dy()))
		op.GeoM.Scale(scaleFactor, scaleFactor)
		op.GeoM.Translate(
			(float64(Width)-(float64(StdImagePool.GetImage(s.imageKey).Bounds().Dx())*scaleFactor))/2,
			(float64(Height)-(float64(StdImagePool.GetImage(s.imageKey).Bounds().Dy())*scaleFactor))/2,
		)
		op.ColorScale.ScaleAlpha(float32(alpha)) //  // 调整透明度
	}
	screen.DrawImage(StdImagePool.GetImage(s.imageKey), op)
}
