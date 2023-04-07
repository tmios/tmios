package config

import (
	"context"
	"crypto/tls"
	db_sql "database/sql"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
	logutil "tmios/lib/log"
	"tmios/lib/sql"
)

type Config struct {
	context.Context
	Db    *gorm.DB
	Conf  *CnfFile
	Rc    *resty.Client
	Cache *cache.Cache
}

type Option func(conf *Config)

var (
	cnf     *Config
	once    sync.Once
	options []Option
)

func (s *Config) Run() error {
	return nil
}
func NewConfig(ops ...Option) *Config {
	once.Do(func() {
		cnf = &Config{
			Context: context.Background(),
		}
		for _, op := range ops {
			op(cnf)
			options = append(options, op)
		}
	})
	return cnf
}

// WithConf 加载配置文件，排序第一
func WithConf(path string, watch bool) Option {
	return func(conf *Config) {
		conf.Conf = LoadConfigFile(path)
		if watch {
			go WatchConfigFile(path, func() {
				for _, op := range options {
					op(cnf)
				}
			})
		}
	}
}

// 在conf之后
func WithLog() Option {
	return func(conf *Config) {
		if err := logutil.Init(conf.Conf.Log.Level, conf.Conf.Log.Path); err != nil {
			logrus.Fatal(err)
		}
	}
}
func WithMysql() Option {
	return func(conf *Config) {
		//CreateDatabase(conf.Conf.MySQL.DatabaseName, conf.Conf.MySQL.Username, conf.Conf.MySQL.Password, conf.Conf.MySQL.Ip, conf.Conf.MySQL.Port)
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.Conf.MySQL.Username, conf.Conf.MySQL.Password, conf.Conf.MySQL.Ip, conf.Conf.MySQL.Port, conf.Conf.MySQL.DatabaseName))
		if err != nil {
			logrus.Fatal(err)
		}
		conf.Db = db.DB
	}
}

func WithResty() Option {
	return func(conf *Config) {
		conf.Rc = resty.New().SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).SetTimeout(time.Second * 20)
	}
}

func WithCache() Option {
	return func(conf *Config) {
		conf.Cache = cache.New(5*time.Minute, 10*time.Minute)
	}
}

// CreateDatabase 创建数据库
func CreateDatabase(dbName, username, password, ip string, port int) {

	db, err := db_sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, ip, port))
	if err != nil {
		logrus.Panic(err)
	}
	defer func(db *db_sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}(db)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci")
	if err != nil {
		logrus.Fatal(err)
	}
}
