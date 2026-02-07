variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "my_project"
}

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
variable "google_project_id" {
  description = "Google Cloud project ID"
  type        = string
}

variable "google_region" {
  description = "Google Cloud region"
  type        = string
  default     = "us-central1"
}
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
variable "vercel_api_token" {
  description = "Vercel API Token"
  type        = string
  sensitive   = true
}
variable "cloudflare_api_token" {
  description = "Cloudflare API Token"
  type        = string
  sensitive   = true
}
