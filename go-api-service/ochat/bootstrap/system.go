package bootstrap

import (
	"fmt"
	"net/url"
	"ochat/comm"
)

var (
	SystemConf      systemConfT
	HOST_NAME       string
	HTTP_HOST       string
	HTTP_URL        *url.URL
	HTTP_Avatar_URI string
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
}

// system->app config struct type
type appConfT struct {
	Name     string `yaml:"name"`
	Document string `yaml:"document"`
	Env      string `yaml:"env"`
}

// system->serv config struct type
type servConfT struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// system->avatar config struct type
type avatarConfT struct {
	FileDir       string `yaml:"file_dir"`
	Uri           string `yaml:"uri"`
	SuffixName    string `yaml:"suffix_name"`
	DefaultAvatar string `yaml:"default_avatar"`
	UploadDir     string `yaml:"upload_dir"`
}

// system->login_qrcode config struct type
type loginQRCodeT struct {
	FileDir    string `yaml:"file_dir"`
	Uri        string `yaml:"uri"`
	SuffixName string `yaml:"suffix_name"`
}

// read  system.yaml config and set var
func initSystemConfig() {
	SystemConf = comm.ReadYamlConfig[systemConfT]("system.yaml")
	servConf := SystemConf.Serv
	HTTP_URL := &url.URL{
		Scheme: servConf.Protocol,
		Host:   fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
	}
	HOST_NAME = fmt.Sprintf("%s:%s", HTTP_URL.Hostname(), HTTP_URL.Port())
	HTTP_HOST = HTTP_URL.String()

	HTTP_Avatar_URI = fmt.Sprintf("%s%s", HTTP_HOST, SystemConf.Avatar.Uri)
}