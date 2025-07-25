/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"io"

	"github.com/spf13/cobra"
	"path/filepath"
)

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

// copyTemplate walks through a template directory and copies its structure and files.
func copyTemplate(templatePath, projectPath string) error {
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
			return copyFile(path, destPath)
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
		err = os.Mkdir("./"+projectName, 0755)
		if err != nil {
			fmt.Printf("Error creating project directory '%s': %v\n", projectName, err)
			return
		}

		// 3. Copy the entire template structure.
		err = copyTemplate(templatePath, projectName)
		if err != nil {
			fmt.Printf("Error creating project from template: %v\n", err)
			return
		}
		
		fmt.Println("✅ Project created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
