package providers

import (
	"testing"
)

func TestGetLatestVersion_RealRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	// Using a very stable provider to ensure test reliability
	version, err := client.GetLatestVersion("hashicorp/aws")
	if err != nil {
		t.Fatalf("Failed to fetch real version from Terraform Registry: %v", err)
	}

	if version == "" {
		t.Error("Received empty version from real registry")
	}
}
