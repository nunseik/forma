package cmd

import (
	"embed"
	"os"
	"github.com/spf13/cobra"
)

// This package-level variable will hold the filesystem passed from main.go.
var embeddedTemplates embed.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "forma",
	Short: "FORMA is a smart project initializer.",
	Long: `FORMA is a CLI tool designed to quickly scaffold and initialize new projects with best practices and templates.
It helps developers set up consistent project structures, apply templates, and automate repetitive setup tasks.
Use FORMA to boost productivity and maintain standardization across your projects.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(fs embed.FS) {
	// Store the passed filesystem in our package variable so other funcs can use it.
	embeddedTemplates = fs

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}


