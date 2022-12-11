package app

import (
	"github.com/GritselMaks/BT_API/internal/store/postgresql"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Host       string `env:"HOST,required"`
	Port       string `env:"PORT,required"`
	LogLevel   string `env:"LOGLEVEL,required"`
	LogPath    string `env:"LOGPATH,required"`
	Store      *postgresql.DBConfig
	LocalStore string `env:"LOCAL_STORE,required"`
}

// Initialize config
func LoadConfig() (*Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	var dbcfg postgresql.DBConfig
	err = env.Parse(&dbcfg)
	if err != nil {
		return nil, err
	}
	cfg.Store = &dbcfg
	return &cfg, nil
}
