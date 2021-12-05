package config

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/chenhqchn/gotools/logs"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	global *Config
	logger *zap.SugaredLogger
	db     *sql.DB
	gormdb *gorm.DB
)

func G() *Config {
	if global == nil {
		L().Error("Getting global config failed")
		return nil
	}
	return global
}

func L() *zap.SugaredLogger {
	if logger == nil {
		L().Error("Getting logger failed")
		return nil
	}
	return logger
}

func D() *sql.DB {
	if db == nil {
		L().Error("Unable to connect to the database")
		return nil
	}
	return db
}

func DG() *gorm.DB {
	if gormdb == nil {
		L().Error("Unable to connect to the database")
		return nil
	}
	return gormdb
}

type Config struct {
	Mysql *Mysql `toml:"mysql"`
	Log   *Log   `toml:"log"`
	Jwt   *Jwt   `toml:"jwt"`
}

type Mysql struct {
	Username          string `toml:"username"`
	Password          string `toml:"password"`
	Database          string `toml:"database"`
	Host              string `toml:"host"`
	Port              string `toml:"port"`
	MaxIdleConnection int    `toml:"max_idle_connection"`
	MaxOpenConnection int    `toml:"max_open_connection"`
	ShowType          string `toml:"sql_type"`
	lock              sync.Mutex
}

// 原生 sql
func (ms *Mysql) getDBConn() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true&parseTime=true",
		ms.Username, ms.Password, ms.Host, ms.Port, ms.Database)
	db, err := sql.Open(ms.ShowType, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(ms.MaxOpenConnection)
	db.SetMaxIdleConns(ms.MaxIdleConnection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// gorm 连接 mysql
func (ms *Mysql) gormConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true&parseTime=true",
		ms.Username, ms.Password, ms.Host, ms.Port, ms.Database)
	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormdb, nil
}

func (ms *Mysql) GetDB() (*sql.DB, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	if db == nil {
		conn, err := ms.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil
}

type Log struct {
	Level      string `toml:"level"`
	Filepath   string `toml:"filepath"`
	Format     string `toml:"format"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

type Jwt struct {
	ExpireTime int    `toml:"expire_time"`
	Secret     string `toml:"secret"`
	Issuer     string `toml:"issuer"`
}

// 加载 toml 配置文件，解析至 Config struct
func LoadConfig() error {
	var config Config

	if _, err := toml.DecodeFile("config/base.toml", &config); err != nil {
		return err
	}
	global = &config

	return nil
}

// 导入自定义 logger 模块，并初始化
func InitLogger() {
	logs.SetLevel(global.Log.Level)
	logs.SetFilepath(global.Log.Filepath)
	logs.SetFormat(global.Log.Format)
	logs.SetMaxSize(global.Log.MaxSize)
	logs.SetMaxAge(global.Log.MaxAge)
	logs.SetCompress(global.Log.Compress)

	logger = logs.InitLogger()
}

// 初始化 mysql 引擎
func InitDB() error {
	tmp_db, err := global.Mysql.GetDB()
	if err != nil {
		return err
	}

	db = tmp_db

	tmp_gorm, err := global.Mysql.gormConn()
	if err != nil {
		return err
	}
	gormdb = tmp_gorm

	return nil
}
