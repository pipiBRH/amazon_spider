package schema

var BlockCategory = map[string]bool{
	"Prime Video":    true,
	"Amazon Music":   true,
	"Android アプリストア": true,
	"Echo & Alexa":   true,
	"Fireタブレット":      true,
	"Fire TV":        true,
	"Kindle 本＆電子書籍リーダー":   true,
	"クレジットカード＆Amazonポイント": true,
	"DVD・ミュージック・ゲーム":      true,
}

var Config Conf

type Conf struct {
	SSDB   SSDBConf
	Spider SpiderConf
}
type SSDBConf struct {
	IP   string
	Port int
}

type SpiderConf struct {
	Znum           int
	CategoryLevel  int
	EnableCategory bool
	EnableProduct  bool
}
