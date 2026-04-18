package yoshino

import (
	"encoding/json"
	"maps"

	"github.com/hajimehoshi/ebiten/v2"
)

// //用户存档信息
//
//	type Player interface {
//		Load(path string)
//		Save(path string)
//		Bytes() []byte  //序列化后的存档
//		String() string //作为日志的打印
//	}
type SelectMap map[string]int //type:int
type Player struct {
	ID                string
	ScreenData        []byte    //存档的保存截图
	Data              SelectMap //type:val
	screenEbitenImage *ebiten.Image
}

func (p *Player) Load(path string) {

}

func (p *Player) Save(path string) {

}
func (p *Player) Bytes() []byte {
	d, _ := json.Marshal(p)
	return d
}
func (p *Player) String() string {
	return string(p.Bytes())
}

func (p *Player) GetSelectNext(c caseMap) string {
	for _, val := range c {
		if maps.Equal(p.Data, val.Key) {
			return val.Next
		}
	}
	return ""
}
