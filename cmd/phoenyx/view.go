package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

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
