package sockets

import (
	"fmt"
	"log"
)

type Action struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

type Register struct {
	client *Client
	username string
}

func handleActions(client *Client, action *Action) {
	clientUsername, isRegistered := client.hub.clients[client]
	switch action.Event {
	case "join":
		username := action.Content
		userNameInHub := client.hub.findUsernameInHub(username)
		if userNameInHub {
			res := &Action{
				Event: "alreadyExists",
				Content: fmt.Sprintf("%s already exists", username),
			} 
			client.actions <- res
			log.Println(username, "already exits")
		} else {
			if isRegistered {
				client.hub.unregister <- client
			}
			client.hub.register <- &Register{client: client, username: action.Content}
		}		

	case "message":
		if isRegistered {
			formatted := fmt.Sprintf("%s: %s", clientUsername, action.Content)
			client.sendToEveryone(formatted)
			// TODO: private messages
		}
	default:
		// TODO: handle command not found
	}
}