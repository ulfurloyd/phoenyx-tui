package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os/exec"
)

type model struct {
	choices []string
	cursor  int
	status  string
}

type shellFinishedMsg struct {
	err error
}

func shell() tea.Cmd {
	return tea.ExecProcess(
		exec.Command("bash"),
		func(err error) tea.Msg {
			return shellFinishedMsg{err: err}
		},
	)
}

func initialModel() model {
	return model{
		choices: []string{
			"ssh nyx",
			"ssh hermes",
			"shell",
		},
		status: "",
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.choices[m.cursor] == "shell" {
				return m, shell()
			}

			m.status = "Selected:" + m.choices[m.cursor]
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case shellFinishedMsg:
		if msg.err != nil {
			m.status = "shell exited with an error"
		} else {
			m.status = "shell exited successfully"
		}

	}

	return m, nil
}

func (m model) View() tea.View {
	s := "phoenyx\n\n"

	for i, choice := range m.choices {
		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		s += cursor + " " + choice + "\n"
	}

	s += "\n" + m.status

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
