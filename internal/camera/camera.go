package camera

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/stone/timelapser/internal/config"
	"github.com/stone/timelapser/internal/utils"
)

// HTTPClient interface for mocking in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Camera struct {
	Config config.CameraConfig
	Client HTTPClient
}

func (c *Camera) GetSnapshot() ([]byte, error) {
	req, err := c.prepareRequest()
	if err != nil {
		return nil, fmt.Errorf("error preparing request: %v", err)
	}

	resp, err := c.Client.Do(req)
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
	baseURL := c.Config.SnapshotURL

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}

	if c.Config.Insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Apply authentication
	switch c.Config.Auth.Type {
	case "basic":
		auth := base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s",
				c.Config.Auth.Username,
				c.Config.Auth.Password)))
		req.Header.Add("Authorization", "Basic "+auth)

	case "bearer":
		req.Header.Add("Authorization", "Bearer "+c.Config.Auth.Token)
	}

	return req, nil
}

func ListCameras(config *config.Config) {
	for _, camConfig := range config.Cameras {
		name := utils.ToCamelCase(camConfig.Name)
		fmt.Printf("%s [%s]\n", name, camConfig.Name)
	}
}
