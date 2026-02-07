# tfinit - A Terraform Project Scaffolder & Updater

`tfinit` is a command-line interface (CLI) tool to quickly scaffold and manage your Terraform projects. It provides an interactive TUI to select common providers and also includes commands to keep your provider versions up-to-date.

## Features

*   **Interactive Scaffolding:** Interactively select from a list of popular Terraform providers (AWS, Google Cloud, Azure, etc.) to generate your initial project files.
*   **Version Management:** Automatically fetches the latest provider versions from the Terraform Registry.
*   **Automated Updates:** A simple `update` command to parse your existing `provider.tf` and update versions to the latest available.
*   **Standard File Generation:** Creates `provider.tf`, `variables.tf`, `main.tf`, and `terraform.tfvars` with sensible defaults.

## Installation

### Homebrew (macOS & Linux)

Once the project is published, you will be able to install it via Homebrew:
```bash
brew install warike/tools/tfinit
```

### From Source (Go)

Ensure you have a working Go environment (Go 1.18+).

```bash
go install github.com/warike/terraform-files@latest
```

## Usage

### 1. Create a New Project

The `create` command launches an interactive terminal UI to help you scaffold a new Terraform project.

**To create a project in the current directory:**

```bash
tfinit create --name .
```

**To create a project in a new directory called `my-infra`:**

```bash
tfinit create --name my-infra
```

**Interactive Selection:**

*   Use the **arrow keys** (`↑`/`↓` or `j`/`k`) to navigate.
*   Press **Spacebar** to select or deselect providers.
*   Press **`y`** to confirm and generate the files.
*   Press **`q`** or `Ctrl+C` to quit.

### 2. Update Provider Versions

The `update` command checks for newer versions of the providers listed in your `provider.tf` file and updates them automatically.

**To update a project in the current directory:**

```bash
tfinit update --name .
```

**To update a project in the `my-infra` directory:**

```bash
tfinit update --name my-infra
```

The tool will print a list of providers that were updated.

## Contributing

Contributions are welcome! Please see the [Contributing Guidelines](CONTRIBUTING.md) for more details on how to set up your development environment and submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
