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
	content := ""

	if m.width < 40 || m.height < 20 {
		return tea.NewView("terminal too small")
	}

	if m.help {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			m.renderHelp(),
		)
	} else {
		commands := m.renderCommands()
		description := m.renderDescription()
		search := ""
		if m.searching {
			search = "/ " + m.search
		}
		help := helpStyle.Render("press `?` for keybinds")

		contentHeight := lipgloss.Height(title) +
			lipgloss.Height(commands) +
			lipgloss.Height(help) +
			lipgloss.Height(description) +
			lipgloss.Height(search) +
			2
		remaining := m.height - contentHeight - 4

		spacer := ""
		if remaining > 0 {
			spacer = strings.Repeat("\n", remaining)
		}

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			search,
			commands,
			description,
			spacer,
			help,
		)
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

func (m model) renderDescription() string {
	filtered := m.filteredCommands()

	if len(filtered) <= 0 {
		return ""
	}

	description := filtered[m.cursor].description

	return helpStyle.Render(description)
}

func (m model) renderHelp() string {
	const helpText = "normal mode keybinds:" +
		"\n" +
		"j/k or ↑/↓\tnavigate up/down" +
		"\n" +
		"g\t\t\t jump to top" +
		"\n" +
		"G\t\t\t jump to bottom" +
		"\n" +
		"enter\t\t select menu entry" +
		"\n" +
		"?\t\t\t toggle help mode" +
		"\n" +
		"q\t\t\t quit" +
		"\n" +
		"/\t\t\t enter search mode" +
		"\n" +
		"\n" +
		"search mode keybinds:" +
		"\n" +
		"↑/↓ or j/k\tnavigate up/down" +
		"\n" +
		"enter\t\t select highlighted result" +
		"\n" +
		"esc/q\t\t exit search"
	return helpStyle.Render(helpText)
}
