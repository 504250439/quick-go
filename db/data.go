package db

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"quick-go/conf"
	"quick-go/log"
	"quick-go/utils"
)

var (
	// mysql 连接
	DBLocalhost *gorm.DB

	// redis连接
	RedisLocal *redis.Client
)

func InitRedis() (err error) {
	RedisLocal, err = redisConnect("redisLocal")
	if err != nil {
		return err
	}
	return nil
}

// 初始化连接
func redisConnect(key string) (rdb *redis.Client, err error) {
	addr := conf.Env.GetString(key+".host") + ":" + conf.Env.GetString(key+".port")
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Env.GetString(key + ".password"),
		DB:       0,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		log.ErrorLogger.Error("redis连接异常", zap.Error(err), zap.String("connectInfo", key),
			zap.String("addr", addr))
		return nil, err
	}
	return rdb, nil
}

func InitMysql() (err error) {
	// 建立 MySQL 连接
	DBLocalhost, err = mysqlConnect("TEST", "dbLocal")
	if err != nil {
		return err
	}

	return nil
}

func mysqlConnect(dbName string, key string) (db *gorm.DB, err error) {
	username := conf.Env.GetString(key + ".user")
	pw := conf.Env.GetString(key + ".pwd")
	host := conf.Env.GetString(key + ".host")
	port := conf.Env.GetString(key + ".port")
	dsn := utils.StringConcat("", username, ":", pw, "@tcp(", host, ":", port, ")/", dbName, "?timeout=5s&readTimeout=5s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.ErrorLogger.Info("", zap.Error(err), zap.String("connect info", dsn))
		return nil, err
	}
	return db, nil
}