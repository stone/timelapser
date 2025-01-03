package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func createTimelapse(config *Config) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %v", err)
	}

	for _, camConfig := range config.Cameras {
		name := toCamelCase(camConfig.Name)
		folderPath := filepath.Join(config.OutputDir, name)

		// Skip if camera directory doesn't exist
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			fmt.Printf("No images found for camera %s, skipping\n", camConfig.Name)
			continue
		}

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			fmt.Printf("Error reading directory for camera %s: %v\n", camConfig.Name, err)
			continue
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
			fmt.Printf("No images found for camera %s\n", camConfig.Name)
			continue
		}

		// Sort files by modification time
		sort.Slice(imageFiles, func(i, j int) bool {
			iInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[i]))
			jInfo, _ := os.Stat(filepath.Join(folderPath, imageFiles[j]))
			return iInfo.ModTime().Before(jInfo.ModTime())
		})

		// Create FFmpeg input file list
		timestamp := time.Now().Format("20060102-150405")
		listPath := filepath.Join(config.OutputDir, fmt.Sprintf("%s-%s.txt", name, timestamp))
		outputPath := filepath.Join(config.OutputDir, fmt.Sprintf("%s-%s.mp4", name, timestamp))

		var fileList strings.Builder
		for _, file := range imageFiles {
			// fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(folderPath, file)))
			fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(name, file)))
			// fileList.WriteString(fmt.Sprintf("file '%s'\n", file))
			fileList.WriteString(fmt.Sprintf("duration 0.0416667\n")) // 1/24 for 24fps
		}
		// Add last frame one more time to ensure last image is visible
		fileList.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(folderPath, imageFiles[len(imageFiles)-1])))

		if err := os.WriteFile(listPath, []byte(fileList.String()), 0o644); err != nil {
			return fmt.Errorf("failed to write file list for %s: %v", camConfig.Name, err)
		}

		fmt.Printf("Creating timelapse for %s (%d images)\n", camConfig.Name, len(imageFiles))

		// Run FFmpeg command
		cmd := exec.Command("ffmpeg",
			"-f", "concat",
			"-safe", "0",
			"-i", listPath,
			"-vf", "fps=24,format=yuv420p", // Ensure compatibility
			"-c:v", "libx264",
			"-preset", "medium",
			"-crf", "23",
			"-movflags", "+faststart",
			"-y",
			outputPath,
		)

		// fmt.Println(cmd)

		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("FFmpeg error for %s: %v\nOutput: %s\n", camConfig.Name, err, output)
			continue
		}

		// Cleanup
		if err := os.Remove(listPath); err != nil {
			fmt.Printf("Warning: Failed to remove temporary file list for %s: %v\n", camConfig.Name, err)
		}

		fmt.Printf("Timelapse created for %s at %s\n", camConfig.Name, outputPath)
		// Optionally remove original images
		if camConfig.Delete {
			for _, file := range imageFiles {
				if err := os.Remove(filepath.Join(folderPath, file)); err != nil {
					fmt.Printf("Warning: Failed to remove %s: %v\n", file, err)
				}
			}
			fmt.Printf("Original images removed for %s\n", camConfig.Name)
		}
	}

	return nil
}
