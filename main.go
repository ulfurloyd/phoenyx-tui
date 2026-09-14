package main

import (
	"fmt"
	"os"
	"os/exec"

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

var titleStyle = lipgloss.NewStyle().Bold(true)
var selectedStyle = lipgloss.NewStyle().Bold(true)
var helpStyle = lipgloss.NewStyle().Faint(true)
var boxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).MarginLeft(1).PaddingLeft(2)

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
	content := titleStyle.Render("phoenyx") + "\n\n"

	content += m.renderCommands()
	content += "\n" + helpStyle.Render("↑/↓ or j/k · enter · q")

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

		if i == m.cursor {
			cursor = ">"
		}

		if i == m.cursor {
			s += cursor + " " + selectedStyle.Render(command.name) + "\n"
		} else {
			s += cursor + " " + command.name + "\n"
		}
	}

	return s
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
