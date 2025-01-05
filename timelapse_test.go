package main

import (
	"reflect"
	"testing"
)

// Sample configuration for testing
var sampleConfig = &CameraConfig{
	FfmpegTemplate: "ffmpeg -f concat -safe 0 -i {{.ListPath}} -vf fps=24,format=yuv420p -c:v libx264 -preset medium -crf 23 -movflags +faststart -y {{.OutputPath}}",
}

func TestBuildFfmpegCommand(t *testing.T) {
	tests := []struct {
		name       string
		listPath   string
		outputPath string
		expected   []string
	}{
		{
			name:       "Basic command generation",
			listPath:   "/path/to/input.txt",
			outputPath: "/path/to/output.mp4",
			expected: []string{
				"ffmpeg", "-f", "concat", "-safe", "0", "-i", "/path/to/input.txt",
				"-vf", "fps=24,format=yuv420p", "-c:v", "libx264", "-preset", "medium",
				"-crf", "23", "-movflags", "+faststart", "-y", "/path/to/output.mp4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := buildFfmpegCommand(sampleConfig, tt.listPath, tt.outputPath)
			if err != nil {
				t.Fatalf("buildTimelapseCommand() error = %v", err)
			}

			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("buildTimelapseCommand() = %v, expected %v", args, tt.expected)
			}
		})
	}
}
