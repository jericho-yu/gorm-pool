package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/jericho-yu/gorm-pool/gormPool"
)

type (
	DbSetting struct {
		Common *Common                `yaml:"common"`
		MySql  *gormPool.MySqlSetting `yaml:"mysql"`
	}

	Common struct {
		Driver string `yaml:"driver"`
	}
)

// NewDbSetting 初始化
func NewDbSetting(path string) *DbSetting {
	var (
		file      []byte
		err       error
		dbSetting *DbSetting
	)
	file, err = os.ReadFile(path + "/db.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件（db.yaml）失败：%s", err.Error()))
	}

	err = yaml.Unmarshal(file, &dbSetting)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件（db.yaml）失败：%s", err.Error()))
	}

	return dbSetting
}
