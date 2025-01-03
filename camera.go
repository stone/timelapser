package main

import (
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

type Camera struct {
    config CameraConfig
    client *http.Client
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

    // Handle query parameters
    if c.config.Auth.Type == "query" {
        parsedURL, err := url.Parse(baseURL)
        if err != nil {
            return nil, err
        }
        q := parsedURL.Query()
        for key, value := range c.config.Auth.Params {
            q.Add(key, value)
        }
        parsedURL.RawQuery = q.Encode()
        baseURL = parsedURL.String()
    }

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

