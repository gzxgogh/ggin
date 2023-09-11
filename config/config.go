package config

// Cfg 全局的Config配置，解析dns.yaml的结果
var Cfg *Configs

type App struct {
	Name   string `json:"name" yaml:"name"`
	Listen string `json:"listen" yaml:"listen"`
}

type Logger struct {
	Path             string `json:"path" yaml:"path"`
	Type             string `json:"type" yaml:"type"`
	RequestTableName string `json:"requestTableName" yaml:"requestTableName"`
}

type Mysql struct {
	Conn string `json:"conn" yaml:"conn"`
	Used bool   `json:"used" yaml:"used"`
}

type Mongo struct {
	Conn string `json:"conn" yaml:"conn"`
	Db   string `json:"db" yaml:"db"`
	Used bool   `json:"used" yaml:"used"`
}

type Redis struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
	Used     bool   `json:"used" yaml:"used"`
}

type RabbitMq struct {
	Conn string `json:"conn" yaml:"conn"`
	Used bool   `json:"password" yaml:"password"`
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
	RabbitMq     RabbitMq     `json:"rabbitMq" yaml:"rabbitMq"`
	Spec         string       `json:"spec" yaml:"spec"`
	RuntimeParam runtimeParam `json:"-" yaml:"-"`
}
