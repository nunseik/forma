/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// TemplateConfig matches the structure of the template.yaml file.
type TemplateConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available project templates.",
	Long:  `Scans the templates directory and lists all available project templates found.`,
	Run: func(cmd *cobra.Command, args []string) {
		templatesPath := "./templates"

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
				fmt.Printf("    └─ Description: %s\n\n", config.Description)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
