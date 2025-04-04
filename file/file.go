package file

import (
	"embed"
	"io/fs"
)

//go:embed  MaokenZhuyuanTi.ttf
var MaoKenTTF []byte

//go:embed STLITI.TTF
var StlitiTTF []byte

//go:embed material
var Material embed.FS

func ReadMaterial(name string) (data []byte) {
	data, _ = Material.ReadFile("material/" + name)
	return
}
func OpenMaterial(name string) (fs fs.File) {
	fs, _ = Material.Open("material/" + name)
	return
}
