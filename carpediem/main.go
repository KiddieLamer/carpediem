package main

import (
	"fmt"
	"os"

	"github.com/KiddieLamer/carpediem/cmd"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	logoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	subStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#A7A7A7")).Margin(0, 0, 1, 0)
	itemStyle    = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#FFFFFF"))
	selectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#00FF87")).Bold(true)
	quitStyle    = lipgloss.NewStyle().Margin(1, 0, 0, 0).Foreground(lipgloss.Color("#A7A7A7"))
)

type model struct {
	choices  []string
	cursor   int
	selected string
	quitting bool
}

func initialModel() model {
	return model{
		choices: []string{"▶ Run Automation", "  Input Accounts", "  Init File", "  About", "  Exit"},
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := "\n"
	s += logoStyle.Render(` ▓▓▓   ▓▓▓  ▓▓▓▓  ▓▓▓▓  ▓▓▓▓▓    ▓▓▓▓  ▓▓▓ ▓▓▓▓▓ ▓   ▓`) + "\n"
	s += logoStyle.Render(`▓     ▓   ▓ ▓   ▓ ▓   ▓ ▓        ▓   ▓  ▓  ▓     ▓▓ ▓▓`) + "\n"
	s += logoStyle.Render(`▓     ▓▓▓▓▓ ▓▓▓▓  ▓▓▓▓  ▓▓▓▓     ▓   ▓  ▓  ▓▓▓▓  ▓ ▓ ▓`) + "\n"
	s += logoStyle.Render(`▓     ▓   ▓ ▓  ▓  ▓     ▓        ▓   ▓  ▓  ▓     ▓   ▓`) + "\n"
	s += logoStyle.Render(` ▓▓▓  ▓   ▓ ▓   ▓ ▓     ▓▓▓▓▓    ▓▓▓▓  ▓▓▓ ▓▓▓▓▓ ▓   ▓`) + "\n\n"
	s += subStyle.Render("     ChatGPT Auth Automator") + "\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "▸ "
			s += selectedStyle.Render(fmt.Sprintf("%s%s", cursor, choice)) + "\n"
		} else {
			s += itemStyle.Render(fmt.Sprintf("%s%s", cursor, choice)) + "\n"
		}
	}

	s += quitStyle.Render("\n  ↑/↓ navigate • enter select • q quit\n")
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("\n  ─────────────────────────────\n")

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result := m.(model)
	if result.quitting {
		return
	}

	switch result.cursor {
	case 0:
		cmd.RunInteractive()
	case 1:
		cmd.InputAccounts()
	case 2:
		cmd.Init()
	case 3:
		fmt.Println("\nCarpeDiem v1.0")
		fmt.Println("Go + Rod browser automation")
		fmt.Println("github.com/KiddieLamer/carpediem\n")
	case 4:
		return
	}
}
