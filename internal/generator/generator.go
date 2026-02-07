package generator

import (
	"bytes"
	"os"
	"text/template"
)

type ProviderConfig struct {
	Name          string
	Source        string
	LatestVersion string
}

type GeneratorData struct {
	ProjectName string
	Providers   []ProviderConfig
}

const providerTemplate = `terraform {
  required_providers {
{{- range .Providers }}
    {{ .Name }} = {
      source  = "{{ .Source }}"
      version = "{{ .LatestVersion }}"
    }
{{- end }}
  }
}

{{ range .Providers -}}
provider "{{ .Name }}" {
{{- if eq .Name "aws" }}
  region  = local.aws_region
  profile = local.aws_profile
  default_tags {
   tags = local.tags
  }
{{- else if eq .Name "google" }}
  project = local.gcp_project_id
  region  = local.gcp_region
{{- else if eq .Name "azurerm" }}
  subscription_id = local.azure_subscription_id
  features {}
{{- else if eq .Name "github" }}
  owner = local.gh_owner
  token = local.gh_token
{{- else if eq .Name "vercel" }}
  api_token = local.vercel_api_token
{{- else if eq .Name "cloudflare" }}
  api_token = local.cloudflare_api_token
{{- end }}
}

{{ end -}}

locals {
  project_name = var.project_name
{{- range .Providers }}
{{- if eq .Name "aws" }}
  aws_region  = var.aws_region
  aws_profile = var.aws_profile
{{- else if eq .Name "google" }}
  gcp_project_id = var.google_project_id
  gcp_region     = var.google_region
{{- else if eq .Name "azurerm" }}
  azure_location            = var.azure_location
  azure_subscription_id     = var.azure_subscription_id
{{- else if eq .Name "github" }}
  gh_owner           = var.gh_owner
  gh_token           = var.gh_token
{{- else if eq .Name "vercel" }}
  vercel_api_token = var.vercel_api_token
{{- else if eq .Name "cloudflare" }}
  cloudflare_api_token = var.cloudflare_api_token
{{- end }}
{{- end }}

  tags = {
    project     = local.project_name
    environment = "dev"
    owner       = "warike"
    cost-center  = "development"
    terraform   = "true"
  }
}
`

const variablesTemplate = `variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "my_project"
}
{{ range .Providers }}
{{- if eq .Name "aws" }}
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "aws_profile" {
  description = "AWS profile name"
  type        = string
  default     = "default"
}
{{- else if eq .Name "google" }}
variable "google_project_id" {
  description = "Google Cloud project ID"
  type        = string
}

variable "google_region" {
  description = "Google Cloud region"
  type        = string
  default     = "us-central1"
}
{{- else if eq .Name "azurerm" }}
variable "azure_location" {
  description = "Azure location"
  type        = string
  default     = "East US"
}

variable "azure_subscription_id" {
  description = "Azure subscription ID"
  type        = string
  sensitive   = true
}
{{- else if eq .Name "github" }}
variable "gh_owner" {
  description = "GitHub owner (user or organization)"
  type        = string
  default     = "warike"
}

variable "gh_token" {
  description = "GitHub token"
  type        = string
  sensitive   = true
}
{{- else if eq .Name "vercel" }}
variable "vercel_api_token" {
  description = "Vercel API Token"
  type        = string
  sensitive   = true
}
{{- else if eq .Name "cloudflare" }}
variable "cloudflare_api_token" {
  description = "Cloudflare API Token"
  type        = string
  sensitive   = true
}
{{- end }}
{{- end }}
`

const tfvarsTemplate = `project_name = "{{ .ProjectName }}"
{{ range .Providers }}
{{- if eq .Name "aws" }}
aws_region   = "us-west-2"
aws_profile  = "default"
{{- else if eq .Name "google" }}
google_project_id = "gcp-project-id-goes-here"
google_region     = "us-central1"
{{- else if eq .Name "azurerm" }}
azure_location = "East US"
azure_subscription_id = "azure-subscription-id-goes-here"
{{- else if eq .Name "github" }}
gh_owner = "warike"
gh_token = "your-github-token"
{{- else if eq .Name "vercel" }}
vercel_api_token = "your-vercel-token"
{{- else if eq .Name "cloudflare" }}
cloudflare_api_token = "your-cloudflare-token"
{{- end }}
{{- end }}
`

const mainTemplate = `// main.tf
`

func GenerateProviderFile(data GeneratorData) ([]byte, error) {
	return generateFromTemplate("provider", providerTemplate, data)
}

func GenerateVariablesFile(data GeneratorData) ([]byte, error) {
	return generateFromTemplate("variables", variablesTemplate, data)
}

func GenerateTfvarsFile(data GeneratorData) ([]byte, error) {
	return generateFromTemplate("tfvars", tfvarsTemplate, data)
}

func GenerateMainFile(data GeneratorData) ([]byte, error) {
	return generateFromTemplate("main", mainTemplate, data)
}

func generateFromTemplate(name, text string, data GeneratorData) ([]byte, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func WriteFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}
