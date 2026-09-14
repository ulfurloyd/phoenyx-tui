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
	commands []command
	cursor   int
	width    int
	height   int
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
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down", "ctrl+n":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			return m, tea.ExecProcess(
				m.commands[m.cursor].cmd,
				nil,
			)
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
	help := helpStyle.Render("↑/↓ or j/k · enter · q")

	contentHeight := lipgloss.Height(title) +
		lipgloss.Height(commands) +
		lipgloss.Height(help) +
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

	for i, command := range m.commands {
		cursor := " "
		name := command.name

		if i == m.cursor {
			cursor = cursorStyle.Render(">")
			name = selectedStyle.Render(name)
		}

		s += cursor + " " + name + "\n\n"
	}

	return s
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
