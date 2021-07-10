package sockets

import (
	"fmt"
	"log"
)
type Hub struct {
	clients    map[*Client]string
	register   chan *Register
	unregister chan *Client
}

func CreateHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]string),
		register:   make(chan *Register),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) findUsernameInHub(username string) bool  {
	for client := range hub.clients {
		if username == hub.clients[client] {
			return true
		}
	}
	return false
}

// TODO: list of active users

// listens to register and unregister channel
func (hub *Hub) RunHub() {
	for {
		select {
		case register := <-hub.register:
			hub.clients[register.client] = register.username
			action := &Action{
				Event: "joined",
			}
			register.client.actions <- action

			res := fmt.Sprintf("%s joined", register.username)
			log.Println(res)
			register.client.sendToEveryone(res)

		case client := <-hub.unregister:
			if username, registered := hub.clients[client]; registered {
				res := fmt.Sprintf("%s left", username)
				log.Println(res)
				// send a message to everyone except the client 
				for c := range client.hub.clients {
					if c != client {
						c.feed <- res
					}
				}
				
				delete(hub.clients, client)
			}
		}
	}
}
