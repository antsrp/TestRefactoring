package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	CONFIGS_RELATIVE_PATH = "configs"
	CONFIG_FILENAME       = "config.yaml"
)

type Config struct {
	URL struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"url"`
	Database struct {
		Source string `yaml:"source"`
	} `yaml:"db"`
}

func getPathToConfigsFolder() string {
	curPath, _ := os.Getwd()
	return filepath.Join(curPath, CONFIGS_RELATIVE_PATH)
}

func ParseConfig(logger *zap.Logger) *Config {
	f, err := os.Open(fmt.Sprintf("%s//%s", getPathToConfigsFolder(), CONFIG_FILENAME))
	if err != nil {
		logger.Sugar().Fatal("Can't read config of db: ", err)
	}
	defer f.Close()

	var cfg Config

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Sugar().Fatal("Can't parse config of db: ", err)
	}

	return &cfg
}
