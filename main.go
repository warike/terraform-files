package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var Version string = "dev" // Default value for development

var (
	// Styles
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	checkboxStyle  = lipgloss.NewStyle().PaddingLeft(2)
	checkedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	uncheckedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(2)
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
)

// Provider represents a Terraform provider with its name and version info
type Provider struct {
	Name            string
	Source          string
	Version         string
	LatestVersion   string
	IsVersionLatest bool
}

// model represents the state of the application
type model struct {
	providers      []Provider
	selected       []bool
	cursor         int
	spinner        spinner.Model
	loading        bool
	versionsLoaded bool
	error          string
	filesGenerated bool
}

// initialModel creates the initial state of the application
func initialModel() model {
	providers := []Provider{
		{Name: "aws", Source: "hashicorp/aws"},
		{Name: "google", Source: "hashicorp/google"},
		{Name: "azurerm", Source: "hashicorp/azurerm"},
		{Name: "github", Source: "integrations/github"},
		{Name: "vercel", Source: "vercel/vercel"},
		{Name: "cloudflare", Source: "cloudflare/cloudflare"},
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return model{
		providers: providers,
		selected:  make([]bool, len(providers)),
		spinner:   s,
		loading:   true,
	}
}

// Init is the first function that will be called. It returns a command.
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchAllVersions(m.providers))
}

// Update handles all incoming messages and updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading || m.filesGenerated {
			if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "enter" {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.providers)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "y", "Y":
			m.filesGenerated = true
			err := m.generateProviderFile()
			if err != nil {
				m.error = fmt.Sprintf("Error generating provider.tf: %v", err)
			}
			err = m.generateVariablesFile()
			if err != nil {
				m.error = fmt.Sprintf("Error generating variables.tf: %v", err)
			}

			err = m.generateTfvarsFile()
			if err != nil {
				m.error = fmt.Sprintf("Error generating terraform.tfvars: %v", err)
			}

			err = m.generateMainFile()
			if err != nil {
				m.error = fmt.Sprintf("Error generating main.tf: %v", err)
			}

			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case versionsFetchedMsg:
		m.loading = false
		m.versionsLoaded = true
		m.providers = msg.providers
		if msg.err != nil {
			m.error = msg.err.Error()
		}

	case tea.WindowSizeMsg:
	}

	return m, nil
}

// View renders the UI.
func (m model) View() string {
	if m.error != "" {
		return errorStyle.Render(fmt.Sprintf("Error: %s\n\nPress Enter to exit.", m.error))
	}

	if m.filesGenerated {
		s := successStyle.Render("Terraform files (provider.tf, variables.tf, main.tf, terraform.tfvars) generated successfully!") + "\n\n"

		s += "\n\nPress Enter to exit."
		return s
	}

	if m.loading {
		return fmt.Sprintf("%s Fetching latest provider versions...", m.spinner.View())
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Select Terraform Providers"))
	sb.WriteString("\n\n")

	for i, p := range m.providers {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		style := uncheckedStyle
		if m.selected[i] {
			checked = "x"
			style = checkedStyle
		}

		versionInfo := fmt.Sprintf("latest: %s", p.LatestVersion)
		if !p.IsVersionLatest && p.Version != "" {
			versionInfo = fmt.Sprintf("current: %s, latest: %s", p.Version, p.LatestVersion)
		}

		sb.WriteString(fmt.Sprintf("%s [%s] %s (%s)\n", cursor, style.Render(checked), p.Name, versionInfo))
	}

	sb.WriteString(helpStyle.Render("\nUse arrow keys to navigate, space/enter to select, 'y' to generate files, 'q' to quit.\n"))

	return sb.String()
}

// --- Commands and Messages ---

type versionsFetchedMsg struct {
	providers []Provider
	err       error
}

func fetchAllVersions(providers []Provider) tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		var mu sync.Mutex
		updatedProviders := make([]Provider, len(providers))
		copy(updatedProviders, providers)
		var firstErr error

		for i := range providers {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				version, err := getLatestProviderVersion(updatedProviders[i].Source)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = fmt.Errorf("failed to fetch version for %s: %w", updatedProviders[i].Name, err)
					}
					return
				}
				updatedProviders[i].LatestVersion = version
				// In a real app, you'd compare against a version from a file
				updatedProviders[i].IsVersionLatest = true
			}(i)
		}

		wg.Wait()
		return versionsFetchedMsg{providers: updatedProviders, err: firstErr}
	}
}

// --- Helper Functions ---

func getLatestProviderVersion(source string) (string, error) {
	url := fmt.Sprintf("https://registry.terraform.io/v1/providers/%s", source)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	var result struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Version, nil
}

func (m *model) generateProviderFile() error {
	var sb strings.Builder

	sb.WriteString("terraform {\n")
	sb.WriteString("  required_providers {\n")
	for i, p := range m.providers {
		if m.selected[i] {
			sb.WriteString(fmt.Sprintf("    %s = {\n", p.Name))
			sb.WriteString(fmt.Sprintf("      source  = \"%s\"\n", p.Source))
			sb.WriteString(fmt.Sprintf("      version = \"%s\"\n", p.LatestVersion))
			sb.WriteString("    }\n")
		}
	}
	sb.WriteString("  }\n")
	sb.WriteString("}\n")

	// Add provider configurations based on selected providers
	for i, p := range m.providers {
		if m.selected[i] {
			sb.WriteString("\n")
			switch p.Name {
			case "aws":
				sb.WriteString("provider \"aws\" {\n")
				sb.WriteString("  region  = local.aws_region\n")
				sb.WriteString("  profile = local.aws_profile\n")
				sb.WriteString("  default_tags {\n")
				sb.WriteString("   tags = local.tags\n")
				sb.WriteString("  }\n")
				sb.WriteString("}\n")
			case "google":
				sb.WriteString("provider \"google\" {\n")
				sb.WriteString("  project = local.gcp_project_id\n")
				sb.WriteString("  region  = local.gcp_region\n")
				sb.WriteString("}\n")
			case "azurerm":
				sb.WriteString("provider \"azurerm\" {\n")
				sb.WriteString("  subscription_id = local.azure_subscription_id\n")
				sb.WriteString("  features {}\n")
				sb.WriteString("}\n")
			case "github":
				sb.WriteString("provider \"github\" {\n")
				sb.WriteString("  owner = local.gh_owner\n")
				sb.WriteString("  token = local.gh_token\n")
				sb.WriteString("}\n")
			case "vercel":
				sb.WriteString("provider \"vercel\" {\n")
				sb.WriteString("  api_token = local.vercel_api_token\n")
				sb.WriteString("}\n")
			case "cloudflare":
				sb.WriteString("provider \"cloudflare\" {\n")
				sb.WriteString("  api_token = local.cloudflare_api_token\n")
				sb.WriteString("}\n")
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString("locals {\n")
	sb.WriteString("  project_name = var.project_name\n")
	for i, p := range m.providers {
		if m.selected[i] {
			switch p.Name {
			case "aws":
				sb.WriteString("  aws_region  = var.aws_region\n")
				sb.WriteString("  aws_profile = var.aws_profile\n")
			case "google":
				sb.WriteString("  gcp_project_id = var.google_project_id\n")
				sb.WriteString("  gcp_region     = var.google_region\n")
			case "azurerm":
				sb.WriteString("  azure_location            = var.azure_location\n")
				sb.WriteString("  azure_subscription_id     = var.azure_subscription_id\n")
			case "github":
				sb.WriteString("  gh_owner           = var.gh_owner\n")
				sb.WriteString("  gh_token           = var.gh_token\n")
			case "vercel":
				sb.WriteString("  vercel_api_token = var.vercel_api_token\n")
			case "cloudflare":
				sb.WriteString("  cloudflare_api_token = var.cloudflare_api_token\n")
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString("  tags = {\n")
	sb.WriteString("    project     = local.project_name\n")
	sb.WriteString("    environment = \"dev\"\n")
	sb.WriteString("    owner       = \"warike\"\n")
	sb.WriteString("    cost-center  = \"development\"\n")
	sb.WriteString("    terraform   = \"true\"\n")
	sb.WriteString("  }\n")

	sb.WriteString("}\n")

	return os.WriteFile("provider.tf", []byte(sb.String()), 0644)
}

func (m *model) generateVariablesFile() error {
	var sb strings.Builder

	sb.WriteString("variable \"project_name\" {\n")
	sb.WriteString("  description = \"Name of the project\"\n")
	sb.WriteString("  type        = string\n")
	sb.WriteString("  default     = \"my_project\"\n")
	sb.WriteString("}\n")

	for i, p := range m.providers {
		if m.selected[i] {
			switch p.Name {
			case "aws":
				sb.WriteString("\n")
				sb.WriteString("variable \"aws_region\" {\n")
				sb.WriteString("  description = \"AWS region\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  default     = \"us-west-2\"\n")
				sb.WriteString("}\n")
				sb.WriteString("\n")
				sb.WriteString("variable \"aws_profile\" {\n")
				sb.WriteString("  description = \"AWS profile name\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  default     = \"default\"\n")
				sb.WriteString("}\n")
			case "google":
				sb.WriteString("\n")
				sb.WriteString("variable \"google_project_id\" {\n")
				sb.WriteString("  description = \"Google Cloud project ID\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("}\n")

				sb.WriteString("\n")
				sb.WriteString("variable \"google_region\" {\n")
				sb.WriteString("  description = \"Google Cloud region\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  default     = \"us-central1\"\n")
				sb.WriteString("}\n")
			case "azurerm":
				sb.WriteString("\n")
				sb.WriteString("variable \"azure_location\" {\n")
				sb.WriteString("  description = \"Azure location\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  default     = \"East US\"\n")
				sb.WriteString("}\n")

				sb.WriteString("\n")
				sb.WriteString("variable \"azure_subscription_id\" {\n")
				sb.WriteString("  description = \"Azure subscription ID\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  sensitive   = true\n")
				sb.WriteString("}\n")
			case "github":
				sb.WriteString("\n")
				sb.WriteString("variable \"gh_owner\" {\n")
				sb.WriteString("  description = \"GitHub owner (user or organization)\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  default     = \"warike\"\n")
				sb.WriteString("}\n")
				sb.WriteString("\n")
				sb.WriteString("variable \"gh_token\" {\n")
				sb.WriteString("  description = \"GitHub token\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  sensitive   = true\n")
				sb.WriteString("}\n")
			case "vercel":
				sb.WriteString("\n")
				sb.WriteString("variable \"vercel_api_token\" {\n")
				sb.WriteString("  description = \"Vercel API Token\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  sensitive   = true\n")
				sb.WriteString("}\n")
			case "cloudflare":
				sb.WriteString("\n")
				sb.WriteString("variable \"cloudflare_api_token\" {\n")
				sb.WriteString("  description = \"Cloudflare API Token\"\n")
				sb.WriteString("  type        = string\n")
				sb.WriteString("  sensitive   = true\n")
				sb.WriteString("}\n")
			}
		}
	}

	return os.WriteFile("variables.tf", []byte(sb.String()), 0644)
}

// Generate main.tf
func (m *model) generateMainFile() error {
	var sb strings.Builder

	sb.WriteString("// main.tf\n")
	sb.WriteString("\n")

	return os.WriteFile("main.tf", []byte(sb.String()), 0644)
}

// Generate terraform.tfvars based on selected providers
func (m *model) generateTfvarsFile() error {
	var sb strings.Builder

	sb.WriteString("project_name = \"my_project\"\n")

	for i, p := range m.providers {
		if m.selected[i] {
			sb.WriteString("\n")
			switch p.Name {
			case "aws":
				sb.WriteString("aws_region   = \"us-west-2\"\n")
				sb.WriteString("aws_profile  = \"default\"\n")
			case "google":
				sb.WriteString("google_project_id = \"gcp-project-id-goes-here\"\n")
				sb.WriteString("google_region     = \"us-central1\"\n")
			case "azurerm":
				sb.WriteString("azure_location = \"East US\"\n")
				sb.WriteString("azure_subscription_id = \"azure-subscription-id-goes-here\"\n")
			case "github":
				sb.WriteString("gh_owner = \"warike\"\n")
				sb.WriteString("gh_token = \"your-github-token\"\n")
			case "vercel":
				sb.WriteString("vercel_api_token = \"your-vercel-token\"\n")
			case "cloudflare":
				sb.WriteString("cloudflare_api_token = \"your-cloudflare-token\"\n")
			}
		}
	}

	return os.WriteFile("terraform.tfvars", []byte(sb.String()), 0644)
}

// Extract major.minor version from full version string
func extractMajorMinor(version string) string {
	re := regexp.MustCompile(`^(\d+\.\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		return matches[1]
	}
	return version
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
