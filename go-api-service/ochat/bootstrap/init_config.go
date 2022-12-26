package bootstrap

import (
	"fmt"
	"ochat/comm"
)

var (
	SystemConf      systemConfT
	HTTP_HOST       string
	HTTP_Avatar_URI string
)

// init config centent
func InitConfig() {
	systemConfig()
	readDatabaseConfig()
}

// system config struct type
type systemConfT struct {
	App    appConfT    `yaml:"app"`
	Serv   servConfT   `yaml:"server"`
	Avatar avatarConfT `yaml:"avatar"`
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

// system->avatar config struct type
type avatarConfT struct {
	FileDir       string `yaml:"file_dir"`
	Uri           string `yaml:"uri"`
	SuffixName    string `yaml:"suffix_name"`
	DefaultAvatar string `yaml:"default_avatar"`
	UploadDir     string `yaml:"upload_dir"`
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
	HTTP_HOST = fmt.Sprintf("http://%s:%d",
		SystemConf.Serv.Host, SystemConf.Serv.Port)

	HTTP_Avatar_URI = fmt.Sprintf("%s%s", HTTP_HOST, SystemConf.Avatar.Uri)
}

// read  database.yaml config and set var
func readDatabaseConfig() {
	mysqlConf = comm.ReadYamlConfig[mysqlConfT]("database.yaml")
}
