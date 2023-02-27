package bootstrap

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

var (
	SystemConf     systemConfT
	HOST_AUTHORITY string   // domain:port e.g 127.0.0.1:8080
	HTTP_HOST      string   // scheme://domain:port
	HTTP_URL       *url.URL // return: url.URL struct
	// system_init_once sync.Once
	UploadDirs map[string]string
)

// init config centent
func InitSysConfig() {
	// loading system config
	initSystemConfig()
	// system_init_once.Do(initSystemConfig)

	initUploadDirectory()
}

// system config struct type
type systemConfT struct {
	App        appConfT    `toml:"app"`
	Serv       servConfT   `toml:"server"`
	UploadPath uploadPathT `toml:"upload_file_path"`
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

type uploadPathT struct {
	Image       string `toml:"image"`
	Video       string `toml:"video"`
	Voice       string `toml:"voice"`
	UserAvatar  string `toml:"user_avatar"`
	LoginQrcode string `toml:"login_qrcode"`
	UserQrcode  string `toml:"user_qrcode"`
	GroupQrcode string `toml:"group_qrcode"`
	GroupAvatar string `toml:"group_avatar"`
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

// init file upload directorys
func initUploadDirectory() {
	UploadDirs = map[string]string{
		"image":        SystemConf.UploadPath.Image,
		"video":        SystemConf.UploadPath.Video,
		"voice":        SystemConf.UploadPath.Voice,
		"user_avatar":  SystemConf.UploadPath.UserAvatar,
		"user_qrcode":  SystemConf.UploadPath.UserQrcode,
		"login_qrcode": SystemConf.UploadPath.LoginQrcode,
		"group_avatar": SystemConf.UploadPath.GroupAvatar,
		"group_qrcode": SystemConf.UploadPath.GroupQrcode,
	}

	for _, f := range UploadDirs {
		n, err := os.Stat(f)
		if err != nil && n == nil {
			go os.MkdirAll(f, os.ModePerm)
		}
	}
}
