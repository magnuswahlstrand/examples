package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
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

func (m *MessageHandler) formatInput(s string) string {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 3 {
		return s
	}

	if parts[0] == m.clientId {
		return fmt.Sprintf("%s: %s", color.HiGreenString(parts[1]), color.GreenString(parts[2]))
	}

	return fmt.Sprintf("%s: %s", color.HiWhiteString(parts[1]), color.WhiteString(parts[2]))
}

func (m *MessageHandler) updateNickname(newNickname string) {
	m.nickname = newNickname
}

func main() {
	// Flags for nickname and room
	var nickname string
	flag.StringVar(&nickname, "nickname", "anonymous", "Nickname to use in chat")
	var room string
	flag.StringVar(&room, "room", "", "Room to join, e.g. 'stockholm'")
	var clientId string
	flag.StringVar(&clientId, "clientId", "", "Client ID to use in chat")
	flag.Parse()

	if room == "" {
		fmt.Println("Please provide a room")
		return
	}

	if clientId == "" {
		clientId = uuid.New().String()
	}

	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8787", Path: "/ws", RawQuery: "room=" + room}
	fmt.Printf("connecting to %s\n", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	defer resp.Body.Close()

	//var result map[string]interface{}
	//if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	//	log.Fatal("decode:", err)
	//}

	messageHandler := MessageHandler{
		clientId: clientId,
		nickname: nickname,
	}

	//When the program closes the connection
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			fmt.Printf("%s\n", messageHandler.formatInput(string(message)))
		}

	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			trimmed := strings.Trim(text, "\n")

			if strings.HasPrefix(trimmed, "/") {
				switch {
				case strings.HasPrefix(trimmed, "/nickname "):
					newNickname := strings.TrimPrefix(trimmed, "/nickname ")
					messageOut <- fmt.Sprintf("<%s> changed name to <%s>", nickname, newNickname)
					messageHandler.updateNickname(newNickname)
					continue
				default:
					fmt.Printf(color.RedString("Unknown command %q\n"), trimmed)
					continue
				}
			}

			messageOut <- messageHandler.formatOutput(trimmed)
		}
	}()

	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
