package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AuthConfig struct {
	Type     string            `yaml:"type"`     // basic, bearer, or query
	Username string            `yaml:"username"` // for basic auth
	Password string            `yaml:"password"` // for basic auth
	Token    string            `yaml:"token"`    // for bearer auth
	Params   map[string]string `yaml:"params"`   // for query parameters
}

type CameraConfig struct {
	Name        string     `yaml:"name"`
	SnapshotURL string     `yaml:"snapshotUrl"`
	Auth        AuthConfig `yaml:"auth"`
	Delete      bool       `yaml:"delete"`
}

type Config struct {
	Interval  int            `yaml:"interval"` // snapshot interval in seconds
	OutputDir string         `yaml:"outputDir"`
	Cameras   []CameraConfig `yaml:"cameras"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}
