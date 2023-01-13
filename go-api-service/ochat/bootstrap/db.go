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

type mysqlListConfT struct {
	DBs []mysqlConfT `toml:"mysql"`
}

// db mysql config struct type
type mysqlConfT struct {
	DriverName  string `toml:"driver_name"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	Charset     string `toml:"charset"`
	MaxLineNums int    `toml:"max_line_nums"`
	ShowSQL     bool   `toml:"show_sql"`
}

func DBOnceInit() *xorm.Engine {

	db_init_once.Do(initDbConnect)
	// fmt.Printf("%#v\n", DB_Engine)

	return DB_Engine
}

// read  database.yaml config and set var
func loadDatabaseConfig() {
	mysqlListConf := readTomlConfig[mysqlListConfT]("database.toml")
	mysqlConf = mysqlListConf.DBs[0]
}

func initDbConnect() {
	// load database config info
	loadDatabaseConfig()

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
