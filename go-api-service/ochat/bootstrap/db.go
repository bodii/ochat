package bootstrap

import (
	"fmt"
	"log"
	"ochat/comm"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	DB_Engine    *xorm.Engine
	db_init_once sync.Once
	mysqlConf    mysqlConfT
)

// db mysql config struct type
type mysqlConfT struct {
	DriverName  string `yaml:"driver_name"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	Charset     string `yaml:"charset"`
	MaxLineNums int    `yaml:"max_line_nums"`
	ShowSQL     bool   `yaml:"show_sql"`
}

func DBOnceInit() *xorm.Engine {

	db_init_once.Do(initDbConnect)
	// fmt.Printf("%#v\n", DB_Engine)

	return DB_Engine
}

// read  database.yaml config and set var
func loadDatabaseConfig() {
	mysqlConf = comm.ReadYamlConfig[mysqlConfT]("database.yaml")
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
