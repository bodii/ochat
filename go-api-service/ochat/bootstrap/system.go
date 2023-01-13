package bootstrap

import (
	"fmt"
	"log"
	"net/url"
)

var (
	SystemConf     systemConfT
	HOST_AUTHORITY string   // domain:port e.g 127.0.0.1:8080
	HTTP_HOST      string   // scheme://domain:port
	HTTP_URL       *url.URL // return: url.URL struct
	// system_init_once sync.Once
)

// init config centent
func InitSysConfig() {
	// loading system config
	initSystemConfig()
	// system_init_once.Do(initSystemConfig)

}

// system config struct type
type systemConfT struct {
	App  appConfT  `toml:"app"`
	Serv servConfT `toml:"server"`
}

// system->app config struct type
type appConfT struct {
	Name     string `toml:"name"`
	Document string `toml:"document"`
	Env      string `toml:"env"`
}

// system->serv config struct type
type servConfT struct {
	Scheme string `toml:"scheme"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
}

// read  system.yaml config and set var
func initSystemConfig() {
	SystemConf = readTomlConfig[systemConfT]("system.toml")
	servConf := SystemConf.Serv
	HTTP_URL := &url.URL{
		Scheme: servConf.Scheme,
		Host:   fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
	}
	HOST_AUTHORITY = fmt.Sprintf("%s:%s", HTTP_URL.Hostname(), HTTP_URL.Port())
	HTTP_HOST = HTTP_URL.String()

	log.Println("init system config success!")
}
