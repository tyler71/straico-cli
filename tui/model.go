package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"straico-cli.tylery.com/m/v2/cmd"
	p "straico-cli.tylery.com/m/v2/prompt"
	"strings"
)

const gap = "\n\n"

// LLMResponseMsg represents a message containing the LLM response
type LLMResponseMsg struct {
	response string
	err      error
}

type Model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	config      *cmd.ConfigFile
	prompt      p.Prompt
}

func NewModel(config *cmd.ConfigFile) Model {
	ta := textarea.New()
	ta.Placeholder = "Ask the LLM... (" + config.Model + ")"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 2000

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	// Initialize viewport with scrolling enabled
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to Straico Chat!
Type a message and press Enter to send.

Use ↑/↓ arrows or mouse wheel to scroll through chat history.`)

	// Enable mouse wheel scrolling
	vp.MouseWheelEnabled = true

	// Add a subtle style to indicate scrollable area with full border
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		config:      config,
		prompt:      cmd.Config.Prompt,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		// Leave more room for the chat history
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap) - 2 // -2 for viewport borders

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width - 2).Render(strings.Join(m.messages, "\n"))) // -2 for viewport padding
		}
		m.viewport.GotoBottom()

	case LLMResponseMsg:
		if msg.err != nil {
			m.messages = append(m.messages, m.senderStyle.Render("Error: ")+msg.err.Error())
		} else {
			m.messages = append(m.messages, m.senderStyle.Render("Assistant: ")+msg.response)
		}
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width - 2).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			userMessage := m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+userMessage)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width - 2).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			return m, func() tea.Msg {
				response, err := m.prompt.Request(m.config.Key, userMessage)
				if err != nil {
					return LLMResponseMsg{err: err}
				}
				llmResponse := response.Data.Completions[m.prompt.Model[0]].Completion.Choices[0].Message.Content
				return LLMResponseMsg{response: llmResponse}
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
