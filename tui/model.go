package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"straico-cli.tylery.com/m/v2/cmd"
	"strconv"
	"strings"
)

const gap = "\n\n"

// LLMResponseMsg represents a message containing the LLM response
type LLMResponseMsg struct {
	response string
	err      error
}
type Messages []string

func (m Messages) Render(width int) string {
	return lipgloss.NewStyle().Width(width).Render(strings.Join(m, "\n"))
}

type Model struct {
	viewport      viewport.Model
	convSelection int
	Conversations Conversations
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	err           error
	config        *cmd.ConfigFile
}

func NewModel(config *cmd.ConfigFile) Model {
	ta := textarea.New()
	ta.Placeholder = "Ask the LLM... (" + config.Model + ")" + " "
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

	// Add a subtle style to indicate scrollable area with full border
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	conversations := make(Conversations, 9)
	for i := 0; i < 9; i++ {
		conversations.InitConversation(i)
	}
	conversations.LoadConversations()

	return Model{
		textarea:      ta,
		convSelection: 0,
		Conversations: conversations,
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:           nil,
		config:        config,
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

	c := &m.Conversations[m.convSelection]

	if len(c.Messages) == 0 {
		m.viewport.SetContent(`Welcome to Straico Cli!
Type a message and press Enter to send.


Use ↑/↓ arrows to scroll through chat history.`)
	}

	m.textarea.Placeholder = "Ask the LLM... (" + m.config.Model + ")" +
		" " + "(%" + strconv.Itoa(int(m.viewport.ScrollPercent()*100)) + ")" +
		" " + "(" + strconv.Itoa(m.convSelection+1) + ")"

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		// Leave more room for the chat history
		h := m.viewport.Style.GetVerticalFrameSize()
		m.viewport.Height = msg.Height - m.textarea.Height() - h - 1

		if len(c.Messages) > -1 {
			m.viewport.SetContent(c.Messages.Render(m.viewport.Width - 6))
		}

	case LLMResponseMsg:
		if msg.err != nil {
			c.Messages = append(c.Messages, m.senderStyle.Render("Error: ")+msg.err.Error())
		} else {
			c.Messages = append(c.Messages, m.senderStyle.Render("LLM: ")+msg.response)
		}
		m.viewport.SetContent(c.Messages.Render(m.viewport.Width - 6))
		m.viewport.GotoBottom()
		err := m.Conversations.SaveConversations()
		if err != nil {
			m.err = err
			return m, nil
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			userMessage := m.textarea.Value()
			c.PromptHistory = append(c.PromptHistory, userMessage)
			c.Messages = append(c.Messages, m.senderStyle.Render("You: ")+userMessage)
			m.viewport.SetContent(c.Messages.Render(m.viewport.Width - 6))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			return m, func() tea.Msg {
				response, err := m.config.Prompt.Request(m.config.Key, userMessage, c.PromptHistory)
				if err != nil {
					return LLMResponseMsg{err: err}
				}
				llmResponse := response.Data.Completions[m.config.Prompt.Model[0]].Completion.Choices[0].Message.Content
				return LLMResponseMsg{response: llmResponse}
			}
		case tea.KeyF1, tea.KeyF2, tea.KeyF3, tea.KeyF4, tea.KeyF5, tea.KeyF6, tea.KeyF7, tea.KeyF8, tea.KeyF9:
			m.convSelection = int(tea.KeyF1 - msg.Type)
			c = &m.Conversations[m.convSelection]
			m.viewport.SetContent(c.Messages.Render(m.viewport.Width - 6))
		case tea.KeyF12:
			m.Conversations.InitConversation(m.convSelection)
			m.Conversations.SaveConversations()
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return m.viewport.View() + gap + m.textarea.View()
}
