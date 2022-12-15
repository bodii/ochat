package bootstrap

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	DB_Engine    *xorm.Engine
	db_init_once sync.Once
	mysqlConf    mysqlConfT
)

func DBOnceInit() *xorm.Engine {

	db_init_once.Do(initConnect)
	// fmt.Printf("%#v\n", DB_Engine)

	return DB_Engine
}

func initConnect() {
	dataSource := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=%s",
		mysqlConf.User,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.Database,
		mysqlConf.Charset,
	)

	db_engine, err := xorm.NewEngine(
		mysqlConf.DriverName, dataSource)
	if err != nil {
		panic(err.Error())
	}

	// show sql
	db_engine.ShowSQL(mysqlConf.ShowSQL)

	// max connect number
	db_engine.SetMaxOpenConns(mysqlConf.MaxLineNums)

	// auto sync
	// db_engine.Sync2()

	DB_Engine = db_engine

	log.Println("init database success!")
}
