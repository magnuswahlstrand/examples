package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type (
	errMsg error
)

type model struct {
	websocket      *websocket.Conn
	err            error
	messageHandler MessageHandler
	textInput      textinput.Model
	viewport       viewport.Model
	ready          bool
	title          string
	messages       []string
}

func initialModel(c *websocket.Conn, handler MessageHandler, title string, messages []string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:      ti,
		err:            nil,
		websocket:      c,
		messageHandler: handler,
		messages: append([]string{
			"Welcome to the chat!",
			"Type /nickname <name> to change your nickname.",
			"",
		}, messages...),
		title: title,
	}
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

type WebSocketMessage struct {
	Content string
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:

			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleInput()
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.separator())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			// TODO: How do I make this a common method?
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	case WebSocketMessage:
		m.messages = append(m.messages, m.messageHandler.formatInput(msg.Content))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, nil

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// Inspired by https://github.com/charmbracelet/bubbletea/blob/master/examples/chat/main.go
func (m model) handleInput() (tea.Model, tea.Cmd) {
	v := m.textInput.Value()
	if v == "" {
		// Don't send empty messages.
		return m, nil
	}

	var message string
	if strings.HasPrefix(v, "/") {
		switch {
		case strings.HasPrefix(v, "/nickname "):
			newNickname := strings.TrimPrefix(v, "/nickname ")
			// TODO: Refactor
			message = fmt.Sprintf("%s changed nickname to %s", m.messageHandler.nickname, newNickname)
			m.messageHandler.updateNickname(newNickname)
		default:
			// TODO: Handle in a better way
			fmt.Printf(color.RedString("Unknown command %q\n"), v)
		}
	} else {
		message = m.messageHandler.formatOutput(v)
	}
	err := m.websocket.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
		return m, nil
	}

	m.textInput.Reset()
	return m, nil
}

func (m model) redrawContent() {
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()
}

func (m model) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) separator() string {
	line := strings.Repeat("─", max(0, m.viewport.Width))
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), m.separator(), m.textInput.View())
}
