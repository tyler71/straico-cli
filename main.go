/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tyler71/straico-cli/m/v0/cmd"
	"github.com/tyler71/straico-cli/m/v0/tui"
	"log"
	"os"
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
	if state.CoinUsage > 0 {
		os.Stderr.Write([]byte(strconv.FormatFloat(state.CoinUsage, 'f', 2, 64) + " coins used during session.\n"))
	}
}
