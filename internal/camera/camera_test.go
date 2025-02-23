package camera

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stone/timelapser/internal/config"
)

// mockHTTPClient implements HTTPClient interface
type mockHTTPClient struct {
	doFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestCamera_getSnapshot(t *testing.T) {
	tests := []struct {
		name          string
		config        config.CameraConfig
		mockResponse  *http.Response
		mockErr       error
		expectedData  []byte
		expectedError string
	}{
		{
			name: "successful snapshot",
			config: config.CameraConfig{
				SnapshotURL: "http://example.com/snapshot",
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("image data"))),
			},
			expectedData: []byte("image data"),
		},
		{
			name: "http error",
			config: config.CameraConfig{
				SnapshotURL: "http://example.com/snapshot",
			},
			mockErr:       http.ErrAbortHandler,
			expectedError: "error making request",
		},
		{
			name: "non-200 status code",
			config: config.CameraConfig{
				SnapshotURL: "http://example.com/snapshot",
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			expectedError: "unexpected status code: 401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camera := &Camera{
				Config: tt.config,
				Client: &mockHTTPClient{
					doFunc: func(*http.Request) (*http.Response, error) {
						return tt.mockResponse, tt.mockErr
					},
				},
			}

			data, err := camera.GetSnapshot()

			// Check error cases
			if tt.expectedError != "" {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !bytes.Equal(data, tt.expectedData) {
				t.Errorf("data = %v, want %v", data, tt.expectedData)
			}
		})
	}
}
