package generator

import (
	"strings"
	"testing"
)

func TestGenerateProviderFile_AWS(t *testing.T) {
	data := GeneratorData{
		ProjectName: "test-project",
		Providers: []ProviderConfig{
			{Name: "aws", Source: "hashicorp/aws", LatestVersion: "5.30.0"},
		},
	}

	got, err := GenerateProviderFile(data)
	if err != nil {
		t.Fatalf("GenerateProviderFile() error = %v", err)
	}

	content := string(got)
	
	expectedStrings := []string{
		`aws = {`,
		`source  = "hashicorp/aws"`,
		`version = "5.30.0"`,
		`provider "aws" {`,
		`region  = local.aws_region`,
		`aws_region  = var.aws_region`,
	}

	for _, s := range expectedStrings {
		if !strings.Contains(content, s) {
			t.Errorf("Expected content to contain %q", s)
		}
	}
}
