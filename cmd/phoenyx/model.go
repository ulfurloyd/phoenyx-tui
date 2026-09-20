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
	cpuInitialized bool
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

func readCPUStat() (int64, int64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, fmt.Errorf("could not read /proc/stat: %w", err)
	}

	found := false
	var total int64
	var idle int64

	for line := range strings.Lines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if fields[0] == "cpu" {
			found = true
			for i, field := range fields[1:] {
				fieldVal, err := strconv.Atoi(field)
				if err != nil {
					return 0, 0, fmt.Errorf("could not parse CPU stat: %w", err)
				}

				total += int64(fieldVal)

				if i == 3 {
					idle = int64(fieldVal)
				}
			}
		}
	}

	if !found {
		return 0, 0, fmt.Errorf("aggregate CPU stats not found")
	}

	return total, idle, nil
}

func calculateCPUUsage(previousTotal int64, previousIdle int64, currentTotal int64, currentIdle int64) float64 {
	totalDelta := currentTotal - previousTotal
	if totalDelta == 0 {
		return 0
	}
	idleDelta := currentIdle - previousIdle
	cpuUsage := (1 - float64(idleDelta)/float64(totalDelta)) * 100

	return cpuUsage
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
