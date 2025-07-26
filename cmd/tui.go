package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
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
		errorStyle  lipgloss.Style
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
		errorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
	}
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

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.step == stepChooseTemplate {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.templates)-1 {
					m.cursor++
				}
			case "enter":
				m.template = m.templates[m.cursor]
				m.step = stepEnterProjectName // Move to next step
				return m, nil
			}
		}
		return m, nil
	}
	// Handle the other steps which use text input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.step {
			case stepEnterProjectName:
				m.projectName = m.textInput.Value()
				if !isValidName(m.projectName) {
					m.err = fmt.Errorf("invalid project name: %s\nNames should start with a letter or number, only contains letters, numbers, hyphens, or underscores", m.projectName)
					return m, nil
				}
				m.err = nil // Reset error
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
				if !isValidName(m.author) {
					m.err = fmt.Errorf("invalid author name: %s\nNames should start with a letter or number, only contains letters, numbers, hyphens, or underscores", m.author)
					return m, nil
				}
				m.err = nil // Reset error
				return m, tea.Quit
			}
			return m, nil
		}
	}

	// Handle text input updates
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func isValidName(name string) bool {
    if name == "" {
        return false
    }
    // Regex to match a valid project name: starts and ends with a letter or number,
    // and contains only letters, numbers, hyphens, or underscores.
    re := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
    return re.MatchString(name)
}

// View renders the UI.
func (m model) View() string {
	var s string
	switch m.step {
	case stepChooseTemplate:
		s = "Which template would you like to use?\n\n"
		for i, tpl := range m.templates {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, tpl)
		}
	case stepEnterProjectName:
		s = fmt.Sprintf("What is the name of your project?\n\n%s\n\n(press enter to confirm)", m.textInput.View())
	case stepEnterAuthorName:
		s = fmt.Sprintf("What is your GitHub username?\n\n%s\n\n(press enter to confirm)", m.textInput.View())
	}
	
	if m.err != nil {
		s += fmt.Sprintf("\n\n%s", m.errorStyle.Render(m.err.Error()))
	}

	s += "\n(press ctrl+c to quit)\n"
	return s
}
