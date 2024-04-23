package gormPool

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	DbSetting struct {
		Common *Common       `yaml:"common"`
		MySql  *MySqlSetting `yaml:"mysql"`
	}

	Common struct {
		Driver string `yaml:"driver"`
	}

	Dsn struct {
		Name    string
		Content string
	}

	MySqlSetting struct {
		MaxOpenConns int                         `yaml:"maxOpenConns"`
		MaxIdleConns int                         `yaml:"maxIdleConns"`
		MaxLifetime  int                         `yaml:"maxLifetime"`
		MaxIdleTime  int                         `yaml:"maxIdleTime"`
		Rws          bool                        `yaml:"rws"`
		Main         *MySqlConnection            `yaml:"main"`
		Sources      map[string]*MySqlConnection `yaml:"sources"`
		Replicas     map[string]*MySqlConnection `yaml:"replicas"`
	}

	MySqlConnection struct {
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Host      string `yaml:"host"`
		Port      uint16 `yaml:"port"`
		Database  string `yaml:"database"`
		Charset   string `yaml:"charset"`
		Collation string `yaml:"collation"`
	}
)

// NewDbSetting 初始化
func NewDbSetting(path string) *DbSetting {
	var (
		file      []byte
		err       error
		dbSetting *DbSetting
	)
	file, err = os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件（db.yaml）失败：%s", err.Error()))
	}

	err = yaml.Unmarshal(file, &dbSetting)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件（db.yaml）失败：%s", err.Error()))
	}

	return dbSetting
}
