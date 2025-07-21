package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Writer  *struct {
			Topic string `yaml:"topic"`
		} `yaml:"writer"`
		Reader *struct {
			Topic   string `yaml:"topic"`
			GroupID string `yaml:"groupID,omitempty"`
		} `yaml:"reader"`
	} `yaml:"kafka"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
