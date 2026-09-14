package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m model) Init() tea.Cmd {
	return nil
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
					func(err error) tea.Msg {
						return commandFinishedMsg{err: err}
					},
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
	case commandFinishedMsg:
		m.searching = false
		m.search = ""
		m.cursor = m.previousCursor
	}
	return m, nil
}
