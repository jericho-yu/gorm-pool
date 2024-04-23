package gormPool

type (
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
