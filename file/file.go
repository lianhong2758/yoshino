package file

import (
	"embed"
)

//go:embed  MaokenZhuyuanTi.ttf
var MaoKenTTF []byte

//go:embed material
var Material embed.FS

func ReadMaterial(name string) (data []byte) {
	data, _ = Material.ReadFile("material/" + name)
	return
}
