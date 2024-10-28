package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/magnuswahlstrand/cloudflare-examples/durable-object-websockets-chatroom/client"
)

func main() {
	var nickname string
	flag.StringVar(&nickname, "nickname", "anonymous", "Nickname to use in chat")
	var room string
	flag.StringVar(&room, "room", "", "Room to join, e.g. 'stockholm'")
	var clientId string
	flag.StringVar(&clientId, "clientId", "", "Client ID to use in chat")
	var host string
	flag.StringVar(&host, "host", "localhost:8787", "Host to connect to")
	flag.Parse()

	if room == "" {
		fmt.Println("Please provide a room")
		return
	}

	if clientId == "" {
		clientId = uuid.New().String()
	}

	client.Run(host, room, clientId, nickname)
}
