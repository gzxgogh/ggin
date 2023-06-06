package ggin

import (
	"github.com/gzxgogh/ggin/config"
	"github.com/gzxgogh/ggin/db"
	"github.com/gzxgogh/ggin/logs"
)

func Init(configFile string) {
	config.Init(configFile)
	logs.InitLogger()
	db.DBObj.Init()
	return
}
