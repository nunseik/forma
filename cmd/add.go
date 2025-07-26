package cmd

import (
	"fmt"
	"strings"
	"path/filepath"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <git_repo_url>",
	Short: "Add a new template from a Git repository",
	Long:  `Clones a Git repository into the FORMA templates directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]

		templatesPath, err := getTemplatesPath()
		if err != nil {
			fmt.Printf("Error getting templates path: %v\n", err)
			return
		}

		repoName := strings.TrimSuffix(filepath.Base(repoURL), ".git")
		destPath := filepath.Join(templatesPath, repoName)

		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("Template '%s' already exists in '%s'.\n", repoName, templatesPath)
			return
		}
		fmt.Printf("Cloning template from '%s' into '%s'...\n", repoURL, destPath)

		gitCmd := exec.Command("git", "clone", repoURL, destPath)

		// Run the command and capture the combined output (stdout and stderr).
		output, err := gitCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error cloning repository: %v\nOutput: %s\n", err, output)
			return
		}	
		fmt.Printf("Successfully added template '%s'.\n", repoName)
		fmt.Println("You can now use this template with the 'new' command.")
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <template_name>",
	Short: "Remove a template",
	Long:  `Removes a template from the FORMA templates directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify a template name to remove.")
			return
		}

		templateName := args[0]
		templatesPath, err := getTemplatesPath()
		if err != nil {
			fmt.Printf("Error getting templates path: %v\n", err)
			return
		}

		templatePath := filepath.Join(templatesPath, templateName)

		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			fmt.Printf("Template '%s' does not exist in '%s'.\n", templateName, templatesPath)
			return
		}

		fmt.Printf("Are you sure you want to remove the template '%s'? (y/n): ", templateName)
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			fmt.Println("Aborted.")
			return
		}

		err = os.RemoveAll(templatePath)
		if err != nil {
			fmt.Printf("Error removing template: %v\n", err)
			return
		}

		fmt.Printf("Successfully removed template '%s'.\n", templateName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
}
