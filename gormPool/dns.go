package gormPool

type (
	Dsn struct {
		Name    string
		Content string
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
