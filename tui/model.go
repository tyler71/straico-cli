package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tyler71/straico-cli/m/v0/cmd"
	"strconv"
	"strings"
)

const gap = "\n\n"

// LLMResponseMsg represents a message containing the LLM response
type LLMResponseMsg struct {
	response  string
	coinUsage float64
	err       error
}
type Messages []string

func (m Messages) Render(width int) string {
	return lipgloss.NewStyle().Width(width).Render(strings.Join(m, "\n"))
}

func NewModel(config *cmd.ConfigFile, state *State) *State {
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

	state.Textarea = ta
	state.ConvSelection = 0
	state.Conversations = conversations
	state.Viewport = vp
	state.SenderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	state.Config = *config
	return state

}

func (s State) Init() tea.Cmd {
	return textarea.Blink
}

func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	c := &s.Conversations[s.ConvSelection]

	s.Textarea, _ = s.Textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Viewport.Width = msg.Width
		s.Textarea.SetWidth(msg.Width)
		// Leave more room for the chat history
		h := s.Viewport.Style.GetVerticalFrameSize()
		s.Viewport.Height = msg.Height - s.Textarea.Height() - h - 1

		if len(c.Messages) > -1 {
			s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
		}

	case LLMResponseMsg:
		if msg.err != nil {
			c.Messages = append(c.Messages, s.SenderStyle.Render("Error: ")+msg.err.Error())
		} else {
			c.Messages = append(c.Messages, s.SenderStyle.Render("LLM: ")+msg.response)
			s.CoinUsage += msg.coinUsage
		}
		s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
		s.Viewport.HalfViewDown()
		//s.Viewport.GotoBottom()
		err := s.Conversations.SaveConversations()
		if err != nil {
			s.Err = err
			return s, nil
		}

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelDown, tea.MouseButtonWheelUp:
			s.Viewport, _ = s.Viewport.Update(msg)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			//coinUsageMessage := strconv.FormatFloat(s.CoinUsage, 'f', 2, 64) + " coins used during session"
			return s, tea.Quit
		case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
			s.Viewport, _ = s.Viewport.Update(msg)
		case tea.KeyEnd:
			s.Viewport.GotoBottom()
			s.Viewport, _ = s.Viewport.Update(msg)
		case tea.KeyHome:
			s.Viewport.GotoTop()
			s.Viewport, _ = s.Viewport.Update(msg)
		case tea.KeyEnter:
			userMessage := s.Textarea.Value()
			if strings.Trim(userMessage, " ") == "" {
				s.Textarea.Reset()
				return s, nil
			}
			c.PromptHistory = append(c.PromptHistory, userMessage)
			c.Messages = append(c.Messages, s.SenderStyle.Render("You: ")+userMessage)
			s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
			s.Textarea.Reset()
			s.Viewport.GotoBottom()
			s.Textarea.Placeholder = "Loading..."

			return s, func() tea.Msg {
				response, err := s.Config.Prompt.Request(s.Config.Key, userMessage, c.PromptHistory)
				if err != nil {
					return LLMResponseMsg{err: err}
				}
				llmResponse := response.Data.Completions[s.Config.Prompt.Model[0]].Completion.Choices[0].Message.Content
				coins := response.Data.OverallPrice.Total
				return LLMResponseMsg{response: llmResponse, coinUsage: coins}
			}
		case tea.KeyF1, tea.KeyF2, tea.KeyF3, tea.KeyF4, tea.KeyF5, tea.KeyF6, tea.KeyF7, tea.KeyF8, tea.KeyF9:
			s.ConvSelection = int(tea.KeyF1 - msg.Type)
			c = &s.Conversations[s.ConvSelection]
			s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
		//	ShiftLeft and ShiftRight used to
		case tea.KeyShiftLeft:
			if s.ConvSelection-1 >= 0 {
				t := s.Conversations[s.ConvSelection-1]
				s.Conversations[s.ConvSelection-1] = s.Conversations[s.ConvSelection]
				s.Conversations[s.ConvSelection] = t
				s.ConvSelection--
				c = &s.Conversations[s.ConvSelection]
				s.Conversations.SaveConversations()
				s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
			}
		case tea.KeyShiftRight:
			if s.ConvSelection+1 < len(s.Conversations) {
				t := s.Conversations[s.ConvSelection+1]
				s.Conversations[s.ConvSelection+1] = s.Conversations[s.ConvSelection]
				s.Conversations[s.ConvSelection] = t
				s.ConvSelection++
				c = &s.Conversations[s.ConvSelection]
				s.Conversations.SaveConversations()
				s.Viewport.SetContent(c.Messages.Render(s.Viewport.Width - 6))
			}
		case tea.KeyF12:
			s.Conversations.InitConversation(s.ConvSelection)
			s.Conversations.SaveConversations()
		default:
			var command tea.Cmd
			//s.textarea, command = s.textarea.Update(msg)
			return s, command
		}

	default:
		return s, nil
	case error:
		s.Err = msg
		return s, nil
	}

	s.Textarea.Placeholder = "Ask the LLM... (" + s.Config.Model + ")" +
		" " + "(%" + strconv.Itoa(int(s.Viewport.ScrollPercent()*100)) + ")" +
		" " + "(" + strconv.Itoa(s.ConvSelection+1) + ")" +
		" " + "(" + strconv.FormatFloat(s.CoinUsage, 'f', 2, 64) + ")"
	if len(c.Messages) == 0 {
		s.Viewport.SetContent(`Welcome to Straico Cli!
Type a message and press Enter to send.


Use ↑/↓ arrows to scroll through chat history.`)
	}

	return s, nil
	//return s, tea.Batch(tiCmd, vpCmd)
}

func (s State) View() string {
	return s.Viewport.View() + gap + s.Textarea.View()
}
