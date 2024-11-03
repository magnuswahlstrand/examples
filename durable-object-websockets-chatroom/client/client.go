package client

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/goodies/stringslice"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

type MessageHandler struct {
	clientId string
	nickname string
}

func (m *MessageHandler) formatOutput(message string) string {
	return fmt.Sprintf("%s:%s: %s", m.clientId, m.nickname, message)
}

var styleOwnMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
var styleOwnNickname = styleOwnMessage.Bold(true).Width(10)

var styleMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var styleOtherNickname = styleMessage.Bold(true).Width(10)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func (m *MessageHandler) formatInput(s string) string {
	//return s
	if strings.HasPrefix(s, "err") {
		return errorStyle.Render(strings.TrimPrefix(s, "err:"))
	}

	parts := strings.SplitN(s, ":", 3)

	if len(parts) < 3 {
		return s
	}
	//return fmt.Sprintf("%s ||| %s", parts[1], parts[2]) //, len(parts), parts[0], parts[1])

	if parts[0] == m.clientId {
		return styleOwnNickname.Render(parts[1]) + "|" + styleOwnMessage.Render(parts[2])
	}

	return styleOtherNickname.Render(parts[1]) + "|" + styleMessage.Render(parts[2])
}

func (m *MessageHandler) updateNickname(newNickname string) {
	m.nickname = newNickname
}

func Run(host string, room string, clientId string, nickname string) {
	//messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Yes, this is silly :-)
	schemePostfix := "s"
	if strings.Contains(host, "localhost") {
		schemePostfix = ""
	}
	u1 := url.URL{Scheme: "http" + schemePostfix, Host: host, Path: "/history", RawQuery: "room=" + room}
	u := url.URL{Scheme: "ws" + schemePostfix, Host: host, Path: "/ws", RawQuery: "room=" + room}

	resp0, err := http.DefaultClient.Get(u1.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp0.Body.Close()

	// Parse the response
	var messages []string
	if err := json.NewDecoder(resp0.Body).Decode(&messages); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("connecting to %s\n", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	defer resp.Body.Close()

	messageHandler := MessageHandler{
		clientId: clientId,
		nickname: nickname,
	}

	formattedMessages := stringslice.Map(messages, func(s string) string {
		return messageHandler.formatInput(s)
	})

	//When the program closes the connection
	defer c.Close()

	p := tea.NewProgram(initialModel(c, messageHandler, fmt.Sprintf("# %s", room), formattedMessages))

	// TODO: Refactor
	//done := make(chan struct{})
	go func() {
		//defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			p.Send(ReceivedWebSocketMessage{Content: string(message)})
		}

	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
