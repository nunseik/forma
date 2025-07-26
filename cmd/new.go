/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var author string
var goVersion string

type TemplateData struct {
	ProjectName string
	Author      string
	GoVersion   string
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

		if d.IsDir() {
			// It's a directory, so create it in the destination.
			// MkdirAll is used to create parent directories if they don't exist.
			return os.MkdirAll(destPath, d.Type().Perm())
		} else {
			// It's a file, so copy it.
			// (Assuming you have the copyFile function from our previous conversation)
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
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: You must specify a template and a project name.")
			return
		}
		template := args[0]
		projectName := args[1]

		fmt.Printf("Creating a new project '%s' from template '%s'\n", projectName, template)

		// --- TODO: DAY 1 LOGIC GOES HERE ---

		templatePath := "./templates/" + template
		// 1. Check if the template exists in the "./templates" directory.
		//    If it does not exist, print an error message and exit.
		_, err := os.Stat(templatePath)
		if err != nil {
			fmt.Printf("Template '%s' not found at '%s'.\n", template, templatePath)
			return
		}

		_, err = os.Stat("./" + projectName)
		if err == nil {
			// 2. If the project directory already exists, prompt the user for confirmation to overwrite it.
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
		// 3. Create the new project directory.
		projectPath := "./" + projectName
		data := TemplateData{
			ProjectName: projectName,
			Author:      author,
			GoVersion:   goVersion,
			Timestamp:   time.Now().Format(time.RFC822),
		}

		// 3. Copy the entire template structure.
		err = copyTemplate(templatePath, projectPath, data)
		if err != nil {
			fmt.Printf("Error creating project from template: %v\n", err)
			return
		}

		fmt.Println("✅ Project created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Add flags here
	newCmd.Flags().StringVarP(&author, "author", "a", "Your Name", "Author of the project")
	newCmd.Flags().StringVarP(&goVersion, "go-version", "g", "1.24", "Go version for the go.mod file")
}
