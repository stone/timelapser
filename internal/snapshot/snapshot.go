package snapshot

import (
	"errors"
	"fmt"
	"log/slog"
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
	cam := camera.NewCamera(*camconfig)
	name := utils.ToCamelCase(camconfig.Name)
	logger.Debug("Retrieving snapshot", "name", camconfig.Name)

	cameraDir := filepath.Join(outdir, name)
	if err := os.MkdirAll(cameraDir, 0o755); err != nil {
		return fmt.Errorf("failed to create camera directory: %v", err)
	}

	snapshot, err := cam.GetSnapshot()
	if err != nil {
		return fmt.Errorf("snapshot error for: %s error: %s", camconfig.Name, err)
	}

	// Write atomically: temp file in the same directory → rename.
	// This prevents the timelapse job from reading a partially-written file.
	filename := filepath.Join(cameraDir, fmt.Sprintf("%d.png", time.Now().UnixNano()))
	tmp := filename + ".tmp"
	if err := os.WriteFile(tmp, snapshot, 0o644); err != nil {
		return fmt.Errorf("failed to write snapshot: %v", err)
	}
	if err := os.Rename(tmp, filename); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("failed to rename snapshot: %v", err)
	}

	logger.Info("Snapshot saved for", "name", camconfig.Name, "file", filename)

	return nil
}

func TakeSnapshot(config *config.Config) error {
	logger := config.Logger
	var errs []error
	for _, camConfig := range config.Cameras {
		if err := TakeCameraSnapshot(&camConfig, config.OutputDir, logger); err != nil {
			logger.Error("Snapshot error for", "name", camConfig.Name, "err", err, "continue", true)
			errs = append(errs, fmt.Errorf("%s: %w", camConfig.Name, err))
		}
	}
	return errors.Join(errs...)
}
