package store

type DBConfig struct {
	Host     string `yaml:"hots"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}
