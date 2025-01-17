package configuration

import (
	"context"
	"errors"
	"github.com/vitortenor/lead-stream-service/internal/tools"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
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

func InitConfig(_ context.Context, fileName string) (*Config, error) {
	log.Println("Loading configuration file...")

	path, err := tools.FindProjectRoot()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(path, "resources", fileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("configuration file does not exist")
	}

	data, err := os.ReadFile(configPath)
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
