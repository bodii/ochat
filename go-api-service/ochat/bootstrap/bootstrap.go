package bootstrap

import (
	"os"
	"path"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

var (
	PROJECT_DIR string = projectDIR()
)

func Init() {
	// 初始化日志
	InitLog()

	// 初始化系统配置项
	InitSysConfig()

	// 创建数据库链接
	RedisOnceInit()

	// 创建数据库链接
	DBOnceInit()
}

// get the project root directory
//
// param:
//
// return:
//   - [string]: dir path
func projectDIR() string {
	dir, err := os.Getwd()

	if err != nil {
		return os.Args[0]
	}

	return dir
}

// Read the configuration content of a yaml file type
//
//	@this is a generic func [T]
//
// param:
//   - file [string]: a string file path
//
// return:
//   - [T]: type
func readYamlConfig[T any](file string) T {
	if 1 > strings.LastIndex(file, ".yaml") {
		panic("input file name not is yaml file")
	}

	configPath := path.Join(projectDIR(), "config/"+file)
	content, _ := os.ReadFile(configPath)

	var t T
	err := yaml.Unmarshal(content, &t)
	if err != nil {
		panic(err.Error())
	}

	filename := strings.TrimSuffix(file, ".yaml")

	Log.Info("read %s config succuee!", filename)

	return t
}

// Read the configuration content of a toml file type
//
//	@this is a generic func [T]
//
// param:
//   - file [string]: a string file path
//
// return:
//   - [T]: type
func readTomlConfig[T any](file string) T {
	if 1 > strings.LastIndex(file, ".toml") {
		panic("input file name not is toml file")
	}

	configPath := path.Join(projectDIR(), "config", file)
	content, _ := os.ReadFile(configPath)

	var t T
	err := toml.Unmarshal(content, &t)
	// fmt.Printf("%#v\n", t)
	if err != nil {
		panic(err.Error())
	}

	filename := strings.TrimSuffix(file, ".toml")

	Log.Info("read %s config succuee!", filename)

	return t
}
