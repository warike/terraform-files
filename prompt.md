
# Prompts.md

## Core Task: Go CLI Tool Development & Deployment

**Persona:** You are a senior software engineer specializing in Go and DevOps, adept at rapid prototyping and robust deployment strategies. You prioritize efficiency, best practices, and secure, maintainable solutions. Your responses must be concise, accurate, and directly actionable.

-----

### Phase 1: Go CLI Development

**Goal:** Generate the complete `main.go` file for the "wtf" CLI, exactly matching the provided complex structure and functionality for Terraform provider selection and file generation.

**Prompt:**

````
Generate the complete `main.go` file for a Go CLI application named "wtf".
The application is a Text User Interface (TUI) for selecting Terraform providers and generating `.tf` configuration files.

**Exact Code Structure and Functionality Requirements:**

1.  **Package and Imports:**
    - `package main`
    - **All necessary imports** as follows:
        ```go
        import (
            "encoding/json"
            "fmt"
            "io/ioutil"
            "net/http"
            "os"
            "regexp"
            "strings"
            "sync"

            "[github.com/charmbracelet/bubbles/spinner](https://github.com/charmbracelet/bubbles/spinner)"
            tea "[github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)"
            "[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)"
        )
        ```

2.  **Global Version Variable:**
    - Declare `var Version string = "dev"` immediately after imports.

3.  **Global Lipgloss Styles:**
    - Define all `lipgloss.NewStyle()` variables as provided: `titleStyle`, `checkboxStyle`, `checkedStyle`, `uncheckedStyle`, `helpStyle`, `spinnerStyle`, `errorStyle`, `successStyle`.

4.  **`Provider` Struct:**
    - Define the `Provider` struct with fields: `Name`, `Source`, `Version`, `LatestVersion`, `IsVersionLatest`.

5.  **`model` Struct:**
    - Define the `model` struct with fields: `providers`, `selected`, `cursor`, `spinner`, `loading`, `versionsLoaded`, `error`, `filesGenerated`.

6.  **`initialModel()` Function:**
    - Initialize the `model` with a predefined slice of `Provider` structs for "aws", "google", "azurerm", "github", "vercel", "cloudflare".
    - Configure `spinner.New()` using `spinner.Dot` and `spinnerStyle`.
    - Set `loading` to `true`.

7.  **`Init()` Method (on `model`):**
    - Return `tea.Batch` with `m.spinner.Tick` and `fetchAllVersions(m.providers)`.

8.  **`Update()` Method (on `model`):**
    - Handle `tea.KeyMsg` for navigation (`up`/`k`, `down`/`j`), selection (`enter`/` `), and exit (`ctrl+c`/`q`).
    - Handle 'y'/'Y' key to set `m.filesGenerated = true` and call `m.generateProviderFile()`, `m.generateVariablesFile()`, `m.generateTfvarsFile()`, `m.generateMainFile()`. Capture and set `m.error` if any generation fails.
    - Handle `spinner.TickMsg` to update the spinner.
    - Handle `versionsFetchedMsg` to set `m.loading = false`, `m.versionsLoaded = true`, update `m.providers`, and set `m.error` if `msg.err` exists.
    - Include `tea.WindowSizeMsg` case (empty).

9.  **`View()` Method (on `model`):**
    - Render error messages if `m.error` is not empty.
    - Render success message if `m.filesGenerated` is true.
    - Render loading state with spinner if `m.loading` is true.
    - Otherwise, build the main TUI view:
        - Display "Select Terraform Providers" title.
        - Iterate through `m.providers`, showing cursor, checkbox state, provider name, and version info (current/latest).
        - Include help text for navigation and actions.
        - Use `strings.Builder` for efficient string concatenation.

10. **`versionsFetchedMsg` Type:**
    - Define `type versionsFetchedMsg struct { providers []Provider; err error }`.

11. **`fetchAllVersions()` Function:**
    - Implement concurrent fetching of latest provider versions using `sync.WaitGroup` and `sync.Mutex`.
    - Call `getLatestProviderVersion()` for each provider.
    - Aggregate errors and return `versionsFetchedMsg`.

12. **`getLatestProviderVersion()` Function:**
    - Fetch latest version from `https://registry.terraform.io/v1/providers/<source>`.
    - Handle HTTP errors and JSON decoding errors.
    - Return the `version` string.

13. **File Generation Methods (on `*model`):**
    - `generateProviderFile()`: Generate `provider.tf` with `terraform` block, `required_providers`, and `provider` configurations (aws, google, azurerm, github, vercel, cloudflare) using `local` variables. Include `locals` block with `project_name` and `tags`.
    - `generateVariablesFile()`: Generate `variables.tf` with `project_name` and variables for selected providers (aws_region, aws_profile, google_project_id, google_region, azure_location, gh_owner, gh_token, vercel_api_token, cloudflare_api_token). Include `sensitive = true` where appropriate.
    - `generateMainFile()`: Generate a simple `main.tf` with a comment.
    - `generateTfvarsFile()`: Generate `terraform.tfvars` with default values for `project_name` and selected provider variables.
    - All generation functions must use `strings.Builder` and `ioutil.WriteFile`.

14. **`extractMajorMinor()` Function:**
    - Implement `func extractMajorMinor(version string) string` to extract major.minor version using `regexp`.

15. **`main()` Function:**
    - Initialize `tea.NewProgram(initialModel())`.
    - Run the program and handle any fatal errors by printing and `os.Exit(1)`.

**Output:** Provide only the complete, runnable `main.go` file.
````

-----

### Phase 2: Build Script Generation

**Goal:** Create a Bash script (`build.sh`) to compile the Go CLI for multiple platforms and prepare release assets.

**Prompt:**

```
Generate a Bash script named `build.sh` for building the "wtf" Go CLI.
- **Supported OS/Arch:** Must compile for:
    - macOS (Apple Silicon - `darwin_arm64`)
    - Linux (AMD64 - `linux_amd64`)
- **Version Retrieval:**
    - Use `git describe --tags --abbrev=0` to get the Git tag as the build version.
    - **Crucial:** If `git describe` fails (no tags), **default the version to "1.0.0"**.
    - **No verbose output** for version detection; only the final derived version.
- **Version Injection:** Inject the determined version into the `main.Version` variable of the Go binary using Go's `-ldflags="-X main.Version=${VERSION}"`.
- **Output Directory:** All compiled binaries must be placed in a `build/` directory. Create it if it doesn't exist.
- **Binary Naming Convention:**
    - Compiled binaries within `build/` must be named `wtf_<OS>_<ARCH>` (e.g., `wtf_darwin_arm64`). No version in binary name.
- **Release Asset Packaging:**
    - After compilation, for each binary, create a **single ZIP archive**.
    - ZIP file names must follow: `wtf_<VERSION_NUMBER_ONLY>_<OS>_<ARCH>.zip`.
    - **Crucial:** `VERSION_NUMBER_ONLY` must be derived from the Git tag by **removing any 'v' prefix** (e.g., `v1.0.0` -> `1.0.0` in ZIP name).
    - Ensure the `zip` command places the binary directly at the root of the ZIP archive, not inside subdirectories.
- **Script Executability:** Conclude the script with clear instructions for making it executable (`chmod +x build.sh`).
- **Output:** Provide only the complete `build.sh` script.
```

-----

### Phase 3: GitHub Release Instructions

**Goal:** Provide precise, step-by-step instructions for creating a GitHub Release for the compiled assets.

**Prompt:**

```
Outline the exact manual steps to create a GitHub Release for the "wtf" CLI.
- **Prerequisite:** Assume compiled and zipped binaries exist (e.g., `wtf_1.0.0_darwin_arm64.zip` in the `build/` directory).
- **Step 1: Code Push:** Detail commands to ensure all local code changes are committed and pushed to the remote GitHub repository.
- **Step 2: Git Tag Creation & Push:** Provide precise Git commands for:
    - Creating an **annotated Git tag** (`git tag -a`) corresponding to the release version (e.g., `v1.0.0`). The tag message should be brief and descriptive.
    - Pushing this specific Git tag to the remote GitHub repository (`git push origin <tag>`).
- **Step 3: GitHub UI Workflow:** Provide a step-by-step walkthrough of the GitHub website interface:
    - Navigation: "Your_Repository" -> "Releases".
    - Action: "Draft a new release" or "Create a new release" button.
    - Configuration:
        - "Choose a tag": Select the recently pushed Git tag.
        - "Release title": Suggest a clear title (e.g., "Warike Terraform v1.0.0").
        - "Description": Advise on content for release notes (new features, fixes).
        - **Critical:** "Attach binaries": How to upload the generated `.zip` files.
        - **Critical:** "Set as the latest release": Explicitly mention ensuring this option is selected.
    - Final Action: "Publish release".
- **Output:** Provide only the step-by-step instructions, concise and action-oriented.
```

-----

### Phase 4: Installation Script Generation

**Goal:** Create a Bash script (`install.sh`) to download and install the CLI from GitHub Releases.

**Prompt:**

```
Create a robust Bash script named `install.sh` for installing the "wtf" CLI.
- **Installation Directory:** Hardcode installation to `/usr/local/bin`.
- **Supported OS/Arch:** Must correctly auto-detect and install only on:
    - macOS (Apple Silicon - `darwin_arm64`)
    - Linux (AMD64 - `linux_amd64`)
    - **Crucial:** For unsupported OS/arch (e.g., Intel Mac, ARM Linux, Windows, unknown), print a clear error and `exit 1`.
- **GitHub Repository:**
    - Use `GITHUB_ORG_REPO="warike/terraform-files"` as the source repository.
    - **Crucial:** Hardcode `RELEASE_TAG="v1.0.0"` for the specific version to download.
- **Download Logic:**
    - Construct the `DOWNLOAD_URL` using the standard GitHub Releases pattern: `https://github.com/OWNER/REPO/releases/download/TAG/ASSET_NAME.zip`.
    - **Critical:** The `ASSET_NAME` in the URL must match `wtf_<VERSION_NO_V_PREFIX>_<OS>_<ARCH>.zip`.
        - Derive `VERSION_NO_V_PREFIX` from `RELEASE_TAG` (e.g., `v1.0.0` -> `1.0.0`).
    - Use `curl -L -f -o` for robust downloading to a temporary directory.
    - **Error Handling:** If download fails (e.g., 404), print specific error message including the URL, and `exit 1`.
- **Extraction:**
    - Use `unzip` to extract the binary from the downloaded `.zip` into the temporary directory.
    - **Error Handling:** If `unzip` fails, print an error and `exit 1`.
    - **Critical:** Verify the expected binary name (`wtf_<OS>_<ARCH>`) exists directly at the root of the unzipped temporary directory. If not found, print error and `exit 1`.
- **Installation:**
    - Use `sudo mv` to move the extracted binary to `INSTALL_DIR`.
    - Use `sudo chmod +x` to make the binary executable.
    - **Error Handling:** If `mv` or `chmod` fails (e.g., permissions), print error and `exit 1`.
- **Cleanup:** Delete the temporary directory and its contents upon success or failure.
- **User Feedback:**
    - Include ASCII art for "Warike" at the start.
    - Use concise, clear messages for success and error states. Avoid verbose explanations.
    - Final success message should instruct the user how to run the installed CLI.
- **Output:** Provide only the complete `install.sh` script.
```

-----

### Phase 5: Uninstallation Script Generation

**Goal:** Create a Bash script (`uninstall.sh`) to remove the installed CLI.

**Prompt:**

```
Generate a Bash script named `uninstall.sh` for uninstalling the "wtf" CLI.
- **Binary Location:** Assume the binary is installed at `/usr/local/bin/wtf`.
- **Uninstallation Logic:**
    - Check if the binary exists at the expected `INSTALL_DIR`. If not, print a message and exit.
    - Use `sudo rm` to remove the binary.
    - **Error Handling:** If `rm` fails (e.g., permissions), print an error and exit.
- **User Feedback:**
    - Include ASCII art for "Warike" at the start.
    - Use concise success/failure messages.
- **Output:** Provide only the complete `uninstall.sh` script.
```