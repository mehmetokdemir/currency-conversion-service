package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Db     DbConfig `mapstructure:"db"`
	Server Server   `mapstructure:"server"`
}

type DbConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type Server struct {
	Port string
}

var vp *viper.Viper

func LoadConfig() (*Config, error) {
	vp = viper.New()
	var config *Config

	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./config")
	vp.AddConfigPath(".")
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := vp.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
