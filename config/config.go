package config

import (
	"github.com/go-xorm/xorm"
	"github.com/iphuket/gowc/config/cache"
	_ "github.com/mattn/go-sqlite3" // go-sqlite3
)

// Cache 配置相关
type Cache struct {
	cache.Cache
}

// DB ...
type DB struct {
}

// NewEngine xorm
func (db *DB) NewEngine() (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("sqlite3", "./gowc.db")
	return
}
