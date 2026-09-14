package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type command struct {
	name string
	cmd  *exec.Cmd
}

type model struct {
	commands       []command
	cursor         int
	width          int
	height         int
	searching      bool
	search         string
	previousCursor int
}

type shellFinishedMsg struct {
	err error
}

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("6")).
	PaddingLeft(1)
var cursorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("13"))
var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("12")).
	MarginLeft(1).
	Padding(1)
var helpStyle = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("7"))
var selectedStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("15"))

func initialModel() model {
	return model{
		commands: []command{
			{
				name: "open phoenyx configs",
				cmd:  exec.Command("nvim", os.ExpandEnv("$HOME/.local/share/chezmoi")),
			},
			{
				name: "open homelab",
				cmd:  exec.Command("nvim", os.ExpandEnv("$HOME/Projects/phoenyxlab")),
			},
			{
				name: "shell",
				cmd:  exec.Command("bash"),
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	filtered := m.filteredCommands()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.searching {
			switch msg.String() {
			case "backspace":
				if len(m.search) > 0 {
					m.search = m.search[:len(m.search)-1]
					m.cursor = 0
				}
			case "esc":
				m.searching = false
				m.search = ""
				m.cursor = m.previousCursor
			case "enter":
				if len(filtered) == 0 {
					return m, nil
				}
				m.searching = false
				return m, tea.ExecProcess(
					filtered[m.cursor].cmd,
					nil,
				)
			case "down", "ctrl+n":
				if m.cursor < len(filtered)-1 {
					m.cursor++
				}
			case "up", "ctrl+p":
				if m.cursor > 0 {
					m.cursor--
				}
			case "ctrl+c":
				return m, tea.Quit
			default:
				if msg.Text != "" {
					m.search += msg.Text
					m.cursor = 0
				}
			}
		}
		switch msg.String() {
		case "j", "down", "ctrl+n":
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if len(filtered) == 0 {
				return m, nil
			}

			return m, tea.ExecProcess(
				filtered[m.cursor].cmd,
				nil,
			)
		case "g":
			m.cursor = 0
		case "G":
			if len(filtered) > 0 {
				m.cursor = len(filtered) - 1
			}
		case "/":
			m.searching = true
			m.search = ""
			m.previousCursor = m.cursor
			m.cursor = 0
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() tea.View {
	title := titleStyle.Render("phoenyx")
	commands := m.renderCommands()
	search := ""
	if m.searching {
		search = "/ " + m.search
	}
	help := helpStyle.Render("↑/↓ or j/k · enter · q")

	contentHeight := lipgloss.Height(title) +
		lipgloss.Height(commands) +
		lipgloss.Height(help) +
		lipgloss.Height(search) +
		2
	remaining := m.height - contentHeight - 4

	spacer := ""
	if remaining > 0 {
		spacer = strings.Repeat("\n", remaining)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		search,
		commands,
		spacer,
		help,
	)

	if m.width < 40 || m.height < 20 {
		return tea.NewView("terminal too small")
	}

	box := boxStyle.Width(m.width - 2).Height(m.height - 1)
	content = box.Render(content)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m model) renderCommands() string {
	s := ""

	for i, command := range m.filteredCommands() {
		cursor := " "
		name := command.name

		if i == m.cursor {
			cursor = cursorStyle.Render(">")
			name = selectedStyle.Render(name)
		}

		s += cursor + " " + name + "\n"
	}

	return s
}

func (m model) filteredCommands() []command {
	if m.search == "" {
		return m.commands
	}

	var filtered []command
	for _, command := range m.commands {
		if strings.Contains(strings.ToLower(command.name), strings.ToLower(m.search)) {
			filtered = append(filtered, command)
		}
	}

	return filtered
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
