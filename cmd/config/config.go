package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Database DatabaseSection `yaml:"database"`
	App      AppSection      `yaml:"application"`
}

type DatabaseSection struct {
	Name     string `yaml:"POSTGRES_DB"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Port     string `yaml:"PORT"`
	Host     string `yaml:"HOST"`
}

type AppSection struct {
	Port   string `yaml:"port"`
	JWTKey string `yaml:"jwtkey"`
}

func GetConfig() (*AppConfig, error) {
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла YAML: %v", err)
		return nil, err
	}

	var config AppConfig

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Ошибка разбора YAML: %v", err)
		return nil, err
	}

	return &config, nil
}
