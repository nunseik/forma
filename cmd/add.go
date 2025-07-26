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

func init() {
	rootCmd.AddCommand(addCmd)
}
