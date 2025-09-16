package yoshino

import (
	"bytes"
	"encoding/gob"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/lianhong2758/yoshino/file"
	"github.com/quasilyte/gdata"
)

// 存档系统初始化
var StdFileSystem = NewFileSystem()

type FileSystem struct{ *gdata.Manager }

func NewFileSystem() *FileSystem {
	m, err := gdata.Open(gdata.Config{
		AppName: "yoshino",
	})
	if err != nil {
		panic(err)
	}
	return &FileSystem{m}
}

// 返回12个项目,用于绘制存档界面
func (f *FileSystem) LoadPlayers() ([]Player, error) {
	if !f.ItemExists("players.gob") {
		return make([]Player, 12), nil
	}
	data, err := f.LoadItem("players.gob")
	if err != nil {
		log.Println("load err: ", err)
		return nil, err
	}
	players := make([]Player, 0)
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&players)
	if err != nil {
		log.Println("load err: ", err)
		return nil, err
	}
	for k := range len(players) {
		// players[k].screenContent
		players[k].screenEbitenImage, _, _ = ebitenutil.NewImageFromReader(bytes.NewReader(players[k].ScreenData))
	}
	//如果不够12就补齐
	for range 12 - len(players) {
		players = append(players, Player{})
	}
	return players, nil
}

func (f *FileSystem) SavePlayers(p []Player) error {
	for range 12 - len(p) {
		p = append(p, Player{})
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(p)
	if err != nil {
		log.Println("save err: ", err)
		return err
	}
	err = f.SaveItem("players.gob", buff.Bytes())
	if err != nil {
		log.Println("loaderr: ", err)
		return err
	}
	return nil
}

// 字体初始化
var StdFonts []*FontFaceSource

type FontFaceSource text.GoTextFaceSource

func LoadFont(fns fs.File) (f *FontFaceSource, err error) {
	defer fns.Close()
	gf := new(text.GoTextFaceSource)
	gf, err = text.NewGoTextFaceSource(fns)
	f = (*FontFaceSource)(gf)
	return
}

func LoadFontFromByte(ttf []byte) (f *FontFaceSource, err error) {
	gf := new(text.GoTextFaceSource)
	gf, err = text.NewGoTextFaceSource(bytes.NewReader(ttf))
	f = (*FontFaceSource)(gf)
	return
}

func (f *FontFaceSource) FacePointer(size float64) *text.Face {
	var gt text.Face = f.Face(size)
	return &gt
}
func (f *FontFaceSource) Face(size float64) *text.GoTextFace {
	return &text.GoTextFace{
		Source: (*text.GoTextFaceSource)(f),
		Size:   size,
	}
}

// 字体初始化入口
func LoadFontsFromFs(names ...string) error {
	for _, v := range names {
		f, err := LoadFont(file.OpenMaterial(v))
		if err != nil {
			return err
		}
		StdFonts = append(StdFonts, f)
	}
	return nil
}
