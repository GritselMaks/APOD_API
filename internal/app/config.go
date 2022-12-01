package app

import (
	"io/ioutil"

	"github.com/GritselMaks/BT_API/internal/store/postgresql"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Http       Http                 `yaml:"http"`
	LogLevel   string               `yaml:"log_level"`
	LogPath    string               `yaml:"log_path"`
	Store      *postgresql.DBConfig `yaml:"store"`
	LocalStore string               `yaml:"local_store"`
}
type Http struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Initialize config
func LoadConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
