package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	step int
	model struct {
		step       step
		templates  []string
		cursor     int
		template   string
		textInput  textinput.Model
		err        error
	}
)

const (
	stepChooseTemplate step = iota
	stepEnterProjectName
)
	
// Initialize the model with available templates.
func initialModel() model {
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepChooseTemplate:
		return updateChooseTemplate(msg, m)
	case stepEnterProjectName:
		return updateEnterProjectName(msg, m)
	}
	return m, cmd
}

func updateChooseTemplate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
		case "down", "j":
			if m.cursor < len(m.templates)-1 { m.cursor++ }
		case "enter":
			m.template = m.templates[m.cursor]
			m.step = stepEnterProjectName // Move to the next step
			return m, nil
		}
	}
	return m, nil
}

func updateEnterProjectName(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return m, tea.Quit // Done!
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n(press q to quit)", m.err)
	}

	if m.step == stepChooseTemplate {
		s := "Which template would you like to use?\n\n"
		for i, tpl := range m.templates {
			cursor := " "
			if m.cursor == i { cursor = ">" }
			s += fmt.Sprintf("%s %s\n", cursor, tpl)
		}
		s += "\n(press q to quit)\n"
		return s
	}

	return fmt.Sprintf(
		"What is the name of your project?\n\n%s\n\n(press enter to confirm)",
		m.textInput.View(),
	)
}
