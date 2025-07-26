package cmd

import (
	"fmt"
    "io/fs"
    "os"
    "path/filepath"
)

// getTemplatesPath ensures the config path exists and returns it.
// It handles the first-run-only copying of embedded templates.
func getTemplatesPath() (string, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user config dir: %w", err)
    }

    templatesPath := filepath.Join(configDir, "forma", "templates")

    // Check if the directory exists. If not, it's the first run.
    if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
        fmt.Println("Performing first-time setup, creating templates folder...")

        // Create the full path, e.g., ~/.config/forma/templates
        if err := os.MkdirAll(templatesPath, 0755); err != nil {
            return "", fmt.Errorf("failed to create templates directory: %w", err)
        }

        // Copy the embedded templates
        templatesRoot, _ := fs.Sub(embeddedTemplates, "templates")
        err := fs.WalkDir(templatesRoot, ".", func(path string, d fs.DirEntry, err error) error {
            if err != nil {
                return err
            }
            destPath := filepath.Join(templatesPath, path)
            if d.IsDir() {
                return os.MkdirAll(destPath, 0755)
            }
            content, _ := fs.ReadFile(templatesRoot, path)
            return os.WriteFile(destPath, content, 0644)
        })

        if err != nil {
            return "", fmt.Errorf("failed to copy embedded templates: %w", err)
        }
    }

    return templatesPath, nil
}