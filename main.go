package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
	"os/exec"
)

type command struct {
	name string
	cmd  *exec.Cmd
}

type model struct {
	commands []command
	cursor   int
}

type shellFinishedMsg struct {
	err error
}

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
	}
	return m, nil
}

func (m model) View() tea.View {
	s := "phoenyx\n\n"

	for i, command := range m.commands {
		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		s += cursor + " " + command.name + "\n"
	}

	s += "\nq: quit"

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
