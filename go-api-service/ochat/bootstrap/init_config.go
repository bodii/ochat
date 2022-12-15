package bootstrap

import (
	"ochat/comm"
)

var (
	SystemConf systemConfT
)

// init config centent
func InitConfig() {
	systemConfig()
	readDatabaseConfig()
}

// system config struct type
type systemConfT struct {
	App  appConfT  `yaml:"app"`
	Serv servConfT `yaml:"server"`
}

// system->app config struct type
type appConfT struct {
	Name     string `yaml:"name"`
	Document string `yaml:"document"`
	Env      string `yaml:"env"`
}

// system->serv config struct type
type servConfT struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

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

// read  system.yaml config and set var
func systemConfig() {
	SystemConf = comm.ReadYamlConfig[systemConfT]("system.yaml")
}

// read  database.yaml config and set var
func readDatabaseConfig() {
	mysqlConf = comm.ReadYamlConfig[mysqlConfT]("database.yaml")
}
