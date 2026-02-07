package main

import (
	"os"
	"path/filepath"
	"testing"
	"warike/base/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func TestE2E_GenerateFiles_InSubdirectory(t *testing.T) {
	// Setup a temporary CWD for the test to run in
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	testDir := t.TempDir()
	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	defer os.Chdir(originalWd) // Clean up by changing back

	// Define a relative subdirectory where the project should be created
	targetProjectDir := "my-cool-project"

	// Initialize the model, passing the RELATIVE path, just like the real app.
	m := ui.InitialModel(targetProjectDir)

	// BYPASS LOADING STATE for a synchronous test
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

	// 3. Verify Model State
	if !m.FilesGenerated {
		t.Error("Expected filesGenerated to be true after pressing 'g'")
	}

	// 4. Verify Files on Disk
	expectedFiles := []string{"provider.tf", "variables.tf", "main.tf", "terraform.tfvars"}
	for _, f := range expectedFiles {
		// Check that the file exists inside the RELATIVE subdirectory
		expectedPath := filepath.Join(targetProjectDir, f)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("CRITICAL: Expected file '%s' was NOT generated in the target directory", expectedPath)
		}
	}

	// 5. Verify a file was NOT created in the base directory (CWD of the test)
	unexpectedPath := "provider.tf"
	if _, err := os.Stat(unexpectedPath); err == nil {
		t.Errorf("CRITICAL: A provider.tf file was incorrectly generated in the base directory")
	}
}
