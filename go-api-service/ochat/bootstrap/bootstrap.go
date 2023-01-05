package bootstrap

import (
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	PROJECT_DIR string = projectDIR()
)

func Init() {
	// 初始化系统配置项
	InitSysConfig()

	// 创建数据库链接
	RedisOnceInit()

	// 创建数据库链接
	DBOnceInit()
}

func projectDIR() string {
	dir, err := os.Getwd()

	if err != nil {
		return os.Args[0]
	}

	return dir
}

// Read the configuration content of a yaml file type
// Parmas: [file string] filename
// Returns: Specifies a structure of type [T]
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

	log.Printf("read %s config succuee!", filename)

	return t
}
