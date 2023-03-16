package config

import (
	"errors"

	"github.com/spf13/viper"
)

// TODO: implement config validation and default values

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Environment string
	Port        int
	BaseUrl     string `mapstructure:"base_url"`
}

// LoadConfig takes in a filename and attempts to load in a config file using viper from the current directy, "./etc/config", and "/etc/config"
func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return parseConfig()
}

// parseConfig
func parseConfig() (*Config, error) {
	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, errors.New("unable to decode config file into Config struct")
	}

	return config, nil
}
