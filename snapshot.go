package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
)

func takeSnapshot(config *Config) error {
    if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %v", err)
    }

    for _, camConfig := range config.Cameras {
        camera := NewCamera(camConfig)
        name := toCamelCase(camConfig.Name)
        fmt.Printf("Retrieving snapshot for %s\n", camConfig.Name)

        cameraDir := filepath.Join(config.OutputDir, name)
        if err := os.MkdirAll(cameraDir, 0755); err != nil {
            return fmt.Errorf("failed to create camera directory: %v", err)
        }

        snapshot, err := camera.getSnapshot()
        if err != nil {
            fmt.Printf("Snapshot error for %s: %v\n", camConfig.Name, err)
            continue
        }

        filename := filepath.Join(cameraDir, fmt.Sprintf("%d.png", time.Now().UnixNano()))
        if err := os.WriteFile(filename, snapshot, 0644); err != nil {
            return fmt.Errorf("failed to write snapshot: %v", err)
        }

        fmt.Printf("Snapshot for %s saved\n", camConfig.Name)
    }

    return nil
}

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

