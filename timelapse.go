package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// ErrNoSnapshots indicates that no snapshot images were found for processing
var ErrNoSnapshots = fmt.Errorf("no snapshots found for camera")

// buildFFmpegCommand creates the ffmpeg command string using the provided template
func buildFFmpegCommand(cfg *CameraConfig, listPath, outputPath string) (string, error) {
	tmpl, err := template.New("ffmpeg").Parse(cfg.FFmpegTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing ffmpeg template: %w", err)
	}

	var cmdBuffer bytes.Buffer
	data := map[string]string{
		"ListPath":   listPath,
		"OutputPath": outputPath,
	}

	if err := tmpl.Execute(&cmdBuffer, data); err != nil {
		return "", fmt.Errorf("executing ffmpeg template: %w", err)
	}

	return cmdBuffer.String(), nil
}

// CreateTimelapse generates a timelapse video from a sequence of images
func CreateTimelapse(cfg *CameraConfig, outputDir string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	name := toCamelCase(cfg.Name)
	folderPath := filepath.Join(outputDir, name)

	// Verify camera directory exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("camera directory not found: %w", err)
	}

	imageFiles, err := collectImageFiles(folderPath)
	if err != nil {
		return fmt.Errorf("collecting image files: %w", err)
	}

	if len(imageFiles) == 0 {
		return ErrNoSnapshots
	}

	timestamp := time.Now().Format("20060102-150405")
	listPath := filepath.Join(outputDir, fmt.Sprintf("%s-%s.txt", name, timestamp))
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s-%s.mp4", name, timestamp))

	if err := writeFileList(listPath, name, imageFiles, cfg.FrameDuration); err != nil {
		return fmt.Errorf("writing file list: %w", err)
	}
	defer cleanupFile(listPath)

	if err := executeFFmpeg(cfg, listPath, outputPath); err != nil {
		return err
	}

	logger.Info("timelapse created",
		"camera", cfg.Name,
		"output", outputPath,
		"snapshots", len(imageFiles),
	)

	if cfg.Delete {
		if err := cleanupImages(folderPath, imageFiles); err != nil {
			logger.Info("failed to cleanup some images",
				"camera", cfg.Name,
				"error", err,
			)
		}
	}

	return nil
}

func collectImageFiles(folderPath string) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	var imageFiles []string
	for _, entry := range entries {
		if isImageFile(entry) {
			imageFiles = append(imageFiles, entry.Name())
		}
	}

	// Sort files by modification time
	sort.Slice(imageFiles, func(i, j int) bool {
		iInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[i]))
		jInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[j]))
		return iInfo.ModTime().Before(jInfo.ModTime())
	})

	return imageFiles, nil
}

func isImageFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := strings.ToLower(entry.Name())
	return strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png")
}

func writeFileList(listPath, name string, imageFiles []string, frameDuration float64) error {
	var fileList strings.Builder
	for _, file := range imageFiles {
		fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(name, file)))
		fileList.WriteString(fmt.Sprintf("duration %f\n", frameDuration))
	}

	// Add last frame again to ensure visibility
	lastFrame := filepath.Join(name, imageFiles[len(imageFiles)-1])
	fileList.WriteString(fmt.Sprintf("file '%s'\n", lastFrame))

	return os.WriteFile(listPath, []byte(fileList.String()), 0o644)
}

func executeFFmpeg(cfg *CameraConfig, listPath, outputPath string) error {
	cmdStr, err := buildFFmpegCommand(cfg, listPath, outputPath)
	if err != nil {
		return fmt.Errorf("building ffmpeg command: %w", err)
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	logger.Debug("executing ffmpeg", "command", cmd.String())

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg execution failed: %w (%s)", err, output)
	}

	return nil
}

func cleanupFile(path string) {
	if err := os.Remove(path); err != nil {
		logger.Info("failed to remove temporary file",
			"path", path,
			"error", err,
		)
	}
}

func cleanupImages(folderPath string, imageFiles []string) error {
	var errs []string
	for _, file := range imageFiles {
		if err := os.Remove(filepath.Join(folderPath, file)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to remove some images: %s", strings.Join(errs, "; "))
	}
	return nil
}

func createAllTimelapse(config *Config) error {
	for _, camConfig := range config.Cameras {
		// we do not want to delete the original images when manually creating timelapse.
		camConfig.Delete = false
		if err := CreateTimelapse(&camConfig, config.OutputDir); err != nil {
			return err
		}
	}
	return nil
}
