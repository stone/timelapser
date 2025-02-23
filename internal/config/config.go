package config

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultOutputDir         = "/tmp"
	defaultInterval          = "*/5 * * * *"
	defaultTimelapseInterval = "*/60 * * * *"
	defaultFrameDuration     = 0.0416667
	defaultFFmpegTemplate    = "ffmpeg -f concat -safe 0 -i {{.ListPath}} -vf fps=24,format=yuv420p -c:v libx264 -preset medium -crf 23 -movflags +faststart -y {{.OutputPath}}"
)

type AuthConfig struct {
	Type     string `yaml:"type,omitempty"`     // basic, bearer
	Username string `yaml:"username,omitempty"` // for basic auth
	Password string `yaml:"password,omitempty"` // for basic auth
	Token    string `yaml:"token,omitempty"`    // for bearer auth
}

type CameraConfig struct {
	Name              string     `yaml:"name"`
	SnapshotURL       string     `yaml:"snapshotUrl"`
	Insecure          bool       `yaml:"insecure"`
	Auth              AuthConfig `yaml:"auth,omitempty"`
	Delete            bool       `yaml:"delete"`
	Interval          string     `yaml:"interval,omitempty"`
	TimelapseInterval string     `yaml:"timelapseInterval,omitempty"`
	FrameDuration     float64    `yaml:"frameDuration,omitempty"`
	FFmpegTemplate    string     `yaml:"ffmpeg_template,omitempty"`
}

type Config struct {
	OutputDir         string         `yaml:"outputDir"`
	Cameras           []CameraConfig `yaml:"cameras"`
	Interval          string         `yaml:"interval"`
	TimelapseInterval string         `yaml:"timelapseInterval"`
	FrameDuration     float64        `yaml:"frameDuration"`
	FFmpegTemplate    string         `yaml:"ffmpeg_template"`
	Logger            *slog.Logger
}

func newDefaultConfig() Config {
	// Create a new Config struct with default values
	return Config{
		OutputDir:         defaultOutputDir,
		Interval:          defaultInterval,
		TimelapseInterval: defaultTimelapseInterval,
		FrameDuration:     defaultFrameDuration,
		FFmpegTemplate:    defaultFFmpegTemplate,
	}
}

func NewExampleConfig() string {
	config := newDefaultConfig()
	// Create example CameraConfig struct with default values
	cameraConfig := CameraConfig{
		Name:        "camera1",
		SnapshotURL: "http://localhost:8080/snapshot",
		Auth:        AuthConfig{Type: "basic", Username: "admin", Password: "admin"},
		Delete:      false,
	}
	config.Cameras = append(config.Cameras, cameraConfig)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Sprintf("error creating example configuration: %s", err.Error())
	}

	return string(data)
}

func LoadConfig(path string, logger *slog.Logger) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Set defaults
	config := newDefaultConfig()
	// Attach logger to our config, not amazing but sometimes you get to be lazy
	config.Logger = logger

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	if err := applyDefaultsToCameras(&config); err != nil {
		return nil, fmt.Errorf("error applying defaults to cameras: %v", err)
	}

	return &config, nil
}

func applyDefaultsToCameras(config *Config) error {
	for i := range config.Cameras {
		camConfig := &config.Cameras[i]
		if camConfig.TimelapseInterval == "" {
			config.Logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.TimelapseInterval,
				"to", config.TimelapseInterval)
			camConfig.TimelapseInterval = config.TimelapseInterval
		}

		if camConfig.FrameDuration == 0 {
			config.Logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.FrameDuration,
				"to", config.FrameDuration)
			camConfig.FrameDuration = config.FrameDuration
		}

		if camConfig.FFmpegTemplate == "" {
			config.Logger.Debug("Setting defaults for camera", "name", camConfig.Name,
				"from", camConfig.FFmpegTemplate,
				"to", config.FFmpegTemplate)
			camConfig.FFmpegTemplate = config.FFmpegTemplate
		}

	}
	return nil
}
