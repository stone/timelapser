package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func takeCameraSnapshot(camconfig *CameraConfig, outdir string) error {
	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	camera := NewCamera(*camconfig)
	name := toCamelCase(camconfig.Name)
	logger.Info("Retrieving snapshot", "name", camconfig.Name)

	cameraDir := filepath.Join(outdir, name)
	if err := os.MkdirAll(cameraDir, 0o755); err != nil {
		return fmt.Errorf("failed to create camera directory: %v", err)
	}

	snapshot, err := camera.getSnapshot()
	if err != nil {
		return fmt.Errorf("snapshot error for: %s error: %s", camconfig.Name, err)
	}

	filename := filepath.Join(cameraDir, fmt.Sprintf("%d.png", time.Now().UnixNano()))
	if err := os.WriteFile(filename, snapshot, 0o644); err != nil {
		return fmt.Errorf("failed to write snapshot: %v", err)
	}

	logger.Info("Snapshot saved for", "name", camconfig.Name, "file", filename)

	return nil
}

func takeSnapshot(config *Config) error {
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for _, camConfig := range config.Cameras {
		camera := NewCamera(camConfig)
		name := toCamelCase(camConfig.Name)
		logger.Info("Retrieving snapshot", "name", camConfig.Name)

		cameraDir := filepath.Join(config.OutputDir, name)
		if err := os.MkdirAll(cameraDir, 0o755); err != nil {
			return fmt.Errorf("failed to create camera directory: %v", err)
		}

		snapshot, err := camera.getSnapshot()
		if err != nil {
			logger.Error("Snapshot error for", "name", camConfig.Name, "err", err, "continue", true)
			continue
		}

		filename := filepath.Join(cameraDir, fmt.Sprintf("%d.png", time.Now().UnixNano()))
		if err := os.WriteFile(filename, snapshot, 0o644); err != nil {
			return fmt.Errorf("failed to write snapshot: %v", err)
		}

		logger.Info("Snapshot saved for", "name", camConfig.Name, "file", filename)
	}

	return nil
}

// Convert names to "safer" paths using camelCase
// It could be a bit confusing, but it is a simple way to handle
// spaces in paths.
// "This is a Test" -> "thisIsATest"
func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(strings.ToLower(words[0]))

	for _, word := range words[1:] {
		result.WriteString(strings.Title(word))
	}

	return result.String()
}
