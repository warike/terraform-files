# Terraform Quickstart Generator

A command-line interface (CLI) tool to quickly scaffold your Terraform projects by interactively selecting common cloud and service providers. This utility fetches the latest provider versions and generates essential Terraform configuration files, giving you a head start on your infrastructure-as-code journey.

## Features

* **Interactive Provider Selection:** Choose from a list of popular Terraform providers (AWS, Google Cloud, Azure, GitHub, Vercel, Cloudflare).

* **Latest Version Fetching:** Automatically retrieves the most recent stable version for each selected provider from the Terraform Registry.

* **Automated File Generation:** Creates the following foundational Terraform files:

  * `provider.tf`: Defines required providers and their configurations.

  * `variables.tf`: Declares input variables for your project.

  * `main.tf`: A placeholder for your primary Terraform resources.

  * `terraform.tfvars`: An example file to set default values for your variables, including sensitive ones.

## Getting Started

### Prerequisites

* [Go](https://golang.org/doc/install) (version 1.16 or higher recommended)

### How to Use

1. **Clone the repository (or save the code):**

`git clone [https://github.com/warike/terraform-files.git](https://github.com/warike/terraform-files.git)`

*(If you just have the `.go` file, navigate to its directory.)*

2. **Run the application:**

`go run .`



3. **Interactive Selection:**

* Use the **arrow keys** (or `j`/`k`) to navigate the provider list.

* Press **Spacebar** (or `Enter`) to select or deselect a provider.

* Once you've made your selections, press **`y`** to generate the Terraform files.

* Press **`q`** or `Ctrl+C` to quit at any time.

## Generated Files Overview

* `provider.tf`: Configures the `terraform` block with `required_providers` and sets up the `provider` blocks with `locals` for common variables like regions, project IDs, and API tokens.

* `variables.tf`: Defines the input variables used in `provider.tf` and other potential configurations, including sensitive variables where appropriate.

* `main.tf`: An empty file ready for you to add your infrastructure resources.

* `terraform.tfvars`: Provides example values for the variables defined in `variables.tf`. **Remember to replace placeholder values (e.g., `your-github-token`) with your actual credentials or desired settings.**

## Customization

The generated files are a solid starting point. Feel free to modify them to fit your specific project requirements, add more resources, or integrate with additional Terraform modules