package config

import (
	"os"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App App
		DB  DB
	}

	App struct {
		Host string `mapstructure:"APP_HOST"`
		Port int    `mapstructure:"APP_PORT"`
	}

	DB struct {
		Host       string `mapstructure:"DB_HOST"`
		Port       int    `mapstructure:"DB_PORT"`
		Name       string `mapstructure:"DB_NAME"`
		User       string `mapstructure:"DB_USER"`
		Password   string `mapstructure:"DB_PASSWORD"`
		YelpApiKey string `mapstructure:"YELP_API_KEY"`
	}
)

func Load(configFile string) (*Config, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}

	var (
		app App
		db  DB
	)
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(configFile)

	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = v.Unmarshal(&app); err != nil {
		return nil, err
	}

	if err = v.Unmarshal(&db); err != nil {
		return nil, err
	}

	return &Config{app, db}, nil
}
