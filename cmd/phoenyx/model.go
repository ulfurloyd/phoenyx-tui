package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	cpuUsage       float64
	memUsage       float64
	ssdUsage       float64
	hddUsage       float64
	prevTotal      int64
	prevIdle       int64
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

func readCPUStat() ([]byte, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, fmt.Errorf("could not read /proc/stat: %w", err)
	}

	for line := range strings.Lines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if fields[0] == "cpu" {
			for _, field := range fields[1:] {
				fieldVal, err := strconv.Atoi(field)
				if err != nil {
				}
			}
		}
	}

	return data, nil
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
