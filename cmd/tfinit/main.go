package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"warike/base/internal/ui"
	"warike/base/internal/updater"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	
	switch cmd {
	case "create":
		handleCreate(os.Args[2:])
	case "update":
		handleUpdate(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func handleCreate(args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	nameFlag := createCmd.String("name", "", "Name of the project directory (optional, positional argument takes precedence)")

	createCmd.Parse(args)

	targetDir := *nameFlag

	// If a positional argument is provided, it's our target directory.
	if createCmd.NArg() > 0 {
		targetDir = createCmd.Arg(0)
	}

	// Default to current directory if no flag or argument is provided.
	if targetDir == "" {
		targetDir = "."
	}

	// Validate the target directory.
	if targetDir != "." {
		info, err := os.Stat(targetDir)
		if err == nil && info.IsDir() {
			// Directory exists, check if it's empty.
			entries, readErr := os.ReadDir(targetDir)
			if readErr != nil {
				fmt.Printf("Error reading directory '%s': %v\n", targetDir, readErr)
				os.Exit(1)
			}
			if len(entries) > 0 {
				fmt.Printf("Error: directory '%s' already exists and is not empty\n", targetDir)
				os.Exit(1)
			}
		}
		// os.IsNotExist(err) is the happy path for a new directory, so we don't check for it.
	}

	p := tea.NewProgram(ui.InitialModel(targetDir))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func handleUpdate(args []string) {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	name := updateCmd.String("name", ".", "Name of the project directory to update")
	
	updateCmd.Parse(args)
	
	targetDir := *name
	if targetDir == "" {
		targetDir = "."
	}
	
	u := updater.NewUpdater()
	updates, err := u.UpdateProject(targetDir)
	if err != nil {
		fmt.Printf("Error updating project: %v\n", err)
		os.Exit(1)
	}
	
	if len(updates) == 0 {
		fmt.Println("No updates available.")
	} else {
		for _, update := range updates {
			fmt.Println(update)
		}
	}
}

func printHelp() {
	fmt.Println("Usage: tfinit <command> [directory]")
	fmt.Println("\nCommands:")
	fmt.Println("  create [name]   Create a new Terraform project in the specified directory (defaults to current dir)")
	fmt.Println("  update [name]   Update providers in an existing project (defaults to current dir)")
}