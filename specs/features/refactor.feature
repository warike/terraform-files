Feature: Modular CLI Architecture
  As a developer
  I want to refactor the monolithic main.go into a modular architecture
  So that the code is testable, maintainable, and scalable

  @e2e
  Scenario: Generate Terraform files for selected providers
    Given the CLI is installed
    And the user runs the "wtf" command
    When the user selects "aws" and "github" providers
    And the user confirms generation
    Then the file "provider.tf" should exist
    And the file "variables.tf" should exist
    And the file "main.tf" should exist
    And the file "terraform.tfvars" should exist
    And "provider.tf" should contain "provider "aws""
    And "provider.tf" should contain "provider "github""

  @integration
  Scenario: Provider Client fetches latest versions
    Given the Terraform Registry is accessible
    When I fetch the latest version for "hashicorp/aws"
    Then I should receive a valid semantic version string (e.g., "5.0.0")
    And it should not be empty

  @unit
  Scenario: Provider Client handles network errors gracefully
    Given the HTTP client returns a 500 error
    When I fetch the latest version for "hashicorp/aws"
    Then I should receive an error matching "bad status: 500"

  @unit
  Scenario: File Generator creates correct AWS provider block
    Given I have a provider configuration for "aws" with version "5.30.0"
    When I generate the "provider.tf" content
    Then it should contain:
      """
      provider "aws" {
        region  = local.aws_region
        profile = local.aws_profile
      """
