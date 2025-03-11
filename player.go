package yoshino

//用户存档信息
type Player interface {
	Load(path string)
	Save(path string)
	Bytes() []byte  //序列化后的存档
	String() string //作为日志的打印
}
