package main

import (
	"fmt"
	"os"

	"github.com/aymanbagabas/hknui/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hknui",
		Short: "Hacker News TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}

			defer f.Close()
			m := ui.NewModel()
			p := tea.NewProgram(m)
			p.EnterAltScreen()
			defer p.ExitAltScreen()
			// p.EnableMouseCellMotion()
			// defer p.DisableMouseCellMotion()
			if err := p.Start(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
			return cmd.Help()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
