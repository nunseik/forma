package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var author string

type TemplateData struct {
	ProjectName string
	Author      string
	Timestamp   string
}

var newCmd = &cobra.Command{
	Use:   "new <template> <project_name>",
	Short: "Creates a new project from a specified template.",
	Long: `Creates a new project directory based on a template.
For example:
forma new go-api my-awesome-project`,
	Example: `  forma new go-api my-awesome-project
  forma new python-app my-python-project --author "Jane Doe"`,
	// This makes sure the user provides exactly two arguments.
	Run: func(cmd *cobra.Command, args []string) {
		var templateName, projectName, finalAuthor string

		// If we have all required info, run directly.
		if len(args) == 2 && author != "" {
			templateName = args[0]
			projectName = args[1]
			finalAuthor = author
		} else {
			// No arguments, launch the TUI!
			m := initialModel(author)
			p := tea.NewProgram(m)
			finalModel, err := p.Run()
			if err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

			// Cast the final model to our model type
			final, ok := finalModel.(model)
			if !ok {
				fmt.Println("Error: unexpected model type returned from TUI.")
				return
			}

			// Check if there was an error in the TUI
			if final.err != nil {
				fmt.Println("\nProject creation aborted due to invalid input.")
				return
			}

			// Check if the user quit without confirming
			if final.template == "" || final.projectName == "" {
				fmt.Println("Aborted.")
				return
			}

			templateName = final.template
			projectName = final.projectName
			finalAuthor = final.author
		}

		systemTemplatesPath, err := getTemplatesPath()
		if err != nil {
			fmt.Printf("Error getting templates path: %v\n", err)
			return
		}
		templatePath := filepath.Join(systemTemplatesPath, templateName)

		// 1. Read and parse the template.yaml file to get hook info
		configPath := filepath.Join(templatePath, "template.yaml")
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading template config: %v\n", err)
			return
		}

		var templateConfig TemplateConfig
		if err := yaml.Unmarshal(yamlFile, &templateConfig); err != nil {
			fmt.Printf("Error parsing template config: %v\n", err)
			return
		}

		fmt.Printf("Creating a new project '%s' from template '%s'\n", projectName, templateName)

		// Create the new project directory.
		projectPath := filepath.Join(".", projectName)
		_, err = os.Stat(projectPath)
		if err == nil {
			// If the project directory already exists, prompt the user for confirmation to overwrite it.
			fmt.Printf("Project directory '%s' already exists. Do you want to overwrite it? (y/n): ", projectName)
			var response string
			fmt.Scanln(&response)
			normalized := strings.ToLower(strings.TrimSpace(response))
			if normalized != "y" {
				fmt.Println("Project creation aborted.")
				return
			} else {
				// If the user confirms, remove the existing directory.
				err = os.RemoveAll(projectPath)
				if err != nil {
					fmt.Printf("Error removing existing project directory: %v\n", err)
					return
				}
			}
		} else if !os.IsNotExist(err) {
			// If there was an error other than "not found", print it and exit.
			fmt.Printf("Error checking project directory: %v\n", err)
			return
		}

		data := TemplateData{
			ProjectName: projectName,
			Author:      finalAuthor,
			Timestamp:   time.Now().Format(time.RFC822),
		}


		// Copy the entire template structure.
		err = copyTemplate(templatePath, projectPath, data)
		if err != nil {
			fmt.Printf("Error creating project from template: %v\n", err)
			return
		}

		// 2. Run the post-create hooks
		if len(templateConfig.Hooks.PostCreate) > 0 {
			if err := runHooks(templateConfig.Hooks.PostCreate, projectPath, data); err != nil {
				fmt.Printf("Error running post-create hooks: %v\n", err)
				return
			}
		}

		fmt.Println("✅ Project created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&author, "author", "a", "", "Author of the project")
}
