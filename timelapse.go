package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func buildFfmpegCommand(camConfig *CameraConfig, listPath, outputPath string) (string, error) {
	logger.Debug("buildFfmpegCommand", "camConfig", camConfig.FfmpegTemplate, "list", listPath, "out", outputPath)
	tmpl, err := template.New("ffmpeg").Parse(camConfig.FfmpegTemplate)
	if err != nil {
		return "", err
	}

	var cmdBuffer bytes.Buffer
	data := map[string]string{
		"ListPath":   listPath,
		"OutputPath": outputPath,
	}
	if err := tmpl.Execute(&cmdBuffer, data); err != nil {
		return "", err
	}

	return cmdBuffer.String(), nil
}

func createTimelapse(camConfig *CameraConfig, outputdir string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	name := toCamelCase(camConfig.Name)
	folderPath := filepath.Join(outputdir, name)

	// Skip if camera directory doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("no snapshots found for camera: %s", err)
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory for camera: %s", err)
	}

	// Filter and sort image files
	var imageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(strings.ToLower(entry.Name()), ".jpg") ||
			strings.HasSuffix(strings.ToLower(entry.Name()), ".png")) {
			imageFiles = append(imageFiles, entry.Name())
		}
	}

	if len(imageFiles) == 0 {
		return fmt.Errorf("no snapshots found for camera")
	}

	// Sort files by modification time
	sort.Slice(imageFiles, func(i, j int) bool {
		iInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[i]))
		jInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[j]))
		return iInfo.ModTime().Before(jInfo.ModTime())
	})

	// Create FFmpeg input file list
	timestamp := time.Now().Format("20060102-150405")
	listPath := filepath.Join(outputdir, fmt.Sprintf("%s-%s.txt", name, timestamp))
	outputPath := filepath.Join(outputdir, fmt.Sprintf("%s-%s.mp4", name, timestamp))

	var fileList strings.Builder
	for _, file := range imageFiles {
		fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(name, file)))
		fileList.WriteString(fmt.Sprintf("duration %f\n", camConfig.FrameDuration)) // 1/24 for 24fps
	}
	// Add last frame one more time to ensure last image is visible
	fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(folderPath, imageFiles[len(imageFiles)-1])))

	if err := os.WriteFile(listPath, []byte(fileList.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write file list for %s: %v", camConfig.Name, err)
	}

	logger.Info("Creating timelapse for camera", "name", camConfig.Name, "snapshots", len(imageFiles))

	cmdBuf, err := buildFfmpegCommand(camConfig, listPath, outputPath)
	if err != nil {
		return fmt.Errorf("error building ffmpeg command: %s", err)
	}
	cmd := exec.Command("sh", "-c", cmdBuf)

	logger.Debug("ffmpeg command", "exec", cmd.String())

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error creating ffmpeg timelapse: %s (%s)", err, string(output))
	}

	// Cleanup
	if err := os.Remove(listPath); err != nil {
		logger.Info("Failed to remove temporary timelampse file list", "name", camConfig.Name, "err", err)
	}

	logger.Info("Timelampse created for camera", "name", camConfig.Name, "output", outputPath)

	// Optionally remove original images
	if camConfig.Delete {
		logger.Info("Removing snapshots", "name", camConfig.Name)
		for _, file := range imageFiles {
			if err := os.Remove(filepath.Join(folderPath, file)); err != nil {
				logger.Info("Failed to remove snapshot for camera", "name", camConfig.Name, "err", err)
			}
		}
		logger.Info("Snapshot images removed for camera", "name", camConfig.Name)
	}

	return nil
}

func createAllTimelapse(config *Config) error {
	for _, camConfig := range config.Cameras {
		if err := createTimelapse(&camConfig, config.OutputDir); err != nil {
			return err
		}
	}
	return nil
}
