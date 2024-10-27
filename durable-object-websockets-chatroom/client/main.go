package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

func main() {
	//Create Message Out
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	const room = "stockholm"

	u := url.URL{Scheme: "ws", Host: "localhost:8787", Path: "/ws", RawQuery: "room=" + room}
	log.Printf("connecting to %s", u.String())
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
			log.Printf("recv: %s", message)
		}

	}()

	go func() {
		username := "Magnus"
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			trimmed := strings.Trim(text, "\n")

			if strings.HasPrefix(trimmed, "/nickname ") {
				newUsername := strings.TrimPrefix(trimmed, "/nickname ")
				messageOut <- fmt.Sprintf("<%s> changed name to <%s>", username, newUsername)
				username = newUsername
				continue
			}

			messageOut <- fmt.Sprintf("<%s>: %s", username, trimmed)
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
