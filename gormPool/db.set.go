package gormPool

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	DbSetting struct {
		Common    *Common           `yaml:"common"`
		MySql     *MySqlSetting     `yaml:"mysql"`
		Postgres  *PostgresSetting  `yaml:"postgres"`
		SqlServer *SqlServerSetting `yaml:"sqlServer"`
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

	PostgresSetting struct {
		MaxOpenConns int                 `yaml:"maxOpenConns"`
		MaxIdleConns int                 `yaml:"maxIdleConns"`
		MaxLifetime  int                 `yaml:"maxLifetime"`
		MaxIdleTime  int                 `yaml:"maxIdleTime"`
		Main         *PostgresConnection `yaml:"main"`
	}

	PostgresConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
		TimeZone string `yaml:"timezone"`
		SslMode  string `yaml:"sslmode"`
	}

	SqlServerSetting struct {
		MaxOpenConns int                  `yaml:"maxOpenConns"`
		MaxIdleConns int                  `yaml:"maxIdleConns"`
		MaxLifetime  int                  `yaml:"maxLifetime"`
		MaxIdleTime  int                  `yaml:"maxIdleTime"`
		Main         *SqlServerConnection `yaml:"main"`
	}

	SqlServerConnection struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Database string `yaml:"database"`
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
		panic(fmt.Sprintf("读取配置文件（数据库）失败：%s", err.Error()))
	}

	err = yaml.Unmarshal(file, &dbSetting)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件（yaml）失败：%s", err.Error()))
	}

	return dbSetting
}
