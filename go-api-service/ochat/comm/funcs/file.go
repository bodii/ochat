package funcs

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func GetProjectDIR() string {
	dir, err := os.Getwd()

	if err != nil {
		return os.Args[0]
	}

	return dir
}

// Read the configuration content of a yaml file type
// Parmas: [file string] filename
// Returns: Specifies a structure of type [T]
func ReadYamlConfig[T any](file string) T {
	if 1 > strings.LastIndex(file, ".yaml") {
		panic("input file name not is yaml file")
	}

	configPath := path.Join(GetProjectDIR(), "config/"+file)
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

// RandFileName func
//
// param:
//   - suffix the affix of the file name to be generated.
//
// returns a file name randomly generated by suffix.
func RandFileName(suffix string) string {
	return fmt.Sprintf("%d%04d%s",
		time.Now().Unix(),
		GetUnixNanoRandSeed().Int31(),
		suffix)
}
