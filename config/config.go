package config

// Cfg 全局的Config配置，解析dns.yaml的结果
var Cfg *Configs

type App struct {
	Name   string `json:"name" yaml:"name"`
	Listen string `json:"listen" yaml:"listen"`
}

type Logger struct {
	Path             string `json:"path" yaml:"path"`
	RequestTableName string `json:"requestTableName" yaml:"requestTableName"`
}

type Mysql struct {
	User     string `json:"user" yaml:"user"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
	Used     bool   `json:"used" yaml:"used"`
}

type Mongo struct {
	User     string `json:"user" yaml:"user"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
	Used     bool   `json:"used" yaml:"used"`
}

type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
	Used     bool   `json:"used" yaml:"used"`
}

type runtimeParam struct {
	RootDir string `json:"-" yaml:"-"` // 此软件运行后的工作目录
}

type Configs struct {
	App          App          `json:"app" yaml:"app"`
	Logger       Logger       `json:"logger" yaml:"logger"`
	Mysql        Mysql        `json:"mysql" yaml:"mysql"`
	Mongo        Mongo        `json:"mongo" yaml:"mongo"`
	Redis        Redis        `json:"redis" yaml:"redis"`
	Spec         string       `json:"spec" yaml:"spec"`
	RuntimeParam runtimeParam `json:"-" yaml:"-"`
}
