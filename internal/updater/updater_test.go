package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProviderFile(t *testing.T) {
	content := `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.10.0"
    }
  }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "provider.tf")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater()
	providers, _, err := u.ParseProviderFile(filePath)
	if err != nil {
		t.Fatalf("ParseProviderFile failed: %v", err)
	}

	if len(providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(providers))
	}

	if providers[0].Source != "hashicorp/aws" {
		t.Errorf("Expected source hashicorp/aws, got %s", providers[0].Source)
	}
	if providers[0].Version != "4.10.0" {
		t.Errorf("Expected version 4.10.0, got %s", providers[0].Version)
	}
}

func TestUpdateProject_MissingFile(t *testing.T) {
	u := NewUpdater()
	tmpDir := t.TempDir()
	
	_, err := u.UpdateProject(tmpDir)
	if err == nil {
		t.Error("Expected error for missing provider.tf, got nil")
	}
	
	expected := "provider.tf not found"
	if err != nil && err.Error() != expected && !contains(err.Error(), expected) {
		t.Errorf("Expected error containing %q, got %q", expected, err.Error())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
