package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const DefaultRegistryURL = "https://registry.terraform.io/v1/providers"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: DefaultRegistryURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetLatestVersion(source string) (string, error) {
	// Construct URL. If BaseURL is the default, we append the source.
	// If it's a test server (implied by not matching default), we might just hit the root 
	// or append strictly if the test server expects it.
	// For simplicity in this refactor, we just concatenate.
	
	// Handle trailing slash in BaseURL just in case
	url := fmt.Sprintf("%s/%s", c.BaseURL, source)
	// If fetching from a test server that doesn't handle paths, this might be tricky, 
	// but httptest.Server handles paths fine. 
	// In the test: BaseURL = server.URL. Code does: server.URL + "/" + source.
	// The test handler needs to match anything or we don't care about the path in the test handler.

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	var result struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Version, nil
}
