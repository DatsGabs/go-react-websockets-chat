package sockets

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
 ReadBufferSize:  1024,
 WriteBufferSize: 1024,
 CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	feed chan string
	actions chan *Action
	conn *websocket.Conn
	hub  *Hub
}


func (client *Client) sendToEveryone(message string)  {
	for c := range client.hub.clients {
		c.feed <- message
	}
}

// listens to any changes in the feed and action channel
func (client *Client) feedHandler() {
	for {
		fmt.Println("feed")
		select {
		case action := <- client.actions:
			err := client.conn.WriteJSON(action)

			if err != nil {
				log.Printf("err: %s", err)
				return
			}

		case message  := <- client.feed:
			res := &Action {
				Event: "message",
				Content: message,
			}
			err := client.conn.WriteJSON(res)
			if err != nil {
				log.Printf("err: %s", err)
				return
			}
		}
	}
}

// reads if there was a message sent in the connection
func (client *Client) reader() {
	defer func() {	
		client.hub.unregister <- client
		client.conn.Close()	
	}()

	for {
		var action Action

		err := client.conn.ReadJSON(&action)
		if err != nil {
			log.Printf("err: %s", err)
			break
		}

		handleActions(client, &action)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request, hub *Hub) {
	// upgrades the HTTP server connection to the WebSocket protocol
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("err: %s", err)
		return
	}

	client := &Client{feed: make(chan string), actions: make(chan *Action), conn: ws, hub: hub}
		
	go client.reader()
	go client.feedHandler()
}