package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gzxgogh/ggin/config"
	"github.com/gzxgogh/ggin/logs"
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
	redisConn *redis.Client
	lock      sync.Mutex
}

// Init Init
func (i *InitDB) Init() {
	i.initMysql()
	i.initMongo()
	i.initRedis()
}

func (i *InitDB) initMysql() (done bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if config.Cfg.Mysql.Used == false {
		return false
	}
	conn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`,
		config.Cfg.Mysql.User, config.Cfg.Mysql.Password, config.Cfg.Mysql.Host, config.Cfg.Mysql.Port, config.Cfg.Mysql.Db)
	mysqlConn := mysql.Open(conn)
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
	if config.Cfg.Mongo.Used == false {
		return false
	}
	conn := fmt.Sprintf(`mongodb://%s:%s@%s:%d`, config.Cfg.Mongo.User, config.Cfg.Mongo.Password, config.Cfg.Mongo.Host, config.Cfg.Mongo.Port)
	db, err := mgo.Dial(conn)
	if err != nil {
		logs.Log.Error(err)
		return
	}
	i.mongoConn = db.Copy().DB(config.Cfg.Mongo.Db)
	return true
}

func (i *InitDB) initRedis() (done bool) {
	if config.Cfg.Redis.Used == false {
		return false
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%d`, config.Cfg.Redis.Host, config.Cfg.Redis.Port),
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.Db,
	})
	i.redisConn = rdb
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

// GetRedisConn 得到数据库连接实例
func (i *InitDB) GetRedisConn() *redis.Client {
	if i.redisConn == nil {
		if !i.initRedis() {
			return nil
		}
	}
	return i.redisConn
}
