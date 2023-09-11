package db

import (
	"github.com/go-redis/redis"
	"github.com/gzxgogh/ggin/config"
	"github.com/gzxgogh/ggin/logs"
	"github.com/streadway/amqp"
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
	mysqlConn  *gorm.DB
	mongoConn  *mgo.Database
	redisConn  *redis.Client
	rabbitConn *amqp.Channel
	lock       sync.Mutex
}

// Init Init
func (i *InitDB) Init() {
	i.initMysql()
	i.initMongo()
	i.initRedis()
	i.initRabbit()
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
	mysqlConn := mysql.Open(config.Cfg.Mysql.Conn)
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
	db, err := mgo.Dial(config.Cfg.Mongo.Conn)
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
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.Db,
	})
	i.redisConn = rdb
	return true
}

func (i *InitDB) initRabbit() (done bool) {
	if config.Cfg.RabbitMq.Used == false {
		return false
	}
	conn, err := amqp.Dial(config.Cfg.RabbitMq.Conn)
	if err != nil {
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		return false
	}
	i.rabbitConn = ch
	return true
}

// GetMysqlConn 得到Mysql数据库连接实例
func (i *InitDB) GetMysqlConn() *gorm.DB {
	if i.mysqlConn == nil {
		if !i.initMysql() {
			return nil
		}
	}
	return i.mysqlConn
}

// GetMongoConn 得到Mongo数据库连接实例
func (i *InitDB) GetMongoConn() *mgo.Database {
	if i.mongoConn == nil {
		if !i.initMongo() {
			return nil
		}
	}
	return i.mongoConn
}

// GetRedisConn 得到Redis数据库连接实例
func (i *InitDB) GetRedisConn() *redis.Client {
	if i.redisConn == nil {
		if !i.initRedis() {
			return nil
		}
	}
	return i.redisConn
}

// GetRabbitConn 得到Rabbit数据库连接实例
func (i *InitDB) GetRabbitConn() *amqp.Channel {
	if i.rabbitConn == nil {
		if i.initRabbit() {
			return nil
		}
	}
	return i.rabbitConn
}
