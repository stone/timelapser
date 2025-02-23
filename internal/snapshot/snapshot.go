package snapshot

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/stone/timelapser/internal/camera"
	"github.com/stone/timelapser/internal/config"
	"github.com/stone/timelapser/internal/utils"
)

func TakeCameraSnapshot(camconfig *config.CameraConfig, outdir string, logger *slog.Logger) error {
	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	// camera := camera.NewCamera(*camconfig)
	camera := &camera.Camera{
		Config: *camconfig,
		Client: &http.Client{},
	}
	name := utils.ToCamelCase(camconfig.Name)
	logger.Debug("Retrieving snapshot", "name", camconfig.Name)

	cameraDir := filepath.Join(outdir, name)
	if err := os.MkdirAll(cameraDir, 0o755); err != nil {
		return fmt.Errorf("failed to create camera directory: %v", err)
	}

	snapshot, err := camera.GetSnapshot()
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

func TakeSnapshot(config *config.Config) error {
	logger := config.Logger
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for _, camConfig := range config.Cameras {
		camera := &camera.Camera{
			Config: camConfig,
			Client: &http.Client{},
		}
		name := utils.ToCamelCase(camConfig.Name)
		logger.Debug("Retrieving snapshot", "name", camConfig.Name)

		cameraDir := filepath.Join(config.OutputDir, name)
		if err := os.MkdirAll(cameraDir, 0o755); err != nil {
			return fmt.Errorf("failed to create camera directory: %v", err)
		}

		snapshot, err := camera.GetSnapshot()
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
