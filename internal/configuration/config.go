package configuration

import (
	"context"
	"errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		API struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
		} `yaml:"api"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URI        string            `yaml:"uri"`
		Name       string            `yaml:"name"`
		Collection map[string]string `yaml:"collection"`
	} `yaml:"database"`
}

func InitConfig(_ context.Context, path string) (*Config, error) {
	log.Println("Loading configuration file...")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("configuration file does not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host is required")
	}
	if config.Server.Port == 0 {
		return errors.New("server port is required")
	}
	if config.Database.URI == "" {
		return errors.New("database URI is required")
	}
	if config.Database.Name == "" {
		return errors.New("database name is required")
	}
	if len(config.Database.Collection) == 0 {
		return errors.New("at least one database collection is required")
	}
	return nil
}
