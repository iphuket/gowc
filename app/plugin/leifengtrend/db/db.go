package db

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3" // go-sqlite3
)

// NewEngine xorm
func NewEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", "./leifengtrend.db")
	return engine, err
}
