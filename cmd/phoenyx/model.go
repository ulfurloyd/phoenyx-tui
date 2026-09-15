package main

import (
	"os"
	"os/exec"
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

type commandFinishedMsg struct {
	err error
}

func newNvimCommand(dir string) *exec.Cmd {
	cmd := exec.Command("nvim", dir)
	cmd.Dir = dir
	return cmd
}

func initialModel() model {
	return model{
		commands: []command{
			{
				name: "open phoenyx configs",
				cmd:  newNvimCommand(os.ExpandEnv("$HOME/.local/share/chezmoi")),
			},
			{
				name: "open homelab",
				cmd:  newNvimCommand(os.ExpandEnv("$HOME/Projects/phoenyxlab")),
			},
			{
				name: "shell",
				cmd:  exec.Command("bash"),
			},
		},
	}
}
