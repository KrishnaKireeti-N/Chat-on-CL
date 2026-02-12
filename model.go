package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// API Stuff (what other parts of the program (Client) can(technically should) access

// Color enum
type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
	Purple
	Cyan
	White
)

func (c Color) String() string {
	var color string
	switch c {
	case Red:
		color = "#FF0000"
	case Green:
		color = "#00FF00"
	case Blue:
		color = "#0000FF"
	case Yellow:
		color = "FFFF00"
	case Purple:
		color = "800080"
	case Cyan:
		color = "48D1CC" // mediumturquoise
	case White:
		color = "FFFFFF"
	}
	return color
}

// ----------------------------------------------------------------------------------

const gap = "\n\n"

type (
	errMsg  error
	message struct {
		sender  string
		content string
	}
	quit struct {
		end string
	}
)

type model struct {
	user        User
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel(u User) model {
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
		user:        u,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(user.senderstyle.String())),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, recieveMsg(m.user))
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

			return m, tea.Batch(tea.Quit, sendMsg(m.user, "0"))
		case tea.KeyEnter:
			body := m.textarea.Value()

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+body)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// go sendMsg(m.user, body)

			return m, tea.Batch(tiCmd, vpCmd, sendMsg(m.user, body))
		}
	case message:
		if msg.content == "0" {
			return m, tea.Quit
		}
		m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%v: ", msg.sender))+msg.content)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()

		return m, tea.Batch(tiCmd, vpCmd, recieveMsg(m.user))

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
func sendMsg(user User, msg string) func() tea.Msg {
	return func() tea.Msg {
		user.Send(msg)
		return nil
	}
}

func recieveMsg(user User) tea.Cmd {
	return func() tea.Msg {
		msg := user.Recieve()
		if msg == "0" {
			return quit{end: "Connection Ended!"}
		}
		return message{
			sender:  user.conn.RemoteAddr().String(),
			content: msg,
		}
	}
}
