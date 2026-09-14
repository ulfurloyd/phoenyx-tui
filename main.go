package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
)

type model struct {
	choices []string
	cursor  int
	status  string
}

func initialModel() model {
	return model{
		choices: []string{
			"SSH nyx",
			"SSH hermes",
			"Shell",
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
			m.status = "selected: " + m.choices[m.cursor]
		case "q", "ctrl+c":
			return m, tea.Quit
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
