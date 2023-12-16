package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Server struct {
		Host     string        `yaml:"host" validate:"required"`
		Secret   string        `yaml:"secret" validate:"required"`
		TokenTTL time.Duration `yaml:"tokenTTL" validate:"required"`
		Storage  struct {
			Postgres *Postgres `yaml:"postgres" validate:"required"`
		} `yaml:"storage" validate:"required"`
	}

	Postgres struct {
		DSN              string `yaml:"dsn" validate:"required"`
		SourceMigrations string `yaml:"sourceMigrations"`
	}
)

func NewServer(configFile string) (*Server, error) {
	if configFile == "" {
		return nil, ErrInvalidConfigFile
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	s := new(Server)

	if err := viper.Unmarshal(&s); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return nil, err
	}

	return s, nil
}
