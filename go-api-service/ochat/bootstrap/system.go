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
	App         appConfT     `yaml:"app"`
	Serv        servConfT    `yaml:"server"`
	Avatar      avatarConfT  `yaml:"avatar"`
	LoginQRCode loginQRCodeT `yaml:"login_qrcode"`
	UserQRCode  UserQRCodeT  `yaml:"user_qrcode"`
}

// system->app config struct type
type appConfT struct {
	Name     string `yaml:"name"`
	Document string `yaml:"document"`
	Env      string `yaml:"env"`
}

// system->serv config struct type
type servConfT struct {
	Scheme string `yaml:"scheme"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

// system->avatar config struct type
type avatarConfT struct {
	FileDir       string `yaml:"file_dir"`
	SuffixName    string `yaml:"suffix_name"`
	DefaultAvatar string `yaml:"default_avatar"`
	UploadDir     string `yaml:"upload_dir"`
}

// system->login_qrcode config struct type
type loginQRCodeT struct {
	FileDir    string `yaml:"file_dir"`
	SuffixName string `yaml:"suffix_name"`
}

// system->user_qrcode config struct type
type UserQRCodeT struct {
	FileDir    string `yaml:"file_dir"`
	SuffixName string `yaml:"suffix_name"`
}

// read  system.yaml config and set var
func initSystemConfig() {
	SystemConf = readYamlConfig[systemConfT]("system.yaml")
	servConf := SystemConf.Serv
	HTTP_URL := &url.URL{
		Scheme: servConf.Scheme,
		Host:   fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
	}
	HOST_AUTHORITY = fmt.Sprintf("%s:%s", HTTP_URL.Hostname(), HTTP_URL.Port())
	HTTP_HOST = HTTP_URL.String()

	log.Println("init system config success!")
}
