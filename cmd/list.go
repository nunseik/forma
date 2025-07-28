package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"bytes"
	"os/exec"
	"text/template"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// HooksConfig holds commands to be run at different stages.
type HooksConfig struct {
	PostCreate []string `yaml:"post_create"`
}

// TemplateConfig matches the structure of the template.yaml file.
type TemplateConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Hooks       HooksConfig `yaml:"hooks"`
}

// listTemplatesCmd represents the list command
func runHooks(commands []string, projectPath string, data TemplateData) error {
	if len(commands) == 0 {
		return nil
	}

	fmt.Println("--- The following post-creation hooks will be executed ---")
	for i, commandStr := range commands {
		// Process the command string as a template for preview
		tmpl, err := template.New("hook").Parse(commandStr)
		if err != nil {
			fmt.Printf("  [%d] (template parse error): %s\n", i+1, commandStr)
			continue
		}
		var processedCmd bytes.Buffer
		if err := tmpl.Execute(&processedCmd, data); err != nil {
			fmt.Printf("  [%d] (template exec error): %s\n", i+1, commandStr)
			continue
		}
		fmt.Printf("  [%d] %s\n", i+1, processedCmd.String())
	}

	fmt.Print("Do you want to proceed with executing all hooks? [y/n]: ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil || (strings.ToLower(strings.TrimSpace(response)) != "y") {
		fmt.Println("Aborted running hooks.")
		return nil
	}

	fmt.Println("--- Running post-creation hooks ---")
	for _, commandStr := range commands {
		tmpl, err := template.New("hook").Parse(commandStr)
		if err != nil {
			return fmt.Errorf("failed to parse hook command template: %w", err)
		}

		var processedCmd bytes.Buffer
		if err := tmpl.Execute(&processedCmd, data); err != nil {
			return fmt.Errorf("failed to execute hook command template: %w", err)
		}

		command := processedCmd.String()
		fmt.Printf("▶️ Running: %s\n", command)

		cmd := exec.Command("sh", "-c", command)
		cmd.Dir = projectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook command '%s' failed: %w", command, err)
		}
	}

	fmt.Println("--- Hooks finished successfully ---")
	return nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available project templates.",
	Long:  `Scans the templates directory and lists all available project templates found.`,
	Run: func(cmd *cobra.Command, args []string) {
		templatesPath, err := getTemplatesPath()
		if err != nil {
			fmt.Printf("Error getting templates path: %v\n", err)
			return
		}

		entries, err := os.ReadDir(templatesPath)
		if err != nil {
			fmt.Printf("Error reading templates directory: %v\n", err)
			return
		}

		fmt.Println("Available templates:")
		fmt.Println("---------------------")

		foundTemplates := 0

		for _, entry := range entries {
			if entry.IsDir() {
				templateName := entry.Name()
				configPath := filepath.Join(templatesPath, templateName, "template.yaml")

				// Check if template.yaml exists
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					continue // No yaml file, skip this directory
				}

				// Read and parse the yaml file
				yamlFile, err := os.ReadFile(configPath)
				if err != nil {
					fmt.Printf("! Error reading config for '%s': %v\n", templateName, err)
					continue
				}

				var config TemplateConfig
				err = yaml.Unmarshal(yamlFile, &config)
				if err != nil {
					fmt.Printf("! Error parsing config for '%s': %v\n", templateName, err)
					continue
				}

				// Print the details
				fmt.Printf("  %s\n", config.Name)
				fmt.Printf("    └─ ID: %s\n", templateName)
				// Print only the first line of the description, trimmed to 100 characters
				desc := strings.SplitN(config.Description, "\n", 2)[0]
				if len(desc) > 100 {
					desc = desc[:97] + "..."
				}
				fmt.Printf("    └─ Description: %s\n\n", desc)
				foundTemplates++
			}
		}

		if foundTemplates == 0 {
			fmt.Println("No templates found in the './templates' directory.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
