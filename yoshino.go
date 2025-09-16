package yoshino

import (
	"github.com/joelschutz/stagehand"
)

var (
	Width, Height = 1600, 900
)

const (
	PageTitle   int = iota //开屏动画等
	PageMenu               //主菜单
	PageGame               //游戏内容
	PageCG                 //可能存在的CG场景,用于播放视频
	PageSetting            //设置界面
	PageSave               //保存存档界面
	PageLoad               //加载存档界面
	PageTree               //流程树界面
)

type State struct {
	Page   int //页面管理用的状态
	Count  int //如果有需要可以内部使用的计数器
	Count2 int
}

type BaseScene struct {
	State
	sm *stagehand.SceneManager[State]
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	return Width, Height
}
