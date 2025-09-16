package yoshino

import (
	"errors"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lianhong2758/yoshino/file"
)

var (
	NilImage = ebiten.NewImage(1, 1)

	StdImagePool = NewImagePool()
)

// ImagePool 图池数据结构
type ImagePool struct {
	images map[string]*ebiten.Image
}

// NewImagePool 创建一个新的图池实例
func NewImagePool() *ImagePool {
	return &ImagePool{
		images: make(map[string]*ebiten.Image),
	}
}

func (p *ImagePool) LoadImage(key, fileName string) error {
	f := file.OpenMaterial(fileName)
	defer f.Close()
	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		return err
	}
	p.images[key] = img
	return nil
}

// key, fileName的格式输入,必须为对子
func (p *ImagePool) LoadImageArray(arg ...string) error {
	if len(arg)%2 != 0 {
		return errors.New("bad number with arg")
	}
	for key, fileName := range rangeimagekv(arg) {
		f := file.OpenMaterial(fileName)
		defer f.Close()
		img, _, err := ebitenutil.NewImageFromReader(f)
		if err != nil {
			return err
		}
		p.images[key] = img
	}
	return nil
}

// 存入图片
func (p *ImagePool) PostImage(key string, img *ebiten.Image) {
	p.images[key] = img
}

// 返回图片,如果不存在则返回NilImage
func (p *ImagePool) GetImage(key string) *ebiten.Image {
	if image, ok := p.images[key]; ok {
		return image
	}
	return NilImage
}

// RemoveImage 从图池中移除指定键的图像
func (p *ImagePool) RemoveImage(key string) {
	delete(p.images, key)
}

// Clear 清空整个图池
func (p *ImagePool) Clear() {
	p.images = make(map[string]*ebiten.Image)
}

// Size 返回图池中图像的数量
func (p *ImagePool) Size() int {
	return len(p.images)
}

func rangeimagekv(arg []string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for i := 0; i < len(arg); i += 2 {
			if !yield(arg[i], arg[i+1]) {
				return
			}
		}
	}
}
