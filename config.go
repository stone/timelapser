package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AuthConfig struct {
	Type     string `yaml:"type"`     // basic, bearer
	Username string `yaml:"username"` // for basic auth
	Password string `yaml:"password"` // for basic auth
	Token    string `yaml:"token"`    // for bearer auth
}

type CameraConfig struct {
	Name              string     `yaml:"name"`
	SnapshotURL       string     `yaml:"snapshotUrl"`
	Auth              AuthConfig `yaml:"auth"`
	Delete            bool       `yaml:"delete"`
	Interval          string     `yaml:"interval"`
	TimelapseInterval string     `yaml:"timelapseInterval"`
	FrameDuration     float64    `yaml:"frameDuration"`
	FfmpegTemplate    string     `yaml:"ffmpeg_template"`
}

type Config struct {
	OutputDir         string         `yaml:"outputDir"`
	Cameras           []CameraConfig `yaml:"cameras"`
	Interval          string         `yaml:"interval"`
	TimelapseInterval string         `yaml:"timelapseInterval"`
	FrameDuration     float64        `yaml:"frameDuration"`
	FfmpegTemplate    string         `yaml:"ffmpeg_template"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Set defaults
	config := Config{
		OutputDir:         "/tmp",
		Interval:          "*/5 * * * *",
		TimelapseInterval: "*/60 * * * *",
		FrameDuration:     0.0416667,
		FfmpegTemplate:    "ffmpeg -f concat -safe 0 -i {{.ListPath}} -vf fps=24,format=yuv420p -c:v libx264 -preset medium -crf 23 -movflags +faststart -y {{.OutputPath}}",
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Add defaults if not set for Cameras
	for _, camConfig := range config.Cameras {
		if camConfig.TimelapseInterval == "" {
			logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.TimelapseInterval,
				"to", config.TimelapseInterval)
			camConfig.TimelapseInterval = config.TimelapseInterval
		}

		if camConfig.FrameDuration == 0 {
			logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.FrameDuration,
				"to", config.FrameDuration)
			camConfig.FrameDuration = config.FrameDuration
		}

		if camConfig.FfmpegTemplate == "" {
			logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.FfmpegTemplate,
				"to", config.FfmpegTemplate)
			camConfig.FfmpegTemplate = config.FfmpegTemplate
		}
	}

	return &config, nil
}
