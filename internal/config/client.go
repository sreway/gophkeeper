package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Client struct {
		Server         string `yaml:"server" validate:"required"`
		MasterPassword string `yaml:"masterPassword"`
		Storage        struct {
			SQLite *SQLite `yaml:"sqlite" validate:"required"`
		} `yaml:"storage" validate:"required"`
	}

	SQLite struct {
		DSN string `yaml:"dsn" validate:"required"`
	}
)

func NewClient(configFile string) (*Client, error) {
	if configFile == "" {
		return nil, ErrInvalidConfigFile
	}

	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("KEEPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.BindEnv("client.secretKey")
	if err != nil {
		return nil, err
	}

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := new(Client)

	if err = viper.Unmarshal(c); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err = validate.Struct(c); err != nil {
		return nil, err
	}

	return c, nil
}
