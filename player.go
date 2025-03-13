package yoshino

import (
	"encoding/json"
	"image"
)

// //用户存档信息
// type Player interface {
// 	Load(path string)
// 	Save(path string)
// 	Bytes() []byte  //序列化后的存档
// 	String() string //作为日志的打印
// }

type Player struct {
	Token         int
	ID            string
	CGUnlock      string
	ScreenData    []byte
	screenContent image.Image
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
