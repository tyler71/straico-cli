package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/tyler71/straico-cli/m/v0/cmd"
)

type State struct {
	Viewport      viewport.Model
	ConvSelection int
	Conversations Conversations
	Textarea      textarea.Model
	SenderStyle   lipgloss.Style
	Err           error
	Config        cmd.ConfigFile
	CoinUsage     float64
}
