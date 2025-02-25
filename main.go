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
	"strconv"
)

func main() {
	configFile := cmd.Init()

	state := tui.State{}
	p := tea.NewProgram(
		tui.NewModel(configFile, &state),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Capture mouse events
	)

	if _, err := p.Run(); err != nil {
		log.Fatalln("Error running program:", err)
	}
	os.Stdout.Write([]byte(strconv.FormatFloat(state.CoinUsage, 'f', 2, 64) + " coins used during session.\n"))
}
