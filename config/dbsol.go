package config

import (
	"ggin/logs"
	"gopkg.in/mgo.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

// DBObj DBObj
var DBObj = NewInitDB()

func NewInitDB() *InitDB {
	idb := &InitDB{}
	idb.lock = sync.Mutex{}
	return idb
}

// InitDB 初始化数据库的连接
type InitDB struct {
	// DBConn 连接实例
	mysqlConn *gorm.DB
	mongoConn *mgo.Database
	lock      sync.Mutex
}

// Init Init
func (i *InitDB) Init() {
	i.initMysql()
	i.initMongo()
}

func (i *InitDB) initMysql() (done bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	mysqlConn := mysql.Open(Cfg.Mysql.Conn)
	db, err := gorm.Open(mysqlConn, &gorm.Config{
		Logger: logs.LogForDB(),
	})
	if err != nil {
		logs.Log.Error(err)
		return
	}
	i.mysqlConn = db
	return true
}

func (i *InitDB) initMongo() (done bool) {
	conn, err := mgo.Dial(Cfg.Mongo.Conn)
	if err != nil {
		return
	}
	i.mongoConn = conn.Copy().DB(Cfg.Mongo.Db)
	return true
}

// GetMysqlConn 得到数据库连接实例
func (i *InitDB) GetMysqlConn() *gorm.DB {
	if i.mysqlConn == nil {
		if !i.initMysql() {
			return nil
		}
	}
	return i.mysqlConn
}

// GetMongoConn 得到数据库连接实例
func (i *InitDB) GetMongoConn() *mgo.Database {
	if i.mongoConn == nil {
		if !i.initMongo() {
			return nil
		}
	}
	return i.mongoConn
}
