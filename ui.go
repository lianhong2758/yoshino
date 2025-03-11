package yoshino

import "github.com/hajimehoshi/ebiten/v2"

// UI
type UI interface {
	Update(g *Game)
	Draw(g *Game, screen *ebiten.Image)
	Init(g *Game)
	Clear(g *Game)
}

// 一个界面为一个状态,框架示例是这么写得
type GameStatus int

// 按照一般开启顺序排列
const (
	StatusTitle   GameStatus = iota //开屏动画等
	StatusMenu                      //主菜单
	StatusGame                      //游戏内容
	StatusCG                        //可能存在的CG场景,用于播放视频
	StatusSetting                   //设置界面
	StatusSave                      //保存存档界面
	StatusLoad                      //加载存档界面
	StatusTree                      //流程树界面
)
