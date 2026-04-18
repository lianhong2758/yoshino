package yoshino

import (
	"encoding/json"
	"log"
	"os"

	"github.com/lianhong2758/yoshino/file"
)

// 剧本解析
type Script map[string]*Repertoire

type Repertoire struct {
	ID             string
	Types          string      //剧目类型 A 对话,B CG, C 选择,D 个人线判断
	Role           string      //对话用角色名
	Text           string      //对话内容
	Avatar         string      //左下角头像,非必须
	Creation       [3]Creation //立绘,分别对应左中右,非必须
	Background     string      //背景
	BackgroundType string      `json:"backgroundtype"` //image,mpg
	Music          string      //背景音乐
	Voice          string      //角色语音
	Select         []selects
	Next           string
	//option
	CaseMap caseMap //用于types == d 时,用于选择个人线或者后续线路的判断.
}

// 选择分支
type selects struct {
	Text string
	Next string
	Key  SelectMap //作为赋值加入
}

type Creation struct {
	Role   string
	Action string //动画效果
	A      int    //透明度?
}
type caseMap []struct {
	Key  SelectMap
	Next string
}

var script Script

func ScriptInit() {
	log.Println("加载剧本...")
	err := json.Unmarshal(file.ReadMaterial("script.json"), &script)
	if err != nil {
		log.Println("加载剧本错误,err:", err)
		os.Exit(1)
	}
}

func LoadRepertoire(id string) *Repertoire {
	return script[id]
}
