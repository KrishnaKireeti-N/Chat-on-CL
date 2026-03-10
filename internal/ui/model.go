package ui

import (
	"fmt"
	"strings"

	"github.com/KrishnaKireeti-N/Chat-on-CL/internal/node"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type (
	errMsg error

	user_message struct {
		header  *node.Header
		content string
	}
	user_left struct {
		header *node.Header
	}

	quit struct {
		end string
	}
)

type model struct {
	client      *node.Client
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func InitialModel(u *node.Client) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	// ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		client:      u,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(u.Senderstyle)),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, recieveMsg(m.client))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case quit:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())

			return m, tea.Batch(tea.Quit, sendMsg(m.client, "\r"))
		case tea.KeyEnter:
			body := m.textarea.Value()

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+body)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// go sendMsg(m.user, body)

			return m, tea.Batch(tiCmd, vpCmd, sendMsg(m.client, body))
		}
	case user_message:
		h := msg.header

		senderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.client.Users[h.Username].Senderstyle))
		m.messages = append(m.messages, senderStyle.Render(fmt.Sprintf("%v: ", h.Username))+msg.content)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()

		return m, tea.Batch(tiCmd, vpCmd, recieveMsg(m.client))

	case user_left:
		h := msg.header

		m.messages = append(m.messages, fmt.Sprintf("--- %v left the chat", h.Username))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()

		return m, tea.Quit
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

/* >>> Helper Functions <<< */
func sendMsg(client *node.Client, msg string) func() tea.Msg {
	return func() tea.Msg {
		switch msg {
		case "":
			return nil
		case "\r":
			client.SendTimeout(msg)
			return nil
		}
		client.Send(msg)
		return nil
	}
}

func recieveMsg(client *node.Client) tea.Cmd {
	return func() tea.Msg {
		h, msg := client.Recieve()
		switch msg {
		case "\r":
			return user_left{header: h}
		}

		return user_message{
			header:  h,
			content: msg,
		}
	}
}
