/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"straico-cli.tylery.com/m/v2/cmd"
	"straico-cli.tylery.com/m/v2/tui"
)

func main() {
	cmd.Init()
	configFile, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Create program with alternate screen buffer to prevent terminal scrolling
	p := tea.NewProgram(
		tui.NewModel(configFile),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Capture mouse events
	)

	if err := p.Start(); err != nil {
		log.Fatal("Error running program:", err)
		os.Exit(1)
	}
}
