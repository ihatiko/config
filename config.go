package config

import (
	"github.com/ihatiko/config/parser"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	lg "log"
)

const (
	configPath = "./config/config"
)

type Default struct {
	Path string
}

type Option func(p *Default)

func (d *Default) Validate() error {
	if d.Path == "" {
		return errors.New("empty project path")
	}
	return nil
}

func WithPath(path string) Option {
	return func(p *Default) {
		p.Path = path
	}
}

func GetConfig[T any](opts ...Option) (*T, error) {
	d := &Default{Path: configPath}
	for _, opt := range opts {
		opt(d)
	}

	cfgFile, err := LoadConfig(d.Path)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig[T](cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadConfig(filename string) (*parser.Config, error) {
	cfg := parser.New(viper.New())
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

func ParseConfig[T any](v *parser.Config) (*T, error) {
	var c T
	err := v.Unmarshal(&c)
	if err != nil {
		lg.Printf("unable to decode into struct, %v", err)
		return nil, err
	}
	return &c, nil
}
