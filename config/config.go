package config

type App struct {
	Name   string `json:"name" yaml:"name"`
	Listen string `json:"listen" yaml:"listen"`
}

type Logger struct {
	Path             string `json:"path" yaml:"path"`
	RequestTableName string `json:"requestTableName" yaml:"requestTableName"`
}

type Mysql struct {
	Conn string `json:"conn" yaml:"conn"`
}

type Mongo struct {
	Conn string `json:"conn" yaml:"conn"`
	Db   string `json:"db" yaml:"db"`
}

type runtimeParam struct {
	RootDir string `json:"-" yaml:"-"` // 此软件运行后的工作目录
}

type Configs struct {
	App          App          `json:"app" yaml:"app"`
	Logger       Logger       `json:"logger" yaml:"logger"`
	Mysql        Mysql        `json:"mysql" yaml:"mysql"`
	Mongo        Mongo        `json:"mongo" yaml:"mongo"`
	Spec         string       `json:"spec" yaml:"spec"`
	RuntimeParam runtimeParam `json:"-" yaml:"-"`
}

// Cfg 全局的Config配置，解析dns.yaml的结果
var Cfg *Configs

var configFileName = "example.yml"
