package base

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var (
	Env      string
	ConfPath string
)

func ParseConfEnvPath(confPath string) {
	pathArr := strings.Split(confPath, "/")
	ConfPath = strings.Join(pathArr[:len(pathArr)-1], "/")
	Env = pathArr[len(pathArr)-2]
}

func GetConfEnv() string {
	return Env
}

func GetConfPath(fileName string) string {
	return ConfPath + "/" + fileName + ".toml"
}

func GetConfFilePath(fileName string) string {
	return ConfPath + "/" + fileName
}

func ParseConfig(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %v fail, %v", path, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail, %v", err)
	}

	v := viper.New()
	v.SetConfigType("toml")
	err = v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("parse config fail, config:%v, err:%v", string(data), err)
	}
	return nil
}
