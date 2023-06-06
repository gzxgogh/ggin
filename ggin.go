package ggin

import (
	"github.com/gzxgogh/ggin/db"
	"github.com/gzxgogh/ggin/logs"
)

func Init(configFile string) {
	logs.InitLogger()
	db.DBObj.Init()
	return
}
