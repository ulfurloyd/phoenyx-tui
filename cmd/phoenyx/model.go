package main

import (
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
)

type command struct {
	name        string
	description string
	cmd         *exec.Cmd
}

type model struct {
	commands       []command
	cursor         int
	width          int
	height         int
	searching      bool
	help           bool
	search         string
	previousCursor int
	currentTime    time.Time
}

type commandFinishedMsg struct {
	err error
}

type tickMsg time.Time

func newNvimCommand(dir string) *exec.Cmd {
	cmd := exec.Command("nvim", dir)
	cmd.Dir = dir
	return cmd
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	return model{
		commands: []command{
			{
				name:        "open phoenyx configs",
				cmd:         newNvimCommand(os.ExpandEnv("$HOME/.local/share/chezmoi")),
				description: "open the phoenyx dotfiles directory in nvim",
			},
			{
				name:        "open homelab",
				cmd:         newNvimCommand(os.ExpandEnv("$HOME/Projects/phoenyxlab")),
				description: "open the phoenyxlab directory in nvim",
			},
			{
				name:        "shell",
				cmd:         exec.Command("bash"),
				description: "open an interactive shell",
			},
		},
		currentTime: time.Now(),
	}
}
