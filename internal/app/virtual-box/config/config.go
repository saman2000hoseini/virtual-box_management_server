package config

import (
	"time"

	"github.com/saman2000hoseini/virtual-box_management_server/pkg/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

const (
	app       = "virtual-box"
	cfgFile   = "config.yaml"
	cfgPrefix = "virtual-box"
)

type (
	Config struct {
		Server Server `mapstructure:"server"`
	}

	// Server represents server configuration struct.
	Server struct {
		Address      string        `mapstructure:"address"`
		ReadTimeout  time.Duration `mapstructure:"read-timeout"`
		WriteTimeout time.Duration `mapstructure:"write-timeout"`
	}
)

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgFile, &cfg, defaultConfig, cfgPrefix)

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("failed to validate configurations: %s", err.Error())
	}

	return cfg
}
