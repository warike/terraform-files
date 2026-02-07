package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"warike/base/internal/generator"
	"warike/base/internal/providers"
)

type Provider struct {
	Name            string
	Source          string
	LatestVersion   string
	IsVersionLatest bool
}

type Model struct {
	Providers      []Provider
	Selected       []bool
	Cursor         int
	Spinner        spinner.Model
	Loading        bool
	VersionsLoaded bool
	Error          string
	FilesGenerated bool
	Client         *providers.Client
	TargetDir      string
}

func InitialModel(targetDir string) Model {
	p := []Provider{
		{Name: "aws", Source: "hashicorp/aws"},
		{Name: "google", Source: "hashicorp/google"},
		{Name: "azurerm", Source: "hashicorp/azurerm"},
		{Name: "github", Source: "integrations/github"},
		{Name: "vercel", Source: "vercel/vercel"},
		{Name: "cloudflare", Source: "cloudflare/cloudflare"},
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	return Model{
		Providers: p,
		Selected:  make([]bool, len(p)),
		Spinner:   s,
		Loading:   true,
		Client:    providers.NewClient(),
		TargetDir: targetDir,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.fetchAllVersions())
}

type versionsFetchedMsg struct {
	providers []Provider
	err       error
}

func (m Model) fetchAllVersions() tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		var mu sync.Mutex
		updatedProviders := make([]Provider, len(m.Providers))
		copy(updatedProviders, m.Providers)
		var firstErr error

		for i := range m.Providers {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				version, err := m.Client.GetLatestVersion(updatedProviders[i].Source)
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = fmt.Errorf("failed to fetch version for %s: %w", updatedProviders[i].Name, err)
					}
					return
				}
				updatedProviders[i].LatestVersion = version
				updatedProviders[i].IsVersionLatest = true
			}(i)
		}

		wg.Wait()
		return versionsFetchedMsg{providers: updatedProviders, err: firstErr}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Loading || m.FilesGenerated {
			if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "enter" {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Providers)-1 {
				m.Cursor++
			}
		case "enter", " ":
			m.Selected[m.Cursor] = !m.Selected[m.Cursor]
		case "g", "G":
			m.FilesGenerated = true
			if err := m.generateFiles(); err != nil {
				m.Error = err.Error()
			}
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case versionsFetchedMsg:
		m.Loading = false
		m.VersionsLoaded = true
		m.Providers = msg.providers
		if msg.err != nil {
			m.Error = msg.err.Error()
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Error != "" {
		return ErrorStyle.Render(fmt.Sprintf("Error: %s\n\nPress Enter to exit.", m.Error))
	}

	if m.FilesGenerated {
		return SuccessStyle.Render("Terraform files generated successfully!") + "\n\nPress Enter to exit."
	}

	if m.Loading {
		return fmt.Sprintf("%s Fetching latest provider versions...", m.Spinner.View())
	}

	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Select Terraform Providers"))
	sb.WriteString("\n\n")

	for i, p := range m.Providers {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		style := UncheckedStyle
		if m.Selected[i] {
			checked = "x"
			style = CheckedStyle
		}

		versionInfo := fmt.Sprintf("latest: %s", p.LatestVersion)
		sb.WriteString(fmt.Sprintf("%s [%s] %s (%s)\n", cursor, style.Render(checked), p.Name, versionInfo))
	}

	sb.WriteString(HelpStyle.Render("\n[space/enter] select | [g] generate | [q] quit\n"))

	return sb.String()
}

func (m Model) generateFiles() error {
	// Ensure target directory exists
	if m.TargetDir != "" && m.TargetDir != "." {
		if err := os.MkdirAll(m.TargetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", m.TargetDir, err)
		}
	}

	genData := generator.GeneratorData{
		ProjectName: "my_project",
	}
	
	if m.TargetDir != "" && m.TargetDir != "." {
		genData.ProjectName = filepath.Base(m.TargetDir)
	}

	for i, p := range m.Providers {
		if m.Selected[i] {
			genData.Providers = append(genData.Providers, generator.ProviderConfig{
				Name:          p.Name,
				Source:        p.Source,
				LatestVersion: p.LatestVersion,
			})
		}
	}

	files := []struct {
		name string
		gen  func(generator.GeneratorData) ([]byte, error)
	}{
		{"provider.tf", generator.GenerateProviderFile},
		{"variables.tf", generator.GenerateVariablesFile},
		{"terraform.tfvars", generator.GenerateTfvarsFile},
		{"main.tf", generator.GenerateMainFile},
	}

	for _, f := range files {
		content, err := f.gen(genData)
		if err != nil {
			return err
		}
		
		path := f.name
		if m.TargetDir != "" {
			path = filepath.Join(m.TargetDir, f.name)
		}
		
		if err := generator.WriteFile(path, content); err != nil {
			return err
		}
	}
	return nil
}
