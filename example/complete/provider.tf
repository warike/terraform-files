terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.31.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "7.18.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.59.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.11.0"
    }
    vercel = {
      source  = "vercel/vercel"
      version = "4.5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.16.0"
    }
  }
}

provider "aws" {
  region  = local.aws_region
  profile = local.aws_profile
  default_tags {
   tags = local.tags
  }
}

provider "google" {
  project = local.gcp_project_id
  region  = local.gcp_region
}

provider "azurerm" {
  subscription_id = local.azure_subscription_id
  features {}
}

provider "github" {
  owner = local.gh_owner
  token = local.gh_token
}

provider "vercel" {
  api_token = local.vercel_api_token
}

provider "cloudflare" {
  api_token = local.cloudflare_api_token
}

locals {
  project_name = var.project_name
  aws_region  = var.aws_region
  aws_profile = var.aws_profile
  gcp_project_id = var.google_project_id
  gcp_region     = var.google_region
  azure_location            = var.azure_location
  azure_subscription_id     = var.azure_subscription_id
  gh_owner           = var.gh_owner
  gh_token           = var.gh_token
  vercel_api_token = var.vercel_api_token
  cloudflare_api_token = var.cloudflare_api_token

  tags = {
    project     = local.project_name
    environment = "dev"
    owner       = "warike"
    cost-center  = "development"
    terraform   = "true"
  }
}
