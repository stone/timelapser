package timelapse

import (
	"reflect"
	"testing"

	"github.com/stone/timelapser/internal/config"
)

// Sample configuration for testing
var sampleConfig = &config.CameraConfig{
	FFmpegTemplate: "-i {{.ListPath}} -y {{.OutputPath}}",
}

func TestBuildFFmpegCommand(t *testing.T) {
	tests := []struct {
		name       string
		listPath   string
		outputPath string
		expected   string
	}{
		{
			name:       "Basic command generation",
			listPath:   "/path/to/input.txt",
			outputPath: "/path/to/output.mp4",
			expected:   "-i /path/to/input.txt -y /path/to/output.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := buildFFmpegCommand(sampleConfig, tt.listPath, tt.outputPath)
			if err != nil {
				t.Fatalf("buildTimelapseCommand() error = %v", err)
			}

			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("buildTimelapseCommand() = %v, expected %v", args, tt.expected)
			}
		})
	}
}
