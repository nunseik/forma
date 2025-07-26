package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

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

// processAndCopyFile reads a source file, processes it as a Go template,
// and writes the output to the destination file.
func processAndCopyFile(src, dst string, data TemplateData) error {
	// Read the source file content
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", src, err)
	}

	// Create a new template and parse the file content
	tmpl, err := template.New(filepath.Base(src)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", src, err)
	}

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Execute the template, writing the output to the destination file
	err = tmpl.Execute(destFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// copyTemplate walks through a template directory and copies its structure and files.
func copyTemplate(templatePath, projectPath string, data TemplateData) error {
	// Make sure the destination project directory exists.
	// os.MkdirAll is safe to call even if the directory already exists.
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Walk the template directory.
	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Propagate errors from walking the directory.
		}

		// Get the relative path of the file/dir with respect to the template root.
		relativePath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Create the full destination path.
		destPath := filepath.Join(projectPath, relativePath)

		// Skip the template.yaml file itself.
		if d.Name() == "template.yaml" {
			return nil
		}

		if d.IsDir() {
			// It's a directory, so create it in the destination.
			// MkdirAll is used to create parent directories if they don't exist.
			return os.MkdirAll(destPath, d.Type().Perm())
		} else {
			// It's a file, so copy it.
			return processAndCopyFile(path, destPath, data)
		}
	}

	return filepath.WalkDir(templatePath, walkFunc)
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <template> <project_name>",
	Short: "Creates a new project from a specified template.",
	Long: `Creates a new project directory based on a template.
For example:
forma new go-api my-awesome-project`,
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
			final := finalModel.(model)

			// Check if the user quit without confirming
			if final.template == "" || final.projectName == "" {
				fmt.Println("Aborted.")
				return
			}

			templateName = final.template
			projectName = final.projectName
			finalAuthor = final.author
		}

		fmt.Printf("Creating a new project '%s' from template '%s'\n", projectName, templateName)

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

		_, err = os.Stat("./" + projectName)
		if err == nil {
			// If the project directory already exists, prompt the user for confirmation to overwrite it.
			fmt.Printf("Project directory '%s' already exists. Do you want to overwrite it? (y/n): ", projectName)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Project creation aborted.")
				return
			} else if response == "y" || response == "Y" {
				// If the user confirms, remove the existing directory.
				err = os.RemoveAll("./" + projectName)
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

		// Create the new project directory.
		projectPath := "./" + projectName

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
