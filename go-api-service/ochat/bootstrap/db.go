package bootstrap

import (
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	DB_Engine   *xorm.Engine
	_dbInitOnce sync.Once
	_mysqlConf  mysqlConfT
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

	_dbInitOnce.Do(initDbConnect)

	return DB_Engine
}

// read  database.yaml config and set var
func loadDatabaseConfig() {
	mysqlListConf := readTomlConfig[mysqlListConfT]("database.toml")
	_mysqlConf = mysqlListConf.DBs[0]
}

func initDbConnect() {
	// load database config info
	loadDatabaseConfig()

	dataSource := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=%s",
		_mysqlConf.User,
		_mysqlConf.Password,
		_mysqlConf.Host,
		_mysqlConf.Port,
		_mysqlConf.Database,
		_mysqlConf.Charset,
	)

	dbEngine, err := xorm.NewEngine(
		_mysqlConf.DriverName, dataSource)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// show sql
	dbEngine.ShowSQL(_mysqlConf.ShowSQL)

	// max connect number
	dbEngine.SetMaxOpenConns(_mysqlConf.MaxLineNums)

	// auto sync
	// DB_Engine.Sync2()
	DB_Engine = dbEngine

	log.Println("init database success!")
}
