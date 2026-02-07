Feature: CLI Commands for Create and Update
  As a user
  I want to separate project creation and updating into distinct commands
  So that I can manage my Terraform projects more effectively

  @e2e
  Scenario: Create a new project in the current directory
    Given the CLI is installed
    And I am in an empty directory
    When I run "tfinit create ."
    And I select "aws" provider
    And I confirm generation
    Then the file "provider.tf" should exist
    And the file "provider.tf" should contain "provider "aws""

  @e2e
  Scenario: Create a new project in a specific directory
    Given the CLI is installed
    When I run "tfinit create my-infra"
    And I select "google" provider
    And I confirm generation
    Then the directory "my-infra" should exist
    And the file "my-infra/provider.tf" should exist
    And the file "my-infra/provider.tf" should contain "provider "google""

  @e2e
  Scenario: Update existing project providers
    Given I have a Terraform project in "legacy-infra"
    And the "provider.tf" contains "hashicorp/aws" version "4.0.0"
    And the latest version of "hashicorp/aws" is "5.0.0"
    When I run "tfinit update legacy-infra"
    Then the file "legacy-infra/provider.tf" should show version "5.0.0"
    And the console should display "Updated aws from 4.0.0 to 5.0.0"

  @unit
  Scenario: Update command identifies required providers
    Given a "provider.tf" content:
      """
      terraform {
        required_providers {
          aws = {
            source  = "hashicorp/aws"
            version = "4.10.0"
          }
        }
      }
      """
    When I parse the provider file
    Then I should identify "hashicorp/aws" with version "4.10.0"

  @integration
  Scenario: Update command fetches latest version for parsed provider
    Given a parsed provider "hashicorp/aws" with version "4.0.0"
    When I check for updates
    Then I should receive a version greater than "4.0.0"

  @unit
  Scenario: Create command fails if directory already exists and is not empty
    Given a non-empty directory "existing-project" exists
    When I run "tfinit create existing-project"
    Then it should return an error "directory exists and is not empty"

  @unit
  Scenario: Update command fails if provider.tf is missing
    Given the directory "empty-dir" exists
    And "empty-dir" does not contain "provider.tf"
    When I run "tfinit update empty-dir"
    Then it should return an error "provider.tf not found"
