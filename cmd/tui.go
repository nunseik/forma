package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	step  int
	model struct {
		step        step
		templates   []string
		cursor      int
		template    string
		projectName string
		author      string
		textInput   textinput.Model
		err         error
	}
)

const (
	stepChooseTemplate step = iota
	stepEnterProjectName
	stepEnterAuthorName
)

// Initialize the model with available templates.
func initialModel(flagAuthor string) model {
	templates, err := getAvailableTemplates()
	if err != nil {
		fmt.Println("Error getting templates:", err)
		os.Exit(1)
	}
	ti := textinput.New()
	ti.Placeholder = "my-awesome-app"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		step:      stepChooseTemplate,
		templates: templates,
		author:    flagAuthor,
		textInput: ti,
		err:       err,
	}
}

// getAvailableTemplates scans the templates directory and returns a slice of template names.
func getAvailableTemplates() ([]string, error) {
	templatesPath := "./templates"
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

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.step {
			case stepChooseTemplate:
				m.template = m.templates[m.cursor]
				m.step = stepEnterProjectName // Move to next step
			case stepEnterProjectName:
				m.projectName = m.textInput.Value()
				m.textInput.Reset()
				// If author was not provided by flag, ask for it. Otherwise, we're done.
				if m.author == "" {
					m.textInput.Placeholder = "YourGitHubUsername"
					m.step = stepEnterAuthorName
				} else {
					return m, tea.Quit
				}
			case stepEnterAuthorName:
				m.author = m.textInput.Value()
				return m, tea.Quit // Done!
			}
			return m, nil
		}
	}

	// Handle text input updates
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n(press q to quit)", m.err)
	}

	switch m.step {
	case stepChooseTemplate:
		s := "Which template would you like to use?\n\n"
		for i, tpl := range m.templates {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, tpl)
		}
		s += "\n(press q or ctrl+c to quit)\n"
		return s
	case stepEnterProjectName:
		return fmt.Sprintf("What is the name of your project?\n\n%s\n\n(press enter to confirm)", m.textInput.View())
	case stepEnterAuthorName:
		return fmt.Sprintf("What is your GitHub username?\n\n%s\n\n(press enter to confirm)", m.textInput.View())
	}
	return ""
}
