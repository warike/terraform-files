package main

import (
	"os"
	"path/filepath"
	"testing"
	"warike/base/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func TestE2E_GenerateFiles_InNestedSubdirectory(t *testing.T) {
	// Setup a temporary CWD for the test to run in
	originalWd, _ := os.Getwd()
	testDir := t.TempDir()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Define a nested subdirectory path
	targetProjectDir := filepath.Join("example", "complete")

	// Initialize the model, passing the RELATIVE nested path
	m := ui.InitialModel(targetProjectDir)

	// BYPASS LOADING STATE
	m.Loading = false
	m.Providers = []ui.Provider{
		{Name: "aws", Source: "hashicorp/aws", LatestVersion: "5.0.0"},
	}

	// 1. Simulate selecting AWS
	msg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ := m.Update(msg)
	m = newModel.(ui.Model)

	// 2. Simulate pressing 'g' to generate the files
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	newModel, _ = m.Update(msg)
	m = newModel.(ui.Model)

	// 3. Verify that the nested directory was created
	if _, err := os.Stat(targetProjectDir); os.IsNotExist(err) {
		t.Fatalf("CRITICAL: The nested directory '%s' was not created", targetProjectDir)
	}

	// 4. Verify Files on Disk inside the NESTED directory
	expectedFiles := []string{"provider.tf"}
	for _, f := range expectedFiles {
		expectedPath := filepath.Join(targetProjectDir, f)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("CRITICAL: Expected file '%s' was NOT generated in the nested target directory", expectedPath)
		}
	}
}
