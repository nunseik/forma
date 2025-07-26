/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
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
	rootCmd.AddCommand(removeCmd)
}
