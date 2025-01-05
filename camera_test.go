package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPClient implements HTTPClient interface
type mockHTTPClient struct {
    doFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.doFunc(req)
}

func TestNewCamera(t *testing.T) {
    config := CameraConfig{
        SnapshotURL: "http://example.com/snapshot",
        Auth: AuthConfig{
            Type:     "basic",
            Username: "user",
            Password: "pass",
        },
    }

    camera := NewCamera(config)
    if camera == nil {
        t.Error("NewCamera returned nil")
    }
    if camera.config != config {
        t.Errorf("NewCamera config = %v, want %v", camera.config, config)
    }
    if camera.client == nil {
        t.Error("NewCamera client is nil")
    }
}

func TestCamera_getSnapshot(t *testing.T) {
    tests := []struct {
        name           string
        config        CameraConfig
        mockResponse  *http.Response
        mockErr       error
        expectedData  []byte
        expectedError string
    }{
        {
            name: "successful snapshot",
            config: CameraConfig{
                SnapshotURL: "http://example.com/snapshot",
            },
            mockResponse: &http.Response{
                StatusCode: http.StatusOK,
                Body:      io.NopCloser(bytes.NewReader([]byte("image data"))),
            },
            expectedData: []byte("image data"),
        },
        {
            name: "http error",
            config: CameraConfig{
                SnapshotURL: "http://example.com/snapshot",
            },
            mockErr:       http.ErrAbortHandler,
            expectedError: "error making request",
        },
        {
            name: "non-200 status code",
            config: CameraConfig{
                SnapshotURL: "http://example.com/snapshot",
            },
            mockResponse: &http.Response{
                StatusCode: http.StatusUnauthorized,
                Body:      io.NopCloser(strings.NewReader("")),
            },
            expectedError: "unexpected status code: 401",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            camera := &Camera{
                config: tt.config,
                client: &mockHTTPClient{
                    doFunc: func(*http.Request) (*http.Response, error) {
                        return tt.mockResponse, tt.mockErr
                    },
                },
            }

            data, err := camera.getSnapshot()

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
