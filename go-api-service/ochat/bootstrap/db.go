package bootstrap

import (
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var db_engine *xorm.Engine
var once sync.Once

func DBOnceInit() *xorm.Engine {
	if db_engine == nil {
		once.Do(initConnect)
	}

	return db_engine
}

func initConnect() {
	driverName := "mysql"
	dataSource := "root:123456@(127.0.0.1:3306)/ochat_database?charset=utf8mb4"

	db_engine, err := xorm.NewEngine(driverName, dataSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	// show sql
	db_engine.ShowSQL(true)

	// max connect number
	db_engine.SetMaxOpenConns(2)

	// auto sync
	// db_engine.Sync2()

	log.Println("init database success!")
}
