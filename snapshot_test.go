package main

import (
	"testing"
)

func Test_toCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "two words",
			input:    "hello world",
			expected: "helloWorld",
		},
		{
			name:     "multiple words",
			input:    "the quick brown fox",
			expected: "theQuickBrownFox",
		},
		{
			name:     "multiple spaces",
			input:    "hello    world",
			expected: "helloWorld",
		},
		{
			name:     "mixed case input",
			input:    "Hello World",
			expected: "helloWorld",
		},
		{
			name:     "all caps input",
			input:    "HELLO WORLD",
			expected: "helloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toCamelCase(tt.input)
			if got != tt.expected {
				t.Errorf("toCamelCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}
