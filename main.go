/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"straico-cli.tylery.com/m/v2/cmd"
	"straico-cli.tylery.com/m/v2/tui"
)

func main() {
	configFile := cmd.Init()

	p := tea.NewProgram(
		tui.NewModel(configFile),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Capture mouse events
	)

	if err, _ := p.Run(); err != nil {
		log.Fatalln("Error running program:", err)
	}
}
