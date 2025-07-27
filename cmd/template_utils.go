package cmd

import (
	"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"text/template"
)

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
			// Use standard permissions (0777) to avoid permission issues.
			return os.MkdirAll(destPath, 0777)
		} else {
			// It's a file, so copy it.
			return processAndCopyFile(path, destPath, data)
		}
	}

	return filepath.WalkDir(templatePath, walkFunc)
}

// getAvailableTemplates scans the templates directory and returns a slice of template names.
func getAvailableTemplates() ([]string, error) {
	templatesPath, err := getTemplatesPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get templates path: %w", err)
	}
	var templates []string

	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if a template.yaml exists before adding it to the list
			configPath := filepath.Join(templatesPath, entry.Name(), "template.yaml")
			if _, err := os.Stat(configPath); err == nil {
				templates = append(templates, entry.Name())
			}
		}
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no valid templates found in %s", templatesPath)
	}

	return templates, nil
}