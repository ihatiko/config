package config

import (
	"fmt"
	cfg "github.com/ihatiko/viper-env"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	lg "log"
	"os"
)

func GetConfig[T any]() (*T, error) {
	config := localConfigPath
	stand := os.Getenv(standENV)
	if stand == productionStand {
		config = productionConfigPath
	}
	path := fmt.Sprintf("%s/%s", configFolder, config)
	cfgFile, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig[T](cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadConfig(filename string) (*cfg.Config, error) {
	cfg := cfg.New(viper.New())
	cfg.SetConfigName(filename)
	cfg.AddConfigPath(".")
	cfg.AutomaticEnv()
	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return cfg, nil
}

func ParseConfig[T any](v *cfg.Config) (*T, error) {
	var c T
	err := v.Unmarshal(&c)
	if err != nil {
		lg.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
