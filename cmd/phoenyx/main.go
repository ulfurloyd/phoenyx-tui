package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
