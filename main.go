package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	fetchLimit = 200
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	defer f.Close()
	m := newModel()
	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()
	// p.EnableMouseCellMotion()
	// defer p.DisableMouseCellMotion()
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
