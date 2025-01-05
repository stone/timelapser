package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient interface for mocking in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Camera struct {
	config CameraConfig
	client HTTPClient
}

func NewCamera(config CameraConfig) *Camera {
	return &Camera{
		config: config,
		client: &http.Client{},
	}
}

func (c *Camera) getSnapshot() ([]byte, error) {
	req, err := c.prepareRequest()
	if err != nil {
		return nil, fmt.Errorf("error preparing request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Camera) prepareRequest() (*http.Request, error) {
	baseURL := c.config.SnapshotURL

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}

	// Apply authentication
	switch c.config.Auth.Type {
	case "basic":
		auth := base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s",
				c.config.Auth.Username,
				c.config.Auth.Password)))
		req.Header.Add("Authorization", "Basic "+auth)

	case "bearer":
		req.Header.Add("Authorization", "Bearer "+c.config.Auth.Token)
	}

	return req, nil
}

func listCameras(config *Config) {
	for _, camConfig := range config.Cameras {
		name := toCamelCase(camConfig.Name)
		fmt.Printf("%s [%s]\n", name, camConfig.Name)
	}
}
