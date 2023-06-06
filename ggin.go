package ggin

import (
	"github.com/gzxgogh/ggin/config"
	"github.com/gzxgogh/ggin/logs"
)

func Init(configFile string) {
	logs.InitLogger()
	config.DBObj.Init()
	return
}
